package optimizer

import (
	"github.com/stretchr/testify/assert"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/mocks"
	"testing"
)

func TestRangeOptimizer(t *testing.T) {
	type args struct {
		loadSet []uint64
		rates   []int
		powers  []float64
		cType   curve.CurveType
		dims    uint64
		bits    uint64
	}
	tests := []struct {
		name      string
		args      args
		wantRates []int
		wantErr   bool
	}{
		{
			"equal 4 nodes",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0, 0},
				powers:  []float64{1, 1, 1, 1},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[]int{1024, 1024, 1024, 1024},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgs := mocks.GenerateMockCellGroup(tt.args.loadSet, tt.args.rates, tt.args.powers, nil)
			rgs := mocks.GenerateMockCellGroup(tt.args.loadSet, tt.wantRates, tt.args.powers, nil)
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)
			s := balancer.NewMockSpace(cgs, sfc)

			got, err := RangeOptimizer(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("RangeOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(rgs) {
				t.Errorf("RangeOptimizer() different amounf of cell groups got = %v, want %v", len(got), len(rgs))
				return
			}
			for iter := range got {
				assert.Equal(t, rgs[iter].Range(), got[iter].Range())
				assert.Equal(t, rgs[iter].TotalLoad(), got[iter].TotalLoad())
				//assert.Equal(t, rgs[iter].Cells(), got[iter].Cells())
				//if ok, msg := mocks.CompareCellGroup(got[iter], rgs[iter]); !ok {
				//	t.Errorf("RangeOptimizer() %s", msg)
				//	return
				//}
			}
		})
	}
}

func TestPowerRangeOptimizer(t *testing.T) {
	type args struct {
		loadSet []uint64
		rates   []int
		powers  []float64
		caps    []float64
		cType   curve.CurveType
		dims    uint64
		bits    uint64
	}
	tests := []struct {
		name      string
		args      args
		wantRates []int
		wantErr   bool
	}{
		{
			"equal 4 nodes",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0, 0},
				powers:  []float64{1, 1, 1, 1},
				caps:    []float64{1000, 1000, 1000, 1000},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[]int{1024, 1024, 1024, 1024},
			false,
		},
		{
			"not equal 4 nodes",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0, 0},
				powers:  []float64{1, 2, 3, 4},
				caps:    []float64{1000, 1000, 1000, 1000},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[]int{410, 819, 1229, 1638},
			false,
		},
		{
			"2 equal 1 not",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0},
				powers:  []float64{1, 1, 5},
				caps:    []float64{1000, 1000, 1000},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[]int{585, 585, 2925},
			false,
		},
		{
			"2 equal 1 not(in middle)",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0},
				powers:  []float64{1, 5, 1},
				caps:    []float64{1000, 1000, 1000},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[]int{585, 2925, 585},
			false,
		},
		{
			"11 equal",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				powers:  []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				caps:    []float64{1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000},
				cType:   curve.Morton,
				dims:    2,
				bits:    6,
			},
			[]int{372, 372, 372, 372, 372, 372, 372, 372, 372, 372, 376},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgs := mocks.GenerateMockCellGroup(tt.args.loadSet, tt.args.rates, tt.args.powers, tt.args.caps)
			rgs := mocks.GenerateMockCellGroup(tt.args.loadSet, tt.wantRates, tt.args.powers, tt.args.caps)
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)
			s := balancer.NewMockSpace(cgs, sfc)

			got, err := PowerRangeOptimizer(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("RangeOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(rgs) {
				t.Errorf("RangeOptimizer() different amounf of cell groups got = %v, want %v", len(got), len(rgs))
				return
			}
			for iter := range got {
				assert.Equal(t, rgs[iter].Range(), got[iter].Range())
				assert.Equal(t, rgs[iter].TotalLoad(), got[iter].TotalLoad())
			}
		})
	}
}
