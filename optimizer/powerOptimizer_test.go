package optimizer

import (
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/mocks"
	"testing"
)

//TODO FIX this

func TestPowerOptimizer(t *testing.T) {
	type args struct {
		loadSet []uint64
		rates   []int
		powers  []float64
	}
	//cs := generateCells()

	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "test equal power",
			args: args{
				loadSet: []uint64{0, 0, 10, 20, 0, 0, 80, 0, 60, 0, 40, 0, 90, 0, 0},
				rates:   []int{5, 5, 5},
				powers:  []float64{10, 10, 10},
			},
			want:    []int{7, 4, 4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cgs := mocks.GenerateMockCellGroup(tt.args.loadSet, tt.args.rates, tt.args.powers, nil)
			rgs := mocks.GenerateMockCellGroup(tt.args.loadSet, tt.want, tt.args.powers, nil)

			s := balancer.NewMockSpace(cgs, nil)
			got, err := PowerOptimizer(s)

			if (err != nil) != tt.wantErr {
				t.Errorf("PowerOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(rgs) {
				t.Errorf("PowerOptimizer() different amounf of cell groups got = %v, want %v", len(got), len(rgs))
				return
			}

			for iter := range got {
				assert.Equal(t, rgs[iter].Range(), got[iter].Range())
				assert.Equal(t, rgs[iter].TotalLoad(), got[iter].TotalLoad())
			}
		})
	}
}

//
//func TestPowerOptimizerGreedy(t *testing.T) {
//	type args struct {
//		loadSet []uint64
//		rates   []int
//		powers  []float64
//	}
//	//cs := generateCells()
//
//	tests := []struct {
//		name    string
//		args    args
//		want    []int
//		wantErr bool
//	}{
//		{
//			name: "test equal power",
//			args: args{
//				loadSet: []uint64{0, 0, 10, 20, 0, 0, 80, 0, 60, 0, 40, 0, 90, 0, 0},
//				rates:   []int{5, 5, 5},
//				powers:  []float64{10, 10, 10},
//			},
//			want:    []int{7, 4, 4},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cs := mocks.GenerateMockCells(tt.args.loadSet...)
//			cgs := mocks.GenerateMockCellGroup(cs, tt.args.rates, tt.args.powers)
//			rgs := mocks.GenerateMockCellGroup(cs, tt.want, tt.args.powers)
//
//			s := balancer.NewMockSpace(cgs, cs, nil)
//			got, err := PowerOptimizerGreedy(s)
//
//			if (err != nil) != tt.wantErr {
//				t.Errorf("PowerOptimizerGreedy() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, rgs) {
//				t.Errorf("PowerOptimizerGreedy() got = %v, want %v", got, rgs)
//			}
//		})
//	}
//}
//
//func TestPowerOptimizerBreezy(t *testing.T) {
//	type args struct {
//		loadSet []uint64
//		rates   []int
//		powers  []float64
//	}
//	//cs := generateCells()
//
//	tests := []struct {
//		name    string
//		args    args
//		want    []int
//		wantErr bool
//	}{
//		{
//			name: "test equal power",
//			args: args{
//				loadSet: []uint64{0, 0, 10, 20, 0, 0, 80, 0, 60, 0, 40, 0, 90, 0, 0},
//				rates:   []int{5, 5, 5},
//				powers:  []float64{10, 100, 10},
//			},
//			want:    []int{7, 4, 4},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			cs := mocks.GenerateMockCells(tt.args.loadSet...)
//			cgs := mocks.GenerateMockCellGroup(cs, tt.args.rates, tt.args.powers)
//			cs = mocks.GenerateMockCells(tt.args.loadSet...)
//			rgs := mocks.GenerateMockCellGroup(cs, tt.want, tt.args.powers)
//
//			s := balancer.NewMockSpace(cgs, cs, nil)
//			got, err := PowerOptimizerPerms(s)
//
//			if (err != nil) != tt.wantErr {
//				t.Errorf("PowerOptimizerPerms() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, rgs) {
//				t.Errorf("PowerOptimizerPerms() got = %v, want %v", got, rgs)
//			}
//		})
//	}
//}
//
//func prepareSpace() *balancer.Space {
//	bal, err := balancer.NewBalancer(curve.Morton, 3, 32, transform.SpaceTransform, PowerOptimizer, nil)
//	if err != nil {
//		panic(err)
//	}
//	node0 := balancer.NewMockNode("node-0", 10, 20)
//	if err := bal.AddNode(node0); err != nil {
//		panic(err)
//	}
//	node1 := balancer.NewMockNode("node-1", 10, 20)
//	if err := bal.AddNode(node1); err != nil {
//		panic(err)
//	}
//	node2 := balancer.NewMockNode("node-2", 10, 20)
//	if err := bal.AddNode(node2); err != nil {
//		panic(err)
//	}
//
//	s := bal.Space()
//	return s
//}
//
//func BenchmarkPowerOptimizer(b *testing.B) {
//	log.SetOutput(ioutil.Discard)
//
//	b.Run("orig", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			b.StopTimer()
//			s := prepareSpace()
//			b.StartTimer()
//			x, _ := PowerOptimizer(s)
//			log.Print(x)
//		}
//	})
//
//	b.Run("greed", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			b.StopTimer()
//			s := prepareSpace()
//			b.StartTimer()
//			x, _ := PowerOptimizerGreedy(s)
//			log.Print(x)
//		}
//	})
//	b.Run("permutations", func(b *testing.B) {
//		b.ReportAllocs()
//		for i := 0; i < b.N; i++ {
//			b.StopTimer()
//			s := prepareSpace()
//			b.StartTimer()
//			x, _ := PowerOptimizerPerms(s)
//			log.Print(x)
//		}
//	})
//}
