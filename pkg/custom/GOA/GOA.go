package GOA


import (
	"math"
	"math/rand"
	"time"
	"fmt"
	"sync"

	goamath "github.com/apache/yunikorn-core/pkg/custom/GOA/math"
	"github.com/apache/yunikorn-core/pkg/custom/GOA/math/vector"
	"github.com/apache/yunikorn-core/pkg/custom/GOA/metadata"

	"github.com/apache/yunikorn-core/pkg/scheduler/objects"
	"github.com/apache/yunikorn-core/pkg/log"
)

var goa *GOA
var (
	iterations = 50 
	cMax = 1.0
	cMin = 0.00001
	grasshopperAmount = 30
	GForce = 0.2
	WindForce = 0.2
)

func GetGOA() *GOA{
	return goa
}

func Init(){
	goa = NewGOA();
}

type GOA struct {
	metaData 			*metadata.MetaData
	fairness_vector		*vector.Vector
	
	nodeCount			int
	resourceCount		int
	userCount 			int

	grasshoppers		[]*vector.Vector
	bestGrasshopper		*vector.Vector

	sync.RWMutex
}

func NewGOA() *GOA{
	metaData := metadata.NewMetaData()
	return &GOA{
		metaData: metaData,
		nodeCount: metaData.NodeData.NodeCount,
		resourceCount: metaData.NodeData.ResourceCount,
		userCount: metaData.UserData.UserCount, 
		grasshoppers: make([]*vector.Vector, grasshopperAmount),
	}
}

func (goa *GOA) AddNode(n *objects.Node) {
	goa.Lock()
	defer goa.Unlock()
	goa.nodeCount += 1
	goa.metaData.AddNode(n)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", goa.nodeCount))
}

func (goa *GOA) AddUser(ask *objects.AllocationAsk, app *objects.Application){
	goa.Lock()
	defer goa.Unlock()
	goa.userCount += 1
	goa.metaData.AddUser(ask, app)
	log.Log(log.Custom).Info(fmt.Sprintf("userCount is %v", goa.userCount))
}

func (goa *GOA) RemoveUser(index int) {
	log.Log(log.Custom).Info("removing user")
	goa.userCount -= 1
	goa.metaData.RemoveUser(index)
	log.Log(log.Custom).Info("removed user")

}
// GOA utils
func (goa *GOA) getGravityUnitVector(grasshopper *vector.Vector) *vector.Vector{
	grasshopper_unit := grasshopper.GetUnitVector()
	
	t := goa.fairness_vector.Dot(grasshopper_unit) / grasshopper_unit.Dot(grasshopper_unit)

	D := grasshopper_unit.Multiple(t)

	gravity_vector := vector.Subtract(goa.fairness_vector, D)

	return gravity_vector.GetUnitVector()
}

func (goa *GOA) getWindUnitVector(grasshopper *vector.Vector, best_grasshopper *vector.Vector) *vector.Vector{
	wind_vector := vector.Subtract(best_grasshopper, grasshopper)
	return wind_vector.GetUnitVector()
}

func (goa *GOA) getEffectScore(grasshopper *vector.Vector) float64 {
	userTotal := make([]float64, goa.resourceCount)

	for i := 0; i < goa.nodeCount; i++ {
		for j := 0; j < goa.userCount; j++ {
			amountThatUserJTakeAtNodeI := int(grasshopper.Get(i*goa.userCount+j))
			

			if amountThatUserJTakeAtNodeI < 0 {
				return math.Inf(1)
			} 
			for k := 0; k < goa.resourceCount; k++ {
				userTotal[k] += float64(amountThatUserJTakeAtNodeI) * goa.metaData.GetUserAsks()[j][k]
				if userTotal[k] > goa.metaData.GetTotalLimits()[k] {
					return math.Inf(1)
				}
			}
		}
	}

	totalLimitsVector := vector.NewVector(goa.metaData.GetTotalLimits())
	userTotalVector := vector.NewVector(userTotal)
	return vector.Subtract(totalLimitsVector, userTotalVector).Norm()
}

func (goa *GOA) getFairnessScore(grasshopper *vector.Vector) float64 {
	users := goa.userCount
	usersDR := make([]float64, goa.userCount)


	for i := 0; i < goa.nodeCount; i++ {
		for j := 0; j < users; j++ {
			amountThatUserJTakeAtNodeI := grasshopper.Get(i*users+j)
			if amountThatUserJTakeAtNodeI < 0 {
				return math.Inf(1)
			}
			// fmt.Printf("DRF log: %d times %f\n", amountThatUserJTakeAtNodeI, DRs[i][j])
			usersDR[j] += float64(amountThatUserJTakeAtNodeI) * goa.metaData.DRs[i][j]
		}
	}

	if goamath.Sum(usersDR) == 0 {
		return math.Inf(1)
	}

	jainIndexValue := math.Pow(goamath.Sum(usersDR), 2) / (goamath.SumOfSquares(usersDR) * float64(users))

	return jainIndexValue
}

