package hilbert

import (
	"errors"
	"math/big"
)

const bitSize = 8

type HilbertCurve struct {
	dimensions uint64
	bits       uint64
	length     uint64
}

func New(dims, bits uint64) (*HilbertCurve, error) {
	if bits <= 0 || dims <= 0 {
		return nil, errors.New("Number of bits and dimension must be greater than 0")
	}
	return &HilbertCurve{
		dimensions: dims,
		bits:       bits,
		length:     bits * dims,
	}, nil
}

func (c *HilbertCurve) Decode(d *big.Int) (coords []uint64, err error) {
	coords = make([]uint64, c.dimensions)
	coords, err = c.parseIndex(coords, d)
	if err != nil {
		return
	}
	return c.transpose(coords), nil
}

func (c *HilbertCurve) DecodeWithBuffer(buf []uint64, d *big.Int) (coords []uint64, err error) {
	// TODO Need to decide how to deal with less or more the c.dimensions
	coords, err = c.parseIndex(buf, d)
	if err != nil {
		return
	}
	return c.transpose(coords), nil
}

// TODO OPTIMIZE
func (c *HilbertCurve) parseIndex(coords []uint64, d *big.Int) ([]uint64, error) {
	b := d.Bytes()
	lenB := len(b)

	for idx := 0; idx < bitSize*lenB; idx++ {
		if (b[len(b)-1-idx/bitSize] & (1 << (uint32(idx) % bitSize))) != 0 {
			dim := (c.length - uint64(idx) - 1) % c.dimensions
			shift := (uint64(idx) / c.dimensions) % c.bits
			coords[dim] |= 1 << shift
		}
	}
	return coords, nil
}

//! coords may be altered by method
func (c HilbertCurve) Encode(coords []uint64) (d *big.Int, err error) {
	m := uint64(1 << (c.bits - 1))
	coordsLen := len(coords)
	// Inverse undo excess work
	for q := m; q > 0; q >>= 1 {
		p := q - 1
		for i := 0; i < coordsLen; i++ {
			if (coords[i] & q) != 0 {
				coords[0] ^= p
			} else {
				t := (coords[0] ^ coords[i]) & p
				coords[0] ^= t
				coords[i] ^= t
			}
		}
	}

	for i := 1; i < coordsLen; i++ {
		coords[i] ^= coords[i-1]
	}
	t := uint64(0)
	for q := m; q > 1; q >>= 1 {
		if (coords[c.dimensions-1] & q) != 0 {
			t ^= q - 1
		}
	}
	for i := 0; i < coordsLen; i++ {
		coords[i] ^= t
	}

	//h = self._transpose_to_hilbert_integer(x)
	return c.prepareIndex(coords), nil
}

func (c *HilbertCurve) transpose(coords []uint64) []uint64 {
	N := uint64(2 << (c.bits - 1))
	// Note that x is mutated by this method (as a performance improvement
	// to avoid allocation)
	n := int(c.dimensions)

	// Gray decode by H ^ (H/2)
	t := coords[n-1] >> 1
	// Corrected error in Skilling's paper on the following line. The
	// appendix had i >= 0 leading to negative array index.
	for i := n - 1; i > 0; i-- {
		coords[i] ^= coords[i-1]
	}

	coords[0] ^= t
	// Undo excess work
	for q := uint64(2); q != N; q <<= 1 {
		p := q - 1
		for i := n - 1; i >= 0; i-- {
			if (coords[i] & q) != 0 {
				coords[0] ^= p // invert
			} else {
				t = (coords[0] ^ coords[i]) & p
				coords[0] ^= t
				coords[i] ^= t
			}
		}
	} // exchange
	return coords
}

func (c *HilbertCurve) prepareIndex(coords []uint64) *big.Int {
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

	return new(big.Int).SetBytes(tmpCoords)
}

func (c *HilbertCurve) Len() uint64 {
	return c.length
}
