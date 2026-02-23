package formula

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCostFactor(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		factor   float64
		expected float64
	}{
		{"level 1, factor 1.5", 1, 1.5, 1.0},
		{"level 2, factor 1.5", 2, 1.5, 1.5},
		{"level 3, factor 1.5", 3, 1.5, 2.25},
		{"level 10, factor 1.5", 10, 1.5, 38.44},
		{"level 1, factor 1.8", 1, 1.8, 1.0},
		{"level 2, factor 1.8", 2, 1.8, 1.8},
		{"level 5, factor 1.8", 5, 1.8, 10.50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCostFactor(tt.level, tt.factor)
			assert.InDelta(t, tt.expected, result, 0.01)
		})
	}
}

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name          string
		originG, originS, originP int
		targetG, targetS, targetP int
		expected      int
	}{
		{"same coords", 1, 1, 1, 1, 1, 1, 0},
		{"same system diff position", 1, 1, 1, 1, 1, 5, 2020},
		{"different galaxy", 1, 1, 1, 2, 1, 1, 20000},
		{"max distance", 1, 1, 1, 9, 499, 15, 294530},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDistance(tt.originG, tt.originS, tt.originP, tt.targetG, tt.targetS, tt.targetP)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateFlightTime(t *testing.T) {
	// Test basic functionality - returns positive values
	result := CalculateFlightTime(10000, 10000, 1)
	assert.Greater(t, result, int64(0))
	
	// Speed factor should reduce time
	result1 := CalculateFlightTime(10000, 10000, 1)
	result2 := CalculateFlightTime(10000, 10000, 2)
	assert.Less(t, result2, result1)
}

func TestCalculatePositionBonus(t *testing.T) {
	tests := []struct {
		name     string
		position int
		wantMetal, wantCrystal, wantDeuterium float64
	}{
		{"position 1", 1, 1.0, 1.0, 1.0},
		{"position 4", 4, 1.0, 1.1, 1.0},
		{"position 8", 8, 1.0, 1.2, 1.0},
		{"position 9", 9, 1.0, 1.3, 1.0},
		{"position 10", 10, 1.0, 1.4, 1.0},
		{"position 11", 11, 1.0, 1.6, 1.0},
		{"position 12", 12, 1.0, 1.7, 1.0},
		{"position 15", 15, 1.0, 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metal, crystal, deuterium := CalculatePositionBonus(tt.position)
			assert.Equal(t, tt.wantMetal, metal)
			assert.Equal(t, tt.wantCrystal, crystal)
			assert.Equal(t, tt.wantDeuterium, deuterium)
		})
	}
}
