package balancer

import (
	"sync"
)

type cell struct {
	id   uint64
	mu   sync.Mutex
	load uint64
	off  map[string]uint64 // location of relocated DataItem. DataItem.ID -> cell.ID
	cg   *CellGroup
}

func NewCell(id uint64, cg *CellGroup, load uint64) *cell {
	c := cell{
		id:   id,
		load: load,
		cg:   cg,
		off:  make(map[string]uint64),
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

func (c *cell) remove(d DataItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.load -= d.Size()
	c.cg.removeLoad(d.Size())
	return nil
}

func (c *cell) relocate(did string, ncID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.off[did] = ncID
}

func (c *cell) relocated(did string) (uint64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cid, ok := c.off[did]
	return cid, ok
}
