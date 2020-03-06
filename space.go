package balancer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/struckoff/SFCFramework/curve"
)

type Space struct {
	mu    sync.Mutex
	cells []cell
	cgs   []CellGroup
	sfc   curve.Curve
	tf    TransformFunc
	load  uint64
}

func NewSpace(sfc curve.Curve, tf TransformFunc) *Space {
	l := sfc.Length()
	cells := make([]cell, l)
	for i := range cells {
		cells[i] = newCell()
	}
	return &Space{
		mu:    sync.Mutex{},
		cells: cells,
		cgs:   []CellGroup{},
		sfc:   sfc,
		tf:    tf,
	}
}

//CellGroups returns a slice of all CellGroups in the space.
func (s *Space) CellGroups() []CellGroup {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cgs
}

//Cells returns a slice of all cells in the space.
func (s *Space) Cells() []cell {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cells
}

func (s *Space) TotalLoad() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.load
}

//TotalPower returns the sum of the all node powers in the space.
func (s *Space) TotalPower() (power float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for iter := range s.cgs {
		power += s.cgs[iter].node.Power().Get()
	}
	return
}

// SetGroups replace groups in the space.
func (s *Space) SetGroups(groups []CellGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cgs = groups
}

//Len() returns the number of CellGroups in the space.
func (s *Space) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cgs)
}

//AddNode adds a new node to the space.
func (s *Space) AddNode(n Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.addNode(n); err != nil {
		return err
	}
	return nil
}

func (s *Space) addNode(n Node) error {
	if len(s.cgs) == 0 {
		s.cgs = []CellGroup{NewCellGroup(n)}
		return nil
	}
	s.cgs = append(s.cgs, NewCellGroup(n))
	return nil
}

//AddData add data item to the space.
func (s *Space) AddData(d DataItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addData(d)
}

func (s *Space) addData(d DataItem) error {
	if len(s.cgs) == 0 {
		return errors.New("no nodes in the cluster")
	}
	size := s.sfc.DimSize()
	if s.tf == nil {
		return errors.New("transform function is not set")
	}
	coords, err := s.tf(d.Values(), size)
	cID, err := s.sfc.Encode(coords)
	if err != nil {
		return fmt.Errorf("item encoding error: %w", err)
	}
	if cID > uint64(len(s.cells)-1) {
		return errors.New("cell ID is larger that number of cells in the Space")
	}
	if err = s.cells[cID].add(d); err != nil {
		return err
	}
	s.load += s.cells[cID].load //TODO May be just sum CellGroup.load ?
	return nil
}

//Distribution returns representation of how DataItems distributes per nodes in the space.
func (s *Space) Distribution() DataDistribution {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make(DataDistribution, len(s.cgs))
	for i := range s.cgs {
		nd := NodeData{
			ID:    s.cgs[i].node.ID(),
			Items: []string{},
		}
		for _, c := range s.cgs[i].cells {
			ids := c.itemIDs()
			nd.Items = append(nd.Items, ids...)
		}
		res[i] = nd
	}
	return res
}

//func (s *Space) optimize() error {
//	cgs, err := s.of(s)
//	if err != nil {
//		return err
//	}
//	s.cgs = cgs
//	return nil
//}
