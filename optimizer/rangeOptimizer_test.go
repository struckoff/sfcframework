package optimizer

import (
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
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
			cs := balancer.GenerateMockCells(tt.args.loadSet...)
			cgs := balancer.GenerateMockCellGroup(cs, tt.args.rates, tt.args.powers)
			rgs := balancer.GenerateMockCellGroup(cs, tt.wantRates, tt.args.powers)
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)
			s := balancer.NewMockSpace(cgs, cs, sfc)

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
				if ok, msg := balancer.CompareCellGroup(got[iter], rgs[iter]); !ok {
					t.Errorf("RangeOptimizer() %s", msg)
					return
				}
			}
			//if !reflect.DeepEqual(got, rgs) {
			//	for _, cg := range got{
			//		fmt.Println(cg.)
			//	}
			//	t.Errorf("RangeOptimizer() got = %v, want %v", got, rgs)
			//}
		})
	}
}
