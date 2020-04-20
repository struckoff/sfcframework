package curve

type CurveType int

const (
	Hilbert CurveType = iota
	Morton
)

func (c CurveType) String() string {
	switch c {
	case Hilbert:
		return "Hilbert"
	case Morton:
		return "Morton"
	}
	return ""
}
