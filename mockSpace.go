package balancer

func NewMockSpace(cgs []CellGroup, cs []cell) *Space {
	var load uint64
	for iter := range cs {
		load += cs[iter].load
	}
	return &Space{
		cells: cs,
		cgs:   cgs,
		load:  load,
	}
}
