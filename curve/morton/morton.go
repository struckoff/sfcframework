package morton

import (
	"errors"
	"math/big"
)

type MortonCurve struct {
	dimensions uint64
	bits       uint64
	length     uint64
}

func New(b, n uint64) (*MortonCurve, error){
	if b <= 0 || n <= 0{
		return nil, errors.New("Number of bits and dimension must be greater than 0")
	}
	return &MortonCurve{
		dimensions: n,
		bits:       b,
		length:     (b * n) - b,
	}, nil
}

// TODO USE c.bits
func (c *MortonCurve) Decode(d *big.Int) (coords []uint64, err error){
	coords = make([]uint64, c.dimensions)
	for iter := uint64(0); iter < c.dimensions; iter++ {
		coords[iter] = compact(d.Int64() >> iter)
	}
	return coords, nil
}

func compact(x int64) uint64{
	x &= 0x55555555
	x = (x ^ (x >>  1)) & 0x33333333
	x = (x ^ (x >>  2)) & 0x0f0f0f0f
	x = (x ^ (x >>  4)) & 0x00ff00ff
	x = (x ^ (x >>  8)) & 0x0000ffff
	return uint64(x)
}

func (c MortonCurve) Encode(coords []uint64) (d *big.Int, err error){
	return
}