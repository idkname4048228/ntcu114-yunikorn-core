package GOA

import (
	"fmt"
	"math/rand"
	"time"
	"testing"
	// "reflect"

	"github.com/apache/yunikorn-core/pkg/custom/math/vector"
	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"

	NodeData "github.com/apache/yunikorn-core/pkg/custom/metadata/node"
	UserData "github.com/apache/yunikorn-core/pkg/custom/metadata/user"
	agamath "github.com/apache/yunikorn-core/pkg/custom/math"

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

func randanInitValue(metadata *Metadata.Metadata) []*vector.Vector{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	users := metadata.UserData.UserCount
	nodes := metadata.NodeData.NodeCount

	candidates := make([]*vector.Vector, 0)

	for i := 0; i < 20; i++ {
		candidateArray := make([]int, users*nodes)
		for j := 0; j < len(candidateArray); j++ {
			candidateArray[j] = int(r.Int63n(6))
		}
		candidate := vector.NewVectorByInt(candidateArray)
		candidates = append(candidates, candidate)
	} 
	return candidates
}

func TestStart(t *testing.T) {
	metadata := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metadata))
	goa := &GOA{
		metadata: metadata,
	}
	// GOA hyperParameter
	const (
		GOA_iterations = 10 
		GOA_cMax = 1.0
		GOA_cMin = 0.00001
		GOA_grasshopperAmount = 10
		GOA_GForce = 0.2
		GOA_WindForce = 0.2
	)
	GOAParameter := NewGOAHyperParameter(
		GOA_iterations, 
		GOA_cMax, 
		GOA_cMin, 
		GOA_grasshopperAmount,	
		GOA_GForce, 
		GOA_WindForce, 
	)
	goa.SetHyperParameter(GOAParameter)	
	candidates := randanInitValue(metadata)

	log.Log(log.Custom).Info(fmt.Sprintf("basis goa: %v", goa))
	decision := goa.Start(candidates)

	check := vector.NewVectorByInt(decision)
	fmt.Println("check score: ", agamath.GetScore(metadata, check))
	fmt.Println("check effect: ", agamath.GetEffectScore(metadata, check))
	fmt.Println("check fair: ", agamath.GetFairnessScore(metadata, check))
	
	metadata.CalculateGlobalDRs()
	fmt.Println("check globalDRs: ", metadata.GlobalDRRatioReciprocals)
}

func Test_getGravityUnitVector(t *testing.T) {
	metadata := initBasicParameter()
	goa := &GOA{
		metadata: metadata,
	}
	
	goa.calculateDomainResources()

	except := vector.NewVectorByInt([]int{2, 2})
	other := vector.NewVectorByInt([]int{4, 1})
	v1 := goa.getGravityUnitVector(except)
	v2 := goa.getGravityUnitVector(other)
	log.Log(log.Custom).Info(fmt.Sprintf("there is : %v and %v", v1.Norm(), v2.Norm()))
} 

func Test_getEffectScore(t *testing.T) {
	metadata := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metadata))

	metadata.CalculateDRs()

	a := vector.NewVector([]float64{6, 4})
	fmt.Println("except: ", agamath.GetEffectScore(metadata, a))
	b := vector.NewVector([]float64{5, 5})
	fmt.Println("except: ", agamath.GetEffectScore(metadata, b))	
}

func Test_getFairnessScore(t *testing.T) {
	metadata := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metadata))

	metadata.CalculateDRs()

	a := vector.NewVector([]float64{6, 4})
	fmt.Println("{6, 4} fair: ", agamath.GetFairnessScore(metadata, a))
	
	b := vector.NewVector([]float64{5, 5})
	fmt.Println("{5, 5} fair: ", agamath.GetFairnessScore(metadata, b))
}

func Test_getScore(t *testing.T) {
	metadata := initBasicParameter()
	log.Log(log.Custom).Info(fmt.Sprintf("metadata be like: %v", metadata))

	metadata.CalculateDRs()

	a := vector.NewVector([]float64{6, 4})
	fmt.Println("{6, 4} score: ", agamath.GetScore(metadata, a))
	fmt.Println("{6, 4} effect: ", agamath.GetEffectScore(metadata, a))
	fmt.Println("{6, 4} fair: ", agamath.GetFairnessScore(metadata, a))
	
	b := vector.NewVector([]float64{5, 5})
	fmt.Println("{5, 5} score: ", agamath.GetScore(metadata, b))
	fmt.Println("{5, 5} effect: ", agamath.GetEffectScore(metadata, b))
	fmt.Println("{5, 5} fair: ", agamath.GetFairnessScore(metadata, b))
}