package balancer

import (
	"errors"
	"sync"
)

type cell struct {
	id   uint64
	mu   sync.Mutex
	load uint64
	cg   *CellGroup
}

func newCell(id uint64, cgs []CellGroup) (*cell, error) {
	c := cell{
		id:   id,
		load: 0,
	}
	found := false
	for i := range cgs {
		if id >= cgs[i].cRange.Min && id < cgs[i].cRange.Max {
			found = true
			c.cg = &cgs[i]
			break
		}
	}
	if !found {
		return nil, errors.New("unable to bind cell to cell group")
	}
	return &c, nil
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
