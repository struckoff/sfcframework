package balancer

import (
	"sync"
)

type cell struct {
	id   uint64
	mu   sync.Mutex
	load uint64
	cg   *CellGroup
}

func newCell(id uint64, cg *CellGroup) *cell {
	c := cell{
		id:   id,
		load: 0,
		cg:   cg,
	}
	//found := false
	//for i := range cgs {
	//	if id >= cgs[i].cRange.Min && id < cgs[i].cRange.Max {
	//		found = true
	//		c.cg = &cgs[i]
	//		break
	//	}
	//}
	//if !found { //? May be this could be c.cg == nil
	//	return nil, errors.New("unable to bind cell to cell group")
	//}
	return &c
}

func (c *cell) ID() uint64 {
	return c.id
}

func (c *cell) SetGroup(cg *CellGroup) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cg = cg
}

func (c *cell) Load() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.load
}

func (c *cell) add(d DataItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.load += d.Size()
	c.cg.addLoad(d.Size())
	return nil
}
