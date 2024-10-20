package math

import (
	"fmt"
	"testing"
	// "reflect"

	"github.com/apache/yunikorn-core/pkg/custom/math/vector"
	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/metadata/user"

	"github.com/apache/yunikorn-core/pkg/log"
	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
)

var (
	ResourceTypes = []string{sicommon.CPU, sicommon.Memory}
)

func initBasicParameter() (metaData *Metadata.Metadata) {
	userData := UserData.NewUserData(ResourceTypes)
	userData.AddUserDirectly("userA", []float64{1, 2})
	userData.AddUserDirectly("userB", []float64{1, 1})
	userData.AddUserDirectly("userC", []float64{2, 1})

	nodeData := NodeData.NewNodeData(ResourceTypes)
	nodeData.AddNodeDirectly("nodeA", []float64{100, 250})
	nodeData.AddNodeDirectly("nodeB", []float64{150, 200})

	metaData = &Metadata.Metadata{
		UserData: userData,
		NodeData: nodeData,
	}

	log.Log(log.Custom).Info(fmt.Sprintf("userData be like: %v", userData))
	log.Log(log.Custom).Info(fmt.Sprintf("nodeData be like: %v", nodeData))

	return 
}

func TestGetEffectScore(t *testing.T) {
	metadata := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metadata))

	metadata.CalculateDRs()

	a := vector.NewVector([]float64{6, 4})
	fmt.Println("except: ", GetEffectScore(metadata, a))
	b := vector.NewVector([]float64{5, 5})
	fmt.Println("except: ", GetEffectScore(metadata, b))	
}