package hilbert

import (
	"io/ioutil"
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHilbertCurve_Decode(t *testing.T) {
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
			"3 == [1, 0]",
			fields{
				2,
				1,
			},
			args{
				3,
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
			},
			args{
				96,
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
			},
			args{
				1096,
			},
			[]uint64{
				10, 34,
			},
			false,
		},
		{
			"MaxInt64 == [4095, 4096, 0, 0, 0]",
			fields{
				5,
				64,
			},
			args{
				math.MaxInt64,
			},
			[]uint64{
				4095, 4096, 0, 0, 0,
			},
			false,
		},
		{
			"not valid code",
			fields{
				2,
				1,
			},
			args{
				math.MaxUint64,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New(tt.fields.dimensions, tt.fields.bits)
			gotCoords, err := c.Decode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantCoords, gotCoords)
		})
	}
}

func TestHilbertCurve_DecodeWithBuffer(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
	}
	type args struct {
		code uint64
		buf  []uint64
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
			},
			args{
				3,
				make([]uint64, 2),
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
			},
			args{
				96,
				make([]uint64, 2),
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
			},
			args{
				1096,
				make([]uint64, 2),
			},
			[]uint64{
				10, 34,
			},
			false,
		},
		{
			"MaxInt64 == [4095, 4096, 0, 0, 0, 0, 0, 0, 0, 0]",
			fields{
				5,
				64,
			},
			args{
				math.MaxInt64,
				make([]uint64, 10),
			},
			[]uint64{
				4095, 4096, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			false,
		},
		{
			"buf < dimensions",
			fields{
				5,
				64,
			},
			args{
				math.MaxInt64,
				make([]uint64, 1),
			},
			nil,
			true,
		},
		{
			"code < maxCode",
			fields{
				2,
				2,
			},
			args{
				math.MaxInt64,
				make([]uint64, 2),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New(tt.fields.dimensions, tt.fields.bits)
			gotCoords, err := c.DecodeWithBuffer(tt.args.buf, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantCoords, gotCoords)
		})
	}
}

func TestHilbertCurve_Encode(t *testing.T) {
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
			"[1, 0] == 3",
			fields{
				2,
				1,
			},
			args{
				[]uint64{1, 0},
			},
			3,
			false,
		},
		{
			"[4, 12] == 96",
			fields{
				2,
				10,
			},
			args{
				[]uint64{4, 12},
			},
			96,
			false,
		},
		{
			"[10, 34] == 1096",
			fields{
				2,
				10,
			},
			args{
				[]uint64{10, 34},
			},
			1096,
			false,
		},
		{
			"[4095, 4096, 0, 0, 0] == MaxInt64",
			fields{
				5,
				64,
			},
			args{
				[]uint64{4095, 4096, 0, 0, 0},
			},
			math.MaxInt64,
			false,
		},
		{
			"not valid coords",
			fields{
				5,
				64,
			},
			args{
				[]uint64{math.MaxUint64, math.MaxUint64},
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New(tt.fields.dimensions, tt.fields.bits)
			gotCode, err := c.Encode(tt.args.coords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int(tt.wantCode), int(gotCode))
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

func BenchmarkCurve_Decode_Hilbert(b *testing.B) {
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

func BenchmarkCurve_Encode_Hilbert(b *testing.B) {
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
				length:     8,
				maxSize:    15,
				maxCode:    255,
			},
			false,
		},
		{
			"4x32",
			args{dims: 4, bits: 32},
			&Curve{
				dimensions: 4,
				bits:       32,
				length:     128,
				maxSize:    4294967295,
				maxCode:    18446744073709551615,
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
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCurve_validateCoordinates(t *testing.T) {
	type fields struct {
		dimensions uint64
		bits       uint64
		length     uint64
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
				length:     8,
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
				length:     8,
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
				length:     8,
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
				length:     8,
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
				length:     8,
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
				maxSize:    tt.fields.maxSize,
				maxCode:    tt.fields.maxCode,
			}
			if err := c.validateCode(tt.args.code); (err != nil) != tt.wantErr {
				t.Errorf("validateCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCurve_DimensionSize(t *testing.T) {
	type fields struct {
		maxSize uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "test",
			fields: fields{
				maxSize: 42,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Curve{
				maxSize: tt.fields.maxSize,
			}
			got := c.DimensionSize()
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}

func TestCurve_Bits(t *testing.T) {
	type fields struct {
		bits uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "test",
			fields: fields{
				bits: 42,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Curve{
				bits: tt.fields.bits,
			}
			got := c.Bits()
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}
func TestCurve_Length(t *testing.T) {
	type fields struct {
		maxCode uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "test",
			fields: fields{
				maxCode: 42,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Curve{
				maxCode: tt.fields.maxCode,
			}
			got := c.Length()
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}
func TestCurve_Dimensions(t *testing.T) {
	type fields struct {
		dimensions uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "test",
			fields: fields{
				dimensions: 42,
			},
			want: 42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Curve{
				dimensions: tt.fields.dimensions,
			}
			got := c.Dimensions()
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}
