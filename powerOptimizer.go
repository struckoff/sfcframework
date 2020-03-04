package balancer

func PowerOptimizer(cgs []cellGroup) ([]cellGroup, error) {
	totalLoad := uint64(0)
	cells := []*cell{}
	totalPower := 0.0
	for _, cg := range cgs {
		totalPower += cg.node.Power().Get()
		for _, c := range cg.cells {
			totalLoad += c.load
			cells = append(cells, c)
		}
	}
	res := []cellGroup{}
	i := 0
	cg := newCellGroup(cgs[0].node)
	p := cg.node.Power().Get() / totalPower
	l := uint64(float64(totalLoad) * p)
	for j := range cells {
		cg.addLoad(cells[j].load)
		cg.cells = append(cg.cells, cells[j])
		cells[j].cg = &cg
		if cg.load >= l {
			if i == (len(cgs) - 1) {
				continue
			}
			res = append(res, cg)
			i++
			cg = newCellGroup(cgs[i].node)
			p = cg.node.Power().Get() / totalPower
			l = uint64(float64(totalLoad) * p)
		}
	}
	res = append(res, cg)
	return res, nil
}
