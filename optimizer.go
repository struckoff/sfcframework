package balancer

//OptimizerFunc is a function which responsible for dividing curve into cell groups
//This function should contains realisation of an algorithm of distribution cell ranges per node.
type OptimizerFunc func(s *Space) ([]*CellGroup, error)
