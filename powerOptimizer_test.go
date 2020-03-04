package balancer

import (
	"reflect"
	"testing"
)

func generateCells() []cell {
	return []cell{
		{load: 0},
		{load: 0},
		{load: 10},
		{load: 20},
		{load: 0},
		{load: 0},
		{load: 80},
		{load: 0},
		{load: 60},
		{load: 0},
		{load: 40},
		{load: 0},
		{load: 90},
		{load: 0},
		{load: 0},
	}
}

func generateTestResult(cs []cell) []cellGroup {
	cgs := make([]cellGroup, 3)
	cg := newCellGroup(MockNode{power: MockPower{value: 10.0}})
	cg.cells = append(cg.cells, &cs[0])
	cg.cells = append(cg.cells, &cs[1])
	cg.cells = append(cg.cells, &cs[2])
	cg.cells = append(cg.cells, &cs[3])
	cg.cells = append(cg.cells, &cs[4])
	cg.cells = append(cg.cells, &cs[5])
	cg.cells = append(cg.cells, &cs[6])
	cg.load = 110
	cgs[0] = cg
	cg = newCellGroup(MockNode{power: MockPower{value: 10.0}})
	cg.cells = append(cg.cells, &cs[7])
	cg.cells = append(cg.cells, &cs[8])
	cg.cells = append(cg.cells, &cs[9])
	cg.cells = append(cg.cells, &cs[10])
	cg.load = 100
	cgs[1] = cg
	cg = newCellGroup(MockNode{power: MockPower{value: 10.0}})
	cg.cells = append(cg.cells, &cs[11])
	cg.cells = append(cg.cells, &cs[12])
	cg.cells = append(cg.cells, &cs[13])
	cg.cells = append(cg.cells, &cs[14])
	cg.load = 90
	cgs[2] = cg
	return cgs
}

func generateTestCase(cs []cell) []cellGroup {
	cgs := make([]cellGroup, 3)
	cg := newCellGroup(MockNode{power: MockPower{value: 10.0}})
	cg.cells = append(cg.cells, &cs[0])
	cg.cells = append(cg.cells, &cs[1])
	cg.cells = append(cg.cells, &cs[2])
	cg.cells = append(cg.cells, &cs[3])
	cg.cells = append(cg.cells, &cs[4])
	cgs[0] = cg
	cg = newCellGroup(MockNode{power: MockPower{value: 10.0}})
	cg.cells = append(cg.cells, &cs[5])
	cg.cells = append(cg.cells, &cs[6])
	cg.cells = append(cg.cells, &cs[7])
	cg.cells = append(cg.cells, &cs[8])
	cg.cells = append(cg.cells, &cs[9])
	cgs[1] = cg
	cg = newCellGroup(MockNode{power: MockPower{value: 10.0}})
	cg.cells = append(cg.cells, &cs[10])
	cg.cells = append(cg.cells, &cs[11])
	cg.cells = append(cg.cells, &cs[12])
	cg.cells = append(cg.cells, &cs[13])
	cg.cells = append(cg.cells, &cs[14])
	cgs[2] = cg
	return cgs
}

func TestPowerOptimizer(t *testing.T) {
	type args struct {
		cgs []cellGroup
	}
	cs := generateCells()
	cgs := generateTestCase(cs)
	rgs := generateTestResult(cs)
	tests := []struct {
		name    string
		args    args
		want    []cellGroup
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
			got, err := PowerOptimizer(tt.args.cgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("PowerOptimizer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PowerOptimizer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
