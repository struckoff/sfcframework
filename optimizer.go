package balancer

type OptimizerFunc func(cgs []cellGroup) ([]cellGroup, error)
