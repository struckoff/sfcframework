package balancer

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/struckoff/sfcframework/node"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/struckoff/sfcframework/curve"
	"github.com/struckoff/sfcframework/mocks"
)

func Test_splitCells(t *testing.T) {
	type args struct {
		n int
		l uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []Range
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{
				n: 5,
				l: 5,
			},
			want: []Range{
				{0, 1, 1},
				{1, 2, 1},
				{2, 3, 1},
				{3, 4, 1},
				{4, 5, 1},
			},
			wantErr: false,
		},
		{
			name: "error test",
			args: args{
				n: 50,
				l: 5,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "complex test 1",
			args: args{
				n: 3,
				l: 20,
			},
			want: []Range{
				{0, 7, 7},
				{7, 14, 7},
				{14, 20, 6},
			},
			wantErr: false,
		},
		{
			name: "complex test 2",
			args: args{
				n: 5,
				l: 256,
			},
			want: []Range{
				{0, 52, 52},
				{52, 103, 51},
				{103, 154, 51},
				{154, 205, 51},
				{205, 256, 51},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitCells(tt.args.n, tt.args.l)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSpace_CellGroups(t *testing.T) {
	type fields struct {
		cgs []*CellGroup
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				cgs: []*CellGroup{
					{id: "cg-0"},
					{id: "cg-1"},
					{id: "cg-2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cgs: tt.fields.cgs,
			}

			assert.Equal(t, tt.fields.cgs, s.CellGroups())
		})
	}
}

func TestSpace_Cells(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
	}
	tests := []struct {
		name      string
		fields    fields
		wantCells []*cell
	}{
		{
			name: "",
			fields: fields{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: uint64ptr(11),
					},
					2: {
						id:   2,
						load: uint64ptr(11),
					},
				},
			},
			wantCells: []*cell{
				{
					id:   1,
					load: uint64ptr(11),
				},
				{
					id:   2,
					load: uint64ptr(11),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cells: tt.fields.cells,
			}

			cs := s.Cells()
			sort.Slice(tt.wantCells, func(i, j int) bool { return tt.wantCells[i].ID() < tt.wantCells[j].ID() })
			sort.Slice(cs, func(i, j int) bool { return cs[i].ID() < cs[j].ID() })
			assert.Equal(t, tt.wantCells, cs)
		})
	}
}

func TestSpace_TotalLoad(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
		load  uint64
	}
	tests := []struct {
		name     string
		fields   fields
		wantLoad uint64
	}{
		{
			name: "",
			fields: fields{
				load: 123,
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: uint64ptr(11),
					},
					2: {
						id:   2,
						load: uint64ptr(22),
					},
				},
			},
			wantLoad: 33,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cells: tt.fields.cells,
				load:  tt.fields.load,
			}

			got := s.TotalLoad()
			assert.Equal(t, tt.wantLoad, got)
			assert.Equal(t, tt.wantLoad, s.load)
		})
	}
}

func TestSpace_SetGroups(t *testing.T) {
	type fields struct {
		cgs []*CellGroup
	}
	type args struct {
		groups []*CellGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				cgs: []*CellGroup{
					{id: "cg-0"},
					{id: "cg-1"},
					{id: "cg-2"},
				},
			},
			args: args{
				groups: []*CellGroup{
					{id: "cg-01"},
					{id: "cg-11"},
					{id: "cg-21"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cgs: tt.fields.cgs,
			}
			s.SetGroups(tt.args.groups)
			assert.Equal(t, tt.args.groups, s.cgs)
		})
	}
}

