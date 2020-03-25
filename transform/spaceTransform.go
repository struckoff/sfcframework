package transform

import (
	"errors"
	"github.com/struckoff/SFCFramework/curve"
)

const ft = 1609459200
const latStep = 90.0
const lonStep = 180.0

func SpaceTransform(values []interface{}, sfc curve.Curve) ([]uint64, error) {
	dimSize := sfc.DimensionSize()
	if len(values) != 3 || sfc.Dimensions() != 3 {
		return nil, errors.New("number of dimensions must be 3")
	}
	res := make([]uint64, 3)
	lat, ok := values[0].(float64)
	if !ok {
		return nil, errors.New("first value must be float64 latitude")
	}
	res[0] = uint64((lat + latStep) / (latStep * 2) * float64(dimSize))
	lon, ok := values[1].(float64)
	if !ok {
		return nil, errors.New("second value must be float64 longitude")
	}
	res[1] = uint64((lon + lonStep) / (lonStep * 2) * float64(dimSize))
	ts, ok := values[2].(int64)
	if !ok {
		return nil, errors.New("third value must be int64 timestamp")
	}
	res[2] = uint64(ts-ft) % dimSize
	//res[2] = uint64(float64(ts) / ft * float64(dimSize))
	return res, nil
}
