package metadata

import (
	"fmt"
	"math"
	"sync"

	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/GOA/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/GOA/metadata/user"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
	"github.com/apache/yunikorn-core/pkg/log"
)

var (
	ResourceTypes = []string{sicommon.CPU, sicommon.Memory}
)

type MetaData struct {
	UserData           	*UserData.UserData
	NodeData           	*NodeData.NodeData
	DRs                	[][]float64
	DRRatioReciprocals 	[][]float64

	Nodes 				[]string
	Requests 			[]*objects.AllocationAsk

	sync.RWMutex
}

func NewMetaData() *MetaData {
	return &MetaData{
		UserData:           UserData.NewUserData(ResourceTypes),
		NodeData:           NodeData.NewNodeData(ResourceTypes),
		DRs:                make([][]float64, 0),
		DRRatioReciprocals: make([][]float64, 0),
	}
}

func (metaData *MetaData) GetUserAsks() [][]float64 {
	return metaData.UserData.UserAsks
}

func (metaData *MetaData) GetNodeLimits() [][]float64 {
	return metaData.NodeData.ResourceLimits
}

func (metaData *MetaData) GetTotalLimits() []float64 {
	return metaData.NodeData.TotalLimits
}

func (metaData *MetaData) AddUser(ask *objects.AllocationAsk, app *objects.Application) {
	log.Log(log.Custom).Info("metadata add user")
	metaData.Requests = append(metaData.Requests, ask)
	metaData.UserData.AddUser(ask, app)
	log.Log(log.Custom).Info(fmt.Sprintf("GOA add ask: %v", ask.GetAllocationKey()))
	log.Log(log.Custom).Info(fmt.Sprintf("GOA add ask: %v", ask.GetAllocatedResource()))
}

func (metaData *MetaData) RemoveUser(index int) {
	log.Log(log.Custom).Info(fmt.Sprintf("removing index %v, and length is %v", index, len(metaData.Requests)))
	metaData.Requests = append(metaData.Requests[:index], metaData.Requests[index + 1:]...)
	metaData.UserData.RemoveUser(index)
}

func (metaData *MetaData) AddNode(node *objects.Node) {
	metaData.Nodes = append(metaData.Nodes, node.NodeID)
	metaData.NodeData.AddNode(node)
	log.Log(log.Custom).Info(fmt.Sprintf("add node: %v", node.NodeID))
}

func (metaData *MetaData) CalculateDRs() {
	metaData.DRs = make([][]float64, 0)
	metaData.DRRatioReciprocals = make([][]float64, 0)

	for _, limit := range metaData.NodeData.ResourceLimits {
		DR := make([]float64, 0)
		DRRatioReciprocal := make([]float64, 0)

		for _, askResources := range metaData.UserData.UserAsks {
			maxRatio := 0.0

			for i := 0; i < int(metaData.NodeData.ResourceCount); i++ {
				ratio := askResources[i] / limit[i]
				if ratio > maxRatio {
					maxRatio = ratio
				}
			}

			DR = append(DR, maxRatio)
			DRRatioReciprocal = append(DRRatioReciprocal, math.Pow(maxRatio, -1))
		}
		metaData.DRs = append(metaData.DRs, DR)
		metaData.DRRatioReciprocals = append(metaData.DRRatioReciprocals, DRRatioReciprocal)
	}

}
