package curve

type Curve interface {
	Decode(code uint64) (coords []uint64, err error)
	DecodeWithBuffer(buf []uint64, code uint64) (coords []uint64, err error)
	Encode(coords []uint64) (code uint64, err error)
	Len() uint64
}
