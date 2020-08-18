package balancer

import (
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/SFCFramework/node"
	"github.com/struckoff/SFCFramework/node/mocks"
	"testing"
)

func TestCellGroup_ID(t *testing.T) {
	type fields struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				id: "group-0",
			},
			want: "group-0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				id: tt.fields.id,
			}

			assert.Equal(t, tt.want, cg.ID())
		})
	}
}

func TestCellGroup_Node(t *testing.T) {
	type fields struct {
		node node.Node
	}
	tests := []struct {
		name   string
		fields fields
		want   node.Node
	}{
		{
			name: "",
			fields: fields{
				node: &mocks.Node{},
			},
			want: &mocks.Node{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				node: tt.fields.node,
			}

			assert.Equal(t, tt.want, cg.Node())
		})
	}
}

func TestCellGroup_SetNode(t *testing.T) {
	type args struct {
		n node.Node
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				n: &mocks.Node{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{}

			cg.SetNode(tt.args.n)
			assert.Equal(t, tt.args.n, cg.node)
		})
	}
}

func TestCellGroup_Range(t *testing.T) {
	type fields struct {
		cRange Range
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				cRange: Range{
					Min: 11,
					Max: 111,
					Len: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cRange: tt.fields.cRange,
			}
			assert.Equal(t, tt.fields.cRange, cg.Range())
		})
	}
}

