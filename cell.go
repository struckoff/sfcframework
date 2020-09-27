package balancer

import (
	"sync"
	"sync/atomic"
)

type cell struct {
	id   uint64 //unique id of the cell
	mu   sync.RWMutex
	load *uint64
	off  map[string]uint64 // location of Relocated DataItem. DataItem.ID -> cell.ID
	cg   *CellGroup
}

func NewCell(id uint64, cg *CellGroup) *cell {
	c := cell{
		id:   id,
		cg:   cg,
		load: new(uint64),
		off:  make(map[string]uint64),
	}
	if cg != nil {
		cg.AddCell(&c)
	}
	return &c
}

func (c *cell) ID() uint64 {
	return c.id
}

//SetGroup sets the group of the cell
func (c *cell) SetGroup(cg *CellGroup) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cg = cg
}

//Group returns the group to which the cell attached
func (c *cell) Group() *CellGroup {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cg
}

//Load return load of cell
func (c *cell) Load() (load uint64) {
	return atomic.LoadUint64(c.load)
}

//Truncate - emptifies load of the cell
func (c *cell) Truncate() {
	atomic.StoreUint64(c.load, 0)
}

//AddLoad increase load of the cell
func (c *cell) AddLoad(l uint64) {
	atomic.AddUint64(c.load, l)
}

//Remove removes the DataItem from the cell
func (c *cell) RemoveLoad(l uint64) {
	//TODO: overflow check
	atomic.AddUint64(c.load, ^(l - 1))
}

//Relocate sets the sprecified DataItem as moved and store
//index of the new cell
func (c *cell) Relocate(d DataItem, ncID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.off[d.ID()] = ncID
}

//Relocated returns the id of the cell to which DataItem was relocated
//if DataItem was relocated it also returns the true otherwise false
func (c *cell) Relocated(did string) (uint64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cid, ok := c.off[did]
	return cid, ok
}
