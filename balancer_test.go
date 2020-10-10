package balancer

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/struckoff/sfcframework/curve/morton"

	"github.com/struckoff/sfcframework/mocks"

	"github.com/struckoff/sfcframework/curve"
	"github.com/struckoff/sfcframework/node"

	"github.com/stretchr/testify/assert"
)

func TestLog2(t *testing.T) {
	type args struct {
		n uint64
	}
	tests := []struct {
		name    string
		args    args
		wantP   uint64
		wantErr bool
	}{
		{
			name: "256",
			args: args{
				n: 256,
			},
			wantP:   8,
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				n: 255,
			},
			wantP:   0,
			wantErr: true,
		},
		{
			name: "1",
			args: args{
				n: 1,
			},
			wantP:   0,
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				n: 2,
			},
			wantP:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, err := log2(tt.args.n)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantP, gotP)
			}
		})
	}
}

func TestNewBalancer(t *testing.T) {
	type args struct {
		cType curve.CurveType
		dims  uint64
		size  uint64
		tf    TransformFunc
		of    OptimizerFunc
		nodes []node.Node
	}
	type want struct {
		b          *Balancer
		cType      curve.CurveType
		dims, bits uint64
		err        bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test",
			args: args{
				cType: curve.Morton,
				dims:  2,
				size:  256,
				tf:    nil,
				of:    nil,
				nodes: make([]node.Node, 0),
			},
			want: want{
				b: &Balancer{
					space: &Space{
						cells: make(map[uint64]*cell),
						cgs:   make([]*CellGroup, 0),
						tf:    nil,
						load:  0,
					},
					of: nil,
				},
				cType: curve.Morton,
				dims:  2,
				bits:  8,
				err:   false,
			},
		},
		{
			name: "wrong size",
			args: args{
				cType: curve.Morton,
				dims:  2,
				size:  255,
				tf:    nil,
				of:    nil,
				nodes: make([]node.Node, 0),
			},
			want: want{
				err: true,
			},
		},
		{
			name: "wrong dims",
			args: args{
				cType: curve.Morton,
				dims:  0,
				size:  256,
				tf:    nil,
				of:    nil,
				nodes: []node.Node{&mocks.Node{}},
			},
			want: want{
				err: true,
			},
		},
		{
			name: "nil node",
			args: args{
				cType: curve.Morton,
				dims:  2,
				size:  2,
				tf:    nil,
				of:    nil,
				nodes: make([]node.Node, 42),
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBalancer(tt.args.cType, tt.args.dims, tt.args.size, tt.args.tf, tt.args.of, tt.args.nodes)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				wsfc, _ := curve.NewCurve(tt.want.cType, tt.want.dims, tt.want.bits)
				tt.want.b.space.sfc = wsfc
				assert.NoError(t, err)
				assert.Equal(t, tt.want.b, got)
			}
		})
	}
}

func TestBalancer_Space(t *testing.T) {
	type fields struct {
		space *Space
	}
	tests := []struct {
		name   string
		fields fields
		want   *Space
	}{
		{
			name: "test",
			fields: fields{
				space: &Space{
					cells: make(map[uint64]*cell, 42),
					cgs:   make([]*CellGroup, 42),
				},
			},
			want: &Space{
				cells: make(map[uint64]*cell, 42),
				cgs:   make([]*CellGroup, 42),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
			}
			got := b.Space()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBalancer_AddNode(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type nd struct {
		id string
	}
	type args struct {
		n        *nd
		optimize bool
	}
	type want struct {
		b   *Balancer
		ns  []nd
		err bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "optimize",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   nil,
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return append(s.cgs, s.cgs...), nil
				},
			},
			args: args{
				n:        &nd{"test-node"},
				optimize: true,
			},
			want: want{
				ns: []nd{
					{"test-node"},
					{"test-node"},
				},
				b: &Balancer{
					space: &Space{
						cells: nil,
						cgs:   nil,
						sfc:   nil,
						tf:    nil,
						load:  0,
					},
					of: nil,
				},
				err: false,
			},
		},
		{
			name: "no optimize",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   nil,
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return append(s.cgs, s.cgs...), nil
				},
			},
			args: args{
				n:        &nd{"test-node"},
				optimize: false,
			},
			want: want{
				ns: []nd{
					{"test-node"},
				},
				b: &Balancer{
					space: &Space{
						cells: nil,
						cgs:   nil,
						sfc:   nil,
						tf:    nil,
						load:  0,
					},
					of: nil,
				},
				err: false,
			},
		},
		{
			name: "nil node",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   nil,
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return nil, errors.New("test err")
				},
			},
			args: args{
				n:        nil,
				optimize: false,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "optimize err",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   nil,
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return nil, errors.New("test err")
				},
			},
			args: args{
				n:        &nd{"test-node"},
				optimize: true,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			var n *mocks.Node
			if tt.args.n != nil {
				n = &mocks.Node{}
				n.On("ID").Return(tt.args.n.id)
			}

			err := b.AddNode(n, tt.args.optimize)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.want.b.space.cgs = make([]*CellGroup, len(tt.want.ns))
				for i := range tt.want.ns {
					tt.want.b.space.cgs[i] = &CellGroup{
						id:     tt.want.ns[i].id,
						node:   n,
						cells:  make(map[uint64]*cell),
						load:   0,
						cRange: Range{},
					}
				}
				tt.want.b.space.tf = nil
				tt.want.b.of = nil
				b.space.tf = nil
				b.of = nil
				assert.Equal(t, tt.want.b, b)
			}
		})
	}
}

