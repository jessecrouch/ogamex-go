package service

import (
	"testing"

	"ogamex-go/internal/schema"
)

func TestProductionService_CalculateProduction(t *testing.T) {
	// Create a mock planet with level 1 metal mine and solar plant
	planet := &schema.Planet{
		ID:                          1,
		UserID:                      1,
		MetalMine:                  1,
		CrystalMine:                0,
		DeuteriumSynthesizer:       0,
		SolarPlant:                 1,
		FusionPlant:                0,
		MetalMinePercent:           100,
		CrystalMinePercent:         100,
		DeuteriumSynthesizerPercent: 100,
		SolarPlantPercent:          100,
		FusionPlantPercent:         100,
		TempMax:                    50,
		SolarSatellite:             0,
	}

	tech := &schema.UserTech{
		UserID:             1,
		EnergyTechnology:  0,
	}

	// Create production service with economy speed 8 (matching config)
	productionService := NewProductionService(nil, nil, 8)

	result := productionService.CalculateProduction(planet, tech)

	// Metal mine level 1 produces 33 per hour base
	// With economy speed 8: 33 * 8 = 264
	if result.Metal != 264 {
		t.Errorf("Expected metal production 264, got %d", result.Metal)
	}

	// Solar plant level 1 produces energy
	// With energy tech 0, solar plant produces 20
	if result.Energy <= 0 {
		t.Errorf("Expected positive energy, got %d", result.Energy)
	}

	t.Logf("Production result: Metal=%d, Crystal=%d, Deuterium=%d, Energy=%d",
		result.Metal, result.Crystal, result.Deuterium, result.Energy)
}

func TestProductionService_NilTech(t *testing.T) {
	planet := &schema.Planet{
		ID:              1,
		UserID:          1,
		MetalMine:       1,
		SolarPlant:      1,
		MetalMinePercent: 100,
		SolarPlantPercent: 100,
		TempMax:         50,
	}

	productionService := NewProductionService(nil, nil, 1)

	// Test with nil tech - should not panic
	result := productionService.CalculateProduction(planet, nil)

	if result.Metal == 0 {
		t.Error("Expected non-zero metal production with level 1 mine, got 0")
	}

	t.Logf("Production with nil tech: Metal=%d, Energy=%d", result.Metal, result.Energy)
}

func TestProductionService_ZeroProductionPercentages(t *testing.T) {
	planet := &schema.Planet{
		ID:                    1,
		UserID:                1,
		MetalMine:             1,
		SolarPlant:            1,
		MetalMinePercent:      0, // 0% production
		SolarPlantPercent:     100,
		TempMax:               50,
	}

	tech := &schema.UserTech{
		UserID:            1,
		EnergyTechnology: 0,
	}

	productionService := NewProductionService(nil, nil, 1)

	result := productionService.CalculateProduction(planet, tech)

	// With 0% production, should be 0
	if result.Metal != 0 {
		t.Errorf("Expected 0 metal production with 0%%, got %d", result.Metal)
	}

	t.Logf("Production with 0%%: Metal=%d, Energy=%d", result.Metal, result.Energy)
}
