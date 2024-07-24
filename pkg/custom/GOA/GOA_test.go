package GOA

import (
	"fmt"
	"testing"
	"reflect"

	"github.com/apache/yunikorn-core/pkg/custom/GOA/math/vector"
	"github.com/apache/yunikorn-core/pkg/custom/GOA/metadata"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/GOA/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/GOA/metadata/user"

	"github.com/apache/yunikorn-core/pkg/log"
	sicommon "github.com/apache/yunikorn-scheduler-interface/lib/go/common"
)

var (
	ResourceTypes = []string{sicommon.CPU, sicommon.Memory}
)

func initBasicParameter() (metaData *metadata.MetaData) {
	userData := UserData.NewUserData(ResourceTypes)
	userData.AddUserDirectly("userA", []float64{1, 4})
	userData.AddUserDirectly("userB", []float64{3, 1})
	nodeData := NodeData.NewNodeData(ResourceTypes)
	nodeData.AddNodeDirectly("nodeA", []float64{18, 28})

	metaData = &metadata.MetaData{
		UserData: userData,
		NodeData: nodeData,
	}

	log.Log(log.Custom).Info(fmt.Sprintf("userData be like: %v", userData))
	log.Log(log.Custom).Info(fmt.Sprintf("nodeData be like: %v", nodeData))

	return 
}

func TestStartScheduler(t *testing.T) {
	metaData := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metaData))
	goa = &GOA{
		metaData: metaData,
		nodeCount: metaData.NodeData.NodeCount,
		resourceCount: metaData.NodeData.ResourceCount,
		userCount: metaData.UserData.UserCount, 
		grasshoppers: make([]*vector.Vector, grasshopperAmount),
	}

	log.Log(log.Custom).Info(fmt.Sprintf("basis goa: %v", goa))
	decision := goa.StartScheduler()

    expected := []int{6, 4}

    if !reflect.DeepEqual(decision, expected) {
		t.Errorf("Expected answer wiil be [6, 4], but got %v", decision)
	}
}

func Test_getGravityUnitVector(t *testing.T) {
	metaData := initBasicParameter()
	goa = &GOA{
		metaData: metaData,
		nodeCount: metaData.NodeData.NodeCount,
		resourceCount: metaData.NodeData.ResourceCount,
		userCount: metaData.UserData.UserCount, 
		grasshoppers: make([]*vector.Vector, grasshopperAmount),
	}
	
	goa.calculateDomainResources()

	except := vector.NewVectorByInt([]int{3, 2})
	other := vector.NewVectorByInt([]int{4, 1})
	v1 := goa.getGravityUnitVector(except)
	v2 := goa.getGravityUnitVector(other)
	log.Log(log.Custom).Info(fmt.Sprintf("there is : %v and %v", v1.Norm(), v2.Norm()))
} 

func Test_getEffectScore(t *testing.T) {
	metaData := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metaData))
	goa = &GOA{
		metaData: metaData,
		nodeCount: metaData.NodeData.NodeCount,
		resourceCount: metaData.NodeData.ResourceCount,
		userCount: metaData.UserData.UserCount, 
		grasshoppers: make([]*vector.Vector, grasshopperAmount),
	}

	except := vector.NewVector([]float64{4.852167171869434, 0.7797597317740037})
	fmt.Println("except: ", goa.getEffectScore(except))
	
}
func Test_getScore(t *testing.T) {
	metaData := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metaData))
	goa = &GOA{
		metaData: metaData,
		nodeCount: metaData.NodeData.NodeCount,
		resourceCount: metaData.NodeData.ResourceCount,
		userCount: metaData.UserData.UserCount, 
		grasshoppers: make([]*vector.Vector, grasshopperAmount),
	}

	goa.calculateDomainResources()

	except := vector.NewVector([]float64{4.852167171869434, 0.7797597317740037})
	fmt.Println("except: ", goa.getScore(except))
	fmt.Println("except: ", goa.getEffectScore(except))
	fmt.Println("except: ", goa.getFairnessScore(except))
}