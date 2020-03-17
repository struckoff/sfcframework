package morton

import (
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"testing"
)

func TestMortonCurve_Decode(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
	}
	type args struct {
		code uint64
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
				3,
			},
			[]uint64{
				1, 1,
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
				96,
			},
			[]uint64{
				8, 4,
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
				1096,
			},
			[]uint64{
				40, 2,
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
				math.MaxInt32,
			},
			[]uint64{
				65535, 32767,
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
				math.MaxInt64,
			},
			[]uint64{
				4294967295, 2147483647,
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
				6442450941,
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
			if err != nil {
				t.Fatal(err)
			}
			gotCoords, err := c.Decode(tt.args.code)
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
		name     string
		fields   fields
		args     args
		wantCode uint64
		wantErr  bool
	}{
		{
			"3 == [1, 1]",
			fields{
				2,
				1,
			},
			args{
				[]uint64{
					1, 1,
				},
			},
			3,
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
					8, 4,
				},
			},
			96,
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
					40, 2,
				},
			},
			1096,
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
					65535, 32767,
				},
			},
			math.MaxInt32,
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
					4294967295, 2147483647,
				},
			},
			math.MaxInt64,
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
			6442450941,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.fields.dimensions, tt.fields.bits)
			if err != nil {
				t.Fatal(err)
			}
			gotD, err := c.Encode(tt.args.coords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotD, tt.wantCode) {
				t.Errorf("Encode() gotD = %v, want %v", gotD, tt.wantCode)
			}
		})
	}
}

func BenchmarkCurve_Decode(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	c, err := New(2, 10)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := uint64(0); i < uint64(b.N); i++ {
		log.Print(c.Decode(i))
	}
}

func BenchmarkCurve_Decode_Morton(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	type args struct {
		dims uint64
		bits uint64
	}
	benchmarks := []struct {
		name string
		args args
	}{
		{
			"2x2",
			args{dims: 2, bits: 2},
		},
		{
			"2x10",
			args{dims: 2, bits: 10},
		},
		{
			"32x2",
			args{dims: 32, bits: 2},
		},
		{
			"32x512",
			args{dims: 32, bits: 512},
		},
	}
	for _, bm := range benchmarks {
		c, err := New(bm.args.dims, bm.args.bits)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := uint64(0); i < uint64(b.N); i++ {
				log.Print(c.Decode(i))
			}
		})
	}
}

func BenchmarkCurve_Encode_Morton(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	type args struct {
		dims uint64
		bits uint64
	}
	benchmarks := []struct {
		name string
		args args
	}{
		{
			"2x2",
			args{dims: 2, bits: 2},
		},
		{
			"2x10",
			args{dims: 2, bits: 10},
		},
		{
			"32x2",
			args{dims: 32, bits: 2},
		},
		{
			"32x512",
			args{dims: 32, bits: 512},
		},
	}

	for _, bm := range benchmarks {
		c, err := New(bm.args.dims, bm.args.bits)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(bm.name, func(b *testing.B) {
			b.StopTimer()
			b.ReportAllocs()
			coordsSet := [][]uint64{}
			for i := uint64(0); i < uint64(b.N); i++ {
				coord, _ := c.Decode(i)
				coordsSet = append(coordsSet, coord)
			}
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				log.Print(c.Encode(coordsSet[i]))
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		dims uint64
		bits uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *Curve
		wantErr bool
	}{
		{
			"zero dimensions",
			args{dims: 0, bits: 4},
			nil,
			true,
		},
		{
			"zero bits",
			args{dims: 4, bits: 0},
			nil,
			true,
		},
		{
			"zero dimensions and bits",
			args{dims: 0, bits: 0},
			nil,
			true,
		},
		{
			"2x4",
			args{dims: 2, bits: 4},
			&Curve{
				dimensions: 2,
				bits:       4,
				length:     4,
				masksArray: []uint64{
					0xF,
					0x33,
					0x55,
				},
				maxSize: 15,
				maxCode: 255,
			},
			false,
		},
		{
			"4x32",
			args{dims: 4, bits: 32},
			&Curve{
				dimensions: 4,
				bits:       32,
				length:     96,
				masksArray: []uint64{
					0xffffffff,
					0x3fffff,
					0x3ff800000007ff,
					0xf80007c0003f,
					0xc0380700c03807,
					0x843084308430843,
					0x909090909090909,
					0x1111111111111111,
				},
				maxSize: 4294967295,
				maxCode: 18446744073709551615,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.dims, tt.args.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurve_validateCoordinates(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
		length     uint64
		masksArray []uint64
		maxSize    uint64
		maxCode    uint64
	}
	type args struct {
		coords []uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				dimensions: 2,
				bits:       4,
				length:     4,
				maxSize:    15,
				maxCode:    255,
			},
			args: args{
				coords: []uint64{4, 12},
			},
			wantErr: false,
		},
		{
			name: "not enough coordinates",
			fields: fields{
				dimensions: 2,
				bits:       4,
				length:     4,
				maxSize:    15,
				maxCode:    255,
			},
			args: args{
				coords: []uint64{4},
			},
			wantErr: true,
		},
		{
			name: "coordinate exceeds limit",
			fields: fields{
				dimensions: 2,
				bits:       4,
				length:     8,
				maxSize:    15,
				maxCode:    255,
			},
			args: args{
				coords: []uint64{4, 120},
			},
			wantErr: true,
		},
		{
			name: "all coordinates exceeds limit",
			fields: fields{
				dimensions: 2,
				bits:       4,
				length:     4,
				maxSize:    15,
				maxCode:    255,
			},
			args: args{
				coords: []uint64{400, 120},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Curve{
				dimensions: tt.fields.dimensions,
				bits:       tt.fields.bits,
				length:     tt.fields.length,
				masksArray: tt.fields.masksArray,
				maxSize:    tt.fields.maxSize,
				maxCode:    tt.fields.maxCode,
			}
			if err := c.validateCoordinates(tt.args.coords); (err != nil) != tt.wantErr {
				t.Errorf("validateCoordinates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCurve_validateCode(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
		length     uint64
		masksArray []uint64
		maxSize    uint64
		maxCode    uint64
	}
	type args struct {
		code uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				dimensions: 2,
				bits:       4,
				length:     4,
				maxSize:    15,
				maxCode:    255,
			},
			args: args{
				96,
			},
			wantErr: false,
		},
		{
			name: "code exceeds limit",
			fields: fields{
				dimensions: 2,
				bits:       4,
				length:     4,
				maxSize:    15,
				maxCode:    255,
			},
			args: args{
				412,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Curve{
				dimensions: tt.fields.dimensions,
				bits:       tt.fields.bits,
				length:     tt.fields.length,
				masksArray: tt.fields.masksArray,
				maxSize:    tt.fields.maxSize,
				maxCode:    tt.fields.maxCode,
			}
			if err := c.validateCode(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("validateCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
