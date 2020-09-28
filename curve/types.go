package curve

//CurveType - type of the space-filling curve to use.
type CurveType int

const (
	Hilbert CurveType = iota //Hilbert curve
	Morton                   //Morton curve
)

//String - string representation of the curve type.
func (c CurveType) String() string {
	switch c {
	case Hilbert:
		return "Hilbert"
	case Morton:
		return "Morton"
	}
	return ""
}
