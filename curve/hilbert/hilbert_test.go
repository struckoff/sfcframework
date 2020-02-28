package hilbert

import (
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"testing"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New(tt.fields.dimensions, tt.fields.bits)
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
		//{
		//	"[1, 0] == 3",
		//	fields{
		//		2,
		//		1,
		//	},
		//	args{
		//		[]uint64{1, 0},
		//	},
		//	3,
		//	false,
		//},
		//{
		//	"[4, 12] == 96",
		//	fields{
		//		2,
		//		10,
		//	},
		//	args{
		//		[]uint64{4, 12},
		//	},
		//	96,
		//	false,
		//},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := New(tt.fields.dimensions, tt.fields.bits)
			gotCode, err := c.Encode(tt.args.coords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCode != tt.wantCode {
				t.Errorf("Encode() gotCode = %v, want %v", gotCode, tt.wantCode)
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
