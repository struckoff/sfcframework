package hilbert

import (
	"encoding/binary"
	"errors"
)

const bitSize = 8

 //The Hilbert index is expressed as an array of transposed bits.
 //
 //Example: 5 bits for each of n=3 coordinates.
 //15-bit Hilbert integer = A B C D E F G H I J K L M N O is stored
 //as its Transpose                        ^
 //X[0] = A D G J M                    X[2]|  7
 //X[1] = B E H K N        <------->       | /X[1]
 //X[2] = C F I L O                   axes |/
 //       high low                         0------> X[0]
 //
 //NOTE: This algorithm is derived from work done by John Skilling and published in "Programming the Hilbert curve".
 //(c) 2004 American Institute of Physics.
type Curve struct {
	dimensions uint64
	bits       uint64
	length     uint64
}

func New(dims, bits uint64) (*Curve, error) {
	if bits <= 0 || dims <= 0 {
		return nil, errors.New("number of bits and dimension must be greater than 0")
	}
	return &Curve{
		dimensions: dims,
		bits:       bits,
		length:     bits * dims,
	}, nil
}

func (c Curve) Decode(code uint64) (coords []uint64, err error) {
	coords = make([]uint64, c.dimensions)
	coords, err = c.parseIndex(coords, code)
	if err != nil {
		return
	}
	return c.transpose(coords), nil
}

func (c Curve) DecodeWithBuffer(buf []uint64, code uint64) (coords []uint64, err error) {
	if len(buf) < int(c.dimensions){
		return nil, errors.New("buffer length less then dimensions")
	}
	coords, err = c.parseIndex(buf, code)
	if err != nil {
		return
	}
	coords = c.transpose(coords)
	return coords, nil
}

// TODO OPTIMIZE
func (c Curve) parseIndex(coords []uint64, code uint64) ([]uint64, error) {
	var bitLen int

	b := make([]byte, bitSize)
	binary.LittleEndian.PutUint64(b, code)

	for iter := 0; iter < bitSize; iter++ {
		if (1 << (iter * bitSize)) > code {
			bitLen = iter
			break
		}
	}

	//fmt.Println(b, b[:bitLen], bitLen, new(big.Int).SetUint64(code).Bytes())

	b = b[:bitLen]
	for iter := 0; iter < bitSize*bitLen; iter++ {
		if (b[iter/bitSize] & (1 << (iter % bitSize))) != 0 {
			dim := (c.length - uint64(iter) - 1) % c.dimensions
			shift := (uint64(iter) / c.dimensions) % c.bits
			coords[dim] |= 1 << shift
		}
	}
	return coords, nil
}

//! coords may be altered by method
func (c Curve) Encode(coords []uint64) (code uint64, err error) {
	if len(coords) < int(c.dimensions) {
		return 0, errors.New("number of coordinates less then dimensions")
	}
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
	code = c.prepareIndex(coords)
	return
}

func (c Curve) transpose(coords []uint64) []uint64 {
	m := uint64(2 << (c.bits - 1))
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
	for q := uint64(2); q != m; q <<= 1 {
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

func (c Curve) prepareIndex(coords []uint64) uint64 {
	tmpCoords := make([]byte, bitSize)
	bIndex := c.length - 1
	mask := uint64(1 << (c.bits - 1))

	for iter := uint64(0); iter < c.bits; iter++ {
		for coordsIter := range coords {
			if (coords[coordsIter] & mask) != 0 {
				tmpCoords[c.length-1-bIndex/bitSize] |= 1 << (bIndex % 8)
			}
			bIndex--
		}
		mask >>= 1
	}

	return binary.LittleEndian.Uint64(tmpCoords)
}

func (c Curve) Len() uint64 {
	return c.length
}
