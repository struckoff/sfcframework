package balancer

import (
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"sort"
)

// MockNode is Node implementation used for testing.
type MockNode struct {
	id       string
	power    MockPower
	capacity MockCapacity
	h        uint64
}

func NewMockNode(id string, power float64, capacity float64, h uint64) *MockNode {
	return &MockNode{
		id:       id,
		power:    MockPower{power},
		capacity: MockCapacity{capacity},
		h:        h,
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

func (n MockNode) Hash() uint64 {
	return n.h
}

func GenerateMockCells(loadSet ...uint64) map[uint64]*cell {
	cs := make(map[uint64]*cell)
	for iter, load := range loadSet {
		cs[uint64(iter)] = &cell{load: load, id: uint64(iter)}
	}
	return cs
}

func GenerateMockCellGroup(cs map[uint64]*cell, rates []int, powers []float64) []*CellGroup {
	cgs := make([]*CellGroup, len(rates))
	var min, max uint64
	for iter, rate := range rates {
		var load uint64
		cgs[iter] = NewCellGroup(NewMockNode("node-"+string(iter), powers[iter], 0, uint64(uuid.New().ID())))
		cells := make([]*cell, 0, len(cs))
		for key := range cs {
			cells = append(cells, cs[key])
		}
		sort.Slice(cells, func(i, j int) bool { return cells[i].ID() < cells[j].ID() })
		for iterCell := range cells[:rate] {
			cells[iterCell].cg = cgs[iter]
			cgs[iter].cells[cells[iterCell].ID()] = cells[iterCell]
			load += cells[iterCell].load
		}
		max = min + uint64(rate)
		cells = cells[rate:]
		cgs[iter].load = load
		if err := cgs[iter].SetRange(min, max); err != nil {
			panic(err)
		}
		min = max
	}
	return cgs
}

func CompareCellGroup(cg0, cg1 *CellGroup) (bool, string) {
	if !reflect.DeepEqual(cg0.cRange, cg1.cRange) {
		return false, fmt.Sprintf("Different CellGroup.cRange cg0 = %v, cg1 = %v", cg0.cRange, cg1.cRange)
	}
	return true, ""
}
