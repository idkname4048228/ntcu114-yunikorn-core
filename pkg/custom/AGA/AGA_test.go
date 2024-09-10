package AGA

import (
	"fmt"
	"testing"

	// "reflect"

	agamath "github.com/apache/yunikorn-core/pkg/custom/math"
	"github.com/apache/yunikorn-core/pkg/custom/math/vector"
	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"

	"github.com/apache/yunikorn-core/pkg/custom/ACO"
	"github.com/apache/yunikorn-core/pkg/custom/GOA"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/metadata/user"

	// "github.com/apache/yunikorn-core/pkg/log"
	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
)

var (
	ResourceTypes = []string{sicommon.CPU, sicommon.Memory}
)

func initBasicParameter() (metadata *Metadata.Metadata){
	userData := UserData.NewUserData(ResourceTypes)
	userData.AddUserDirectly("userA", []float64{1, 2})
	userData.AddUserDirectly("userB", []float64{1, 1})
	userData.AddUserDirectly("userC", []float64{2, 1})

	nodeData := NodeData.NewNodeData(ResourceTypes)
	nodeData.AddNodeDirectly("nodeA", []float64{100, 250})
	nodeData.AddNodeDirectly("nodeB", []float64{150, 200})

	metadata = &Metadata.Metadata{
		UserData: userData,
		NodeData: nodeData,
	}

	// log.Log(log.Custom).Info(fmt.Sprintf("userData be like: %v", userData))
	// log.Log(log.Custom).Info(fmt.Sprintf("nodeData be like: %v", nodeData))

	return 
}

func TestStart(t *testing.T) {
	metadata := initBasicParameter()
	// log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metadata))
	aga := &AGA{
		metadata: metadata,
		goa: GOA.NewGOA(),
		aco: ACO.NewACO(),
	}

	// log.Log(log.Custom).Info(fmt.Sprintf("basis aga: %v", aga))
	decision := aga.Start()
	check := vector.NewVectorByInt(decision)

	fmt.Println("check score: ", agamath.GetScore(metadata, check))
	fmt.Println("check effect: ", agamath.GetEffectScore(metadata, check))
	fmt.Println("check fair: ", agamath.GetFairnessScore(metadata, check))
	

	metadata.CalculateGlobalDRs()
	fmt.Println("check globalDRs: ", metadata.GlobalDRRatioReciprocals)
}