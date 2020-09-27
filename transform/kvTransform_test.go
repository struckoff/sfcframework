package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/struckoff/sfcframework/curve"
)

func TestKVTransform(t *testing.T) {
	type args struct {
		values []interface{}
		dims   uint64
		bits   uint64
		cType  curve.CurveType
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "key 3x4",
			args: args{
				[]interface{}{"key"},
				3,
				4,
				curve.Morton,
			},
			want:    []uint64{2, 11, 1},
			wantErr: false,
		},
		{
			name: "key 8x4",
			args: args{
				[]interface{}{"key"},
				8,
				4,
				curve.Morton,
			},
			want:    []uint64{2, 11, 1, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			name: "key 1x4",
			args: args{
				[]interface{}{"key"},
				1,
				4,
				curve.Morton,
			},
			want:    []uint64{14},
			wantErr: false,
		},
		{
			name: "not enough values",
			args: args{
				[]interface{}{},
				1,
				4,
				curve.Morton,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not string",
			args: args{
				[]interface{}{42},
				1,
				4,
				curve.Morton,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)
			got, err := KVTransform(tt.args.values, sfc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, got)
				assert.NoError(t, err)
			}
		})
	}
}
