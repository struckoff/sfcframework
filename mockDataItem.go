package balancer

import (
	"github.com/google/uuid"
	"math/rand"
)

// MockDataItem is DataItem implementation used for testing.
type MockDataItem struct {
	id     string
	size   uint64
	values []interface{}
}

func (m MockDataItem) ID() string {
	return m.id
}

func (m MockDataItem) Size() uint64 {
	return m.size
}

func (m MockDataItem) Values() []interface{} {
	res := make([]interface{}, len(m.values))
	for i := range m.values {
		res[i] = m.values[i]
	}
	return res
}

//func GenerateRandomMockDataItem(dimensions uint64) MockDataItem {
//	coords := make([]uint64, dimensions)
//	for c := range coords {
//		coords[c] = rand.Uint64()
//	}
//	return MockDataItem{
//		id:     uuid.New().String(),
//		size:   rand.Uint64(),
//		values: coords,
//	}
//}

func GenerateRandomMockSpaceItem() MockDataItem {
	coords := make([]interface{}, 3)
	coords[0] = rand.Float64()
	coords[1] = rand.Float64()
	coords[2] = rand.Int63n(31)
	return MockDataItem{
		id:     uuid.New().String(),
		size:   rand.Uint64(),
		values: coords,
	}
}

func NewMockDataItem(id string, size uint64, values []interface{}) DataItem {
	return MockDataItem{
		id:     id,
		size:   size,
		values: values,
	}
}