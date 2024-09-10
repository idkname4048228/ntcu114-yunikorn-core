package math

import (
	"fmt"
	"math"
	"strings"
)

// 資料結構來儲存每個座標的分數總和和出現次數
type ScoreData struct {
	Sum   float64
	Count int
}

// 有序結果的結構體
type CoordScore struct {
	Coord []float64
	Score float64
}


// 計算曼哈頓距離
func ManhattanDistance(coord1, coord2 []float64) int {
	dist := 0
	for i := 0; i < len(coord1); i++ {
		dist += int(math.Abs(coord1[i] - coord2[i]))
	}
	return dist
}

func softmax(scores []CoordScore) []CoordScore {
	if len(scores) == 0 {
		return nil
	}

	// 步驟 1: 找到最大分數，進行數值穩定化處理
	maxScore := scores[0].Score
	for _, item := range scores {
		if item.Score > maxScore {
			maxScore = item.Score
		}
	}

	// 步驟 2: 計算指數和（減去最大分數以避免溢出）
	expSum := 0.0
	for _, item := range scores {
		expSum += math.Exp(item.Score - maxScore)
	}

	// 步驟 3: 計算 Softmax 結果
	softmaxScores := make([]CoordScore, len(scores))
	for i, item := range scores {
		softmaxScores[i] = CoordScore{
			Coord: item.Coord,
			Score: math.Exp(item.Score - maxScore) / expSum,
		}
	}

	return softmaxScores
}

// ori2soft 函數，計算 paths 與 scores 的平均分數並應用 Softmax
func SoftmaxNormalize(paths [][][]float64, scores []float64) []CoordScore {
	scoreMap := make(map[string]*ScoreData)

	// 迭代所有 paths
	for i, path := range paths {
		score := scores[i]
		for _, coord := range path {
			coordKey := fmt.Sprintf("%v", coord) // 將座標轉換成字串作為鍵值

			if _, exists := scoreMap[coordKey]; !exists {
				scoreMap[coordKey] = &ScoreData{}
			}

			scoreMap[coordKey].Sum += score
			scoreMap[coordKey].Count++
		}
	}

	// 計算每個座標的平均分數
	avgScores := []CoordScore{}
	for coordKey, data := range scoreMap {
		cleanCoordKey := strings.Trim(coordKey, "[]") // 移除方括號

		fields := strings.Fields(cleanCoordKey) // 按空格分割字串
		coord := make([]float64, len(fields)) // 根據 fields 的長度創建一個 float64 切片

		for i, field := range fields {
		    fmt.Sscanf(field, "%v", &coord[i]) // 依次將每個部分轉換為 float64
		}
		
		avgScore := data.Sum / float64(data.Count)
		if math.IsInf(avgScore, 0) {
			avgScores = append(avgScores, CoordScore{Coord: coord, Score: math.MaxFloat64})
		} else {
			avgScores = append(avgScores, CoordScore{Coord: coord, Score: avgScore})
		}
	}

	// 計算 Softmax
	softmaxScores := softmax(avgScores)

	return softmaxScores
}