package AGA

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"

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
		for j := 0; j < len(candidateArray); j++ {
			candidateArray[j] = int(r.Int63n(6))
		}
		candidate := vector.NewVectorByInt(candidateArray)
		candidates = append(candidates, candidate)
	} 
	return candidates
}

func (aga *AGA) Start() []int{
	ACOParameter := ACO.NewACOHyperParameter(ACO_NumAnt, ACO_Epochs)
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


	removeIndexs := make([]int, 0)
	visited := make([]int, users)
	for i := 0; i < users; i++ {
		visited[i] = 0
	}

	for nodeIndex, _ := range aga.metadata.Nodes {
		for userIndex, _ := range aga.metadata.Requests {
			if distributeAmount := decision[nodeIndex*users+userIndex]; distributeAmount != 0 {
				nodeId := aga.metadata.Nodes[nodeIndex]
				ask := aga.metadata.Requests[userIndex]
				alloc := objects.NewAllocation(nodeId, ask)
				allocs = append(allocs, alloc)
				if visited[userIndex] == 0 {
					removeIndexs = append([]int{userIndex}, removeIndexs...)
					visited[userIndex] = 1
				}

			}

		}
	}

	for _, index := range removeIndexs {
		aga.RemoveUser(index)
	}
	return allocs
}