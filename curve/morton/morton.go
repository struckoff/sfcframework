/*
	The Morton index is expressed by bit interleaving of each dimension.

	Example: 010 & 011 -> 001101
*/
package morton

import (
	"errors"
	"fmt"
)

//Curve - the representation of Morton curve.
type Curve struct {
	dimensions   uint64 //amount of curve dimensions
	bits         uint64 //size in bits of each dimension
	length       uint64 //(bits * dims) - bits
	maxSize      uint64 //maximum value of each dimension
	maxCode      uint64 //biggest code which could be decoded by curve
	masksArray   []uint64
	lshiftsArray []uint64
}

//New - create new hilbert curve.
//
//dims - amount of curve dimensions.
//
//bits - size in bits of each dimension.
func New(dims, bits uint64) (*Curve, error) {
	if bits <= 0 || dims <= 0 {
		return nil, errors.New("number of bits and dimension must be greater than 0")
	}

	mc := &Curve{
		dimensions: dims,
		bits:       bits,
		length:     (bits * dims) - bits,
		maxSize:    (1 << bits) - 1,
		maxCode:    (1 << (dims * bits)) - 1,
	}
	mc.masksArray, mc.lshiftsArray = mc.masks()

	return mc, nil
}

//Decode returns coordinates for a given code(distance)
//Method will return error if code(distance) exceeds the limit(2 ^ (dims * bits) - 1)
func (c *Curve) Decode(code uint64) (coords []uint64, err error) {
	if err := c.validateCode(code); err != nil {
		return nil, err
	}
	coords = make([]uint64, c.dimensions)
	coords = c.compacter(coords, code)
	return coords, nil
}

//DecodeWithBuffer returns coordinates for a given code(distance).
//Method will return error if:
//  - buffer less than number of dimensions
//	- code(distance) exceeds the limit(2 ^ (dims * bits) - 1)
func (c *Curve) DecodeWithBuffer(buf []uint64, code uint64) (coords []uint64, err error) {
	if len(buf) < int(c.dimensions) {
		return nil, errors.New("buffer length less then dimensions")
	}
	if err := c.validateCode(code); err != nil {
		return nil, err
	}
	buf = c.compacter(buf, code)
	return buf, nil
}

func (c *Curve) validateCode(code uint64) error {
	if code > c.maxCode {
		return fmt.Errorf("code == %v exceeds limit (2^(dimensions * bits) - 1) == %v", code, c.maxSize)
	}
	return nil
}

func (c *Curve) compacter(coords []uint64, code uint64) []uint64 {
	for i := uint64(0); i < c.dimensions; i++ {
		coords[i] = c.compact(code >> i)
	}
	return coords
}

func (c *Curve) compact(x uint64) uint64 {
	//x &= 0x55555555
	//x = (x ^ (x >> 1)) & 0x33333333
	//x = (x ^ (x >> 2)) & 0x0f0f0f0f
	//x = (x ^ (x >> 4)) & 0x00ff00ff
	//x = (x ^ (x >> 8)) & 0x0000ffff

	x &= c.masksArray[len(c.masksArray)-1]
	for i := 0; i < len(c.masksArray)-1; i++ {
		//TODO may be "1 << i" should be pregenerated
		x = (x ^ (x >> (1 << i))) & (c.masksArray[len(c.masksArray)-2-i])
	}

	return x
}

func (c *Curve) masks() (masks []uint64, lshifts []uint64) {
	mask := uint64((1 << c.bits) - 1)

	shift := c.dimensions * (c.bits - 1)
	shift |= shift >> 1
	shift |= shift >> 2
	shift |= shift >> 4
	shift |= shift >> 8
	shift |= shift >> 16
	shift |= shift >> 32
	shift -= shift >> 1

	masks = make([]uint64, 0, 8)
	lshifts = make([]uint64, 1, 8)

	masks = append(masks, mask)

	for ; shift > 0; shift >>= 1 {
		mask = 0
		shifted := uint64(0)

		for bit := uint64(0); bit < c.bits; bit++ {
			distance := (c.dimensions * bit) - bit
			shifted |= shift & distance
			mask |= 1 << bit << (((shift - 1) ^ uint64(0xffffffffffffffff)) & distance)
		}

		if shifted != 0 {
			masks = append(masks, mask)
			lshifts = append(lshifts, shift)
		}
	}
	return masks, lshifts
}

//Encode returns code(distance) for a given set of coordinates
//Method will return error if any of the coordinates exceeds limit(2 ^ bits - 1)
func (c *Curve) Encode(coords []uint64) (code uint64, err error) {
	if err := c.validateCoordinates(coords); err != nil {
		return 0, err
	}
	for i := uint64(0); i < c.dimensions; i++ {
		code |= c.split(coords[i]) << i
	}
	return
}

func (c *Curve) validateCoordinates(coords []uint64) error {
	if len(coords) < int(c.dimensions) {
		return fmt.Errorf("number of coordinates == %v less then dimensions == %v", len(coords), c.dimensions)
	}
	for i := range coords {
		if coords[i] > c.maxSize {
			return fmt.Errorf("coordinate == %v exceeds limit == %v", coords[i], c.maxSize)
		}
	}
	return nil
}

func (c *Curve) split(x uint64) uint64 {
	//shiftIter := len(c.masksArray) - 1
	for i := 0; i < len(c.masksArray); i++ {
		x = (x | (x << c.lshiftsArray[i])) & c.masksArray[i]
		//shiftIter--
	}

	return x
}

// DimensionSize returns the maximum coordinate value in any dimension
func (c *Curve) DimensionSize() uint64 {
	return c.maxSize
}

// Length returns the maximum distance along curve(code value)
//
// 2^(dimensions * bits) - 1
func (c *Curve) Length() uint64 {
	return c.maxCode
}

//Dimensions - amount of curve dimensions
func (c *Curve) Dimensions() uint64 {
	return c.dimensions
}

//Bits - size in bits of each dimension
func (c *Curve) Bits() uint64 {
	return c.bits
}
