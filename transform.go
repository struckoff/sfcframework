package balancer

import "github.com/struckoff/sfcframework/curve"

//TransformFunc is an adapter which purpose to convert values into encodable format
type TransformFunc func(values []interface{}, sfc curve.Curve) ([]uint64, error)