func TestSpace_Len(t *testing.T) {
	type fields struct {
		cgs []*CellGroup
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "",
			fields: fields{
				cgs: []*CellGroup{
					{id: "cg-0"},
					{id: "cg-1"},
					{id: "cg-2"},
				},
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cgs: tt.fields.cgs,
			}
			got := s.Len()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSpace_AddNode(t *testing.T) {
	s := &Space{}

	n0 := &mocks.Node{}
	n0.On("ID").Return("test-node")

	err := s.AddNode(n0)
	assert.NoError(t, err)
	assert.Equal(t, n0, s.cgs[0].node)

	n1 := &mocks.Node{}
	n1.On("ID").Return("test-node")

	err = s.AddNode(n1)
	assert.NoError(t, err)
	assert.Equal(t, n1, s.cgs[0].node)
}

func TestSpace_GetNode(t *testing.T) {
	type nodefields struct {
		ID string
	}
	type fields struct {
		ns []nodefields
	}
	type args struct {
		id string
	}
	type want struct {
		n  *nodefields
		ok bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "exist",
			fields: fields{
				ns: []nodefields{
					{ID: "test-node"},
				},
			},
			args: args{
				id: "test-node",
			},
			want: want{
				n: &nodefields{
					ID: "test-node",
				},
				ok: true,
			},
		},
		{
			name: "not exist",
			fields: fields{
				ns: []nodefields{
					{ID: "test-node"},
				},
			},
			args: args{
				id: "test-no-node",
			},
			want: want{
				n:  nil,
				ok: false,
			},
		},
		{
			name: "empty",
			fields: fields{
				ns: []nodefields{},
			},
			args: args{
				id: "test-no-node",
			},
			want: want{
				n:  nil,
				ok: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgs := make([]*CellGroup, len(tt.fields.ns))
			for i := range tt.fields.ns {
				n := &mocks.Node{}
				n.On("ID").Return(tt.fields.ns[i].ID)
				cgs[i] = &CellGroup{
					id:   tt.fields.ns[i].ID,
					node: n,
				}
			}
			s := &Space{
				cgs: cgs,
			}
			n, ok := s.GetNode(tt.args.id)

			var mn *mocks.Node
			if tt.want.n != nil {
				mn = &mocks.Node{}
				mn.On("ID").Return(tt.want.n.ID)
				assert.Equal(t, mn.ID(), n.ID())
			} else {
				assert.Nil(t, n)
			}
			assert.Equal(t, tt.want.ok, ok)
		})
	}
}

func TestSpace_RemoveNode(t *testing.T) {
	type nodefields struct {
		ID string
	}
	type fields struct {
		ns []nodefields
	}
	type args struct {
		id string
	}
	type want struct {
		ns  []nodefields
		ok  bool
		err bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "exist",
			fields: fields{
				ns: []nodefields{
					{ID: "test-node"},
				},
			},
			args: args{
				id: "test-node",
			},
			want: want{
				ns:  []nodefields{},
				ok:  true,
				err: false,
			},
		},
		{
			name: "not exist",
			fields: fields{
				ns: []nodefields{
					{ID: "test-node"},
				},
			},
			args: args{
				id: "test-no-node",
			},
			want: want{
				ns: []nodefields{
					{ID: "test-node"},
				},
				ok:  false,
				err: true,
			},
		},
		{
			name: "empty",
			fields: fields{
				ns: []nodefields{},
			},
			args: args{
				id: "test-no-node",
			},
			want: want{
				ns:  []nodefields{},
				ok:  false,
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgs := make([]*CellGroup, len(tt.fields.ns))

			for i := range tt.fields.ns {
				n := &mocks.Node{}
				n.On("ID").Return(tt.fields.ns[i].ID)
				cgs[i] = &CellGroup{
					id:   tt.fields.ns[i].ID,
					node: n,
				}
			}

			wantcgs := make([]*CellGroup, len(tt.want.ns))
			for i := range tt.want.ns {
				n := &mocks.Node{}
				n.On("ID").Return(tt.want.ns[i].ID)
				wantcgs[i] = &CellGroup{
					id:   tt.want.ns[i].ID,
					node: n,
				}
			}
			s := &Space{
				cgs: cgs,
			}
			err := s.RemoveNode(tt.args.id)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, wantcgs, s.cgs)
			}
		})
	}
}

