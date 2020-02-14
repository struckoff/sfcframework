package curve

type Curve interface {
	Decode(d int) (coords []uint, err error)
	DecodeWithBuffer(buf []uint, d int) (coords []uint, err error)
	Encode(coords []uint) (d int, err error)
	Len() uint
}


