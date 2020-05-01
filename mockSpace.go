package balancer

import (
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/transform"
)

func NewMockSpace(cgs []*CellGroup, sfc curve.Curve) *Space {
	var load uint64
	cs := make(map[uint64]*cell)
	for _, cg := range cgs {
		for _, c := range cg.cells {
			cs[c.id] = c
			load += c.load
		}
	}
	return &Space{
		cells: cs,
		cgs:   cgs,
		load:  load,
		sfc:   sfc,
		tf:    transform.SpaceTransform,
	}
}
