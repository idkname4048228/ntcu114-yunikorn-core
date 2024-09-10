package GOA

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	goamath "github.com/apache/yunikorn-core/pkg/custom/GOA/math"
	agamath "github.com/apache/yunikorn-core/pkg/custom/math" 
	"github.com/apache/yunikorn-core/pkg/custom/math/vector"
	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"

	"github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
)


type GOAHyperParameter struct {
	iterations        int
	cMax              float64
	cMin              float64
	grasshopperAmount int
	GForce            float64
	WindForce         float64
}

func NewGOAHyperParameter(iterations int, cMax float64, cMin float64, grasshopperAmount int, GForce float64, WindForce float64) *GOAHyperParameter {
	return &GOAHyperParameter{
		iterations:        iterations,
		cMax:              cMax,
		cMin:              cMin,
		grasshopperAmount: grasshopperAmount,
		GForce:            GForce,
		WindForce:         WindForce,
	}
}

type GOA struct {
	metadata        *Metadata.Metadata
	fairness_vector *vector.Vector

	grasshoppers    []*vector.Vector
	bestGrasshopper *vector.Vector

	hyperParameter *GOAHyperParameter

	sync.RWMutex
}

func (goa *GOA) SetHyperParameter(hyperParameter *GOAHyperParameter) {
	goa.hyperParameter = hyperParameter
	goa.grasshoppers =  make([]*vector.Vector, goa.hyperParameter.grasshopperAmount)
}

func(goa *GOA) SetMetaData(metaData *Metadata.Metadata){
	goa.metadata = metaData
}

func NewGOA() *GOA {
	return &GOA{}
}

// GOA utils

func (goa *GOA) getGravityUnitVector(grasshopper *vector.Vector) *vector.Vector {
	grasshopper_unit := grasshopper.GetUnitVector()

	t := goa.fairness_vector.Dot(grasshopper_unit) / grasshopper_unit.Dot(grasshopper_unit)

	D := vector.Multiple(grasshopper_unit, t)

	gravity_vector := vector.Subtract(goa.fairness_vector, D)

	return gravity_vector.GetUnitVector()
}

func (goa *GOA) getWindUnitVector(grasshopper *vector.Vector, best_grasshopper *vector.Vector) *vector.Vector {
	wind_vector := vector.Subtract(best_grasshopper, grasshopper)
	return wind_vector.GetUnitVector()
}

func (goa *GOA) greedyMove(minValue float64, nextPositions *[]*vector.Vector) {
	// goa.bestGrasshopper = nil

	for i := 0; i < goa.hyperParameter.grasshopperAmount; i++ {
		nextPosition := (*nextPositions)[i]

		oldValue := agamath.GetScore(goa.metadata, goa.grasshoppers[i])
		newValue := agamath.GetScore(goa.metadata, nextPosition)

		if newValue != math.Inf(1) {
			if newValue < oldValue {

				goa.grasshoppers[i] = nextPosition
				oldValue = newValue
			}

			if oldValue < minValue || goa.bestGrasshopper == nil {
				minValue = oldValue

				goa.bestGrasshopper = goa.grasshoppers[i]
			}
		}

		// 加入 G_force * calculate_gravity_unit_vector 的結果到 grasshoppers[i]
		gravityVector := goa.getGravityUnitVector(goa.grasshoppers[i])
		goa.grasshoppers[i].Add(vector.Multiple(gravityVector, goa.hyperParameter.GForce))

	}
}

func (goa *GOA) calculateDomainResources() {
	goa.metadata.CalculateDRs()
	fairness_array := make([]float64, 0)
	for _, value := range goa.metadata.DRRatioReciprocals {
		fairness_array = append(fairness_array, value...)
	}
	goa.fairness_vector = vector.NewVector(fairness_array)
	log.Log(log.Custom).Info(fmt.Sprintf("fairness_vector is %v", goa.fairness_vector))
	
}


