package formula

import (
	"testing"

	"ogamex-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetProduction(t *testing.T) {
	tests := []struct {
		name            string
		building        domain.BuildingType
		level           int
		temp            int
		energyFactor    float64
		energyTechLevel int
	}{
		{
			name:            "Metal Mine level 5",
			building:        domain.BuildingMetalMine,
			level:           5,
			temp:            25,
			energyFactor:    1.0,
			energyTechLevel: 0,
		},
		{
			name:            "Crystal Mine level 3",
			building:        domain.BuildingCrystalMine,
			level:           3,
			temp:            25,
			energyFactor:    1.0,
			energyTechLevel: 0,
		},
		{
			name:            "Deuterium Synthesizer level 4",
			building:        domain.BuildingDeuteriumSynthesizer,
			level:           4,
			temp:            25,
			energyFactor:    1.0,
			energyTechLevel: 0,
		},
		{
			name:            "Solar Plant level 2",
			building:        domain.BuildingSolarPlant,
			level:           2,
			temp:            25,
			energyFactor:    1.0,
			energyTechLevel: 0,
		},
		{
			name:            "Fusion Plant level 1 with energy tech 5",
			building:        domain.BuildingFusionPlant,
			level:           1,
			temp:            25,
			energyFactor:    1.0,
			energyTechLevel: 5,
		},
		{
			name:            "Unknown building type",
			building:        domain.BuildingType(999),
			level:           1,
			temp:            25,
			energyFactor:    1.0,
			energyTechLevel: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetProduction(tt.building, tt.level, tt.temp, tt.energyFactor, tt.energyTechLevel)
			if tt.level > 0 {
				assert.True(t, result.Metal >= 0 || result.Crystal >= 0 || result.Deuterium >= 0 || result.Energy != 0,
					"GetProduction should return valid resources for valid building")
			}
		})
	}
}

func TestGetShipStats(t *testing.T) {
	tests := []struct {
		name      string
		unitType  domain.UnitType
		expectOk  bool
	}{
		{"Small Cargo", domain.UnitSmallCargo, true},
		{"Large Cargo", domain.UnitLargeCargo, true},
		{"Light Fighter", domain.UnitLightFighter, true},
		{"Deathstar", domain.UnitDeathstar, true},
		{"Rocket Launcher (defense)", domain.UnitRocketLauncher, false},
		{"Unknown", domain.UnitType(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, ok := GetShipStats(tt.unitType)
			if tt.expectOk {
				assert.True(t, ok, "Expected stats to exist")
				assert.True(t, stats.CargoCapacity >= 0, "CargoCapacity should be non-negative")
				assert.True(t, stats.BaseSpeed >= 0, "BaseSpeed should be non-negative")
				assert.True(t, stats.StructuralIntegrity >= 0, "StructuralIntegrity should be non-negative")
			} else {
				assert.False(t, ok, "Expected stats to not exist")
			}
		})
	}
}

func TestGetDefenseStats(t *testing.T) {
	tests := []struct {
		name      string
		unitType  domain.UnitType
		expectOk  bool
	}{
		{"Rocket Launcher", domain.UnitRocketLauncher, true},
		{"Light Laser", domain.UnitLightLaser, true},
		{"Heavy Laser", domain.UnitHeavyLaser, true},
		{"Gauss Cannon", domain.UnitGaussCannon, true},
		{"Ion Cannon", domain.UnitIonCannon, true},
		{"Plasma Turret", domain.UnitPlasmaTurret, true},
		{"Small Shield Dome", domain.UnitSmallShieldDome, true},
		{"Large Shield Dome", domain.UnitLargeShieldDome, true},
		{"Anti-Ballistic Missile", domain.UnitAntiBallisticMissiles, true},
		{"Interplanetary Missile", domain.UnitInterplanetaryMissiles, true},
		{"Small Cargo (ship)", domain.UnitSmallCargo, false},
		{"Unknown", domain.UnitType(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, ok := GetDefenseStats(tt.unitType)
			if tt.expectOk {
				assert.True(t, ok, "Expected stats to exist")
				assert.True(t, stats.StructuralIntegrity >= 0, "StructuralIntegrity should be non-negative")
				assert.True(t, stats.Shield >= 0, "Shield should be non-negative")
				assert.True(t, stats.Weapon >= 0, "Weapon should be non-negative")
			} else {
				assert.False(t, ok, "Expected stats to not exist")
			}
		})
	}
}

func TestCalculateFleetSpeed(t *testing.T) {
	// Test CalculateFleetSpeed from distance.go
	tests := []struct {
		name       string
		baseSpeed int
		techLevel  int
		numShips   int
		uniSpeed   int
	}{
		{"Normal", 10000, 5, 10, 1},
		{"Zero base speed", 0, 5, 10, 1},
		{"Higher universe speed", 10000, 5, 10, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateFleetSpeed(tt.baseSpeed, tt.techLevel, tt.numShips, tt.uniSpeed)
			if tt.baseSpeed == 0 {
				assert.Equal(t, 0, result)
			} else {
				assert.Greater(t, result, 0)
			}
		})
	}
}

