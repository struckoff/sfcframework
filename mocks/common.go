package mocks

import (
	"github.com/struckoff/SFCFramework"
	"strconv"
)

// MockNode is Node implementation used for testing.
//type MockNode struct {
//	id       string
//	power    MockPower
//	capacity MockCapacity
//	h        uint64
//}
//
//func NewMockNode(id string, power float64, capacity float64, h uint64) *MockNode {
//	return &MockNode{
//		id:       id,
//		power:    MockPower{power},
//		capacity: MockCapacity{capacity},
//		h:        h,
//	}
//}
//
//// ID returns identifier of the node.
//func (n MockNode) ID() string {
//	return n.id
//}
//
//// Power returns MockPower instance of the MockNode.
//func (n MockNode) Power() Power {
//	return n.power
//}
//
//// Capacity returns MockCapacity instance of the MockNode.
//func (n MockNode) Capacity() Capacity {
//	return n.capacity
//}
//
//func (n MockNode) Hash() uint64 {
//	return n.h
//}

//func GenerateMockCells(loadSet ...uint64) map[uint64]*balancer.cell {
//	cs := make(map[uint64]*balancer.cell)
//	for iter, load := range loadSet {
//		cs[uint64(iter)] = &balancer.cell{load: load, id: uint64(iter)}
//	}
//	return cs
//}

func GenerateMockCellGroup(loadSet []uint64, rates []int, powers, caps []float64) []*balancer.CellGroup {
	cgs := make([]*balancer.CellGroup, len(rates))
	var min, max uint64
	for iter, rate := range rates {
		r := &Power{}
		r.On("Get").Return(powers[iter])

		capacity := &Capacity{}
		capacity.On("Get").Return(caps[iter])

		n := &Node{}
		n.On("ID").Return("node-" + strconv.Itoa(iter))
		n.On("Power").Return(r)
		n.On("Capacity").Return(capacity)
		n.On("Hash").Return(uint64(iter))

		//cgs[iter] = NewCellGroup(NewMockNode("node-"+string(iter), powers[iter], 0, uint64(uuid.New().ID())))
		cgs[iter] = balancer.NewCellGroup(n)
		for iterCell := range loadSet[:rate] {
			cell := balancer.NewCell(uint64(iterCell), nil, loadSet[iterCell])
			cgs[iter].AddCell(cell, false)
		}
		max = min + uint64(rate)
		loadSet = loadSet[rate:]
		if err := cgs[iter].SetRange(min, max); err != nil {
			panic(err)
		}
		min = max
	}
	return cgs
}

//func CompareCellGroup(cg0, cg1 *balancer.CellGroup) (bool, string) {
//	if !reflect.DeepEqual(cg0.cRange, cg1.cRange) {
//		return false, fmt.Sprintf("Different CellGroup.cRange cg0 = %v, cg1 = %v", cg0.cRange, cg1.cRange)
//	}
//	return true, ""
//}
