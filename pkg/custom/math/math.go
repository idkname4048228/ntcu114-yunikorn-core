package math

import (
	"math"
)

func SocialInfluence(r float64) float64 {
	f, l := 0.5, 1.5
	return f*math.Exp(-r/l) - math.Exp(-r)
}

func Sum(arr []float64) float64 {
	total := 0.0
	for _, value := range arr {
		total += value
	}
	return total
}

func SumOfSquares(arr []float64) float64 {
	total := 0.0
	for _, value := range arr {
		total += value * value
	}
	return total
}
