package node

// Node is an interface describing storage/processing node in the cluster.
type Node interface {
	ID() string
	Power() Power
	Capacity() Capacity
	Hash() uint64 //unique node hash
}
