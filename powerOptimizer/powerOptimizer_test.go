package powerOptimizer

import (
	"github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/spaceTransform"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
)

//TODO FIX this

func TestPowerOptimizer(t *testing.T) {
	type args struct {
		cgs []balancer.CellGroup
	}
	//cs := generateCells()
	cs := balancer.GenerateMockCells(0, 0, 10, 20, 0, 0, 80, 0, 60, 0, 40, 0, 90, 0, 0)
	cgs := balancer.GenerateMockCellGroup(cs, []int{7, 4, 4})
	rgs := balancer.GenerateMockCellGroup(cs, []int{5, 5, 5})
	tests := []struct {
		name    string
		args    args
		want    []balancer.CellGroup
		wantErr bool
	}{
		{
			name:    "test",
			args:    args{cgs: cgs},
			want:    rgs,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := balancer.NewMockSpace(tt.args.cgs, cs)
			got, err := PowerOptimizer(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("PowerOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PowerOptimizer() got = %v, want %v", got[2].Cells(), tt.want[2].Cells())
			}
		})
	}
}

func BenchmarkPowerOptimizer(b *testing.B) {
	b.ReportAllocs()
	bal, err := balancer.NewBalancer(curve.Morton, 3, 32, spaceTransform.SpaceTransform, PowerOptimizer)
	if err != nil {
		panic(err)
	}
	node0 := balancer.NewMockNode("node-0", 10, 20)
	if err := bal.AddNode(node0); err != nil {
		panic(err)
	}
	node1 := balancer.NewMockNode("node-1", 10, 20)
	if err := bal.AddNode(node1); err != nil {
		panic(err)
	}
	node2 := balancer.NewMockNode("node-2", 10, 20)
	if err := bal.AddNode(node2); err != nil {
		panic(err)
	}

	s := bal.Space()

	log.SetOutput(ioutil.Discard)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		x, _ := PowerOptimizer(s)
		log.Print(x)
	}
	b.StopTimer()
}