func TestSpace_AddData(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		res := make([]uint64, len(values))
		for i := range values {
			res[i] = values[i].(uint64)
		}
		return res, nil
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(1),
		uint64(2),
		uint64(3),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:  n,
		cells: nil,
		load:  0,
		cRange: Range{
			Min: 0,
			Max: math.MaxUint64,
			Len: math.MaxUint64,
		},
	}

	cells := map[uint64]*cell{
		42: {
			id:   42,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	err := s.AddData(0, d)

	assert.NoError(t, err)
	//assert.Equal(t, n, gotn)
	//assert.Equal(t, 42, int(gotcid))
	assert.Equal(t, 111, int(s.load))
}

func TestSpace_LocateData(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		res := make([]uint64, len(values))
		for i := range values {
			res[i] = values[i].(uint64)
		}
		return res, nil
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(1),
		uint64(2),
		uint64(3),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		42: {
			id:   42,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	gotn, gotcid, err := s.LocateData(d)

	assert.NoError(t, err)
	assert.Equal(t, n, gotn)
	assert.Equal(t, 42, int(gotcid))
	assert.Equal(t, 0, int(s.load))
}

func TestSpace_LocateData_Relocated(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		res := make([]uint64, len(values))
		for i := range values {
			res[i] = values[i].(uint64)
		}
		return res, nil
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(1),
		uint64(2),
		uint64(3),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		42: {
			id:   42,
			load: uint64ptr(0),
			cg:   cg,
			off:  map[string]uint64{"test-di": 21},
		},
		21: {
			id:   21,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	gotn, gotcid, err := s.LocateData(d)

	assert.NoError(t, err)
	assert.Equal(t, n, gotn)
	assert.Equal(t, 21, int(gotcid))
	assert.Equal(t, 0, int(s.load))
}

func TestSpace_LocateData_NoNodes(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		res := make([]uint64, len(values))
		for i := range values {
			res[i] = values[i].(uint64)
		}
		return res, nil
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(1),
		uint64(2),
		uint64(3),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		42: {
			id:   42,
			load: uint64ptr(0),
			cg:   cg,
			off:  map[string]uint64{"test-di": 21},
		},
		21: {
			id:   21,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	s := &Space{
		sfc:   sfc,
		cgs:   nil,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	_, _, err := s.LocateData(d)

	assert.Error(t, err)
}

func TestSpace_LocateData_TFNil(t *testing.T) {
	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(1),
		uint64(2),
		uint64(3),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		42: {
			id:   42,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    nil,
		load:  0,
	}

	_, _, err := s.LocateData(d)
	assert.Error(t, err)
}

func TestSpace_LocateData_EncodeErr(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		res := make([]uint64, len(values))
		for i := range values {
			res[i] = values[i].(uint64)
		}
		return res, nil
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(0), errors.New("tests err"))
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(math.MaxUint64),
		uint64(math.MaxUint64),
		uint64(math.MaxUint64),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		42: {
			id:   42,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	_, _, err := s.LocateData(d)
	assert.Error(t, err)
}

func TestSpace_LocateData_NotFound(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		res := make([]uint64, len(values))
		for i := range values {
			res[i] = values[i].(uint64)
		}
		return res, nil
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(math.MaxUint64),
		uint64(math.MaxUint64),
		uint64(math.MaxUint64),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		0: {
			id:   0,
			load: uint64ptr(0),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	_, _, err := s.LocateData(d)
	assert.Error(t, err)
}

func TestSpace_LocateData_TransformErr(t *testing.T) {
	tf := func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
		return nil, errors.New("test err")
	}

	sfc := &mocks.Curve{}
	sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)
	d := &mocks.DataItem{}
	d.On("ID").Return("test-di")
	d.On("Size").Return(uint64(111))
	d.On("Values").Return([]interface{}{
		uint64(math.MaxUint64),
		uint64(math.MaxUint64),
		uint64(math.MaxUint64),
	})
	n := &mocks.Node{}

	cg := &CellGroup{
		node:   n,
		cells:  nil,
		load:   0,
		cRange: Range{},
	}

	cells := map[uint64]*cell{
		42: {
			id:   0,
			load: uint64ptr(42),
			cg:   cg,
		},
	}

	cg.cells = cells

	cgs := []*CellGroup{cg}

	s := &Space{
		sfc:   sfc,
		cgs:   cgs,
		cells: cells,
		tf:    tf,
		load:  0,
	}

	_, _, err := s.LocateData(d)
	assert.Error(t, err)
}

func TestSpace_Nodes(t *testing.T) {
	wantns := make([]node.Node, 11)
	cgs := make([]*CellGroup, len(wantns))
	for i := range wantns {
		n := &mocks.Node{}
		n.On("ID").Return(fmt.Sprintf("test-node-%d", i))
		cgs[i] = &CellGroup{node: n}
		wantns[i] = n
	}

	s := &Space{
		cgs: cgs,
	}

	gotns := s.Nodes()

	assert.Equal(t, wantns, gotns)
}

func TestSpace_FillCellGroup(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
	}
	type args struct {
		cg *CellGroup
	}
	type want struct {
		cells map[uint64]*cell
		load  uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "group fits cells, cell has cg",
			fields: fields{
				cells: map[uint64]*cell{
					0:   {id: 0, load: uint64ptr(1), cg: &CellGroup{load: 100, cells: map[uint64]*cell{0: {id: 0, load: uint64ptr(1)}}}},
					42:  {id: 42, load: uint64ptr(10)},
					111: {id: 111, load: uint64ptr(100)},
				},
			},
			args: args{
				cg: &CellGroup{
					id:    "",
					load:  0,
					cells: make(map[uint64]*cell),
					cRange: Range{
						Min: 0,
						Max: 42,
						Len: 42,
					},
				},
			},
			want: want{
				cells: map[uint64]*cell{
					0: {id: 0, load: uint64ptr(1)},
				},
				load: 1,
			},
		},
		{
			name: "group fits cells",
			fields: fields{
				cells: map[uint64]*cell{
					0:   {id: 0, load: uint64ptr(1)},
					42:  {id: 42, load: uint64ptr(10)},
					111: {id: 111, load: uint64ptr(100)},
				},
			},
			args: args{
				cg: &CellGroup{
					id:    "",
					load:  0,
					cells: make(map[uint64]*cell),
					cRange: Range{
						Min: 0,
						Max: 42,
						Len: 42,
					},
				},
			},
			want: want{
				cells: map[uint64]*cell{
					0: {id: 0, load: uint64ptr(1)},
				},
				load: 1,
			},
		},
		{
			name: "all cells fits",
			fields: fields{
				cells: map[uint64]*cell{
					0:   {id: 0, load: uint64ptr(1)},
					42:  {id: 42, load: uint64ptr(10)},
					111: {id: 111, load: uint64ptr(100)},
				},
			},
			args: args{
				cg: &CellGroup{
					id:    "",
					load:  0,
					cells: make(map[uint64]*cell),
					cRange: Range{
						Min: 0,
						Max: 200,
						Len: 200,
					},
				},
			},
			want: want{
				cells: map[uint64]*cell{
					0:   {id: 0, load: uint64ptr(1)},
					42:  {id: 42, load: uint64ptr(10)},
					111: {id: 111, load: uint64ptr(100)},
				},
				load: 111,
			},
		},
		{
			name: "no cells fits",
			fields: fields{
				cells: map[uint64]*cell{
					0:   {id: 0, load: uint64ptr(1)},
					42:  {id: 42, load: uint64ptr(10)},
					111: {id: 111, load: uint64ptr(100)},
				},
			},
			args: args{
				cg: &CellGroup{
					id:    "",
					load:  0,
					cells: make(map[uint64]*cell),
					cRange: Range{
						Min: 200,
						Max: 400,
						Len: 200,
					},
				},
			},
			want: want{
				cells: map[uint64]*cell{},
				load:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cells: tt.fields.cells,
			}

			cg := tt.args.cg

			for cid := range tt.want.cells {
				tt.want.cells[cid].cg = cg
			}

			s.FillCellGroup(cg)

			assert.Equal(t, tt.want.cells, cg.cells)
			assert.Equal(t, int(tt.want.load), int(cg.load))
		})
	}
}

func TestSpace_TotalPower(t *testing.T) {
	type fields struct {
		powers []float64
	}
	tests := []struct {
		name      string
		fields    fields
		wantPower float64
	}{
		{
			name: "test",
			fields: fields{
				powers: []float64{0.1, 0.01, 0.001},
			},
			wantPower: 0.111,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{}

			cgs := make([]*CellGroup, len(tt.fields.powers))
			for i := range cgs {
				p := &mocks.Power{}
				p.On("Get").Return(tt.fields.powers[i])
				n := &mocks.Node{}
				n.On("Power").Return(p)
				cgs[i] = &CellGroup{node: n}
			}

			s.cgs = cgs

			gotPower := s.TotalPower()
			assert.Equal(t, tt.wantPower, gotPower)
		})
	}
}

