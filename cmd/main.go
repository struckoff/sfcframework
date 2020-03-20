package main

import (
	"fmt"
	"github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
	"math/rand"
)

func main() {
	rates := make(map[string]int)
	ratesDi := make(map[[3]float64]int)

	node0 := balancer.NewMockNode("node-0", 1, 20)
	node1 := balancer.NewMockNode("node-1", 2, 20)
	node2 := balancer.NewMockNode("node-2", 1, 20)

	nodes := []balancer.Node{node0, node1}
	bal, err := balancer.NewBalancer(curve.Morton, 3, 16, transform.SpaceTransform,
		optimizer.RangeOptimizer, nodes)
	if err != nil {
		panic(err)
	}
	if err := bal.AddNode(node2); err != nil {
		panic(err)
	}

	for iter := uint64(0); iter < 10; iter++ {
		vals := make([]interface{}, 3)
		vals[0] = float64(rand.Intn(90))
		vals[1] = float64(rand.Intn(180))
		vals[2] = int64(iter + 1609459200)
		di := balancer.NewMockDataItem(fmt.Sprintf("di-%d", iter), 1024*iter, vals)
		if _, err := bal.AddData(di); err != nil {
			panic(err)
		}
	}
	if err := bal.Optimize(); err != nil {
		panic(err)
	}
	for iter := uint64(0); iter < 100000; iter++ {
		vals := make([]interface{}, 3)
		lon := -90 + rand.Float64()  * 180
		lat := -180 + rand.Float64() * 360
		ts := int64(iter + 1609459200)
		vals[0] = lon
		vals[1] = lat
		vals[2] = ts
		k := [3]float64{lon,lat,float64(ts)}
		ratesDi[k]++
		di := balancer.NewMockDataItem(fmt.Sprintf("di-%d", iter), 1024*iter, vals)
		if n, err := bal.AddData(di); err != nil {
			panic(err)
		} else {
			rates[n.ID()]++
		}
	}

	for key := range ratesDi{
		//fmt.Println(key)
		if (ratesDi[key]==1){
			delete(ratesDi, key)
		}
	}
	fmt.Println(rates, ratesDi)
	//fmt.Println(rates, ratesDi, bal.Space().Rates())

}
