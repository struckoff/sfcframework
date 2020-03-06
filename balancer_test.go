package balancer

import (
	"github.com/struckoff/SFCFramework/curve/hilbert"
	"github.com/struckoff/SFCFramework/curve/morton"
	"reflect"
	"testing"
)

var testNode = MockNode{
	id:       "node-1",
	power:    MockPower{value: 10},
	capacity: MockCapacity{value: 20},
}

func TestBalancer_AddNode(t *testing.T) {
	tests := []struct {
		name     string
		balancer *Balancer
		node     Node
		wantErr  bool
	}{
		{
			name:     "empty balancer",
			balancer: &Balancer{},
			node: MockNode{
				id:       "node-1",
				power:    MockPower{value: 10},
				capacity: MockCapacity{value: 20},
			},
			wantErr: false,
		},
		{
			name: "non-empty balancer",
			balancer: &Balancer{
				nType: reflect.TypeOf(testNode),
				space: &Space{
					cgs: []CellGroup{NewCellGroup(testNode)},
				},
			},
			node: MockNode{
				id:       "node-2",
				power:    MockPower{value: 10},
				capacity: MockCapacity{value: 20},
			},
			wantErr: false,
		},
		{
			name: "nil node",
			balancer: &Balancer{
				nType: reflect.TypeOf(testNode),
				space: &Space{
					cgs: []CellGroup{NewCellGroup(testNode)},
				},
			},
			node:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.balancer.AddNode(tt.node); (err != nil) != tt.wantErr {
				t.Errorf("Balancer.AddNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBalancer_AddData(t *testing.T) {
	type fields struct {
		nType reflect.Type
		space *Space
	}
	type args struct {
		d DataItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Morton curve",
			fields{
				nType: reflect.TypeOf(testNode),
				space: &Space{
					cgs: []CellGroup{NewCellGroup(testNode)},
					sfc: morton.Curve{},
				},
			},
			args{
				MockDataItem{
					id:     "item-0",
					size:   42,
					values: []uint64{4, 12},
				},
			},
			false,
		},
		{
			"Hilbert curve",
			fields{
				nType: reflect.TypeOf(testNode),
				space: &Space{
					cgs: []CellGroup{NewCellGroup(testNode)},
					sfc: hilbert.Curve{},
				},
			},
			args{
				MockDataItem{
					id:     "item-0",
					size:   42,
					values: []uint64{4, 12},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Balancer{
				nType: tt.fields.nType,
				space: tt.fields.space,
			}
			//if err := b.AddNode(testNode); (err != nil) != tt.wantErr {
			//	t.Errorf("Balancer.AddNode() error = %v, wantErr %v", err, tt.wantErr)
			//}
			if err := b.AddData(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("AddData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
