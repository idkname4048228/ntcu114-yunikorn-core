package math

import (
	// "fmt"
	"math"

	// "github.com/apache/yunikorn-core/pkg/log"
	"github.com/apache/yunikorn-core/pkg/custom/math/vector"
	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"
)

func GetEffectScore(metadata *Metadata.Metadata, candidate *vector.Vector) float64 {
	resources := metadata.NodeData.ResourceCount
	nodes := metadata.NodeData.NodeCount
	users := metadata.UserData.UserCount
	
	userTotal := make([]float64, resources)

	for i := 0; i < nodes; i++ {
		totalAtNodeI := make([]float64, resources)
		for j := 0; j < users; j++ {
			
			amountThatUserJTakeAtNodeI := int(candidate.Get(i * users + j))

			if amountThatUserJTakeAtNodeI < 0 {
				return math.Inf(1)
			} 
			for k := 0; k < resources; k++ {
				nodeIOccupied := float64(amountThatUserJTakeAtNodeI) * metadata.GetUserAsks()[j][k]
				totalAtNodeI[k] += nodeIOccupied

				if totalAtNodeI[k] > metadata.GetNodeLimits()[i][k] {
					return math.Inf(1)
				}
				userTotal[k] += nodeIOccupied
			}
		}
	}

	totalLimitsVector := vector.NewVector(metadata.GetTotalLimits())
	userTotalVector := vector.NewVector(userTotal)
	return vector.Subtract(totalLimitsVector, userTotalVector).Norm()
}

func GetFairnessScore(metadata *Metadata.Metadata, candidate *vector.Vector) float64 {
	nodes := metadata.NodeData.NodeCount
	users := metadata.UserData.UserCount

	usersDR := make([]float64, users)

	for i := 0; i < nodes; i++ {
		for j := 0; j < users; j++ {
			amountThatUserJTakeAtNodeI := candidate.Get(i*users+j)
			if amountThatUserJTakeAtNodeI < 0 {
				return math.Inf(1)
			}
			usersDR[j] += float64(amountThatUserJTakeAtNodeI) * metadata.DRs[i][j]
		}
	}
	if Sum(usersDR) == 0 {
		return math.Inf(1)
	}

	jainIndexValue := math.Pow(Sum(usersDR), 2) / (SumOfSquares(usersDR) * float64(users))

	return jainIndexValue
}

func GetScore(metadata *Metadata.Metadata, candidate *vector.Vector) float64 {
	effectScore := GetEffectScore(metadata, candidate)
	if effectScore == math.Inf(1) {
		return math.Inf(1)
	}
	fairnessScore := GetFairnessScore(metadata, candidate)
	if fairnessScore == math.Inf(1) {
		return math.Inf(1)
	}

	if fairnessScore >= 0.9 {
		return effectScore / 1.0
	} else {
		return effectScore / fairnessScore
	}
}