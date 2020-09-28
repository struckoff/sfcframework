package balancer

type Range struct {
	Min uint64
	Max uint64
	Len uint64
}

func NewRange(min, max uint64) Range {
	return Range{
		Min: min,
		Max: max,
		Len: max - min,
	}
}

//Fits <- min <= index < max
func (r *Range) Fits(index uint64) bool {
	return index >= r.Min && index < r.Max
}
