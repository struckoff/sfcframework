package transform

import (
	"errors"

	"github.com/struckoff/sfcframework/curve"
)

//KVTransform is used to transform string to fit SFC.
//It requires one string value.
func KVTransform(values []interface{}, sfc curve.Curve) ([]uint64, error) {
	if len(values) != 1 {
		return nil, errors.New("number of values must be 1")
	}

	ds := sfc.DimensionSize()
	dc := int(sfc.Dimensions())
	key, ok := values[0].(string)

	if !ok {
		return nil, errors.New("value must be string")
	}

	res := make([]uint64, dc)
	cut := len(key) / dc

	if cut == 0 && len(key) > 0 {
		cut = 1
	}

	for i := 0; i < dc; i++ {
		if len(key) == 0 {
			break
		}

		if i < dc-1 {
			res[i] = stringhash(key[:cut], ds)
			key = key[cut:]
			continue
		}

		res[i] = stringhash(key, ds)
	}

	return res, nil
}

func stringhash(key string, limiter uint64) uint64 {
	var sum int32
	for _, rn := range key {
		sum += rn
	}

	return uint64(sum) % limiter
}
