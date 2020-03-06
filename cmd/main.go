package main

import (
	"fmt"
	"github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/powerOptimizer"
	"github.com/struckoff/SFCFramework/spaceTransform"
)

func main() {
	bal, err := SFCFramework.NewBalancer(curve.Morton, 3, 32, spaceTransform.SpaceTransform, powerOptimizer.PowerOptimizer)
	if err != nil {
		panic(err)
	}

	node0 := SFCFramework.NewMockNode("node-0", 10, 20)
	if err := bal.AddNode(node0); err != nil {
		panic(err)
	}
	node1 := SFCFramework.NewMockNode("node-1", 10, 20)
	if err := bal.AddNode(node1); err != nil {
		panic(err)
	}

	for iter := uint64(0); iter < 10; iter++ {
		vals := make([]interface{}, 3)
		vals[0] = float64(iter * 4)
		vals[1] = float64(iter * 12)
		vals[2] = int64(iter*10 + 1609459200)
		di := SFCFramework.NewDefaultDataItem(fmt.Sprintf("di-%d", iter), 1024*iter, vals)
		if err := bal.AddData(di); err != nil {
			panic(err)
		}
	}
	node2 := SFCFramework.NewMockNode("node-2", 10, 20)
	if err := bal.AddNode(node2); err != nil {
		panic(err)
	}
	fmt.Printf("bal.Distribution() == %v\n", bal.Distribution())
}
