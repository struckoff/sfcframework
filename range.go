package balancer

type Range struct {
	Min uint64
	Max uint64
	Len uint64
}

func (r *Range) Fits(index uint64) bool {
	return index >= r.Min && index < r.Max
}
