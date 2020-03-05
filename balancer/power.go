package balancer

// Power is an interface decribing computational ability of the node.
type Power interface {
	Get() float64
}
