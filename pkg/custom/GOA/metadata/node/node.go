package NodeData

import (
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)


type NodeData struct {
	NodeCount      int
	NodeIDs        []string
	ResourceCount  int
	ResourceTypes  []string
	ResourceLimits [][]float64
	TotalLimits    []float64
}

func NewNodeData(ResourceTypes []string) *NodeData {
	ResourceCount := len(ResourceTypes)
	return &NodeData{
		NodeCount: 0,
		NodeIDs: make([]string, 0),
		ResourceCount: len(ResourceTypes),
		ResourceTypes:   ResourceTypes,
		ResourceLimits: make([][]float64, 0),
		TotalLimits:    make([]float64, ResourceCount),
	}
}

// Parse the vcore and memory in node
func (nodeData *NodeData) AddNode(n *objects.Node) {
	nodeData.NodeIDs = append(nodeData.NodeIDs, n.NodeID)
	nodeData.NodeCount += 1;
	
	availableLimit := make([]float64, nodeData.ResourceCount)	
	resources := n.GetAvailableResource().Resources
	for index, targetType := range nodeData.ResourceTypes {
		availableLimit[index] += float64(resources[targetType])
		nodeData.TotalLimits[index] += availableLimit[index]
	}

	nodeData.ResourceLimits = append(nodeData.ResourceLimits, availableLimit)
}

// make test easy 
func (nodeData *NodeData) AddNodeDirectly(nodeName string, resource []float64) {
	nodeData.NodeCount += 1
	
	availableLimit := make([]float64, 2)
	for index := 0; index < len(resource); index++ {
		availableLimit[index] += resource[index]
		nodeData.TotalLimits[index] += availableLimit[index]
	}

	nodeData.ResourceLimits = append(nodeData.ResourceLimits, availableLimit)
}