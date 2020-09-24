package balancer

import "github.com/struckoff/sfcframework/curve"

type TransformFunc func(values []interface{}, sfc curve.Curve) ([]uint64, error)
