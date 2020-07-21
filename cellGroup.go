package balancer

import (
	"github.com/pkg/errors"
	"log"
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

func (cg *CellGroup) SetRange(min, max uint64) error {
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
	return nil
}

func (cg *CellGroup) FitsRange(index uint64) bool {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return index >= cg.cRange.Min && index < cg.cRange.Max
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
	//cg.load += c.load
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
	delete(cg.cells, id)
}

func (cg *CellGroup) TotalLoad() (load uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	for iter := range cg.cells {
		log.Println(cg.id, cg.cells[iter].id, cg.cells[iter].load, load)
		load += cg.cells[iter].load
	}
	return
}

func (cg *CellGroup) addLoad(l uint64) {
	cg.load += l
}
