package formula

import "math"

func GetCostFactor(level int, baseFactor float64) float64 {
	if level <= 1 {
		return 1.0
	}
	return math.Pow(baseFactor, float64(level-1))
}
