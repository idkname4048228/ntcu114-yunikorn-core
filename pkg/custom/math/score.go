package math

import (
	// "fmt"
	"math"

	"github.com/apache/yunikorn-core/pkg/custom/math/vector"
	Metadata "github.com/apache/yunikorn-core/pkg/custom/metadata"
	// "github.com/apache/yunikorn-core/pkg/log"
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
	totalLimitsDistance := totalLimitsVector.Norm()
	userTotalDistance := userTotalVector.Norm()
	
	ratio := userTotalDistance / totalLimitsDistance

	// 計算「空閒資源比例」。使用「使用者資源上限」除以「機器資源上限」，(1 - 所得比例) * 100 為「空閒資源佔比」
	percentage := (1 - ratio) * 100

	return percentage
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

	// log.Log(log.Custom).Info(fmt.Sprintf("effect score is %v", effectScore))

	// effectScore: 0 ~ 100
	// fairnessScore: 1/len(user) ~ 1
	// lenUser: len(user)
	
	if fairnessScore >= 0.9 {
		fairnessScore = 1
	}

	score := effectScore / fairnessScore

	scoreMin := 0.0
	scoreMax := 100.0 * float64(metadata.UserData.UserCount)

	// 正規化 score 到 0 ~ 100
	normalizedScore := (score - scoreMin) / (scoreMax - scoreMin) * 100

	// 保證正規化分數在 0 ~ 100 範圍內
	if normalizedScore < 0 {
	    normalizedScore = 0
	} else if normalizedScore > 100 {
	    normalizedScore = 100
	}

	// log.Log(log.Custom).Info(fmt.Sprintf("score is %v", normalizedScore))
	return normalizedScore
}