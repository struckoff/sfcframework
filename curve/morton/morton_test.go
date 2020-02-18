package morton

import (
	"math/big"
	"reflect"
	"testing"
)

func TestMortonCurve_Decode(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
		length     uint64
	}
	type args struct {
		d *big.Int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantCoords []uint64
		wantErr    bool
	}{
		{
			"3 == [1, 1]",
			fields{
				2,
				1,
				10,
			},
			args{
				big.NewInt(3),
			},
			[]uint64{
				1,1,
			},
			false,
		},
		{
			"96 == [8, 4]",
			fields{
				2,
				10,
				20,
			},
			args{
				big.NewInt(96),
			},
			[]uint64{
				8,4,
			},
			false,
		},
		{
			"1096 == [40, 2]",
			fields{
				2,
				10,
				20,
			},
			args{
				big.NewInt(1096),
			},
			[]uint64{
				40,2,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &MortonCurve{
				dimensions: tt.fields.dimensions,
				bits:       tt.fields.bits,
				length:     tt.fields.length,
			}
			gotCoords, err := c.Decode(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCoords, tt.wantCoords) {
				t.Errorf("Decode() gotCoords = %v, want %v", gotCoords, tt.wantCoords)
			}
		})
	}
}