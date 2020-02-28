package morton

import (
	"errors"
	"fmt"
)

type Curve struct {
	dimensions uint64
	bits       uint64
	length     uint64
	masksArray []uint64
	maxSize    uint64
	maxCode    uint64
}

func New(dims, bits uint64) (*Curve, error) {
	if bits <= 0 || dims <= 0 {
		return nil, errors.New("Number of bits and dimension must be greater than 0")
	}

	mc := &Curve{
		dimensions: dims,
		bits:       bits,
		length:     (bits * dims) - bits,
		maxSize:    (1 << bits) - 1,
		maxCode:    (1 << (dims * bits)) - 1,
	}
	mc.masksArray = mc.masks()

	return mc, nil
}

//Decode returns coordinates for a given code(distance)
func (c Curve) Decode(code uint64) (coords []uint64, err error) {
	if err := c.validateCode(code); err != nil {
		return nil, err
	}
	coords = make([]uint64, c.dimensions)
	coords = c.compacter(coords, code)
	return coords, nil
}

func (c Curve) DecodeWithBuffer(buf []uint64, code uint64) (coords []uint64, err error) {
	if len(buf) < int(c.dimensions) {
		return nil, errors.New("buffer length less then dimensions")
	}
	if err := c.validateCode(code); err != nil {
		return nil, err
	}
	buf = c.compacter(buf, code)
	return buf, nil
}

func (c Curve) validateCode(code uint64) error {
	if code > c.maxCode {
		return errors.New(fmt.Sprintf("code == %v exceeds limit (2^(dimensions * bits) - 1) == %v", code, c.maxSize))
	}
	return nil
}

func (c Curve) compacter(coords []uint64, code uint64) []uint64 {
	for iter := uint64(0); iter < c.dimensions; iter++ {
		coords[iter] = c.compact(code >> iter)
	}
	return coords
}

func (c Curve) compacterAsync(coords []uint64, code uint64) []uint64 {
	ch := make(chan [2]uint64, c.dimensions)
	for iter := uint64(0); iter < c.dimensions; iter++ {
		go func(ch chan [2]uint64, iter uint64) {
			ch <- [2]uint64{iter, c.compact(code >> iter)}
		}(ch, iter)
	}
	for iter := uint64(0); iter < c.dimensions; iter++ {
		pair := <-ch
		coords[pair[0]] = pair[1]
	}
	return coords
}

func (c Curve) compact(x uint64) uint64 {
	//x &= 0x55555555
	//x = (x ^ (x >> 1)) & 0x33333333
	//x = (x ^ (x >> 2)) & 0x0f0f0f0f
	//x = (x ^ (x >> 4)) & 0x00ff00ff
	//x = (x ^ (x >> 8)) & 0x0000ffff

	x &= c.masksArray[len(c.masksArray)-1]
	for iter := 0; iter < len(c.masksArray)-1; iter++ {
		x = (x ^ (x >> (1 << iter))) & (c.masksArray[len(c.masksArray)-2-iter]) //TODO may be "1 << iter" should be pregenerated
	}

	return x
}

func (c Curve) masks() []uint64 {
	mask := uint64((1 << c.bits) - 1)

	shift := c.dimensions * (c.bits - 1)
	shift |= shift >> 1
	shift |= shift >> 2
	shift |= shift >> 4
	shift |= shift >> 8
	shift |= shift >> 16
	shift |= shift >> 32
	shift -= shift >> 1

	masks := make([]uint64, 0, 8)

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
		}

	}

	return masks
}

//Encode returns code(distance) for a given set of coordinates
func (c Curve) Encode(coords []uint64) (code uint64, err error) {
	if err := c.validateCoordinates(coords); err != nil {
		return 0, err
	}
	for iter := uint64(0); iter < c.dimensions; iter++ {
		code |= c.split(coords[iter]) << iter
	}
	return
}

func (c Curve) validateCoordinates(coords []uint64) error {
	if len(coords) < int(c.dimensions) {
		return errors.New(fmt.Sprintf("number of coordinates == %v less then dimensions == %v", len(coords), c.dimensions))
	}
	for iter := range coords {
		if coords[iter] > c.maxSize {
			return errors.New(fmt.Sprintf("coordinate == %v exceeds limit == %v", coords[iter], c.maxSize))
		}
	}
	return nil
}

func (c Curve) split(x uint64) uint64 {
	shiftIter := len(c.masksArray) - 1
	for iter := 0; iter < len(c.masksArray); iter++ {
		x = (x | (x << (1 << shiftIter))) & c.masksArray[iter]
		shiftIter--
	}

	return x
}

// Size returns the maximum coordinate value in any dimension
func (c Curve) Size() uint {
	return uint(c.maxSize)
}

// MaxCode returns the maximum distance along curve(code value)
// 2^(dimensions * bits) - 1
func (c Curve) MaxCode() uint64 {
	return c.maxCode
}