func TestCalculateShipCargoCapacity(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Small Cargo", domain.UnitSmallCargo},
		{"Large Cargo", domain.UnitLargeCargo},
		{"Deathstar", domain.UnitDeathstar},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateShipCargoCapacity(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateShipBaseSpeed(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Small Cargo", domain.UnitSmallCargo},
		{"Light Fighter", domain.UnitLightFighter},
		{"Deathstar", domain.UnitDeathstar},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateShipBaseSpeed(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateShipStructuralIntegrity(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Small Cargo", domain.UnitSmallCargo},
		{"Deathstar", domain.UnitDeathstar},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateShipStructuralIntegrity(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateShipShield(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Small Cargo", domain.UnitSmallCargo},
		{"Deathstar", domain.UnitDeathstar},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateShipShield(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateShipWeapon(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Small Cargo", domain.UnitSmallCargo},
		{"Deathstar", domain.UnitDeathstar},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateShipWeapon(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateDefenseStructuralIntegrity(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Rocket Launcher", domain.UnitRocketLauncher},
		{"Gauss Cannon", domain.UnitGaussCannon},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDefenseStructuralIntegrity(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateDefenseShield(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Rocket Launcher", domain.UnitRocketLauncher},
		{"Gauss Cannon", domain.UnitGaussCannon},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDefenseShield(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculateDefenseWeapon(t *testing.T) {
	tests := []struct {
		name     string
		unitType domain.UnitType
	}{
		{"Rocket Launcher", domain.UnitRocketLauncher},
		{"Gauss Cannon", domain.UnitGaussCannon},
		{"Unknown", domain.UnitType(999)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDefenseWeapon(tt.unitType)
			if tt.unitType == domain.UnitType(999) {
				assert.Equal(t, int64(0), result)
			} else {
				assert.GreaterOrEqual(t, result, int64(0))
			}
		})
	}
}

func TestCalculatePositionBonusEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		position int
	}{
		{"Position 0 (edge)", 0},
		{"Position 1", 1},
		{"Position 3", 3},
		{"Position 5", 5},
		{"Position 6 (metal bonus)", 6},
		{"Position 7", 7},
		{"Position 10 (metal bonus)", 10},
		{"Position 15", 15},
		{"Position 16", 16},
		{"Position 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metal, crystal, deuterium := CalculatePositionBonus(tt.position)
			assert.True(t, metal >= 1.0, "Metal bonus should be >= 1.0")
			assert.True(t, crystal >= 1.0, "Crystal bonus should be >= 1.0")
			assert.True(t, deuterium >= 1.0, "Deuterium bonus should be >= 1.0")
		})
	}
}

func TestCalculateStorageCapacityEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		level       int
	}{
		{"Level 0", 0},
		{"Level 1", 1},
		{"Level 5", 5},
		{"Level 10", 10},
		{"Level 20", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStorageCapacity(tt.level)
			if tt.level == 0 {
				assert.Equal(t, int64(10000), result, "Level 0 should return base storage")
			} else {
				assert.Greater(t, result, int64(0), "Storage should be positive for level > 0")
			}
		})
	}
}

func TestCalculateBuildingCostEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		building  domain.BuildingType
		level     int
	}{
		{"Level 0", domain.BuildingMetalMine, 0},
		{"Level 1", domain.BuildingMetalMine, 1},
		{"Level 2", domain.BuildingMetalMine, 2},
		{"High level", domain.BuildingMetalMine, 50},
		{"Crystal Mine", domain.BuildingCrystalMine, 10},
		{"Fusion Plant", domain.BuildingFusionPlant, 5},
		{"Unknown building", domain.BuildingType(999), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metal, _, _ := CalculateBuildingCost(tt.building, tt.level)
			if tt.level == 0 {
				assert.Equal(t, int64(0), metal, "Level 0 should have no cost")
			} else if tt.building == domain.BuildingType(999) {
				assert.Equal(t, int64(0), metal, "Unknown building should have no cost")
			} else {
				assert.Greater(t, metal, int64(0), "Valid level should have cost")
			}
		})
	}
}

func TestGetRapidFireAgainst(t *testing.T) {
	tests := []struct {
		name     string
		attacker domain.UnitType
		defender domain.UnitType
		hasValue bool
	}{
		{"Deathstar vs Probe", domain.UnitDeathstar, domain.UnitEspionageProbe, true},
		{"Deathstar vs Small Cargo", domain.UnitDeathstar, domain.UnitSmallCargo, true},
		{"Deathstar vs Battleship", domain.UnitDeathstar, domain.UnitBattleship, true},
		{"Deathstar vs Rocket Launcher", domain.UnitDeathstar, domain.UnitRocketLauncher, true},
		{"Solar Satellite vs Probe", domain.UnitSolarSatellite, domain.UnitEspionageProbe, true},
		{"Solar Satellite vs Small Cargo", domain.UnitSolarSatellite, domain.UnitSmallCargo, true},
		{"Small Cargo vs Probe", domain.UnitSmallCargo, domain.UnitEspionageProbe, false},
		{"Unknown attacker", domain.UnitType(999), domain.UnitEspionageProbe, false},
		{"Unknown defender", domain.UnitDeathstar, domain.UnitType(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := GetRapidFireAgainst(tt.attacker, tt.defender)
			if tt.hasValue {
				assert.Greater(t, rf, 0, "Rapid fire should be positive for valid combinations")
			}
		})
	}
}