func (goa *GOA) getScore(grasshopper *vector.Vector) float64 {
	effectScore := goa.getEffectScore(grasshopper)
	if effectScore == math.Inf(1) {
		return math.Inf(1)
	}
	fairnessScore := goa.getFairnessScore(grasshopper)
	if fairnessScore == math.Inf(1) {
		return math.Inf(1)
	}

	return effectScore / fairnessScore
}

func (goa *GOA) greedyMove(minValue float64, nextPositions *[]*vector.Vector) {
	// goa.bestGrasshopper = nil
	
	for i := 0; i < grasshopperAmount; i++ {
		nextPosition := (*nextPositions)[i]

		
		oldValue := goa.getScore(goa.grasshoppers[i])
		newValue := goa.getScore(nextPosition)
		
		if newValue != math.Inf(1) {
			if newValue < oldValue {

				goa.grasshoppers[i] = nextPosition
				oldValue = newValue
			}

			if (oldValue < minValue || goa.bestGrasshopper == nil) {
				minValue = oldValue

				goa.bestGrasshopper = goa.grasshoppers[i]
			}
		}

		// 加入 G_force * calculate_gravity_unit_vector 的結果到 grasshoppers[i]
		gravityVector := goa.getGravityUnitVector(goa.grasshoppers[i])
		goa.grasshoppers[i].Add(gravityVector.Multiple(GForce))

	}
}

func (goa *GOA) calculateDomainResources() {
	goa.metaData.CalculateDRs()
	fairness_array := make([]float64, 0)
	for _, value := range goa.metaData.DRRatioReciprocals {
		fairness_array = append(fairness_array, value...)
	}
	goa.fairness_vector = vector.NewVector(fairness_array) 
}

func (goa *GOA) StartScheduler() (decision []int) {
	users := goa.userCount
	nodes := goa.nodeCount

	// log.Log(log.Custom).Info(fmt.Sprintf("now : %v, %v", users, nodes))

	if users * nodes == 0 {
		return nil
	}

	log.Log(log.Custom).Info(fmt.Sprintf("goa start : %v, %v", users, nodes))
	log.Log(log.Custom).Info(fmt.Sprintf("metadata is %v", goa.metaData))
	log.Log(log.Custom).Info(fmt.Sprintf("nodedata is %v", goa.metaData.NodeData))
	log.Log(log.Custom).Info(fmt.Sprintf("userdata is %v", goa.metaData.UserData))

	
	goa.calculateDomainResources()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	minValue := math.Inf(1)
	goa.grasshoppers = make([]*vector.Vector, 0)
	goa.bestGrasshopper = vector.WithSize(users * nodes)



	for i := 0; i < grasshopperAmount; i++ {
		grasshopperArray := make([]int, users*nodes)
		for j := 0; j < len(grasshopperArray); j++ {
			grasshopperArray[j] = int(r.Int63n(6))
		}
		grasshopper := vector.NewVectorByInt(grasshopperArray)
		value := goa.getScore(grasshopper)
		if minValue > value {
			minValue = value
			goa.bestGrasshopper = grasshopper
		}

		goa.grasshoppers = append(goa.grasshoppers, grasshopper)
	} 

	log.Log(log.Custom).Info(fmt.Sprintf("grasshoppers: %v", goa.grasshoppers))
	log.Log(log.Custom).Info(fmt.Sprintf("best grasshopper: %v", goa.bestGrasshopper))

	// main program
	for i := 0; i < iterations; i++ {
		c := cMax - float64(i) * ((cMax - cMin) / float64(iterations))
		nextPositions := make([]*vector.Vector, 0)
		for grasshopper_i := 0; grasshopper_i < grasshopperAmount; grasshopper_i++ {
			moveVector := vector.WithSize(users * nodes)
			for grasshopper_j := 0; grasshopper_j < grasshopperAmount; grasshopper_j++ {
				if grasshopper_i == grasshopper_j {
					continue
				}
				distance := vector.Subtract(goa.grasshoppers[grasshopper_i], goa.grasshoppers[grasshopper_j])
				if distance.Norm() == 0 {
					continue
				}
				moveVector.Add(distance.GetUnitVector().Multiple(goamath.SocialInfluence(distance.Norm())))
				
			}
			moveVector = moveVector.Multiple(c)

			moveVector.Add(goa.getWindUnitVector(goa.grasshoppers[grasshopper_i], goa.bestGrasshopper).Multiple(WindForce))
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
	
	decision := goa.StartScheduler()
	users := goa.userCount

	removeIndexs := make([]int, 0)
	visited := make([]int, users)
	for i := 0; i < users; i++ {
		visited[i] = 0
	}

	for nodeIndex, _ := range goa.metaData.Nodes {
		for userIndex, _ := range goa.metaData.Requests {
			if distributeAmount := decision[nodeIndex * users + userIndex]; distributeAmount != 0 {
				nodeId := goa.metaData.Nodes[nodeIndex]
				ask := goa.metaData.Requests[userIndex]
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