package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/struckoff/sfcframework/curve"
)

func TestSpaceTransform(t *testing.T) {
	type args struct {
		values []interface{}
		cType  curve.CurveType
		bits   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []uint64
		wantErr bool
	}{
		{
			name: "Hilbert 8 bits",
			args: args{
				values: []interface{}{90.0, 180.0},
				cType:  curve.Hilbert,
				bits:   8,
			},
			want:    []uint64{255, 255},
			wantErr: false,
		},
		{
			name: "not enough values",
			args: args{
				values: []interface{}{90.0},
				cType:  curve.Hilbert,
				bits:   8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "latitude not float64",
			args: args{
				values: []interface{}{"90.0", 180.0},
				cType:  curve.Hilbert,
				bits:   8,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "longitude not float64",
			args: args{
				values: []interface{}{90.0, "180.0"},
				cType:  curve.Hilbert,
				bits:   8,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc, err := curve.NewCurve(tt.args.cType, 2, tt.args.bits)
			if err != nil {
				t.Error(err)
				return
			}
			got, err := SpaceTransform(tt.args.values, sfc)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpaceTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
