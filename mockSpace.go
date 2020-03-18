package balancer

import (
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/transform"
)

func NewMockSpace(cgs []CellGroup, cs map[uint64]*cell, sfc curve.Curve) *Space {
	var load uint64
	for iter := range cs {
		load += cs[iter].load
	}
	return &Space{
		cells: cs,
		cgs:   cgs,
		load:  load,
		sfc:   sfc,
		tf:    transform.SpaceTransform,
	}
}
