package SFCFramework

type TransformFunc func(values []interface{}, dimSize uint64) ([]uint64, error)
