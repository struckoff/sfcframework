package optimizer

import (
	"fmt"
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"math"
)

func RangeOptimizer(s *balancer.Space) (res []balancer.CellGroup, err error) {
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	var check float64
	var max, min uint64

	for iter := 0; iter < len(cgs); iter++ {
		min = max
		p := cgs[iter].Node().Power().Get() / totalPower
		max = min + uint64(math.Round(float64(s.Capacity())*p))
		if err := cgs[iter].SetRange(min, max); err != nil {
			return nil, errors.Wrap(err, "count optimizer error")
		}
		check += p
	}

	if check < 1 {
		if err := cgs[len(cgs)-1].SetRange(min, s.Capacity()); err != nil {
			return nil, errors.Wrap(err, "count optimizer error")
		}
	}
	cells := s.Cells()
	for iter := range cells {
		for cgiter := range cgs {
			if cgs[cgiter].FitsRange(cells[iter].ID()) {
				cgs[cgiter].AddCell(&cells[iter], true)
				break
			}
		}
	}
	fmt.Println(check, len(s.Cells()), len(cells), s.Capacity()+1)
	return cgs, nil
}
