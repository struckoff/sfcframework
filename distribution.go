package SFCFramework

// DataDistribution represents Distribution of data items between nodes.
type DataDistribution []NodeData

// NodeData contains information about what data items should be placed onto node.
//
// ID - identifier of the node.
//
// Items - slice of data items identifiers.
type NodeData struct {
	ID    string
	Items []string
}
