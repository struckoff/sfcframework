package balancer

import (
	"errors"
	"reflect"
)

// Balancer is responsible for distributing load between nodes of the cluster. It stores
// meta-information about the current load in the cells of the space-filling curve. Cells are
// stored in the slice of cell groups, one cell group for every node of the cluster. To balance
// the load between nodes, balancer analyzes power and capacity of nodes, and distributes
// cells between cell groups in such way that all nodes would be equally.
type Balancer struct {
	nType reflect.Type
	space space
}

// AddNode adds node to the space of balancer, and initiates rebalancing of cells
// between cell groups.
func (b *Balancer) AddNode(n Node) error {
	if b.space.len() == 0 {
		b.nType = reflect.TypeOf(n)
		b.space.addNode(n)
		return nil
	}
	if reflect.TypeOf(n) != b.nType {
		return errors.New("incorrect node type")
	}
	b.space.addNode(n)
	return nil
}

// AddData loads data into the space of the balancer.
func (b *Balancer) AddData(d DataItem) error {
	return b.space.addData(d)
}

// Distribution generates distribution of data items accross nodes of the cluster.
func (b *Balancer) Distribution() DataDistribution {
	return b.space.distribution()
}
