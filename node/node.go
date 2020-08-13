package node

import (
	"github.com/struckoff/SFCFramework/capacity"
	"github.com/struckoff/SFCFramework/power"
)

// Node is an interface describing storage/processing node in the cluster.
type Node interface {
	ID() string
	Power() power.Power
	Capacity() capacity.Capacity
	Hash() uint64 //unique node hash
}
