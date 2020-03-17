package balancer

type OptimizerFunc func(s *Space) ([]CellGroup, error)
