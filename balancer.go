package balancer

import (
	"errors"
	"reflect"

	"github.com/struckoff/SFCFramework/curve"
)

// Balancer is responsible for distributing load between nodes of the cluster. It stores
// meta-information about the current load in the cells of the space-filling curve. Cells are
// stored in the slice of cell groups, one cell group for every node of the cluster. To balance
// the load between nodes, balancer analyzes power and capacity of nodes, and distributes
// cells between cell groups in such way that all nodes would be equally.
type Balancer struct {
	nType reflect.Type
	space *Space
	of    OptimizerFunc
}

func NewBalancer(cType curve.CurveType, dims, size uint64, tf TransformFunc, of OptimizerFunc, nodes []Node) (*Balancer, error) {
	if (size & (size - 1)) != 0 {
		return nil, errors.New("size must be a power of 2")
	}
	bits := uint64(0)
	v := uint64(1)
	for v != size {
		v *= 2
		bits++
	}
	sfc, err := curve.NewCurve(cType, dims, bits)
	if err != nil {
		return nil, err
	}
	s, err := NewSpace(sfc, tf, nodes)
	if err != nil {
		return nil, err
	}
	return &Balancer{
		space: s,
		of:    of,
	}, nil
}

func (b *Balancer) Space() *Space {
	return b.space
}

// AddNode adds node to the Space of balancer, and initiates rebalancing of cells
// between cell groups.
func (b *Balancer) AddNode(n Node) error {
	if b.space.Len() == 0 || b.nType == nil {
		b.nType = reflect.TypeOf(n)
		//return b.space.AddNode(n)
	} else if reflect.TypeOf(n) != b.nType {
		return errors.New("incorrect node type")
	}
	if err := b.space.AddNode(n); err != nil {
		return err
	}
	cgs, err := b.of(b.space)
	if err != nil {
		return err
	}
	b.space.SetGroups(cgs)
	return nil
}

// LocateData loads data into the Space of the balancer.
func (b *Balancer) LocateData(d DataItem) (Node, error) {
	return b.space.LocateData(d)
}

func (b *Balancer) Optimize() error {
	ns, err := b.of(b.space)
	if err != nil {
		return err
	}
	b.space.cgs = ns
	return nil
}