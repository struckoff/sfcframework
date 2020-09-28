/*
	Default optimizers provided by the library.
	It covers only the most generic use cases.
	For more appropriate use optimizer should be provided by service from the outside.
*/
package optimizer

import (
	"math"
	"sort"

	"github.com/pkg/errors"
	balancer "github.com/struckoff/sfcframework"
)

//RangeOptimizer - divide curve into segments.
//Length of each segment depends on nodes power.
func RangeOptimizer(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	if len(cgs) == 0 {
		return res, nil
	}
	var max, min uint64

	sort.Slice(cgs, func(i, j int) bool { return cgs[i].Node().Hash() < cgs[j].Node().Hash() })

	for i := 0; i < len(cgs); i++ {
		min = max
		p := cgs[i].Node().Power().Get() / totalPower
		max = min + uint64(math.Round(float64(s.Capacity())*p))
		if err := cgs[i].SetRange(min, max); err != nil {
			return nil, errors.Wrap(err, "range optimizer error")
		}
	}

	if max < s.Capacity() {
		if err := cgs[len(cgs)-1].SetRange(min, s.Capacity()+1); err != nil {
			return nil, errors.Wrap(err, "range optimizer error")
		}
	}
	cells := s.Cells()
	for i := range cells {
		for cgi := range cgs {
			if cgs[cgi].FitsRange(cells[i].ID()) {
				cells[i].Group().RemoveCell(cells[i].ID())
				cgs[cgi].AddCell(cells[i])
				break
			}
		}
	}
	return cgs, nil
}

//PowerRangeOptimizer - divide curve into segments.
//Length of each segment depends on nodes power and capacity.
func PowerRangeOptimizer(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	//TODO: reduce Capacity calls

	cells := s.Cells()
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	if len(cgs) == 0 {
		return res, nil
	}
	var max, min uint64

	caps := make([]float64, len(cgs))
	for i := range cgs {
		caps[i], err = cgs[i].Node().Capacity().Get()
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(cgs, func(i, j int) bool {
		capI, _ := cgs[i].Node().Capacity().Get()
		capJ, _ := cgs[j].Node().Capacity().Get()
		return (capI - float64(cgs[i].TotalLoad())) < (capJ - float64(cgs[j].TotalLoad()))
	})

	for i := 0; i < len(cgs); i++ {
		min = max
		p := cgs[i].Node().Power().Get() / totalPower
		f, err := cgs[i].Node().Capacity().Get()
		if err != nil {
			return nil, err
		}
		max = min + uint64(math.Round(float64(s.Capacity())*p))

		for ci := 0; ci < len(cells); ci++ {
			if cells[ci].ID() > max {
				break
			}
			if cells[ci].ID() >= min {
				f -= float64(cells[ci].Load())
				if f <= 0 {
					c := ci - 1
					if c < 0 {
						c = 0
					}
					max = cells[ci].ID()
					break
				}
				cells[ci].Group().RemoveCell(cells[ci].ID())
				cgs[i].AddCell(cells[ci])
			}
		}
		if err := cgs[i].SetRange(min, max); err != nil {
			return nil, errors.Wrap(err, "power range optimizer error")
		}
	}

	if max <= s.Capacity() {
		if err := cgs[len(cgs)-1].SetRange(min, s.Capacity()+1); err != nil {
			return nil, errors.Wrap(err, "range optimizer error")
		}
		for ci := range cells {
			if cells[ci].ID() >= max {
				cells[ci].Group().RemoveCell(cells[ci].ID())
				cgs[len(cgs)-1].AddCell(cells[ci])
			}
		}
	}

	return cgs, nil
}
