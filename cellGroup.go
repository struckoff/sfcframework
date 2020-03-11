package balancer

import "sync"

type CellGroup struct {
	mu    sync.Mutex
	node  Node
	cells []*cell
	load  uint64
}

func NewCellGroup(n Node) CellGroup {
	return CellGroup{
		node:  n,
		cells: []*cell{},
	}
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

func (cg CellGroup) Cells() []*cell {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cells
}

func (cg *CellGroup) AddCell_(c *cell) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.addLoad(c.load)
	cg.cells = append(cg.cells, c)
	c.cg = cg
}

// AddCell adds a cell to the cell group
// If autoremove flag is true method calls CellGroup.RemoveCell of previous cell group.
// Flag is usefull when CellGroup is altered and not refilled
//! CellGroup.RemoveCell destroys order of CellGroup.cells due to optimizations.
func (cg *CellGroup) AddCell(c *cell, autoremove bool) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	if cg == c.cg {
		return
	}
	cg.load += c.load
	cg.cells = append(cg.cells, c)
	if c.cg != nil && autoremove {
		c.cg.RemoveCell(c)
	}
	c.cg = cg
}

//! Destroys order
// RemoveCell removes a cell from cell group.
func (cg *CellGroup) RemoveCell(c *cell) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	for iter := range cg.cells {
		if cg.cells[iter] == c {
			cg.cells[len(cg.cells)-1], cg.cells[iter] = cg.cells[iter], cg.cells[len(cg.cells)-1]
			cg.cells = cg.cells[:len(cg.cells)-1]
			return
		}
	}
	return
}

func (cg *CellGroup) TotalLoad() uint64 {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.load
}

func (cg *CellGroup) addLoad(l uint64) {
	cg.load += l
}
