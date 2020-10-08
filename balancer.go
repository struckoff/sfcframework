/*
Library provides a hashing mechanism that uses the space-filling curve as a core algorithm for locating data between multiple nodes.
This approach provides high data locality and distribution based on data and nodes features.
*/
package balancer

import (
	"errors"

	"github.com/struckoff/sfcframework/node"

	"github.com/struckoff/sfcframework/curve"
)

//Balancer is responsible for distributing load between nodes of the cluster.
//It stores meta-information about the current load in the cells of the space-filling curve.
//Cells are stored in the slice of cell groups, one cell group for every node of the cluster.
//To balance the load between nodes, the balancer analyzes the power and capacity of nodes
//and distributes cells between cell groups in such a way that all nodes would be equally.
type Balancer struct {
	space *Space
	of    OptimizerFunc
}

//NewBalancer creates a new instance of balancer
//using curve type,  amount of dimension to work with,
//size of the each dimension,
//function which transform DataItem into SFC-readable format,
//optimizer function which distributes cells into groups
//and list of nodes in space(could be nil).
func NewBalancer(cType curve.CurveType, dims, size uint64, tf TransformFunc, of OptimizerFunc, nodes []node.Node) (*Balancer, error) {
	bits, err := log2(size)
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

//Space provides direct access to space inside balancer.
func (b *Balancer) Space() *Space {
	return b.space
}

//AddNode adds the node to the Space of balancer and initiates rebalancing of cells
//between cell groups.
func (b *Balancer) AddNode(n node.Node, optimize bool) error {
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

//GetNode returns the node by the given ID.
func (b *Balancer) GetNode(id string) (node.Node, bool) {
	return b.space.GetNode(id)
}

//Nodes return a list of nodes from space.
func (b *Balancer) Nodes() []node.Node {
	return b.space.Nodes()
}

//SFC returns space-filling curve which space use.
func (b *Balancer) SFC() curve.Curve {
	return b.Space().sfc
}

//RemoveNode removes the node from space by the given ID.
//If optimize is true it will also update cell groups by the optimizer.
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

//LocateData finds data in the Space of the balancer.
func (b *Balancer) LocateData(d DataItem) (node.Node, uint64, error) {
	return b.space.LocateData(d)
}

//AddData loads data into the Space of the balancer.
func (b *Balancer) AddData(cID uint64, d DataItem) error {
	return b.space.AddData(cID, d)
}

//RemoveData removes DataItem from the Space of the balancer
func (b *Balancer) RemoveData(d DataItem) error {
	return b.space.RemoveData(d)
}

//RelocateData tells balancer to locate DataItem in specific cell.
func (b *Balancer) RelocateData(d DataItem, ncID uint64) (node.Node, uint64, error) {
	return b.space.RelocateData(d, ncID)
}

//Optimize updates cell groups by optimizer.
func (b *Balancer) Optimize() error {
	cgs, err := b.of(b.space)
	if err != nil {
		return err
	}
	b.space.SetGroups(cgs)
	return nil
}

func log2(n uint64) (p uint64, err error) {
	if (n & (n - 1)) != 0 {
		return 0, errors.New("number must be a power of 2")
	}
	for n > 1 {
		p++
		n = n >> 1
	}
	return
}
