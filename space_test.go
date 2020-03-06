package SFCFramework

import (
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/powerOptimizer"
	"github.com/struckoff/SFCFramework/spaceTransform"
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

func Test_space_addNode(t *testing.T) {
	type fields struct {
		cells []cell
		cg    []CellGroup
		sfc   curve.Curve
		tf    TransformFunc
		of    OptimizerFunc
	}
	type args struct {
		n Node
	}
	cs := powerOptimizer.generateCells()
	sfc, _ := curve.NewCurve(curve.Hilbert, 3, 32)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test case",
			fields: fields{
				cells: cs,
				cg:    []CellGroup{generateCellGroup(cs, testNode)},
				sfc:   sfc,
				tf:    spaceTransform.SpaceTransform,
				of:    powerOptimizer.PowerOptimizer,
			},
			args: args{
				n: MockNode{power: MockPower{value: 10.0}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Space{
				cells: tt.fields.cells,
				cgs:   tt.fields.cg,
				sfc:   tt.fields.sfc,
				tf:    tt.fields.tf,
				of:    tt.fields.of,
			}
			if err := s.addNode(tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("addNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
