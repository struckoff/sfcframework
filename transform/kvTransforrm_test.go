package transform

import (
	"github.com/struckoff/SFCFramework/curve"
	"reflect"
	"testing"
)

func valuesConv(vals ...interface{}) []interface{} {
	res := make([]interface{}, len(vals))
	for iter := range vals {
		res[iter] = vals[iter]
	}
	return res
}

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
				valuesConv("key"),
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
				valuesConv("key"),
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
				valuesConv("key"),
				1,
				4,
				curve.Morton,
			},
			want:    []uint64{14},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)
			got, err := KVTransform(tt.args.values, sfc)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KVTransform() got = %v, want %v", got, tt.want)
			}
		})
	}
}
