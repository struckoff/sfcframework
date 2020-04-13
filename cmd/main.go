package main

import (
	"fmt"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/serialx/hashring"
	"github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
	"math/rand"
)

func main() {
	kvcompare()
}

func kvcompare() {
	var keys []string
	var rates map[string]int
	for iter := uint64(0); iter < 100_000; iter++ {
		key := fmt.Sprintf("key-%d", iter)
		//key := fmt.Sprintf("key-%d", rand.Int())
		keys = append(keys, key)
	}
	nodes := map[string]int{
		"node-0": 1,
		"node-1": 1,
		"node-2": 1,
		"node-3": 1,
	}

	rates = kv(keys, nodes)
	fmt.Println("SFC", rates)
	rates = kvPower(keys, nodes)
	fmt.Println("SFC Power", rates)
	rates = hring(keys, nodes)
	fmt.Println("serialx/hashring", rates)
	rates = consring(keys, nodes)
	fmt.Println("buraksezer/consistent", rates)
}

func kv(keys []string, nodes map[string]int) map[string]int {
	rand.Seed(42)
	var ns []balancer.Node
	for n, w := range nodes {
		ns = append(ns, balancer.NewMockNode(n, float64(w), 20))
	}

	bal, err := balancer.NewBalancer(curve.Morton, 3, 64, transform.KVTransform,
		optimizer.RangeOptimizer, ns)
	if err != nil {
		panic(err)
	}

	if err := bal.Optimize(); err != nil {
		panic(err)
	}

	rates := make(map[string]int)
	for _, key := range keys {
		vals := make([]interface{}, 1)
		vals[0] = key
		di := balancer.NewMockDataItem(key, 1, vals)

		if n, err := bal.LocateData(di); err != nil {
			panic(err)
		} else {
			rates[n.ID()]++
		}
	}
	return rates
}
func kvPower(keys []string, nodes map[string]int) map[string]int {
	rand.Seed(42)
	var ns []balancer.Node
	for n, w := range nodes {
		ns = append(ns, balancer.NewMockNode(n, float64(w), 20))
	}

	bal, err := balancer.NewBalancer(curve.Morton, 3, 64, transform.KVTransform,
		optimizer.PowerRangeOptimizer, ns)
	if err != nil {
		panic(err)
	}

	if err := bal.Optimize(); err != nil {
		panic(err)
	}

	rates := make(map[string]int)
	for _, key := range keys {
		vals := make([]interface{}, 1)
		vals[0] = key
		di := balancer.NewMockDataItem(key, 1, vals)

		if n, err := bal.LocateData(di); err != nil {
			panic(err)
		} else {
			rates[n.ID()]++
		}
		if err := bal.Optimize(); err != nil {
			panic(err)
		}
	}
	return rates
}

func space() {
	rates := make(map[string]int)
	ratesDi := make(map[[3]float64]int)

	node0 := balancer.NewMockNode("node-0", 1, 20)
	node1 := balancer.NewMockNode("node-1", 1, 20)
	node2 := balancer.NewMockNode("node-2", 1, 20)
	node3 := balancer.NewMockNode("node-2", 1, 20)

	nodes := []balancer.Node{node0, node1, node2, node3}
	bal, err := balancer.NewBalancer(curve.Morton, 3, 16, transform.SpaceTransform,
		optimizer.RangeOptimizer, nodes)
	if err != nil {
		panic(err)
	}
	if err := bal.AddNode(node2, true); err != nil {
		panic(err)
	}

	if err := bal.Optimize(); err != nil {
		panic(err)
	}

	for iter := uint64(0); iter < 10; iter++ {
		vals := make([]interface{}, 3)
		vals[0] = float64(rand.Intn(90))
		vals[1] = float64(rand.Intn(180))
		vals[2] = int64(iter + 1609459200)
		di := balancer.NewMockDataItem(fmt.Sprintf("di-%d", iter), 1024*iter, vals)
		if _, err := bal.LocateData(di); err != nil {
			panic(err)
		}
	}
	if err := bal.Optimize(); err != nil {
		panic(err)
	}
	for iter := uint64(0); iter < 100000; iter++ {
		vals := make([]interface{}, 3)
		lon := -90 + rand.Float64()*180
		lat := -180 + rand.Float64()*360
		ts := int64(iter + 1609459200)
		vals[0] = lon
		vals[1] = lat
		vals[2] = ts
		k := [3]float64{lon, lat, float64(ts)}
		ratesDi[k]++
		di := balancer.NewMockDataItem(fmt.Sprintf("di-%d", iter), 1024*iter, vals)
		if n, err := bal.LocateData(di); err != nil {
			panic(err)
		} else {
			rates[n.ID()]++
		}
	}

	for key := range ratesDi {
		//fmt.Println(key)
		if ratesDi[key] == 1 {
			delete(ratesDi, key)
		}
	}
	fmt.Println(rates, ratesDi)
	//fmt.Println(rates, ratesDi, bal.Space().Rates())
}

func hring(keys []string, nodes map[string]int) map[string]int {
	rand.Seed(42)

	ring := hashring.NewWithWeights(nodes)

	rates := make(map[string]int)
	for _, key := range keys {
		s, _ := ring.GetNode(key)
		rates[s]++
	}
	return rates
}

func consring(keys []string, nodes map[string]int) map[string]int {
	cfg := consistent.Config{
		PartitionCount:    4,
		ReplicationFactor: 1,
		Load:              1.25,
		Hasher:            hasher{},
	}

	var ns []consistent.Member
	for n := range nodes {
		ns = append(ns, myMember(n))
	}

	c := consistent.New(ns, cfg)
	// Add some members to the consistent hash table.
	// Add function calculates average load and distributes partitions over members

	rates := make(map[string]int)
	for _, key := range keys {
		owner := c.LocateKey([]byte(key))
		rates[owner.String()]++
	}

	return rates
}

type myMember string

func (m myMember) String() string {
	return string(m)
}

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}
