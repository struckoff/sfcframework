package powerOptimizer

import (
	"github.com/struckoff/SFCFramework"
)

func PowerOptimizer(s *SFCFramework.Space) (res []SFCFramework.CellGroup, err error) {
	var node SFCFramework.Node

	totalLoad := s.TotalLoad()
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	cells := s.Cells()

	for _, cg := range cgs {
		node = cg.Node()
		totalPower += node.Power().Get()
	}

	i := 0
	node = cgs[0].Node()
	cg := SFCFramework.NewCellGroup(node)
	p := node.Power().Get() / totalPower
	l := uint64(float64(totalLoad) * p)
	for j := range cells {
		cg.AddCell(&cells[j])
		if cg.TotalLoad() >= l {
			if i == (len(cgs) - 1) {
				continue
			}
			res = append(res, cg)
			i++
			cg = SFCFramework.NewCellGroup(cgs[i].Node())
			p = cg.Node().Power().Get() / totalPower
			l = uint64(float64(totalLoad) * p)
		}
	}

	res = append(res, cg)
	return res, nil
}
