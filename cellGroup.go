package balancer

import (
	"github.com/pkg/errors"
	"sync"
)

type CellGroup struct {
	id     string //unique id of the cell group
	mu     sync.Mutex
	node   Node
	cells  map[uint64]*cell
	load   uint64
	cRange Range
}

//NewCellGroup builds a new CellGroup
func NewCellGroup(n Node) *CellGroup {
	return &CellGroup{
		id:    n.ID(),
		node:  n,
		cells: map[uint64]*cell{},
	}
}

//ID returns the ID of the cell group
func (cg *CellGroup) ID() string {
	return cg.id
}

//Node returns the node attached to this cell group
func (cg *CellGroup) Node() Node {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.node
}

//SetNode sets the node attached to the cell group
func (cg *CellGroup) SetNode(n Node) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.node = n
}

//Range returns the range of cell group
func (cg *CellGroup) Range() Range {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cRange
}

//SetRange sets minimum and maximum of the cell group range
//if s is not nil method will tryes to update cells of cell group from the space
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
		for _, c := range s.Cells() {
			if cg.cRange.Fits(c.ID()) {
				cg.load += c.Load()
				cg.cells[c.ID()] = c
				if c.cg != nil && c.cg.id != cg.id {
					c.cg.RemoveCell(c.ID())
				}
				c.SetGroup(cg)
			}
		}
	}
	return nil
}

//FitsRange checks if the index of the cell fits the range of the cell group
func (cg *CellGroup) FitsRange(index uint64) bool {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cRange.Fits(index)
}

//Cells returns map of cells in the cell group
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
	if ocl, ok := cg.cells[c.ID()]; ok {
		cg.load -= ocl.Load()
	}
	cg.load += c.Load()
	cg.cells[c.id] = c
	if c.cg != nil && autoremove && c.cg.id != cg.id {
		c.cg.RemoveCell(c.id)
	}
	c.SetGroup(cg)
}

// RemoveCell removes a cell from cell group.
func (cg *CellGroup) RemoveCell(id uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.removeCell(id)
}

func (cg *CellGroup) removeCell(id uint64) {
	if cell, ok := cg.cells[id]; ok {
		cg.load -= cell.Load()
	}
	delete(cg.cells, id)
}

//TotalLoad returns cumulative load of all cells in the cell group
func (cg *CellGroup) TotalLoad() (load uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	for _, cell := range cg.cells {
		load += cell.Load()
	}
	cg.load = load
	return load
}

// AddLoad increase group load by given argument
func (cg *CellGroup) AddLoad(l uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.addLoad(l)
}

func (cg *CellGroup) addLoad(l uint64) {
	cg.load += l
}

//Truncate call truncate method of each cell in the cell group
//emptifies load of each cell without removing cells from the cell group
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
	cg.load = 0
	if cells == nil {
		cg.cells = make(map[uint64]*cell)
	}
	for iter := range cells {
		cg.load += cells[iter].Load()
	}
}
