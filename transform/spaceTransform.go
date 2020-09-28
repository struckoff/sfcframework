/*
	Default transform functions are provided by the library.
	It covers only the most generic use cases.
	For more appropriate use transform function should be provided by service from the outside.
*/
package transform

import (
	"errors"

	"github.com/struckoff/sfcframework/curve"
)

const latStep = 90.0
const lonStep = 180.0

//SpaceTransform is used to transform geo coordinates to fit SFC.
//It requires two float64 values(latitude, longitude).
func SpaceTransform(values []interface{}, sfc curve.Curve) ([]uint64, error) {
	dimSize := sfc.DimensionSize()
	if len(values) != 2 || sfc.Dimensions() != 2 {
		return nil, errors.New("number of dimensions must be 2")
	}
	res := make([]uint64, 2)
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
	return res, nil
}
