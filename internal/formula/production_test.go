package formula

import (
	"testing"
)

func TestCalculateMetalMineProduction(t *testing.T) {
	tests := []struct {
		name        string
		level       int
		energyFactor float64
		want        int64
	}{
		{"level 0", 0, 1.0, 0},
		{"level 1", 1, 1.0, 33},      // 30 * 1 * 1.1^1 = 33
		{"level 10", 10, 1.0, 779},    // 30 * 10 * 1.1^10 = 779
		{"level 20", 20, 1.0, 4037},   // 30 * 20 * 1.1^20 = 4037
		{"level 20 with 0.5 energy", 20, 0.5, 2019},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateMetalMineProduction(tt.level, tt.energyFactor)
			if got != tt.want {
				t.Errorf("CalculateMetalMineProduction(%d, %v) = %d, want %d", tt.level, tt.energyFactor, got, tt.want)
			}
		})
	}
}

func TestCalculateSolarSatelliteEnergyProduction(t *testing.T) {
	tests := []struct {
		name       string
		temp       int
		count      int
		want       int64
	}{
		{"100 temp, 100 count", 100, 100, 4000}, // floor((100+140)/6) = 40 * 100 = 4000
		{"50 temp, 10 count", 50, 10, 310},      // floor((50+140)/6) = 31 * 10 = 310
		{"0 temp, 10 count", 0, 10, 230},         // floor((0+140)/6) = 23 * 10 = 230
		{"negative temp, 10 count", -50, 10, 150}, // floor((-50+140)/6) = 15 * 10 = 150
		{"zero count", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSolarSatelliteEnergyProduction(tt.temp, tt.count)
			if got != tt.want {
				t.Errorf("CalculateSolarSatelliteEnergyProduction(%d, %d) = %d, want %d", tt.temp, tt.count, got, tt.want)
			}
		})
	}
}

func TestCalculateSolarPlantEnergyProduction(t *testing.T) {
	tests := []struct {
		name  string
		level int
		want  int64
	}{
		{"level 0", 0, 0},
		{"level 1", 1, 22},   // 20 * 1 * 1.1^1 = 22
		{"level 10", 10, 519}, // 20 * 10 * 1.1^10 = 519
		{"level 20", 20, 2691}, // 20 * 20 * 1.1^20 = 2691
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateSolarPlantEnergyProduction(tt.level)
			if got != tt.want {
				t.Errorf("CalculateSolarPlantEnergyProduction(%d) = %d, want %d", tt.level, got, tt.want)
			}
		})
	}
}

func TestCalculateFusionPlantEnergyProduction(t *testing.T) {
	tests := []struct {
		name            string
		level           int
		energyTechLevel int
		want            int64
	}{
		{"level 0", 0, 0, 0},
		{"level 1, no tech", 1, 0, 32},    // 30 * 1 * (1.05)^1 = 31.5 → ceil = 32
		{"level 10, no tech", 10, 0, 489},  // 30 * 10 * 1.05^10 = 489
		{"level 10, level 5 tech", 10, 5, 779}, // 30 * 10 * 1.1^10 = 779
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFusionPlantEnergyProduction(tt.level, tt.energyTechLevel)
			if got != tt.want {
				t.Errorf("CalculateFusionPlantEnergyProduction(%d, %d) = %d, want %d", tt.level, tt.energyTechLevel, got, tt.want)
			}
		})
	}
}

func TestCalculateDeuteriumProduction(t *testing.T) {
	tests := []struct {
		name        string
		level       int
		temp        int
		energyFactor float64
		want        int64
	}{
		{"level 0", 0, 50, 1.0, 0},
		{"level 1, temp 50", 1, 50, 1.0, 14},        // 10 * 1 * 1.1^1 * (1.44 - 0.004*50) = 14
		{"level 10, temp 47", 10, 47, 1.0, 325},      // from test
		{"level 20, temp 47", 20, 47, 1.0, 1685},      // calculated
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDeuteriumProduction(tt.level, tt.temp, tt.energyFactor)
			if got != tt.want {
				t.Errorf("CalculateDeuteriumProduction(%d, %d, %v) = %d, want %d", tt.level, tt.temp, tt.energyFactor, got, tt.want)
			}
		})
	}
}

func BenchmarkMetalMineProduction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateMetalMineProduction(20, 1.0)
	}
}
