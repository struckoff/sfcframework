package balancer

// DataItem is an interface describing data that is loaded into the system and need to be placed
// at some node within cluster.
type DataItem interface {
	ID() string
	Size() uint64
	Coordinates() []uint64
}
