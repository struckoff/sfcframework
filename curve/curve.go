package curve

import (
	"github.com/struckoff/SFCFramework/curve/hilbert"
	"github.com/struckoff/SFCFramework/curve/morton"
)

type Curve interface {
	Decode(d int) (coords []uint, err error)
	DecodeWithBuffer(buf []uint, d int) (coords []uint, err error)
	Encode(coords []uint) (d int, err error)
	Len() uint
	Size() uint
}

func NewCurve(cType CurveType, dims, bits uint64) (Curve, error) {
	switch cType {
	case Hilbert:
		return hilbert.New(dims, bits)
	case Morton:
		return morton.New(dims, bits)
	}
}
