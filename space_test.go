package balancer

import (
	"github.com/struckoff/SFCFramework/curve/hilbert"
	"github.com/struckoff/SFCFramework/powerOptimizer"
	"io/ioutil"
	"log"
	"testing"
)

func generateCellGroup(cs []cell, n Node) CellGroup {
	cg := NewCellGroup(n)
	cg.cells = append(cg.cells, &cs[0])
	cg.cells = append(cg.cells, &cs[1])
	cg.cells = append(cg.cells, &cs[2])
	cg.cells = append(cg.cells, &cs[3])
	cg.cells = append(cg.cells, &cs[4])
	cg.cells = append(cg.cells, &cs[5])
	cg.cells = append(cg.cells, &cs[6])
	cg.cells = append(cg.cells, &cs[7])
	cg.cells = append(cg.cells, &cs[8])
	cg.cells = append(cg.cells, &cs[9])
	cg.cells = append(cg.cells, &cs[10])
	cg.cells = append(cg.cells, &cs[11])
	cg.cells = append(cg.cells, &cs[12])
	cg.cells = append(cg.cells, &cs[13])
	cg.cells = append(cg.cells, &cs[14])
	cg.load = 300
	return cg
}

//func Test_space_addNode(t *testing.T) {
//	type fields struct {
//		cells []cell
//		cg    []CellGroup
//		sfc   curve.Curve
//		tf    TransformFunc
//		of    OptimizerFunc
//	}
//	type args struct {
//		n Node
//	}
//	cs := powerOptimizer.generateCells()
//	sfc, _ := curve.NewCurve(curve.Hilbert, 3, 32)
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{
//			name: "test case",
//			fields: fields{
//				cells: cs,
//				cg:    []CellGroup{generateCellGroup(cs, testNode)},
//				sfc:   sfc,
//				tf:    spaceTransform.SpaceTransform,
//				of:    powerOptimizer.PowerOptimizer,
//			},
//			args: args{
//				n: MockNode{power: MockPower{value: 10.0}},
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &Space{
//				cells: tt.fields.cells,
//				cgs:   tt.fields.cg,
//				sfc:   tt.fields.sfc,
//				tf:    tt.fields.tf,
//				of:    tt.fields.of,
//			}
//			if err := s.addNode(tt.args.n); (err != nil) != tt.wantErr {
//				t.Errorf("addNode() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

//func prepareSpace() *Space {
//	bal, err := NewBalancer(curve.Morton, 3, 32, spaceTransform.SpaceTransform, powerOptimizer.PowerOptimizer)
//	if err != nil {
//		panic(err)
//	}
//	node0 := NewMockNode("node-0", 10, 20)
//	if err := bal.AddNode(node0); err != nil {
//		panic(err)
//	}
//	node1 := NewMockNode("node-1", 10, 20)
//	if err := bal.AddNode(node1); err != nil {
//		panic(err)
//	}
//	node2 := NewMockNode("node-2", 10, 20)
//	if err := bal.AddNode(node2); err != nil {
//		panic(err)
//	}
//
//	return bal.Space()
//}

func BenchmarkAddData(b *testing.B) {
	log.SetOutput(ioutil.Discard)

	//b.Run("Morton:3x5, 1000 dataItems", func(b *testing.B) {
	//	count := 1000
	//
	//	b.ReportAllocs()
	//	for i := 0; i < b.N; i++ {
	//		b.StopTimer()
	//		sfc, _ := morton.New(3, 5)
	//		cs := make([]cell, sfc.Length())
	//		cgs := GenerateMockCellGroup(cs, []int{1, 1, 1, 1}, []float64{10, 10, 10, 10})
	//		s := NewMockSpace(cgs, cs, sfc)
	//		dis := make([]DataItem, count)
	//		for iter := range dis {
	//			dis[iter] = GenerateRandomMockSpaceItem()
	//		}
	//		b.StartTimer()
	//		for iter := range dis {
	//			if err := s.AddData(dis[iter]); err != nil {
	//				b.Fatal(err)
	//			}
	//		}
	//	}
	//})
	//b.Run("Morton:3x5, 10000 dataItems", func(b *testing.B) {
	//	count := 10000
	//
	//	b.ReportAllocs()
	//	for i := 0; i < b.N; i++ {
	//		b.StopTimer()
	//		sfc, _ := morton.New(3, 5)
	//		cs := make([]cell, sfc.Length())
	//		cgs := GenerateMockCellGroup(cs, []int{1, 1, 1, 1}, []float64{10, 10, 10, 10})
	//		s := NewMockSpace(cgs, cs, sfc)
	//		dis := make([]DataItem, count)
	//		for iter := range dis {
	//			dis[iter] = GenerateRandomMockSpaceItem()
	//		}
	//		b.StartTimer()
	//		for iter := range dis {
	//			if err := s.AddData(dis[iter]); err != nil {
	//				b.Fatal(err)
	//			}
	//		}
	//	}
	//})
	//b.Run("Morton:3x5, 100000 dataItems", func(b *testing.B) {
	//	count := 100000
	//
	//	b.ReportAllocs()
	//	for i := 0; i < b.N; i++ {
	//		b.StopTimer()
	//		sfc, _ := morton.New(3, 5)
	//		cs := make([]cell, sfc.Length())
	//		cgs := GenerateMockCellGroup(cs, []int{1, 1, 1, 1}, []float64{10, 10, 10, 10})
	//		s := NewMockSpace(cgs, cs, sfc)
	//		dis := make([]DataItem, count)
	//		for iter := range dis {
	//			dis[iter] = GenerateRandomMockSpaceItem()
	//		}
	//		b.StartTimer()
	//		for iter := range dis {
	//			if err := s.AddData(dis[iter]); err != nil {
	//				b.Fatal(err)
	//			}
	//		}
	//	}
	//})
	b.Run("Hilbert:3x5, 1000 dataItems", func(b *testing.B) {
		count := 1000

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			sfc, _ := hilbert.New(3, 5)
			cs := make([]cell, sfc.Length())
			cgs := GenerateMockCellGroup(cs, []int{1, 1, 1, 1}, []float64{10, 10, 10, 10})
			s := NewMockSpace(cgs, cs, sfc)
			cgs, _ = powerOptimizer.PowerOptimizer(s)
			s.SetGroups(cgs)
			dis := make([]DataItem, count)
			for iter := range dis {
				dis[iter] = GenerateRandomMockSpaceItem()
			}
			b.StartTimer()
			for iter := range dis {
				if err := s.AddData(dis[iter]); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
	b.Run("Hilbert:3x5, 10000 dataItems", func(b *testing.B) {
		count := 10000

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			sfc, _ := hilbert.New(3, 5)
			cs := make([]cell, sfc.Length())
			cgs := GenerateMockCellGroup(cs, []int{1, 1, 1, 1}, []float64{10, 10, 10, 10})
			s := NewMockSpace(cgs, cs, sfc)
			dis := make([]DataItem, count)
			for iter := range dis {
				dis[iter] = GenerateRandomMockSpaceItem()
			}
			b.StartTimer()
			for iter := range dis {
				if err := s.AddData(dis[iter]); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
	b.Run("Hilbert:3x5, 100000 dataItems", func(b *testing.B) {
		count := 100000

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			sfc, _ := hilbert.New(3, 5)
			cs := make([]cell, sfc.Length())
			cgs := GenerateMockCellGroup(cs, []int{1, 1, 1, 1}, []float64{10, 10, 10, 10})
			s := NewMockSpace(cgs, cs, sfc)
			dis := make([]DataItem, count)
			for iter := range dis {
				dis[iter] = GenerateRandomMockSpaceItem()
			}
			b.StartTimer()
			for iter := range dis {
				if err := s.AddData(dis[iter]); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
