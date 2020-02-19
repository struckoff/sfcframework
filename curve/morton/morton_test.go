package morton

import (
	"math"
	"math/big"
	"reflect"
	"testing"
)

func TestMortonCurve_Decode(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
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
			},
			args{
				big.NewInt(1096),
			},
			[]uint64{
				40,2,
			},
			false,
		},
		{
			"math.MaxInt32 == [65535, 32767]",
			fields{
				2,
				100,
			},
			args{
				big.NewInt(math.MaxInt32),
			},
			[]uint64{
				65535,32767,
			},
			false,
		},
		{
			"math.MaxInt64 == [4294967295, 2147483647]",
			fields{
				2,
				100,
			},
			args{
				big.NewInt(math.MaxInt64),
			},
			[]uint64{
				4294967295,2147483647,
			},
			false,
		},
		{
			"6442450941 == [131071, 32766]",
			fields{
				2,
				100,
			},
			args{
				big.NewInt(6442450941),
			},
			[]uint64{
				131071, 32766,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.fields.dimensions, tt.fields.bits)
			if err!=nil{
				t.Fatal(err)
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

func TestMortonCurve_Encode(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
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
			"3 == [1, 1]",
			fields{
				2,
				1,
			},
			args{
				[]uint64{
					1,1,
				},
			},
			big.NewInt(3),
			false,
		},
		{
			"96 == [8, 4]",
			fields{
				2,
				10,
			},
			args{
				[]uint64{
					8,4,
				},
			},
			big.NewInt(96),
			false,
		},
		{
			"1096 == [40, 2]",
			fields{
				2,
				10,
			},
			args{
				[]uint64{
					40,2,
				},
			},
			big.NewInt(1096),
			false,
		},
		{
			"math.MaxInt32 == [65535, 32767]",
			fields{
				2,
				100,
			},
			args{
				[]uint64{
					65535,32767,
				},
			},
			big.NewInt(math.MaxInt32),
			false,
		},
		{
			"math.MaxInt64 == [4294967295, 2147483647]",
			fields{
				2,
				100,
			},
			args{
				[]uint64{
					4294967295,2147483647,
				},
			},
			big.NewInt(math.MaxInt64),
			false,
		},
		{
			"6442450941 == [131071, 32766]",
			fields{
				2,
				100,
			},
			args{
				[]uint64{
					131071, 32766,
				},
			},
			big.NewInt(6442450941),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.fields.dimensions, tt.fields.bits)
			if err!=nil{
				t.Fatal(err)
			}
			gotD, err := c.Encode(tt.args.coords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantD) {
				t.Errorf("Encode() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}