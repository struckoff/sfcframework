package balancer

import (
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/SFCFramework/node/mocks"
	"sort"
	"testing"
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
						load: 11,
						dis:  map[string]uint64{"di-1": 11},
					},
					2: {
						id:   2,
						load: 11,
						dis:  map[string]uint64{"di-2": 11},
					},
				},
			},
			wantCells: []*cell{
				{
					id:   1,
					load: 11,
					dis:  map[string]uint64{"di-1": 11},
				},
				{
					id:   2,
					load: 11,
					dis:  map[string]uint64{"di-2": 11},
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
						load: 11,
						dis:  map[string]uint64{"di-0": 11},
					},
					2: {
						id:   2,
						load: 22,
						dis:  map[string]uint64{"di-1": 22},
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
