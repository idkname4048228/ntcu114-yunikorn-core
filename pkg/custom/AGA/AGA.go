package AGA

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"
	agamath "github.com/apache/yunikorn-core/pkg/custom/math"

	"github.com/apache/yunikorn-core/pkg/custom/ACO"
	"github.com/apache/yunikorn-core/pkg/custom/GOA"

	"github.com/apache/yunikorn-core/pkg/custom/math/vector"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)

var aga *AGA

type AGA struct {
	metadata 	*Metadata.Metadata
	goa			*GOA.GOA
	aco 		*ACO.ACO

	sync.RWMutex
}

func Init() {
	aga = NewAGA()
}

func GetAGA() *AGA {
	return aga
}

func NewAGA() *AGA{
	return &AGA{
		metadata: Metadata.NewMetadata(),
		goa: GOA.NewGOA(),
		aco: ACO.NewACO(),
	}
}

func (aga *AGA) randanInitValue() []*vector.Vector{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	users := aga.metadata.UserData.UserCount
	nodes := aga.metadata.NodeData.NodeCount

	candidates := make([]*vector.Vector, 0)

	for i := 0; i < aga.aco.HyperParameter.AntNum; i++ {
		candidateArray := make([]int, users*nodes)
		distributeAmount := make([]int, users)

		for j := 0; j < users; j++ {
			userName := aga.metadata.UserData.GetName(j)
			count := aga.metadata.UserData.GetUserAskCount(userName)
			distributeAmount[j] = count;
		}

		for j := 0; j < len(candidateArray); j++ {
			userIndex := j % users;
			remains := distributeAmount[userIndex]
			tmp := int64(float64(remains) * float64(i) / float64(aga.aco.HyperParameter.AntNum))
			// log.Log(log.Custom).Info(fmt.Sprintf("%v * %v = %v", remains, float64(i) / float64(aga.aco.HyperParameter.AntNum), tmp))
			if tmp <= 0 {
				tmp = 1
			}
			// log.Log(log.Custom).Info(fmt.Sprintf("%v", tmp))
			amount := r.Int63n(tmp)

			candidateArray[j] = int(amount)
			distributeAmount[userIndex] -= int(amount)
		}

		candidate := vector.NewVectorByInt(candidateArray)
		candidates = append(candidates, candidate)

		// log.Log(log.Custom).Info(fmt.Sprintf("candidate score is %v", agamath.GetScore(aga.metadata, candidate)))

	} 
	return candidates
}

func (aga *AGA) Start() []int{

	users := aga.metadata.UserData.UserCount;
	nodes := aga.metadata.NodeData.NodeCount;

	ACOParameter := ACO.NewACOHyperParameter(ACO_NumAnt, ACO_Epochs, users*nodes * ACO_Steps)
	GOAParameter := GOA.NewGOAHyperParameter(
			GOA_iterations, 
			GOA_cMax, 
			GOA_cMin, 
			GOA_grasshopperAmount,	
			GOA_GForce, 
			GOA_WindForce, 
	)
	aga.aco.SetHyperParameter(ACOParameter)
	aga.aco.SetMetadata(aga.metadata)

	aga.goa.SetHyperParameter(GOAParameter)
	aga.goa.SetMetaData(aga.metadata)

	aga.metadata.CalculateDRs()

	aga.aco.Start(aga.randanInitValue())

	ACOBest := aga.aco.GetCandidate(1)
	
	log.Log(log.Custom).Info(fmt.Sprintf("aco best is %v, and score is %v", ACOBest, agamath.GetScore(aga.metadata, ACOBest[0])))

	desision := aga.goa.Start(aga.aco.GetCandidate(GOA_grasshopperAmount))

	log.Log(log.Custom).Info(fmt.Sprintf("desision is %v", desision))

	return desision
}

func (aga *AGA) GetAllocations() (allocs []*objects.Allocation) {
	
	aga.Lock()
	defer aga.Unlock()
	allocs = make([]*objects.Allocation, 0)

	users := aga.metadata.UserData.UserCount
	nodes := aga.metadata.NodeData.NodeCount

	if users * nodes == 0 {
		return 
	}

	log.Log(log.Custom).Info("AGA start")

	decision := aga.Start()

	result := make([]int, users)

	for nodeIndex, nodeId := range aga.metadata.Nodes {
		for userIndex := 0; userIndex < users; userIndex++ {
			if distributeAmount := decision[nodeIndex*users+userIndex]; distributeAmount != 0 {
				
				name := aga.metadata.UserData.GetName(userIndex)
				asks := aga.metadata.UserData.PopAsks(name, distributeAmount)
				for _, ask := range asks {
					alloc := objects.NewAllocation(nodeId, ask)
					allocs = append(allocs, alloc)
				}
				result[userIndex] += distributeAmount
			}
		}
	}
	
	s := ""
	for i, num := range result {
		if i == 0 {
			s += fmt.Sprintf("user result is: %v", num) 
		} else {
			s += fmt.Sprintf(", %v", num) 
		}
	}

	log.Log(log.Custom).Info(s)

	aga.metadata.UserData.Update()
	return allocs
}