func TestCellGroup_SetRange(t *testing.T) {
	type fields struct {
		cells  map[uint64]*cell
		load   uint64
		cRange Range
	}
	type args struct {
		min uint64
		max uint64
		s   *Space
	}
	type want struct {
		err    bool
		cRange Range
		cells  map[uint64]*cell
		load   uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "has cells, space",
			fields: fields{
				cells: map[uint64]*cell{
					1:  {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
					2:  {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
					10: {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
					15: {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
				},
				load: 1111,
				cRange: Range{
					Min: 0,
					Max: 111,
					Len: 111,
				},
			},
			args: args{
				min: 10,
				max: 121,
				s: &Space{
					cells: map[uint64]*cell{
						1:   {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
						2:   {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
						10:  {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
						15:  {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
						111: {id: uint64(111), dis: map[string]uint64{"di-111": 10000}, load: 10000},
						115: {id: uint64(115), dis: map[string]uint64{"di-115": 100000}, load: 100000},
						121: {id: uint64(121), dis: map[string]uint64{"di-121": 1000000}, load: 1000000},
						122: {id: uint64(122), dis: map[string]uint64{"di-122": 10000000}, load: 10000000},
					},
				},
			},
			want: want{
				err: false,
				cRange: Range{
					Min: 10,
					Max: 121,
					Len: 111,
				},
				load: 111100,
				cells: map[uint64]*cell{
					10:  {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
					15:  {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
					111: {id: uint64(111), dis: map[string]uint64{"di-111": 10000}, load: 10000},
					115: {id: uint64(115), dis: map[string]uint64{"di-115": 100000}, load: 100000},
				},
			},
		},
		{
			name: "no cells, space",
			fields: fields{
				cells: make(map[uint64]*cell),
				load:  0,
				cRange: Range{
					Min: 0,
					Max: 0,
					Len: 0,
				},
			},
			args: args{
				min: 10,
				max: 121,
				s: &Space{
					cells: map[uint64]*cell{
						1:   {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
						2:   {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
						10:  {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
						15:  {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
						111: {id: uint64(111), dis: map[string]uint64{"di-111": 10000}, load: 10000},
						115: {id: uint64(115), dis: map[string]uint64{"di-115": 100000}, load: 100000},
						121: {id: uint64(121), dis: map[string]uint64{"di-121": 1000000}, load: 1000000},
						122: {id: uint64(122), dis: map[string]uint64{"di-122": 10000000}, load: 10000000},
					},
				},
			},
			want: want{
				err: false,
				cRange: Range{
					Min: 10,
					Max: 121,
					Len: 111,
				},
				load: 111100,
				cells: map[uint64]*cell{
					10:  {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
					15:  {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
					111: {id: uint64(111), dis: map[string]uint64{"di-111": 10000}, load: 10000},
					115: {id: uint64(115), dis: map[string]uint64{"di-115": 100000}, load: 100000},
				},
			},
		},
		{
			name: "no space",
			fields: fields{
				cells: map[uint64]*cell{
					1:  {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
					2:  {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
					10: {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
					15: {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
				},
				load: 1111,
				cRange: Range{
					Min: 0,
					Max: 111,
					Len: 111,
				},
			},
			args: args{
				min: 10,
				max: 121,
				s:   nil,
			},
			want: want{
				err: false,
				cRange: Range{
					Min: 10,
					Max: 121,
					Len: 111,
				},
				load: 1111,
				cells: map[uint64]*cell{
					1:  {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
					2:  {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
					10: {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
					15: {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells:  tt.fields.cells,
				load:   tt.fields.load,
				cRange: tt.fields.cRange,
			}

			if tt.args.s != nil {
				for _, cl := range tt.want.cells {
					cl.cg = cg
				}
			}

			err := cg.SetRange(tt.args.min, tt.args.max, tt.args.s)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.cRange, cg.cRange)
			assert.Equal(t, tt.want.cells, cg.cells)
			assert.Equal(t, int(tt.want.load), int(cg.load))
		})
	}
}

func TestCellGroup_Cells(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
	}
	tests := []struct {
		name   string
		fields fields
		want   map[uint64]*cell
	}{
		{
			name: "",
			fields: fields{
				cells: map[uint64]*cell{
					1:  {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
					2:  {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
					10: {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
					15: {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
				},
			},
			want: map[uint64]*cell{
				1:  {id: uint64(1), dis: map[string]uint64{"di-1": 1}, load: 1},
				2:  {id: uint64(2), dis: map[string]uint64{"di-2": 10}, load: 10},
				10: {id: uint64(10), dis: map[string]uint64{"di-10": 100}, load: 100},
				15: {id: uint64(15), dis: map[string]uint64{"di-15": 1000}, load: 1000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells: tt.fields.cells,
			}

			assert.Equal(t, tt.want, cg.Cells())
		})
	}
}

func TestCellGroup_AddCell(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
		load  uint64
	}
	type args struct {
		c *cell
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
			name: "cells not empty",
			fields: fields{
				cells: map[uint64]*cell{1: {
					id:   1,
					load: 1,
					off:  nil,
					dis:  map[string]uint64{"di-0": 1},
					cg:   nil,
				}},
				load: 1,
			},
			args: args{
				c: &cell{
					id:   2,
					load: 1,
					off:  nil,
					dis:  map[string]uint64{"di-1": 1},
				},
			},
			want: want{
				load: 2,
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-0": 1},
						cg:   nil,
					},
					2: {
						id:   2,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-1": 1},
					},
				},
			},
		},
		{
			name: "cells empty",
			fields: fields{
				cells: make(map[uint64]*cell),
				load:  0,
			},
			args: args{
				c: &cell{
					id:   2,
					load: 1,
					off:  nil,
					dis:  map[string]uint64{"di-1": 1},
				},
			},
			want: want{
				load: 1,
				cells: map[uint64]*cell{
					2: {
						id:   2,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-1": 1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells: tt.fields.cells,
				load:  tt.fields.load,
			}
			for _, c := range tt.want.cells {
				c.cg = cg
			}
			for _, c := range cg.cells {
				c.cg = cg
			}
			cg.AddCell(tt.args.c)
			assert.Equal(t, tt.want.cells, cg.cells)
			assert.Equal(t, int(tt.want.load), int(cg.load))
		})
	}
}

func TestCellGroup_RemoveCell(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
		load  uint64
	}
	type args struct {
		cid uint64
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
			name: "cells not empty",
			fields: fields{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-0": 1},
						cg:   nil,
					},
					2: {
						id:   2,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-1": 1},
					},
				},
				load: 2,
			},
			args: args{cid: 2},
			want: want{
				load: 1,
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-0": 1},
						cg:   nil,
					},
				},
			},
		},
		{
			name: "cells empty",
			fields: fields{
				cells: map[uint64]*cell{
					2: {
						id:   2,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-1": 1},
					},
				},
				load: 1,
			},
			args: args{cid: 2},
			want: want{
				load:  0,
				cells: map[uint64]*cell{},
			},
		},
		{
			name: "cell not exists",
			fields: fields{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-0": 1},
					},
				},
				load: 1,
			},
			args: args{cid: 2},
			want: want{
				load: 1,
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						off:  nil,
						dis:  map[string]uint64{"di-0": 1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells: tt.fields.cells,
				load:  tt.fields.load,
			}
			for _, c := range tt.want.cells {
				c.cg = cg
			}
			for _, c := range cg.cells {
				c.cg = cg
			}
			c, ok := cg.cells[tt.args.cid]
			cg.RemoveCell(tt.args.cid)
			if ok {
				assert.Nil(t, c.cg)
			}
			assert.Equal(t, tt.want.cells, cg.cells)
			assert.Equal(t, int(tt.want.load), int(cg.load))
		})
	}
}

func TestCellGroup_TotalLoad(t *testing.T) {
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
			name: "not empty",
			fields: fields{
				cells: map[uint64]*cell{
					1: {
						id:  1,
						dis: map[string]uint64{"di-0": 1},
					},
					2: {
						id:  2,
						dis: map[string]uint64{"di-10": 10},
					},
					3: {
						id:  3,
						dis: map[string]uint64{"di-100": 100},
					},
				},
				load: 1,
			},
			wantLoad: 111,
		},
		{
			name: "empty",
			fields: fields{
				cells: map[uint64]*cell{},
				load:  1,
			},
			wantLoad: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells: tt.fields.cells,
				load:  tt.fields.load,
			}

			assert.Equal(t, int(tt.wantLoad), int(cg.TotalLoad()))
			assert.Equal(t, int(tt.wantLoad), int(cg.load))
		})
	}
}

func TestCellGroup_Truncate(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
		load  uint64
	}
	type want struct {
		cells map[uint64]*cell
		load  uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "not empty",
			fields: fields{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						dis:  map[string]uint64{"di-0": 1},
					},
					2: {
						id:   2,
						load: 0,
						dis:  map[string]uint64{"di-2": 0},
					},
					3: {
						id:   3,
						load: 2,
						dis:  map[string]uint64{"di-3": 2},
					},
				},
				load: 3,
			},
			want: want{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 0,
						dis:  map[string]uint64{},
					},
					2: {
						id:   2,
						load: 0,
						dis:  map[string]uint64{},
					},
					3: {
						id:   3,
						load: 0,
						dis:  map[string]uint64{},
					},
				},
				load: 0,
			},
		},
		{
			name: "empty",
			fields: fields{
				cells: map[uint64]*cell{},
				load:  0,
			},
			want: want{
				cells: map[uint64]*cell{},
				load:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells: tt.fields.cells,
				load:  tt.fields.load,
			}
			cg.Truncate()
			assert.Equal(t, tt.want.cells, cg.cells)
			assert.Equal(t, int(tt.want.load), int(cg.load))
		})
	}
}

func TestCellGroup_SetCells(t *testing.T) {
	type fields struct {
		cells map[uint64]*cell
		load  uint64
	}
	type args struct {
		cells map[uint64]*cell
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
			name: "nil",
			fields: fields{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						dis:  map[string]uint64{},
					},
					2: {
						id:   2,
						load: 10,
						dis:  map[string]uint64{},
					},
					3: {
						id:   3,
						load: 100,
						dis:  map[string]uint64{},
					},
				},
				load: 111,
			},
			args: args{
				cells: nil,
			},
			want: want{
				cells: make(map[uint64]*cell),
				load:  0,
			},
		},
		{
			name: "empty",
			fields: fields{
				cells: make(map[uint64]*cell),
				load:  0,
			},
			args: args{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						dis:  map[string]uint64{"di-1": 1},
					},
					2: {
						id:   2,
						load: 10,
						dis:  map[string]uint64{"di-2": 10},
					},
					3: {
						id:   3,
						load: 100,
						dis:  map[string]uint64{"di-3": 100},
					},
				},
			},
			want: want{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						dis:  map[string]uint64{"di-1": 1},
					},
					2: {
						id:   2,
						load: 10,
						dis:  map[string]uint64{"di-2": 10},
					},
					3: {
						id:   3,
						load: 100,
						dis:  map[string]uint64{"di-3": 100},
					},
				},
				load: 111,
			},
		},
		{
			name: "not empty",
			fields: fields{
				cells: map[uint64]*cell{
					10: {
						id:   10,
						load: 1,
						dis:  map[string]uint64{},
					},
					20: {
						id:   20,
						load: 10,
						dis:  map[string]uint64{},
					},
					30: {
						id:   30,
						load: 100,
						dis:  map[string]uint64{},
					},
				},
				load: 111,
			},
			args: args{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						dis:  map[string]uint64{"di-1": 1},
					},
					2: {
						id:   2,
						load: 10,
						dis:  map[string]uint64{"di-2": 10},
					},
					3: {
						id:   3,
						load: 100,
						dis:  map[string]uint64{"di-3": 100},
					},
				},
			},
			want: want{
				cells: map[uint64]*cell{
					1: {
						id:   1,
						load: 1,
						dis:  map[string]uint64{"di-1": 1},
					},
					2: {
						id:   2,
						load: 10,
						dis:  map[string]uint64{"di-2": 10},
					},
					3: {
						id:   3,
						load: 100,
						dis:  map[string]uint64{"di-3": 100},
					},
				},
				load: 111,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := &CellGroup{
				cells: tt.fields.cells,
				load:  tt.fields.load,
			}
			cg.SetCells(tt.args.cells)
			assert.Equal(t, tt.want.cells, cg.cells)
			assert.Equal(t, int(tt.want.load), int(cg.load))
		})
	}
}

func TestNewCellGroup(t *testing.T) {
	type args struct {
		nid string
	}
	tests := []struct {
		name string
		args args
		want *CellGroup
	}{
		{
			name: "",
			args: args{
				nid: "test-node",
			},
			want: &CellGroup{
				id:    "test-node",
				cells: make(map[uint64]*cell),
				load:  0,
				cRange: Range{
					Min: 0,
					Max: 0,
					Len: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &mocks.Node{}
			n.On("ID").Return(tt.args.nid)
			tt.want.node = n

			cg := NewCellGroup(n)
			assert.Equal(t, tt.want, cg)
		})
	}
}
