package SFCFramework

// Node is an interface describing storage/processing node in the cluster.
type Node interface {
	ID() string
	Power() Power
	Capacity() Capacity
}
