package metadata

import (
	"math"

	"github.com/apache/yunikorn-core/pkg/scheduler/objects"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/GOA/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/GOA/metadata/user"

	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
)

var (
	ResourceTypes = []string{sicommon.CPU, sicommon.Memory}
)

type MetaData struct {
	UserData           *UserData.UserData
	NodeData           *NodeData.NodeData
	DRs                [][]float64
	DRRatioReciprocals [][]float64
}

func NewMetaData() *MetaData {
	return &MetaData{
		UserData: UserData.NewUserData(ResourceTypes),
		NodeData: NodeData.NewNodeData(ResourceTypes),
		DRs: make([][]float64, 0),
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
func (metaData *MetaData) AddUser(app *objects.Application) {
	metaData.UserData.AddUser(app)
}

func (metaData *MetaData) AddNode(node *objects.Node) {
	metaData.NodeData.AddNode(node)
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
