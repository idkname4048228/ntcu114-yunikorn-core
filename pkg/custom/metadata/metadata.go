package metadata

import (
	"fmt"
	"math"
	"sync"

	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/metadata/user"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	"github.com/apache/yunikorn-core/pkg/log"
)

var (
	ResourceTypes = []string{sicommon.CPU, sicommon.Memory}
)

type Metadata struct {
	UserData           	*UserData.UserData
	NodeData           	*NodeData.NodeData
	DRs                	[][]float64
	DRRatioReciprocals 	[][]float64
	GlobalDRs                	[]float64
	GlobalDRRatioReciprocals 	[]float64

	Nodes 				[]string
	Requests 			[]*objects.AllocationAsk

	*sync.RWMutex
}

func NewMetadata() *Metadata {
	return &Metadata{
		UserData:           UserData.NewUserData(ResourceTypes),
		NodeData:           NodeData.NewNodeData(ResourceTypes),
		DRs:                make([][]float64, 0),
		DRRatioReciprocals: make([][]float64, 0),
	}
}

func (metadata *Metadata) GetUserAsks() [][]float64 {
	return metadata.UserData.GetUserAsks()
}

func (metadata *Metadata) GetNodeLimits() [][]float64 {
	return metadata.NodeData.ResourceLimits
}

func (metadata *Metadata) GetTotalLimits() []float64 {
	return metadata.NodeData.TotalLimits
}

func (metadata *Metadata) AddUser(ask *objects.AllocationAsk){ 
	log.Log(log.Custom).Info("metadata add user")
	metadata.Requests = append(metadata.Requests, ask)
	metadata.UserData.AddUser(ask)
	log.Log(log.Custom).Info(fmt.Sprintf("GOA add ask: %v", ask.GetAllocationKey()))
	log.Log(log.Custom).Info(fmt.Sprintf("GOA add ask: %v", ask.GetAllocatedResource()))
}

func (metadata *Metadata) RemoveUser(index int) {
	log.Log(log.Custom).Info(fmt.Sprintf("removing index %v, and length is %v", index, len(metadata.Requests)))
	metadata.Requests = append(metadata.Requests[:index], metadata.Requests[index + 1:]...)
	metadata.UserData.RemoveUser(index)
	log.Log(log.Custom).Info(fmt.Sprintf("now the length is %v", len(metadata.Requests)))
}

func (metadata *Metadata) AddNode(node *objects.Node) {
	metadata.Nodes = append(metadata.Nodes, node.NodeID)
	metadata.NodeData.AddNode(node)
	log.Log(log.Custom).Info(fmt.Sprintf("add node: %v", node.NodeID))
}

func (metadata *Metadata) CalculateDRs() {
	metadata.DRs = make([][]float64, 0)
	metadata.DRRatioReciprocals = make([][]float64, 0)

	for _, limit := range metadata.NodeData.ResourceLimits {
		DR := make([]float64, 0)
		DRRatioReciprocal := make([]float64, 0)

		userAsks := metadata.UserData.GetUserAsks()
		for _, askResources := range userAsks {
			maxRatio := 0.0

			for i := 0; i < int(metadata.NodeData.ResourceCount); i++ {
				ratio := askResources[i] / limit[i]
				if ratio > maxRatio {
					maxRatio = ratio
				}
			}

			DR = append(DR, maxRatio)
			DRRatioReciprocal = append(DRRatioReciprocal, math.Pow(maxRatio, -1))
		}
		metadata.DRs = append(metadata.DRs, DR)
		metadata.DRRatioReciprocals = append(metadata.DRRatioReciprocals, DRRatioReciprocal)
	}

}

func (metadata *Metadata) CalculateGlobalDRs() {
	metadata.GlobalDRs = make([]float64, metadata.UserData.UserCount)
	metadata.GlobalDRRatioReciprocals = make([]float64, metadata.UserData.UserCount)

	totalLimits := metadata.GetTotalLimits()
	userAsks := metadata.UserData.GetUserAsks()
	for userIndex, askResources := range userAsks {
		maxRatio := 0.0
		for resourceIndex, limit := range totalLimits {
			ratio := askResources[resourceIndex] / limit
			if ratio > maxRatio {
				maxRatio = ratio
			}
		}
		metadata.GlobalDRs[userIndex] = maxRatio
		metadata.GlobalDRRatioReciprocals[userIndex] = math.Pow(maxRatio, -1)
	}
}