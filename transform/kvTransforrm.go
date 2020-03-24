package transform

import (
	"errors"
	"github.com/struckoff/SFCFramework/curve"
)

func KVTransform(values []interface{}, sfc curve.Curve) ([]uint64, error) {
	ds := sfc.DimensionSize()
	dc := int(sfc.Dimensions())
	if len(values) != 1 {
		return nil, errors.New("number of values must be 1")
	}
	key, ok := values[0].(string)
	if !ok {
		return nil, errors.New("value must be string")
	}

	res := make([]uint64, dc)
	cut := len(key) / dc
	if cut == 0 && len(key) > 0 {
		cut = 1
	}
	for iter := 0; iter < dc; iter++ {
		if len(key) == 0 {
			break
		}
		if iter < dc-1 {
			res[iter] = stringhash(key[:cut], ds)
			key = key[cut:]
			continue
		}
		res[iter] = stringhash(key, ds)
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
