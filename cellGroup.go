package balancer

import (
	"github.com/pkg/errors"
	"sync"
)

type CellGroup struct {
	id     string
	mu     sync.Mutex
	node   Node
	cells  map[uint64]*cell
	load   uint64
	cRange Range
}

func NewCellGroup(n Node) *CellGroup {
	return &CellGroup{
		id:    n.ID(),
		node:  n,
		cells: map[uint64]*cell{},
	}
}

func (cg *CellGroup) ID() string {
	return cg.id
}

func (cg *CellGroup) Node() Node {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.node
}

func (cg *CellGroup) SetNode(n Node) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.node = n
}

func (cg *CellGroup) Range() Range {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cRange
}

func (cg *CellGroup) SetRange(min, max uint64, s *Space) error {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	if min > max {
		return errors.Errorf("min(%d) should be less or equall then max(%d)", min, max)
	}
	cg.cRange = Range{
		Min: min,
		Max: max,
		Len: max - min,
	}

	if s != nil {
		cg.load = 0
		cg.cells = make(map[uint64]*cell)
		for _, cl := range s.cells {
			if cg.cRange.Fits(cl.ID()) {
				cg.load += cl.Load()
				cg.cells[cl.ID()] = cl
				//cl.cg.RemoveCell(cl.ID())
				cl.cg = cg
			}
		}
	}

	return nil
}

func (cg *CellGroup) FitsRange(index uint64) bool {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cRange.Fits(index)
}

func (cg *CellGroup) Cells() map[uint64]*cell {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cells
}

// AddCell adds a cell to the cell group
// If autoremove flag is true, method calls CellGroup.RemoveCell of previous cell group.
// Flag is useful when CellGroup is altered and not refilled
func (cg *CellGroup) AddCell(c *cell, autoremove bool) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	if cg == c.cg {
		return
	}
	cg.load += c.load
	cg.cells[c.id] = c
	if c.cg != nil && autoremove {
		c.cg.RemoveCell(c.id)
	}
	c.cg = cg
}

// RemoveCell removes a cell from cell group.
func (cg *CellGroup) RemoveCell(id uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	if cell, ok := cg.cells[id]; ok {
		cg.load -= cell.load
	}
	delete(cg.cells, id)
}

func (cg *CellGroup) TotalLoad() (load uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	for _, cell := range cg.cells {
		load += cell.load
	}
	return load
}

func (cg *CellGroup) addLoad(l uint64) {
	cg.load += l
}

func (cg *CellGroup) removeLoad(l uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.load -= l
}

func (cg *CellGroup) truncate() {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	for cid := range cg.cells {
		cg.cells[cid].dis = make(map[string]struct{})
		cg.cells[cid].load = 0
	}
	cg.load = 0
}
