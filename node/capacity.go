package node

// Capacity is an interface describing how much data can the node hold.
type Capacity interface {
	Get() (float64, error)
}
