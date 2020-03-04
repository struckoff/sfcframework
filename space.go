package balancer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/struckoff/SFCFramework/curve"
)

type space struct {
	mu    sync.Mutex
	cells []cell
	cg    []cellGroup
	sfc   curve.Curve
	tf    TransformFunc
	of    OptimizerFunc
}

func newSpace(sfc curve.Curve, tf TransformFunc, of OptimizerFunc) *space {
	l := sfc.Length()
	cells := make([]cell, l)
	for i := range cells {
		cells[i] = newCell()
	}
	return &space{
		mu:    sync.Mutex{},
		cells: cells,
		cg:    []cellGroup{},
		sfc:   sfc,
		tf:    tf,
		of:    of,
	}
}

func (s *space) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.cg)
}

func (s *space) addNode(n Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.cg) == 0 {
		s.cg = []cellGroup{newCellGroup(n)}
		return nil
	}
	s.cg = append(s.cg, newCellGroup(n))
	return s.optimize()
}

func (s *space) addData(d DataItem) error {
	if len(s.cg) == 0 {
		return errors.New("no nodes in the cluster")
	}
	size := s.sfc.DimSize()
	if s.tf == nil {
		return errors.New("transform function is not set")
	}
	coords, err := s.tf(d.Values(), size)
	cID, err := s.sfc.Encode(coords)
	if err != nil {
		return fmt.Errorf("item encoding error: %w", err)
	}
	if cID > uint64(len(s.cells)-1) {
		return errors.New("cell ID is larger that number of cells in the space")
	}
	return s.cells[cID].add(d)
}

func (s *space) distribution() DataDistribution {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make(DataDistribution, len(s.cg))
	for i := range s.cg {
		nd := NodeData{
			ID:    s.cg[i].node.ID(),
			Items: []string{},
		}
		for _, c := range s.cg[i].cells {
			ids := c.itemIDs()
			nd.Items = append(nd.Items, ids...)
		}
		res[i] = nd
	}
	return res
}

func (s *space) optimize() error {
	cgs, err := s.of(s.cg)
	if err != nil {
		return err
	}
	s.cg = cgs
	return nil
}
