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
	//return []cell{
	//	{load: 0},
	//	{load: 0},
	//	{load: 10},
	//	{load: 20},
	//	{load: 0},
	//	{load: 0},
	//	{load: 80},
	//	{load: 0},
	//	{load: 60},
	//	{load: 0},
	//	{load: 40},
	//	{load: 0},
	//	{load: 90},
	//	{load: 0},
	//	{load: 0},
	//}
}

func GenerateMockCellGroup(cs []cell, rates []int) []CellGroup {
	cgs := make([]CellGroup, len(rates))
	for iter, rate := range rates {
		var load uint64
		cgs[iter] = NewCellGroup(NewMockNode("", 10.0, 0))
		for iterCell := range cs[:rate] {
			cgs[iter].cells = append(cgs[iter].cells, &cs[iterCell])
			load += cs[iterCell].load
		}
		cs = cs[rate:]
		cgs[iter].load = load
	}
	return cgs
	//cg := balancer.
	//cg.cells = append(cg.cells, &cs[0])
	//cg.cells = append(cg.cells, &cs[1])
	//cg.cells = append(cg.cells, &cs[2])
	//cg.cells = append(cg.cells, &cs[3])
	//cg.cells = append(cg.cells, &cs[4])
	//cg.cells = append(cg.cells, &cs[5])
	//cg.cells = append(cg.cells, &cs[6])
	//cg.load = 110
	//cgs[0] = cg
	//cg = balancer.NewCellGroup(balancer.NewMockNode("", 10.0, 0))
	//cg.cells = append(cg.cells, &cs[7])
	//cg.cells = append(cg.cells, &cs[8])
	//cg.cells = append(cg.cells, &cs[9])
	//cg.cells = append(cg.cells, &cs[10])
	//cg.load = 100
	//cgs[1] = cg
	//cg = balancer.NewCellGroup(balancer.NewMockNode("", 10.0, 0))
	//cg.cells = append(cg.cells, &cs[11])
	//cg.cells = append(cg.cells, &cs[12])
	//cg.cells = append(cg.cells, &cs[13])
	//cg.cells = append(cg.cells, &cs[14])
	//cg.load = 90
	//cgs[2] = cg
	//return cgs
}

//func GenerateTestCase(cs []balancer.cell) []balancer.CellGroup {
//	cgs := make([]balancer.CellGroup, 3)
//	cg := balancer.NewCellGroup(balancer.NewMockNode("", 10.0, 0))
//	cg.cells = append(cg.cells, &cs[0])
//	cg.cells = append(cg.cells, &cs[1])
//	cg.cells = append(cg.cells, &cs[2])
//	cg.cells = append(cg.cells, &cs[3])
//	cg.cells = append(cg.cells, &cs[4])
//	cgs[0] = cg
//	cg = balancer.NewCellGroup(balancer.NewMockNode("", 10.0, 0))
//	cg.cells = append(cg.cells, &cs[5])
//	cg.cells = append(cg.cells, &cs[6])
//	cg.cells = append(cg.cells, &cs[7])
//	cg.cells = append(cg.cells, &cs[8])
//	cg.cells = append(cg.cells, &cs[9])
//	cgs[1] = cg
//	cg = balancer.NewCellGroup(balancer.NewMockNode("", 10.0, 0))
//	cg.cells = append(cg.cells, &cs[10])
//	cg.cells = append(cg.cells, &cs[11])
//	cg.cells = append(cg.cells, &cs[12])
//	cg.cells = append(cg.cells, &cs[13])
//	cg.cells = append(cg.cells, &cs[14])
//	cgs[2] = cg
//	return cgs
//}