func TestSpace_Capacity(t *testing.T) {
	type fields struct {
		length uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "test",
			fields: fields{
				length: 42,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{}
			sfc := &mocks.Curve{}
			sfc.On("Length").Return(tt.fields.length)
			s.sfc = sfc

			got := s.Capacity()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSpace_RelocateData(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
		cgs   []*CellGroup
		tf    TransformFunc
		load  uint64
	}
	type args struct {
		ncID uint64
	}
	type want struct {
		err   bool
		n     node.Node
		cells map[uint64]*cell
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "test",
			fields: fields{
				cells: map[uint64]*cell{
					42: {
						id:   42,
						load: uint64ptr(100),
						off:  make(map[string]uint64),
						cg: &CellGroup{
							node: &mocks.Node{},
						},
					},
					21: {
						id:   21,
						load: uint64ptr(100),
						off:  make(map[string]uint64),
						cg: &CellGroup{
							node: &mocks.Node{},
						},
					},
				},
				cgs: []*CellGroup{{id: "test-node"}},
				tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
					return []uint64{0, 0}, nil
				},
				load: 0,
			},
			args: args{
				ncID: 21,
			},
			want: want{
				err: false,
				n:   &mocks.Node{},
				cells: map[uint64]*cell{
					42: {
						id:   42,
						load: uint64ptr(99),
						off: map[string]uint64{
							"di-id": 21,
						},
						cg: &CellGroup{
							node: &mocks.Node{},
						},
					},
					21: {
						id:   21,
						load: uint64ptr(101),
						off:  make(map[string]uint64),
						cg: &CellGroup{
							node: &mocks.Node{},
						},
					},
				},
			},
		},
		{
			name: "unable to bind cell to cell group",
			fields: fields{
				cells: make(map[uint64]*cell),
				cgs:   []*CellGroup{{id: "test-node"}},
				tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
					return []uint64{0, 0}, nil
				},
				load: 0,
			},
			args: args{
				ncID: 21,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "no nodes in the cluster",
			fields: fields{
				cells: make(map[uint64]*cell),
				cgs:   nil,
				tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
					return []uint64{0, 0}, nil
				},
				load: 0,
			},
			args: args{
				ncID: 21,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cells: tt.fields.cells,
				cgs:   tt.fields.cgs,
				tf:    tt.fields.tf,
				load:  tt.fields.load,
			}

			sfc := &mocks.Curve{}
			sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(42), nil)

			s.sfc = sfc

			d := &mocks.DataItem{}
			d.On("Values").Return([]interface{}{0, 0})
			d.On("Size").Return(uint64(1))
			d.On("ID").Return("di-id")

			got, gotCode, err := s.RelocateData(d, tt.args.ncID)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.ncID, gotCode)
				assert.Equal(t, tt.want.n, got)
				assert.Equal(t, tt.want.cells, s.cells)
			}
		})
	}
}
