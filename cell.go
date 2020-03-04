package balancer

import "sync"

type cell struct {
	mu    sync.Mutex
	load  uint64
	items []DataItem
	cg    *cellGroup
}

func newCell() cell {
	return cell{
		load:  0,
		items: []DataItem{},
	}
}

func (c *cell) add(d DataItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.load += d.Size()
	c.items = append(c.items, d)
	c.cg.addLoad(d.Size())
	return nil
}

func (c *cell) itemIDs() []string {
	res := make([]string, len(c.items))
	for i := range c.items {
		res[i] = c.items[i].ID()
	}
	return res
}
