package balancer

// MockNode is Node implementation used for testing.
type MockNode struct {
	id       string
	power    MockPower
	capacity MockCapacity
}

func NewMockNode(id string, power float64, capacity float64) *MockNode {
	return &MockNode{
		id:       id,
		power:    MockPower{power},
		capacity: MockCapacity{capacity},
	}
}

// ID returns identifier of the node.
func (n MockNode) ID() string {
	return n.id
}

// Power returns MockPower instance of the MockNode.
func (n MockNode) Power() Power {
	return n.power
}

// Capacity returns MockCapacity instance of the MockNode.
func (n MockNode) Capacity() Capacity {
	return n.capacity
}

func GenerateMockCells(loadSet ...uint64) []cell {
	cs := make([]cell, 0, len(loadSet))
	for _, load := range loadSet {
		cs = append(cs, cell{load: load})
	}
	return cs
}

//func GenerateMockCellGroup(cs []cell, rates []int, powers []float64) []CellGroup {
//	cgs := make([]CellGroup, len(rates))
//	for iter, rate := range rates {
//		var load uint64
//		cgs[iter] = NewCellGroup(NewMockNode("node-"+string(iter), powers[iter], 0))
//		for iterCell := range cs[:rate] {
//			cs[iterCell].cg = &cgs[iter]
//			cgs[iter].cells = append(cgs[iter].cells, &cs[iterCell])
//			load += cs[iterCell].load
//		}
//		cs = cs[rate:]
//		cgs[iter].load = load
//	}
//	return cgs
//
//}
