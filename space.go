package balancer

import (
	"github.com/pkg/errors"
	"math"
	"sort"
	"sync"

	"github.com/struckoff/SFCFramework/curve"
)

type Space struct {
	mu    sync.Mutex
	cells map[uint64]*cell
	cgs   []*CellGroup
	sfc   curve.Curve
	tf    TransformFunc
	load  uint64
}

func NewSpace(sfc curve.Curve, tf TransformFunc, nodes []Node) (*Space, error) {
	s := Space{
		mu:    sync.Mutex{},
		cells: map[uint64]*cell{},
		cgs:   []*CellGroup{},
		sfc:   sfc,
		tf:    tf,
	}
	for _, n := range nodes {
		err := s.addNode(n)
		if err != nil {
			return nil, err
		}
	}
	l := sfc.Length()
	r, err := splitCells(len(nodes), l)
	if err != nil {
		return nil, err
	}
	for i := range s.cgs {
		s.cgs[i].cRange = r[i]
	}
	return &s, nil
}

func splitCells(n int, l uint64) ([]Range, error) {
	if l < uint64(n) {
		return nil, errors.New("curve length must be larger than number of nodes")
	}

	s := float64(l) / float64(n)
	var c float64
	i := 0
	res := make([]Range, n)
	if n == 0 {
		return res, nil
	}
	for c < float64(l) {
		nc := c + s
		res[i] = Range{
			Min: uint64(math.Ceil(c)),
			Max: uint64(math.Ceil(nc)),
		}
		res[i].Len = res[i].Max - res[i].Min
		i++
		c = nc
	}
	return res, nil
}

//CellGroups returns a slice of all CellGroups in the space.
func (s *Space) CellGroups() []*CellGroup {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cgs
}

//Cells returns a slice of all cells in the space.
func (s *Space) Cells() []*cell {
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := make([]uint64, 0, len(s.cells))
	for k := range s.cells {
		ids = append(ids, k)
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
	res := make([]*cell, len(s.cells))
	for i, id := range ids {
		res[i] = s.cells[id]
	}
	return res
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
func (s *Space) SetGroups(groups []*CellGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cgs = groups
}

//Len returns the number of CellGroups in the space.
func (s *Space) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cgs)
}

//Capacity returns maximum number of cell which could be located in space
func (s *Space) Capacity() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sfc.Length()
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
	//TODO May be s.cgs should be map
	for iter := range s.cgs {
		if s.cgs[iter].ID() == n.ID() {
			s.cgs[iter].SetNode(n)
			return nil
		}
	}
	s.cgs = append(s.cgs, NewCellGroup(n))
	return nil
}

func (s *Space) SetNodes(ns []Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.setNodes(ns); err != nil {
		return err
	}
	return nil
}

func (s *Space) setNodes(ns []Node) error {
	s.cgs = nil
	for _, n := range ns {
		s.cgs = append(s.cgs, NewCellGroup(n))
	}
	return nil
}

func (s *Space) RemoveNode(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.removeNode(id); err != nil {
		return err
	}
	return nil
}
func (s *Space) removeNode(id string) error {
	for iter := range s.cgs {
		if s.cgs[iter].ID() == id {
			s.cgs = append(s.cgs[iter:], s.cgs[iter+1:]...)
			return nil
		}
	}
	return errors.Errorf("node(%s) not found", id)
}

//LocateData add data item to the space.
func (s *Space) LocateData(d DataItem) (Node, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.locateData(d)
}

func (s *Space) locateData(d DataItem) (Node, error) {
	if len(s.cgs) == 0 {
		return nil, errors.New("no nodes in the cluster")
	}
	cID, err := s.cellID(d)
	if err != nil {
		return nil, err
	}
	if _, ok := s.cells[cID]; !ok {
		cg, ok := s.findCellGroup(cID)
		if !ok {
			return nil, errors.New("unable to bind cell to cell group")
		}
		s.cells[cID] = newCell(cID, cg)
	}
	if err = s.cells[cID].add(d); err != nil {
		return nil, err
	}
	s.load += s.cells[cID].load //TODO May be just sum CellGroup.load ?
	return s.cells[cID].cg.Node(), nil
}

//cellID calculates the id of cell in space based on transform function and space filling curve.
func (s *Space) cellID(d DataItem) (uint64, error) {
	//size := s.sfc.DimensionSize()
	if s.tf == nil {
		return 0, errors.New("transform function is not set")
	}
	coords, err := s.tf(d.Values(), s.sfc)
	if err != nil {
		return 0, err
	}
	cID, err := s.sfc.Encode(coords)
	if err != nil {
		return 0, errors.Wrap(err, "item encoding error")
	}

	return cID, nil
}

func (s *Space) findCellGroup(cID uint64) (cg *CellGroup, ok bool) {
	for iter := range s.cgs {
		if s.cgs[iter].FitsRange(cID) {
			return s.cgs[iter], true
		}
	}
	return nil, false
}

func (s *Space) Nodes() []Node {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]Node, len(s.cgs))
	for iter := range s.cgs {
		res[iter] = s.cgs[iter].Node()
	}
	return res
}
