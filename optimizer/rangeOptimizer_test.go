package optimizer

import (
	"fmt"
	"testing"

	"github.com/struckoff/sfcframework/node"

	"github.com/stretchr/testify/assert"
	balancer "github.com/struckoff/sfcframework"
	"github.com/struckoff/sfcframework/curve"
	"github.com/struckoff/sfcframework/mocks"
)

func TestRangeOptimizer(t *testing.T) {
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
		name       string
		args       args
		wantRanges [][2]uint64
		wantErr    bool
	}{
		{
			"no nodes",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{},
				powers:  []float64{},
				caps:    []float64{},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[][2]uint64{},
			false,
		},
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
			[][2]uint64{
				{0, 1024},
				{1024, 2048},
				{2048, 3072},
				{3072, 4096},
			},
			false,
		},
		{
			"not equal 4 nodes",
			args{
				loadSet: make([]uint64, 4096),
				rates:   []int{4096, 0, 0, 0},
				powers:  []float64{0, 1, 0, 1},
				caps:    []float64{1000, 1000, 1000, 1000},
				cType:   curve.Morton,
				dims:    3,
				bits:    4,
			},
			[][2]uint64{
				{0, 0},
				{0, 2048},
				{2048, 2048},
				{2048, 4096},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)

			var nodes []node.Node
			var rgs []*balancer.CellGroup

			if len(tt.args.powers) > 0 {
				nodes = make([]node.Node, len(tt.args.powers))
				rgs = make([]*balancer.CellGroup, len(tt.args.powers))
			}

			for i := range tt.args.powers {
				p := &mocks.Power{}
				p.On("Get").Return(tt.args.powers[i])
				c := &mocks.Capacity{}
				c.On("Get").Return(tt.args.caps[i], nil)
				n := &mocks.Node{}
				n.On("Power").Return(p)
				n.On("Capacity").Return(c)
				n.On("Hash").Return(uint64(i))
				n.On("ID").Return(fmt.Sprintf("node-%d", i))
				nodes[i] = n
				rgs[i] = balancer.NewCellGroup(n)
				err := rgs[i].SetRange(tt.wantRanges[i][0], tt.wantRanges[i][1])
				if err != nil {
					t.Fatal(err)
				}
			}

			s, err := balancer.NewSpace(sfc, nil, nodes)
			if err != nil {
				t.Fatal(err)
			}

			got, err := RangeOptimizer(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("RangeOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, rgs, got)
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
		name       string
		args       args
		wantRanges [][2]uint64
		wantErr    bool
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
			[][2]uint64{
				{0, 1024},
				{1024, 2048},
				{2048, 3072},
				{3072, 4096},
			},
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
			[][2]uint64{
				{0, 410},
				{410, 1229},
				{1229, 2458},
				{2458, 4096},
			},
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
			[][2]uint64{
				{0, 585},
				{585, 1170},
				{1170, 4096},
			},
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
			[][2]uint64{
				{0, 585},
				{585, 3510},
				{3510, 4096},
			},
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
			[][2]uint64{
				{0, 372},
				{372, 744},
				{744, 1116},
				{1116, 1488},
				{1488, 1860},
				{1860, 2232},
				{2232, 2604},
				{2604, 2976},
				{2976, 3348},
				{3348, 3720},
				{3720, 4096},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfc, _ := curve.NewCurve(tt.args.cType, tt.args.dims, tt.args.bits)

			nodes := make([]node.Node, len(tt.args.powers))
			rgs := make([]*balancer.CellGroup, len(tt.args.powers))
			for i := range tt.args.powers {
				p := &mocks.Power{}
				p.On("Get").Return(tt.args.powers[i])
				c := &mocks.Capacity{}
				c.On("Get").Return(tt.args.caps[i], nil)
				n := &mocks.Node{}
				n.On("Power").Return(p)
				n.On("Capacity").Return(c)
				n.On("Hash").Return(uint64(i))
				n.On("ID").Return(fmt.Sprintf("node-%d", i))
				nodes[i] = n
				rgs[i] = balancer.NewCellGroup(n)
				err := rgs[i].SetRange(tt.wantRanges[i][0], tt.wantRanges[i][1])
				if err != nil {
					t.Fatal(err)
				}
			}

			s, err := balancer.NewSpace(sfc, nil, nodes)
			if err != nil {
				t.Fatal(err)
			}

			got, err := PowerRangeOptimizer(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("PowerRangeOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, rgs, got)
		})
	}
}
