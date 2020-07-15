package balancer

import "github.com/struckoff/SFCFramework/curve"

type TransformFunc func(values []interface{}, sfc curve.Curve) ([]uint64, error)
