package balancer

type cellGroup struct {
	node      Node
	totalLoad uint64
	cells     []cell
}

func newCellGroup(n Node) cellGroup {
	return cellGroup{
		node:  n,
		cells: []cell{},
	}
}
