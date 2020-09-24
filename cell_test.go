package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/struckoff/sfcframework/mocks"
)

func Test_cell_ID(t *testing.T) {
	type fields struct {
		id uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{"test", fields{42}, 42},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{id: tt.fields.id}
			assert.Equal(t, tt.want, c.ID())
		})
	}
}

func Test_cell_Load(t *testing.T) {
	type fields struct {
		load uint64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "33333",
			fields: fields{load: 33333},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{
				load: &tt.fields.load,
			}

			assert.Equal(t, tt.fields.load, c.Load())
		})
	}
}

func Test_cell_Truncate(t *testing.T) {
	type fields struct {
		load uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   *cell
	}{
		{
			name: "",
			fields: fields{
				load: 1111,
			},
			want: &cell{
				load: uint64ptr(0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{
				load: &tt.fields.load,
			}
			c.Truncate()
			assert.Equal(t, tt.want, c)
		})
	}
}

func Test_cell_Add(t *testing.T) {
	type fields struct {
		load uint64
		cg   *CellGroup
	}
	type args struct {
		Size uint64
	}
	type want struct {
		load uint64
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
				load: 0,
				cg: &CellGroup{
					load: 10,
				},
			},
			args: args{
				Size: 42,
			},
			want: want{
				load: 42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{
				cg:   tt.fields.cg,
				load: &tt.fields.load,
			}

			c.AddLoad(tt.args.Size)

			assert.Equal(t, int(tt.want.load), int(*c.load))
		})
	}
}

func Test_cell_Remove(t *testing.T) {
	type fields struct {
		load uint64
		off  map[string]uint64
		cg   *CellGroup
	}
	type args struct {
		ID   string
		Size uint64
	}
	type want struct {
		off  map[string]uint64
		load uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "dis",
			fields: fields{
				load: 100,
				off:  map[string]uint64{"di-2": 10, "di-3": 20},
				cg: &CellGroup{
					load: 100,
				},
			},
			args: args{
				ID:   "di-0",
				Size: 10,
			},
			want: want{
				off:  map[string]uint64{"di-2": 10, "di-3": 20},
				load: 90,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{
				load: &tt.fields.load,
				off:  tt.fields.off,
				cg:   tt.fields.cg,
			}

			c.RemoveLoad(tt.args.Size)

			assert.Equal(t, tt.want.off, c.off)
			assert.Equal(t, int(tt.want.load), int(*c.load))
		})
	}
}

func Test_cell_Relocate(t *testing.T) {
	type fields struct {
		load uint64
		off  map[string]uint64
		cg   *CellGroup
	}
	type args struct {
		ID   string
		Size uint64
		ncID uint64
	}
	type want struct {
		off    map[string]uint64
		load   uint64
		cgload uint64
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "not exist",
			fields: fields{
				load: 100,
				off:  map[string]uint64{"di-2": 10, "di-3": 20},
				cg: &CellGroup{
					load: 100,
				},
			},
			args: args{
				ID:   "di-0",
				Size: 10,
				ncID: 4242,
			},
			want: want{
				off:    map[string]uint64{"di-0": 4242, "di-2": 10, "di-3": 20},
				load:   100,
				cgload: 100,
			},
		},
		{
			name: "exist",
			fields: fields{
				load: 100,
				off:  map[string]uint64{"di-2": 10, "di-3": 20},
				cg: &CellGroup{
					load: 100,
				},
			},
			args: args{
				ID:   "di-0",
				Size: 10,
				ncID: 4242,
			},
			want: want{
				off:    map[string]uint64{"di-0": 4242, "di-2": 10, "di-3": 20},
				load:   100,
				cgload: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{
				load: &tt.fields.load,
				off:  tt.fields.off,
				cg:   tt.fields.cg,
			}

			di := &mocks.DataItem{}
			di.On("ID").Return(tt.args.ID)
			di.On("Size").Return(tt.args.Size)

			c.Relocate(di, tt.args.ncID)

			assert.Equal(t, tt.want.off, c.off)
			assert.Equal(t, int(tt.want.load), int(*c.load))
			assert.Equal(t, int(tt.want.cgload), int(c.cg.load))
		})
	}
}

func Test_cell_Relocated(t *testing.T) {
	type fields struct {
		off map[string]uint64
	}
	type args struct {
		did string
	}
	type want struct {
		ncID uint64
		ok   bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "true",
			fields: fields{
				off: map[string]uint64{"di-0": 4242, "di-1": 123},
			},
			args: args{
				did: "di-0",
			},
			want: want{
				ncID: 4242,
				ok:   true,
			},
		},
		{
			name: "false",
			fields: fields{
				off: map[string]uint64{"di-0": 4242, "di-1": 123},
			},
			args: args{
				did: "di-3",
			},
			want: want{
				ncID: 0,
				ok:   false,
			},
		},
		{
			name: "empty",
			fields: fields{
				off: make(map[string]uint64),
			},
			args: args{
				did: "di-3",
			},
			want: want{
				ncID: 0,
				ok:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cell{
				off: tt.fields.off,
			}
			ncID, ok := c.Relocated(tt.args.did)
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.ncID, ncID)
		})
	}
}

func TestNewCell(t *testing.T) {
	type args struct {
		id uint64
		cg *CellGroup
	}
	tests := []struct {
		name string
		args args
		want *cell
	}{
		{
			name: "not nil cg",
			args: args{
				id: 11,
				cg: &CellGroup{
					id:    "test-cg",
					cells: make(map[uint64]*cell),
				},
			},
			want: &cell{
				id:   11,
				off:  make(map[string]uint64),
				load: new(uint64),
				cg: &CellGroup{
					id:    "test-cg",
					cells: map[uint64]*cell{},
				},
			},
		},
		{
			name: "nil cg",
			args: args{
				id: 11,
				cg: nil,
			},
			want: &cell{
				id:   11,
				load: new(uint64),
				off:  make(map[string]uint64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCell(tt.args.id, tt.args.cg)
			if tt.args.cg != nil {
				tt.want.cg.cells[c.ID()] = c
			}
			assert.Equal(t, tt.want, c)
		})
	}
}
