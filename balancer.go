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
	bits, err := Log2(size)
	if err != nil {
		return nil, err
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
func (b *Balancer) AddNode(n Node, optimize bool) error {
	if b.space.Len() == 0 || b.nType == nil {
		b.nType = reflect.TypeOf(n)
		//return b.space.AddNode(n)
	}
	//else if reflect.TypeOf(n) != b.nType {
	//	return errors.New("incorrect node type")
	//}
	if err := b.space.AddNode(n); err != nil {
		return err
	}
	if optimize {
		cgs, err := b.of(b.space)
		if err != nil {
			return err
		}
		b.space.SetGroups(cgs)
	}
	return nil
}

func (b *Balancer) GetNode(id string) (Node, bool) {
	return b.space.GetNode(id)
}

func (b *Balancer) Nodes() []Node {
	return b.space.Nodes()
}

func (b *Balancer) SFC() curve.Curve {
	return b.Space().sfc
}

func (b *Balancer) RemoveNode(id string, optimize bool) error {
	if err := b.space.RemoveNode(id); err != nil {
		return err
	}
	if optimize {
		cgs, err := b.of(b.space)
		if err != nil {
			return err
		}
		b.space.SetGroups(cgs)
	}
	return nil
}

// LocateData finds data in the Space of the balancer.
func (b *Balancer) LocateData(d DataItem) (Node, uint64, error) {
	return b.space.LocateData(d)
}

// AddData loads data into the Space of the balancer.
func (b *Balancer) AddData(d DataItem) (Node, uint64, error) {
	return b.space.AddData(d)
}

func (b *Balancer) RemoveData(d DataItem) error {
	return b.space.RemoveData(d)
}

// AddData loads data into the Space of the balancer.
func (b *Balancer) RelocateData(d DataItem, ncID uint64) (Node, uint64, error) {
	return b.space.RelocateData(d, ncID)
}

func (b *Balancer) Optimize() error {
	cgs, err := b.of(b.space)
	if err != nil {
		return err
	}
	b.space.SetGroups(cgs)
	return nil
}

func Log2(n uint64) (p uint64, err error) {
	if (n & (n - 1)) != 0 {
		return 0, errors.New("number must be a power of 2")
	}
	for n > 1 {
		p++
		n = n >> 1
	}
	return
}
