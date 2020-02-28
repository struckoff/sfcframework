package balancer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/struckoff/SFCFramework/curve"
)

type space struct {
	mu  sync.Mutex
	cg  []cellGroup
	sfc curve.Curve
	tf  TransformFunc
}

func (s *space) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cg)
}

func (s *space) addNode(n Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.cg) == 0 {
		s.cg = []cellGroup{newCellGroup(n)}
		return
	}
	s.cg = append(s.cg, newCellGroup(n))
}

func (s *space) addData(d DataItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	size := s.sfc.Size()
	if s.tf == nil {
		return errors.New("transform function is not set")
	}
	coords, err := s.tf(d.Values(), size)
	cID, err := s.sfc.Encode(coords)
	if err != nil {
		return fmt.Errorf("item encoding error: %w", err)
	}
	c := 0
	for _, g := range s.cg {
		for i := range g.cells {
			if uint64(c) == cID {
				err = g.cells[i].add(d)
				if err != nil {
					return fmt.Errorf("error on adding data item to cell: %w", err)
				}
				return nil
			}
		}
	}
	return errors.New("unable to find cell to add data item")
}

func (s *space) distribution() DataDistribution {
	res := make(DataDistribution, len(s.cg))
	for i, g := range s.cg {
		nd := NodeData{
			ID:    s.cg[i].node.ID(),
			Items: []string{},
		}
		for _, c := range g.cells {
			ids := c.itemIDs()
			nd.Items = append(nd.Items, ids...)
		}
		res[i] = nd
	}
	return res
}
