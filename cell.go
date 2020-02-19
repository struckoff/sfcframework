package balancer

type cell struct {
	load  uint64
	items []DataItem
}

func (c *cell) add(d DataItem) error {
	c.load += d.Size()
	c.items = append(c.items, d)
	return nil
}

func (c *cell) itemIDs() []string {
	res := make([]string, len(c.items))
	for i := range c.items {
		res[i] = c.items[i].ID()
	}
	return res
}
