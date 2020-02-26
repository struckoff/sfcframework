package morton

import (
	"errors"
)

type Curve struct {
	dimensions uint64
	bits       uint64
	length     uint64
	masksArray []uint64
}

func New(dims, bits uint64) (*Curve, error) {
	if bits <= 0 || dims <= 0 {
		return nil, errors.New("Number of bits and dimension must be greater than 0")
	}

	mc := &Curve{
		dimensions: dims,
		bits:       bits,
		length:     (bits * dims) - bits,
	}
	mc.masksArray = mc.masks()

	return mc, nil
}

// TODO USE c.bits
func (c Curve) Decode(code uint64) (coords []uint64, err error) {
	coords = make([]uint64, c.dimensions)
	//var wg sync.WaitGroup
	//wg.Add(int(c.dimensions))
	for iter := uint64(0); iter < c.dimensions; iter++ {
		// TODO AM I HERETIC?
		//go func(iter uint64, wg *sync.WaitGroup) {
		//	defer wg.Done() // SLOW PART
		coords[iter] = c.compact(code >> iter)
		//}(iter, &wg)
	}
	//wg.Wait()
	return coords, nil
}

func (c Curve) DecodeWithBuffer(buf []uint64, code uint64) (coords []uint64, err error) {
	// TODO IMPLEMENT
	return nil, nil
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

func (c Curve) Encode(coords []uint64) (code uint64, err error) {
	// TODO ADD ARGUMENTS CHECK
	if len(coords) < int(c.dimensions) {
		return 0, errors.New("number of coordinates less then dimensions")
	}

	code = 0
	for iter := uint64(0); iter < c.dimensions; iter++ {
		code |= c.split(coords[iter]) << iter
	}
	return
}

func (c Curve) split(x uint64) uint64 {
	shiftIter := len(c.masksArray) - 1
	for iter := 0; iter < len(c.masksArray); iter++ {
		x = (x | (x << (1 << shiftIter))) & c.masksArray[iter]
		shiftIter--
	}

	return x
}

func (c Curve) Len() uint64 {
	return c.length
}
