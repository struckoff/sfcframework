package balancer

import (
	"sync"

	"github.com/struckoff/sfcframework/node"

	"github.com/pkg/errors"
)

//CellGroup represents a segment of SFC attached to the node.
//It contains information about the range of cells attached to this group and node.
type CellGroup struct {
	id     string //unique id of the cell group
	mu     sync.RWMutex
	node   node.Node
	cells  map[uint64]*cell
	load   uint64
	cRange Range
}

//NewCellGroup builds a new CellGroup.
func NewCellGroup(n node.Node) *CellGroup {
	return &CellGroup{
		id:    n.ID(),
		node:  n,
		cells: map[uint64]*cell{},
	}
}

//ID returns the ID of the cell group
func (cg *CellGroup) ID() string {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return cg.id
}

//Node returns the node attached to this cell group.
func (cg *CellGroup) Node() node.Node {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return cg.node
}

//SetNode sets the node attached to the cell group
func (cg *CellGroup) SetNode(n node.Node) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.node = n
}

//Range returns the range of cell group
func (cg *CellGroup) Range() Range {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return cg.cRange
}

//SetRange sets the minimum and maximum of the cell group range.
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

//FitsRange checks if the index of the cell fits the range of the cell group.
func (cg *CellGroup) FitsRange(index uint64) bool {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return cg.cRange.Fits(index)
}

//Cells returns map of cells in the cell group
func (cg *CellGroup) Cells() map[uint64]*cell {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return cg.cells
}

// AddCell adds a cell to the cell group.
// If autoremove flag is true, method calls CellGroup.RemoveCell of previous cell group.
// Flag is useful when CellGroup is altered and not refilled
func (cg *CellGroup) AddCell(c *cell) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.addCell(c)
}

func (cg *CellGroup) addCell(c *cell) {
	if ocl, ok := cg.cells[c.ID()]; ok {
		cg.load -= ocl.Load()
	}
	cg.load += c.Load()
	cg.cells[c.id] = c
	c.SetGroup(cg)
}

// RemoveCell removes a cell from cell group.
func (cg *CellGroup) RemoveCell(id uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.removeCell(id)
}

func (cg *CellGroup) removeCell(id uint64) {
	if cg == nil {
		return
	}
	if cell, ok := cg.cells[id]; ok {
		cg.load -= cell.Load()
		cell.cg = nil
	}
	delete(cg.cells, id)
}

//TotalLoad returns cumulative load of all cells in the cell group.
func (cg *CellGroup) TotalLoad() (load uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.totalLoad()
}

func (cg *CellGroup) totalLoad() (load uint64) {
	for _, cell := range cg.cells {
		load += cell.Load()
	}
	cg.load = load
	return load
}

// AddLoad increase group load by given argument
//func (cg *CellGroup) AddLoad(l uint64) {
//	cg.mu.Lock()
//	defer cg.mu.Unlock()
//	cg.addLoad(l)
//}

//func (cg *CellGroup) addLoad(l uint64) {
//	cg.load += l
//}

//Truncate call truncate method of each cell in the cell group
//which empties a load of each cell without removing cells from the cell group.
func (cg *CellGroup) Truncate() {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.truncate()
}

func (cg *CellGroup) truncate() {
	for cid := range cg.cells {
		cg.cells[cid].Truncate()
	}
	cg.load = 0
}

//SetCells replaces cells in the cell group with given
func (cg *CellGroup) SetCells(cells map[uint64]*cell) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.setCells(cells)
}

func (cg *CellGroup) setCells(cells map[uint64]*cell) {
	cg.cells = cells
	if cells == nil {
		cg.cells = make(map[uint64]*cell)
	}
	cg.totalLoad()
}
