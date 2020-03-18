package main

import (
	"fmt"
	"github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
)

func main() {
	node0 := balancer.NewMockNode("node-0", 10, 20)
	node1 := balancer.NewMockNode("node-1", 10, 20)
	nodes := []balancer.Node{node0, node1}
	bal, err := balancer.NewBalancer(curve.Morton, 3, 1024, transform.SpaceTransform,
		optimizer.PowerOptimizer, nodes)
	if err != nil {
		panic(err)
	}

	for iter := uint64(0); iter < 10; iter++ {
		vals := make([]interface{}, 3)
		vals[0] = float64(iter * 4)
		vals[1] = float64(iter * 12)
		vals[2] = int64(iter*10 + 1609459200)
		di := balancer.NewMockDataItem(fmt.Sprintf("di-%d", iter), 1024*iter, vals)
		if _, err := bal.AddData(di); err != nil {
			panic(err)
		}
	}
	node2 := balancer.NewMockNode("node-2", 10, 20)
	if err := bal.AddNode(node2); err != nil {
		panic(err)
	}
}
