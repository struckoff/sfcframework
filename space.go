package balancer

import (
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/node"
	"math"
	"sync"

	"github.com/struckoff/SFCFramework/curve"
)

type Space struct {
	mu    sync.Mutex
	cells map[uint64]*cell //cells in the space
	cgs   []*CellGroup     //cell groups in the space
	sfc   curve.Curve      //encoder(Space filling curve)
	tf    TransformFunc    //TransformFunc - transform DataItem into SFC-readable format
	load  uint64
}

func NewSpace(sfc curve.Curve, tf TransformFunc, nodes []node.Node) (*Space, error) {
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
		if i == n-1 {
			res[i] = Range{
				Min: uint64(math.Ceil(c)),
				Max: l,
			}
			break
		}
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
	//ids := make([]uint64, 0, len(s.cells))
	//for k := range s.cells {
	//	ids = append(ids, k)
	//}
	//sort.Slice(ids, func(i, j int) bool {
	//	return ids[i] < ids[j]
	//})
	res := make([]*cell, 0, len(s.cells))
	for id := range s.cells {
		res = append(res, s.cells[id])
	}
	return res
}

func (s *Space) TotalLoad() (load uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalLoad()
}

func (s *Space) totalLoad() (load uint64) {
	for _, cell := range s.cells {
		load += cell.Load()
	}
	s.load = load
	return load
}

//TotalPower returns the sum of the all node powers in the space.
func (s *Space) TotalPower() (power float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for iter := range s.cgs {
		power += s.cgs[iter].Node().Power().Get()
	}
	return
}

// SetGroups replace groups in the space.
func (s *Space) SetGroups(groups []*CellGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cgs = groups
	//for _, cg := range s.cgs {
	//	cg.SetCells(nil)
	//}
	//for _, c := range s.cells {
	//	if cg, ok := s.findCellGroup(c.ID()); ok {
	//		cg.AddCell(c, true)
	//	}
	//}
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
func (s *Space) AddNode(n node.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.addNode(n); err != nil {
		return err
	}
	return nil
}

func (s *Space) addNode(n node.Node) error {
	//TODO May be s.cgs should be map
	for iter := range s.cgs {
		if s.cgs[iter].ID() == n.ID() {
			//s.cgs[iter].Truncate()
			s.cgs[iter].SetNode(n)
			return nil
		}
	}
	s.cgs = append(s.cgs, NewCellGroup(n))
	return nil
}

//GetNode returns node with given ID.
func (s *Space) GetNode(id string) (node.Node, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getNode(id)
}

func (s *Space) getNode(id string) (node.Node, bool) {
	//TODO May be s.cgs should be map
	for iter := range s.cgs {
		if s.cgs[iter].ID() == id {
			return s.cgs[iter].Node(), true
		}
	}
	return nil, false
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
			s.load -= s.cgs[iter].TotalLoad()
			s.cgs[iter].Truncate()
			s.cgs = append(s.cgs[:iter], s.cgs[iter+1:]...)
			return nil
		}
	}
	return errors.Errorf("node(%s) not found", id)
}

//AddData Add data item to the space.
func (s *Space) AddData(d DataItem) (node.Node, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.locateData(d, true)
}

//RemoveData removes data item from the space.
func (s *Space) RemoveData(d DataItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.removeData(d)
}

//LocateData find data item in the space.
func (s *Space) LocateData(d DataItem) (node.Node, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.locateData(d, false)
}

func (s *Space) locateData(d DataItem, load bool) (node.Node, uint64, error) {
	if len(s.cgs) == 0 {
		return nil, 0, errors.New("no nodes in the cluster")
	}
	cID, err := s.cellID(d)
	if err != nil {
		return nil, 0, err
	}
	c, err := s.getCell(cID)
	if err != nil {
		return nil, 0, err
	}
	if ncID, ok := c.Relocated(d.ID()); ok {
		cID = ncID
		c, err = s.getCell(ncID)
		if err != nil {
			return nil, 0, err
		}
	}
	if load {
		if err = c.Add(d); err != nil {
			return nil, 0, err
		}
		s.load += d.Size()
	}
	return c.cg.Node(), cID, nil
}

func (s *Space) getCell(cID uint64) (*cell, error) {
	if cg, ok := s.cells[cID]; !ok {
		if cg != nil {
			return nil, errors.Errorf("cell exists (cID=%v)", cID)
		}
		cg, ok := s.findCellGroup(cID)
		if !ok {
			return nil, errors.Errorf("unable to bind cell to cell group (cID=%v)", cID)
		}
		c := NewCell(cID, cg, 0)
		s.cells[cID] = c
	}
	return s.cells[cID], nil
}

func (s *Space) removeData(d DataItem) error {
	if len(s.cgs) == 0 {
		return nil
	}
	cID, err := s.cellID(d)
	if err != nil {
		return err
	}
	if _, ok := s.cells[cID]; !ok {
		return nil
	}
	if ncID, ok := s.cells[cID].Relocated(d.ID()); ok {
		if _, ok := s.cells[ncID]; ok {
			if err := s.cells[ncID].Remove(d); err != nil {
				return err
			}
		}
	}
	if err := s.cells[cID].Remove(d); err != nil {
		return err
	}
	s.load -= d.Size()
	return nil
}

//RelocateData moves DataItem to another cell
func (s *Space) RelocateData(d DataItem, ncID uint64) (node.Node, uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.relocateData(d, ncID)
}

func (s *Space) relocateData(d DataItem, ncID uint64) (node.Node, uint64, error) {
	if len(s.cgs) == 0 {
		return nil, 0, errors.New("no nodes in the cluster")
	}
	cID, err := s.cellID(d)
	if err != nil {
		return nil, 0, err
	}
	c, err := s.getCell(cID)
	if err != nil {
		return nil, 0, err
	}
	nc, err := s.getCell(ncID)
	if err != nil {
		return nil, 0, err
	}

	c.Relocate(d, ncID)
	if err = nc.Add(d); err != nil {
		return nil, 0, err
	}

	return nc.cg.Node(), ncID, nil
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

func (s *Space) Nodes() []node.Node {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]node.Node, len(s.cgs))
	for iter := range s.cgs {
		res[iter] = s.cgs[iter].Node()
	}
	return res
}