func TestBalancer_GetNode(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   node.Node
		wantOk bool
	}{
		{
			name: "test",
			fields: fields{
				space: &Space{
					cgs: []*CellGroup{{
						id:    "test-node",
						node:  &mocks.Node{},
						cells: nil,
						load:  0,
						cRange: Range{
							Min: 0,
							Max: 0,
							Len: 0,
						},
					}},
					sfc:  nil,
					tf:   nil,
					load: 0,
				},
				of: nil,
			},
			args: args{
				id: "test-node",
			},
			want:   &mocks.Node{},
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			got, ok := b.GetNode(tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestBalancer_Nodes(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   []node.Node
	}{
		{
			name: "test",
			fields: fields{
				space: &Space{
					cgs: []*CellGroup{
						{node: &mocks.Node{}},
						{node: &mocks.Node{}},
						{node: &mocks.Node{}},
					},
				},
				of: nil,
			},
			want: []node.Node{
				&mocks.Node{}, &mocks.Node{}, &mocks.Node{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			got := b.Nodes()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBalancer_SFC(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	tests := []struct {
		name   string
		fields fields
		want   curve.Curve
	}{
		{
			name: "",
			fields: fields{
				space: &Space{
					sfc: &morton.Curve{},
				},
			},
			want: &morton.Curve{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			got := b.SFC()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBalancer_RemoveNode(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type args struct {
		id       string
		optimize bool
	}
	type want struct {
		b   *Balancer
		err bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "optimize",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs: []*CellGroup{
						{id: "test-node-0"},
						{id: "test-node-1"},
						{id: "test-node-2"},
					},
					sfc:  nil,
					tf:   nil,
					load: 0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return append(s.cgs, s.cgs...), nil
				},
			},
			args: args{
				id:       "test-node-0",
				optimize: true,
			},
			want: want{
				b: &Balancer{
					space: &Space{
						cgs: []*CellGroup{
							{id: "test-node-1"},
							{id: "test-node-2"},
							{id: "test-node-1"},
							{id: "test-node-2"},
						},
					},
					of: nil,
				},
				err: false,
			},
		},
		{
			name: "no optimize",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs: []*CellGroup{
						{id: "test-node-0"},
						{id: "test-node-1"},
						{id: "test-node-2"},
					},
					sfc:  nil,
					tf:   nil,
					load: 0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return append(s.cgs, s.cgs...), nil
				},
			},
			args: args{
				id:       "test-node-0",
				optimize: false,
			},
			want: want{
				b: &Balancer{
					space: &Space{
						cgs: []*CellGroup{
							{id: "test-node-1"},
							{id: "test-node-2"},
						},
					},
					of: nil,
				},
				err: false,
			},
		},
		{
			name: "not exist",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   nil,
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return nil, errors.New("test err")
				},
			},
			args: args{
				id:       "not exist",
				optimize: false,
			},
			want: want{
				err: true,
			},
		},
		{
			name: "optimize err",
			fields: fields{
				space: &Space{
					cgs: []*CellGroup{
						{id: "test-node-0"},
						{id: "test-node-1"},
						{id: "test-node-2"},
					},
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return nil, errors.New("test err")
				},
			},
			args: args{
				id:       "test-node-0",
				optimize: true,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}

			err := b.RemoveNode(tt.args.id, tt.args.optimize)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.want.b.space.tf = nil
				tt.want.b.of = nil
				b.space.tf = nil
				b.of = nil
				assert.Equal(t, tt.want.b, b)
			}
		})
	}
}

func TestBalancer_LocateData(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type args struct {
		d DataItem
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     node.Node
		wantCode uint64
		wantErr  bool
	}{
		{
			name: "",
			fields: fields{
				space: &Space{
					cells: make(map[uint64]*cell),
					cgs: []*CellGroup{
						{
							id:     "test-node",
							cRange: Range{0, math.MaxUint64, math.MaxUint64},
							cells:  make(map[uint64]*cell),
						},
					},
					tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
						return []uint64{0, 0}, nil
					},
					load: 0,
				},
				of: nil,
			},
			args: args{
				d: nil,
			},
			want:     nil,
			wantCode: 0,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			d := &mocks.DataItem{}
			d.On("Values").Return([]interface{}{})
			d.On("ID").Return("test-di")

			sfc := &mocks.Curve{}
			sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(0), nil)

			b.space.sfc = sfc

			got, gotCode, err := b.LocateData(d)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantCode, gotCode)
			}
		})
	}
}

