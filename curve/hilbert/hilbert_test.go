package hilbert

import (
	"math/big"
	"reflect"
	"testing"
)

func TestHilbertCurve_Decode(t *testing.T) {
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
			"3 == [1, 0]",
			fields{
				2,
				1,
				2,
			},
			args{
				big.NewInt(3),
			},
			[]uint64{
				1, 0,
			},
			false,
		},
		{
			"96 == [4, 12]",
			fields{
				2,
				10,
				20,
			},
			args{
				big.NewInt(96),
			},
			[]uint64{
				4, 12,
			},
			false,
		},
		{
			"1096 == [10, 34]",
			fields{
				2,
				10,
				20,
			},
			args{
				big.NewInt(1096),
			},
			[]uint64{
				10,34,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := HilbertCurve{
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

func TestHilbertCurve_Encode(t *testing.T) {
	type fields struct {
		dimentions uint64
		bits       uint64
		length     uint64
	}
	type args struct {
		coords []uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantD   *big.Int
		wantErr bool
	}{
		{
			"[1, 0] == 3",
			fields{
				2,
				1,
				2,
			},
			args{
				[]uint64{1, 0},
			},
			big.NewInt(3),
			false,
		},
		{
			"[4, 12] == 96",
			fields{
				2,
				10,
				20,
			},
			args{
				[]uint64{4, 12},
			},
			big.NewInt(96) ,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := HilbertCurve{
				dimensions: tt.fields.dimentions,
				bits:       tt.fields.bits,
				length:     tt.fields.length,
			}
			gotD, err := c.Encode(tt.args.coords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotD.Cmp(tt.wantD) != 0 {
				t.Errorf("Encode() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}
