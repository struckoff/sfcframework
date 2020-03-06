package powerOptimizer

import (
	"github.com/struckoff/SFCFramework"
)

func PowerOptimizer(s *balancer.Space) (res []balancer.CellGroup, err error) {
	var node balancer.Node

	totalLoad := s.TotalLoad()
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	cells := s.Cells()

	//res := make([]balancer.CellGroup, len(cgs))
	//weights := make([]float64, len(cgs))
	//for iter := range cgs {
	//	weights[iter] = cgs[iter].Node().Power().Get() / totalPower
	//}

	i := 0
	node = cgs[0].Node()
	cg := balancer.NewCellGroup(node)
	p := node.Power().Get() / totalPower
	l := float64(totalLoad) * p
	for j := range cells {
		cg.AddCell(&cells[j])
		if float64(cg.TotalLoad()) >= l {
			if i == (len(cgs) - 1) {
				continue
			}
			res = append(res, cg)
			i++
			cg = balancer.NewCellGroup(cgs[i].Node())
			p = cg.Node().Power().Get() / totalPower
			l = float64(totalLoad) * p
		}
	}

	res = append(res, cg)
	return res, nil
}