func TestBalancer_AddData(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type args struct {
		d DataItem
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     node.Node
		wantCode uint64
		wantErr  bool
	}{
		{
			name: "",
			fields: fields{
				space: &Space{
					cells: make(map[uint64]*cell),
					cgs: []*CellGroup{
						{
							id:     "test-node",
							cRange: Range{0, math.MaxUint64, math.MaxUint64},
							cells:  make(map[uint64]*cell),
						},
					},
					tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
						return []uint64{0, 0}, nil
					},
					load: 0,
				},
				of: nil,
			},
			args: args{
				d: nil,
			},
			want:     nil,
			wantCode: 0,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			d := &mocks.DataItem{}
			d.On("Values").Return([]interface{}{})
			d.On("ID").Return("test-di")
			d.On("Size").Return(uint64(1))

			sfc := &mocks.Curve{}
			sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(0), nil)

			b.space.sfc = sfc

			err := b.AddData(1, d)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				//assert.Equal(t, tt.want, got)
				//assert.Equal(t, tt.wantCode, gotCode)
			}
		})
	}
}

func TestBalancer_RemoveData(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type args struct {
		d DataItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    node.Node
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				space: &Space{
					cells: make(map[uint64]*cell),
					cgs: []*CellGroup{
						{
							id:     "test-node",
							cRange: Range{0, math.MaxUint64, math.MaxUint64},
							cells:  make(map[uint64]*cell),
						},
					},
					tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
						return []uint64{0, 0}, nil
					},
					load: 0,
				},
				of: nil,
			},
			args: args{
				d: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			d := &mocks.DataItem{}
			d.On("Values").Return([]interface{}{})
			d.On("ID").Return("test-di")
			d.On("Size").Return(uint64(1))

			sfc := &mocks.Curve{}
			sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(0), nil)

			b.space.sfc = sfc

			err := b.RemoveData(d)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBalancer_RelocateData(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	type args struct {
		d    DataItem
		ncID uint64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     node.Node
		wantCode uint64
		wantErr  bool
	}{
		{
			name: "",
			fields: fields{
				space: &Space{
					cells: make(map[uint64]*cell),
					cgs: []*CellGroup{
						{
							id:     "test-node",
							cRange: Range{0, math.MaxUint64, math.MaxUint64},
							cells:  make(map[uint64]*cell),
						},
					},
					tf: func(values []interface{}, sfc curve.Curve) ([]uint64, error) {
						return []uint64{0, 0}, nil
					},
					load: 0,
				},
				of: nil,
			},
			args: args{
				d:    nil,
				ncID: 42,
			},
			want:     nil,
			wantCode: 42,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}

			d := &mocks.DataItem{}
			d.On("Values").Return([]interface{}{})
			d.On("ID").Return("test-di")
			d.On("Size").Return(uint64(1))

			sfc := &mocks.Curve{}
			sfc.On("Encode", mock.AnythingOfType("[]uint64")).Return(uint64(0), nil)

			b.space.sfc = sfc

			got, gotCode, err := b.RelocateData(d, tt.args.ncID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantCode, gotCode)
			}
		})
	}
}

func TestBalancer_Optimize(t *testing.T) {
	type fields struct {
		space *Space
		of    OptimizerFunc
	}
	tests := []struct {
		name    string
		fields  fields
		wantCgs []*CellGroup
		wantErr bool
	}{
		{
			name: "optimize",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   []*CellGroup{{id: "test-group"}},
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return []*CellGroup{{id: "altered-test-group"}}, nil
				},
			},
			wantCgs: []*CellGroup{{id: "altered-test-group"}},
			wantErr: false,
		},
		{
			name: "err optimize",
			fields: fields{
				space: &Space{
					cells: nil,
					cgs:   []*CellGroup{{id: "test-group"}},
					sfc:   nil,
					tf:    nil,
					load:  0,
				},
				of: func(s *Space) ([]*CellGroup, error) {
					return nil, errors.New("test err")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				space: tt.fields.space,
				of:    tt.fields.of,
			}
			err := b.Optimize()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCgs, b.space.cgs)
			}
		})
	}
}
