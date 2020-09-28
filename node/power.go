package node

// Power is an interface describing computational ability of the node.
type Power interface {
	Get() float64
}
