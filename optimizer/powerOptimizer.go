package optimizer

import (
	"github.com/struckoff/SFCFramework"
	node2 "github.com/struckoff/SFCFramework/node"
)

func PowerOptimizer(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	var node node2.Node

	totalLoad := s.TotalLoad()
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	cells := s.Cells()

	i := 0
	node = cgs[i].Node()
	cg := balancer.NewCellGroup(node)
	p := node.Power().Get() / totalPower
	l := float64(totalLoad) * p
	var min, max uint64
	for iter := range cells {
		cg.AddCell(cells[iter], false)
		max = cells[iter].ID()
		if float64(cg.TotalLoad()) >= l {
			if i == (len(cgs) - 1) {
				continue
			}
			if err := cg.SetRange(min, max, s); err != nil {
				return nil, err
			}
			res = append(res, cg)
			min = max
			i++
			cg = balancer.NewCellGroup(cgs[i].Node())
			p = cg.Node().Power().Get() / totalPower
			l = float64(totalLoad) * p
		}
	}
	res = append(res, cg)
	return res, nil
}

// PowerOptimizerGreedy use prefilling of results slise.
func PowerOptimizerGreedy(s *balancer.Space) (res []balancer.CellGroup, err error) {
	totalLoad := float64(s.TotalLoad())
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	cells := s.Cells()

	lastCgIndex := len(cgs) - 1
	res = make([]balancer.CellGroup, len(cgs))
	ws := make([]float64, len(cgs))
	i := 0

	for iter := range res {
		node := cgs[iter].Node()
		res[iter].SetNode(node)
		ws[iter] = totalLoad * (node.Power().Get() / totalPower)
	}

	for iter := range cells {
		res[i].AddCell(cells[iter], false)
		ws[i] -= float64(cells[iter].Load())
		if ws[i] <= 0 && i < lastCgIndex {
			i++
		}
	}

	return res, nil
}

// PowerOptimizerPerms fills last CellGroup will cells.
// Function mutates cellGroups in space.
func PowerOptimizerPerms(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	totalLoad := float64(s.TotalLoad())
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	cells := s.Cells()

	lastCgIndex := len(cgs) - 1

	l := totalLoad * (cgs[lastCgIndex].Node().Power().Get() / totalPower)
	l -= float64(cgs[lastCgIndex].TotalLoad())

	for iter := range cells {
		if l <= 0 {
			break
		}
		cl := float64(cells[iter].Load())
		if cl > l {
			continue
		}
		l -= cl
		cgs[lastCgIndex].AddCell(cells[iter], true)
	}

	return cgs, nil
}
