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

func (cg CellGroup) Cells() []*cell {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.cells
}

func (cg *CellGroup) AddCell(c *cell) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.addLoad(c.load)
	cg.cells = append(cg.cells, c)
	c.cg = cg
}

func (cg CellGroup) TotalLoad() uint64 {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	//for iter := range cg.cells {
	//	load += cg.cells[iter].load
	//}
	//return load
	return cg.load
}

func (cg *CellGroup) addLoad(l uint64) {
	cg.load += l
}
