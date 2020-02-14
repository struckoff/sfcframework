package hilbert

import (
	"encoding/binary"
	"errors"
)

type HilbertCurve struct {
	dimentions uint64
	bits uint64
	length uint64
}

func New(b, n uint64) (*HilbertCurve, error){
	if b <= 0 || n <= 0{
		return nil, errors.New("Number of bits and dimension must be greater than 0")
	}
	return &HilbertCurve{
		dimentions: n,
		bits: b,
		length: b * n,
	}, nil
}

func (c HilbertCurve) Decode(d uint64) (coords []uint64, err error){
	return
}

func (c HilbertCurve) DecodeWithBuffer(buf []uint64, d uint64) (coords []uint64, err error){
	return
}

func (c HilbertCurve) Encode(coords []uint64) (d uint64, err error){
	m := uint64(1 << (c.bits - 1))
	// Inverse undo excess work
	for q:= m; q > 0; q >>=1{
		p := q-1
		for i:=uint64(0); i < c.dimentions; i++{
			if (coords[i] & q) != 0{
				coords[i]^=p
			} else {
				t := (coords[0] ^ coords[i]) & p
				coords[0] ^= t
				coords[i] ^= t
			}
		}
	}

	for i:=uint64(1); i < c.dimentions; i++{
		coords[i] ^= coords[i-1]
	}
	t := uint64(0)
	for q:= m; q > 1; q >>=1{
		if (coords[c.dimentions - 1] & q) != 0{
			t ^= q -1
		}
	}
	for i:=uint64(0); i < c.dimentions; i++{
		coords[i]^=t
	}

	//h = self._transpose_to_hilbert_integer(x)
	return c.intoNumeric(coords), nil
}

func (c *HilbertCurve) intoNumeric(coords []uint64) uint64 {
	tmpCoords := make([]byte, c.length)
	bIndex := c.length - 1
	mask := uint64(1 << (c.bits - 1))

	for iter := uint64(0); iter < c.bits; iter++ {
		for coordsIter := range coords {
			if (coords[coordsIter] & mask) != 0 {
				tmpCoords[c.length-1-bIndex/8] |= 1 << (bIndex % 8)
			}
			bIndex--
		}
		mask >>= 1
	}

	//switch len(tmpCoords) {}

	switch{
		case c.length >= 8:
			return binary.BigEndian.Uint64(tmpCoords)
		case c.length >= 4:
			return uint64(binary.BigEndian.Uint32(tmpCoords))
		default:
			return uint64(binary.BigEndian.Uint16(tmpCoords))
	}
}

func (c *HilbertCurve) Len() uint64 {
	return c.length
}