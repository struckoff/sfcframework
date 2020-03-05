package balancer

// DataItem is an interface describing data that is loaded into the system and need to be placed
// at some node within cluster.
type DataItem interface {
	ID() string
	Size() uint64
	Values() []interface{}
}

func NewDefaultDataItem(id string, size uint64, values []interface{}) DefaultDataItem {
	return DefaultDataItem{
		id:     id,
		size:   size,
		values: values,
	}
}

type DefaultDataItem struct {
	id     string
	size   uint64
	values []interface{}
}

func (di DefaultDataItem) ID() string {
	return di.id
}

func (di DefaultDataItem) Size() uint64 {
	return di.size
}

func (di DefaultDataItem) Values() []interface{} {
	return di.values
}
