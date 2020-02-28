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

func TestCurve_Size(t *testing.T) {
	type fields struct {
		bits uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint
	}{
		{
			"2 == 3",
			fields{2},
			3,
		},
		{
			"2 == 3",
			fields{2},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Curve{
				bits: tt.fields.bits,
			}
			if got := c.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