func (goa *GOA) Start(candidates []*vector.Vector) (decision []int) {
	users := goa.metadata.UserData.UserCount
	nodes := goa.metadata.NodeData.NodeCount

	// log.Log(log.Custom).Info(fmt.Sprintf("now : %v, %v", users, nodes))

	if users*nodes == 0 {
		return nil
	}

	// log.Log(log.Custom).Info(fmt.Sprintf("goa start : %v, %v", users, nodes))
	// log.Log(log.Custom).Info(fmt.Sprintf("metadata is %v", goa.metaData))
	// log.Log(log.Custom).Info(fmt.Sprintf("nodedata is %v", goa.metaData.NodeData))
	// log.Log(log.Custom).Info(fmt.Sprintf("userdata is %v", goa.metaData.UserData))

	goa.calculateDomainResources()

	minValue := math.Inf(1)
	goa.grasshoppers = candidates
	goa.bestGrasshopper = vector.WithSize(users * nodes)

	for _, candidate := range(candidates) {
		value := agamath.GetScore(goa.metadata, candidate)
		if minValue > value {
			minValue = value
			goa.bestGrasshopper = candidate
		}
	}

	// log.Log(log.Custom).Info(fmt.Sprintf("grasshoppers: %v", goa.grasshoppers))
	// log.Log(log.Custom).Info(fmt.Sprintf("best grasshopper: %v", goa.bestGrasshopper))

	// main program
	for i := 0; i < goa.hyperParameter.iterations; i++ {
		c := goa.hyperParameter.cMax - float64(i)*((goa.hyperParameter.cMax-goa.hyperParameter.cMin)/float64(goa.hyperParameter.iterations))
		nextPositions := make([]*vector.Vector, 0)
		for grasshopper_i := 0; grasshopper_i < goa.hyperParameter.grasshopperAmount; grasshopper_i++ {
			moveVector := vector.WithSize(users * nodes)
			for grasshopper_j := 0; grasshopper_j < goa.hyperParameter.grasshopperAmount; grasshopper_j++ {
				if grasshopper_i == grasshopper_j {
					continue
				}
				distance := vector.Subtract(goa.grasshoppers[grasshopper_i], goa.grasshoppers[grasshopper_j])
				if distance.Norm() == 0 {
					continue
				}
				moveVector.Add(vector.Multiple(distance.GetUnitVector(), goamath.SocialInfluence(distance.Norm())))

			}
			moveVector = vector.Multiple(moveVector, c)

			wind_unit_vector := goa.getWindUnitVector(goa.grasshoppers[grasshopper_i], goa.bestGrasshopper)
			moveVector.Add(vector.Multiple(wind_unit_vector, goa.hyperParameter.WindForce))
			nextPositions = append(nextPositions, vector.Add(goa.grasshoppers[grasshopper_i], moveVector))
		}
		goa.greedyMove(math.Inf(1), &nextPositions)

		// log.Log(log.Custom).Info(fmt.Sprintf("grasshoppers: %v", goa.grasshoppers))
		// log.Log(log.Custom).Info(fmt.Sprintf("best grasshopper: %v", goa.bestGrasshopper))
	}

	log.Log(log.Custom).Info(fmt.Sprintf("Best solution: %v", goa.bestGrasshopper))

	decision = goa.bestGrasshopper.ToIntArray()
	log.Log(log.Custom).Info(fmt.Sprintf("decision: %v", decision))
	return
}

func (goa *GOA) GetAllocations() (allocs []*objects.Allocation) {
	goa.Lock()
	defer goa.Unlock()
	allocs = make([]*objects.Allocation, 0)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	users := goa.metadata.UserData.UserCount
	nodes := goa.metadata.NodeData.NodeCount

	if users*nodes == 0 {
		return nil
	}

	for i := 0; i < goa.hyperParameter.grasshopperAmount; i++ {
		grasshopperArray := make([]int, users*nodes)
		for j := 0; j < len(grasshopperArray); j++ {
			grasshopperArray[j] = int(r.Int63n(6))
		}
		grasshopper := vector.NewVectorByInt(grasshopperArray)
		goa.grasshoppers = append(goa.grasshoppers, grasshopper)
	}


	candidates := make([]*vector.Vector, 0)

	for i := 0; i <  goa.hyperParameter.grasshopperAmount; i++ {
		candidateArray := make([]int, users*nodes)
		for j := 0; j < len(candidateArray); j++ {
			candidateArray[j] = int(r.Int63n(6))
		}
		candidate := vector.NewVectorByInt(candidateArray)
		candidates = append(candidates, candidate)

	}

	decision := goa.Start(candidates)

	removeIndexs := make([]int, 0)
	visited := make([]int, users)
	for i := 0; i < users; i++ {
		visited[i] = 0
	}

	for nodeIndex, _ := range goa.metadata.Nodes {
		for userIndex, _ := range goa.metadata.Requests {
			if distributeAmount := decision[nodeIndex*users+userIndex]; distributeAmount != 0 {
				nodeId := goa.metadata.Nodes[nodeIndex]
				ask := goa.metadata.Requests[userIndex]
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
		goa.RemoveUser(index)
	}
	return allocs
}