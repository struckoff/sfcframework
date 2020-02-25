package morton

import (
	"errors"
	"math/big"
	"sync"
)

type MortonCurve struct {
	dimensions uint64
	bits       uint64
	length     uint64
	masksArray []uint64
}

func New(dims, bits uint64) (*MortonCurve, error){
	if bits <= 0 || dims <= 0{
		return nil, errors.New("Number of bits and dimension must be greater than 0")
	}

	mc := &MortonCurve{
		dimensions: dims,
		bits:       bits,
		length:     (bits * dims) - bits,
	}
	mc.masksArray = mc.masks()

	return mc, nil
}

// TODO USE c.bits
func (c *MortonCurve) Decode(d *big.Int) (coords []uint64, err error){
	coords = make([]uint64, c.dimensions)
	var wg sync.WaitGroup
	wg.Add(int(c.dimensions))
	for iter := uint64(0); iter < c.dimensions; iter++ {
		// TODO AM I HERETIC
		go func (iter uint64, wg *sync.WaitGroup){
			defer wg.Done()
		coords[iter] = c.compact(d.Uint64() >> iter)
		}(iter, &wg)
	}
	wg.Wait()
	return coords, nil
}

func (c *MortonCurve)compact(x uint64) uint64{
	//x &= 0x55555555
	//x = (x ^ (x >> 1)) & 0x33333333
	//x = (x ^ (x >> 2)) & 0x0f0f0f0f
	//x = (x ^ (x >> 4)) & 0x00ff00ff
	//x = (x ^ (x >> 8)) & 0x0000ffff

	x &= c.masksArray[len(c.masksArray)-1]
	for iter := 0; iter < len(c.masksArray)-1; iter++{
		x = (x ^ (x >> (1 << iter))) & (c.masksArray[len(c.masksArray)-2-iter]) //TODO may be "1 << iter" should be pregenerated
	}

	return x
}

func (c MortonCurve) masks() []uint64{
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

	for ;shift > 0; shift>>=1 {
		mask = 0
		shifted := uint64(0)

		for bit := uint64(0); bit < c.bits; bit++ {
			distance := (c.dimensions * bit) - bit
			shifted |= shift & distance
			mask |= 1 << bit << (((shift - 1) ^ uint64(0xffffffffffffffff)) & distance)
		}

		if shifted != 0{
			masks = append(masks, mask)
		}

	}

	return masks
}

func (c MortonCurve) Encode(coords []uint64) (d *big.Int, err error){
	// TODO ADD ARGUMENTS CHECK

	code := uint64(0)
	for iter := uint64(0); iter < c.dimensions; iter++ {
		code |= c.split(coords[iter]) << iter
	}
	return new(big.Int).SetUint64(code), nil
}

func (c MortonCurve) split(x uint64) uint64 {
	shiftIter := len(c.masksArray) - 1
	for iter := 0; iter < len(c.masksArray); iter++ {
		x = (x | (x << (1 << shiftIter))) & c.masksArray[iter]
		shiftIter--
	}

	return x
}