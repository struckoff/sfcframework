package balancer

// Node is an interface describing storage/processing node in the cluster.
type Node interface {
	ID() string
	Power() Power
	Capacity() Capacity
}

func NewDefaultNode(id string, power float64, capacity float64) *DefaultNode {
	return &DefaultNode{}
}

type DefaultNode struct {
	id       string
	power    DefaultPower
	capacity DefaultCapacity
}

func (n DefaultNode) ID() string {
	return n.id
}

func (n DefaultNode) Power() Power {
	return n.power
}

func (n DefaultNode) Capacity() Capacity {
	return n.capacity
}

type DefaultPower struct {
	value float64
}

func (p DefaultPower) Get() float64 {
	return p.value
}

type DefaultCapacity struct {
	value float64
}

func (c DefaultCapacity) Get() float64 {
	return c.value
}
