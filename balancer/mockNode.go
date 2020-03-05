package balancer

// MockNode is Node implementation used for testing.
type MockNode struct {
	id       string
	power    MockPower
	capacity MockCapacity
}

// ID returns identifier of the node.
func (n MockNode) ID() string {
	return n.id
}

// Power returns MockPower instance of the MockNode.
func (n MockNode) Power() Power {
	return n.power
}

// Capacity returns MockCapacity instance of the MockNode.
func (n MockNode) Capacity() Capacity {
	return n.capacity
}
