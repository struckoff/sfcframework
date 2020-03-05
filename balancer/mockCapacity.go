package balancer

// MockCapacity is implementation of Capacity used for testing.
type MockCapacity struct {
	value float64
}

// Get returns the value of MockCapacity instance.
func (c MockCapacity) Get() float64 {
	return c.value
}
