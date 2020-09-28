/*
	The internal implementation of the space-filling curve algorithm responsible for encoding and decoding of dimensions.
	The dimensions and index represented by unsigned int.
*/
package curve

import (
	"errors"

	"github.com/struckoff/sfcframework/curve/hilbert"
	"github.com/struckoff/sfcframework/curve/morton"
)

//Curve is an interface of space filling curve realisation.
type Curve interface {
	Decode(code uint64) (coords []uint64, err error) //Decode returns coordinates for a given code(distance)
	DecodeWithBuffer(buf []uint64, code uint64) (coords []uint64, err error)
	Encode(coords []uint64) (code uint64, err error) //Encode returns code(distance) for a given set of coordinates
	DimensionSize() uint64                           // DimensionSize returns the maximum coordinate value in any dimension
	Length() uint64                                  // Length returns the maximum distance along curve
	Dimensions() uint64                              // Dimensions - amount of curve dimensions
	Bits() uint64                                    // Bits - size in bits of each dimension
}

//NewCurve - create a curve by given type
//
//cType - curve type(Hilbert, Morton)
//
//dims - amount of curve dimensions.
//
//bits - size in bits of each dimension.
func NewCurve(cType CurveType, dims, bits uint64) (Curve, error) {
	switch cType {
	case Hilbert:
		return hilbert.New(dims, bits)
	case Morton:
		return morton.New(dims, bits)
	default:
		return nil, errors.New("unknown curve type")
	}
}
