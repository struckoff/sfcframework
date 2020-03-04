package balancer

import "sync"

type cellGroup struct {
	mu    sync.Mutex
	node  Node
	cells []*cell
	load  uint64
}

func newCellGroup(n Node) cellGroup {
	return cellGroup{
		node:  n,
		cells: []*cell{},
	}
}

func (cg *cellGroup) addLoad(l uint64) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.load += l
}
