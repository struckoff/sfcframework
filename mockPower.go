package SFCFramework

// MockPower is implementation of Power used for testing.
type MockPower struct {
	value float64
}

// Get returns the value of MockPower instance.
func (p MockPower) Get() float64 {
	return p.value
}
