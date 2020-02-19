package balancer

import (
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
				space: space{
					cg: []cellGroup{newCellGroup(testNode)},
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
				space: space{
					cg: []cellGroup{newCellGroup(testNode)},
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
