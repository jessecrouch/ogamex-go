package formula

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"ogamex-go/internal/domain"
)

type WikiShipTestData struct {
	Ships   map[string]WikiShipData `json:"ships"`
	Engines map[string]EngineData   `json:"engines"`
}

type WikiShipData struct {
	Cost                   ShipCostData       `json:"cost"`
	CargoCapacity          int64              `json:"cargo_capacity"`
	BaseSpeed              int64              `json:"base_speed"`
	StructuralIntegrity    int64              `json:"structural_integrity"`
	Shield                 int64              `json:"shield"`
	Weapon                 int64              `json:"weapon"`
	EngineType             string             `json:"engine_type"`
	RapidFireAgainst       map[string]int     `json:"rapid_fire_against"`
}

type ShipCostData struct {
	Metal     int `json:"metal"`
	Crystal   int `json:"crystal"`
	Deuterium int `json:"deuterium"`
}

type EngineData struct {
	SpeedBonusPerLevel float64 `json:"speed_bonus_per_level"`
	Description         string   `json:"description"`
}

func getWikiShipDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiShipData(t *testing.T) WikiShipTestData {
	data, err := os.ReadFile(getWikiShipDataPath("ships.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki ship test data: %v", err)
	}

	var wikiData WikiShipTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki ship test data: %v", err)
	}

	return wikiData
}

func TestWiki_ShipCargoCapacity(t *testing.T) {
	wikiData := loadWikiShipData(t)

	shipMap := map[string]domain.UnitType{
		"small_cargo":     domain.UnitSmallCargo,
		"large_cargo":     domain.UnitLargeCargo,
		"light_fighter":   domain.UnitLightFighter,
		"heavy_fighter":   domain.UnitHeavyFighter,
		"cruiser":         domain.UnitCruiser,
		"battleship":      domain.UnitBattleship,
		"battlecruiser":   domain.UnitBattlecruiser,
		"bomber":          domain.UnitBomber,
		"destroyer":       domain.UnitDestroyer,
		"deathstar":       domain.UnitDeathstar,
		"recycler":        domain.UnitRecycler,
		"espionage_probe": domain.UnitEspionageProbe,
		"solar_satellite": domain.UnitSolarSatellite,
		"colony_ship":     domain.UnitColonyShip,
		"crawler":         domain.UnitCrawler,
		"reaper":          domain.UnitReaper,
		"pathfinder":      domain.UnitPathfinder,
	}

	for shipName, unitType := range shipMap {
		wikiShip, exists := wikiData.Ships[shipName]
		if !exists {
			continue
		}

		t.Run(shipName, func(t *testing.T) {
			got := CalculateShipCargoCapacity(unitType)
			if got != wikiShip.CargoCapacity {
				t.Logf("DISCREPANCY: CalculateShipCargoCapacity(%s): expected=%d, got=%d", shipName, wikiShip.CargoCapacity, got)
			}
		})
	}
}

func TestWiki_ShipBaseSpeed(t *testing.T) {
	wikiData := loadWikiShipData(t)

	shipMap := map[string]domain.UnitType{
		"small_cargo":     domain.UnitSmallCargo,
		"large_cargo":     domain.UnitLargeCargo,
		"light_fighter":   domain.UnitLightFighter,
		"heavy_fighter":   domain.UnitHeavyFighter,
		"cruiser":         domain.UnitCruiser,
		"battleship":      domain.UnitBattleship,
		"battlecruiser":   domain.UnitBattlecruiser,
		"bomber":          domain.UnitBomber,
		"destroyer":       domain.UnitDestroyer,
		"deathstar":       domain.UnitDeathstar,
		"recycler":        domain.UnitRecycler,
		"espionage_probe": domain.UnitEspionageProbe,
		"solar_satellite": domain.UnitSolarSatellite,
		"colony_ship":     domain.UnitColonyShip,
		"crawler":         domain.UnitCrawler,
		"reaper":          domain.UnitReaper,
		"pathfinder":      domain.UnitPathfinder,
	}

	for shipName, unitType := range shipMap {
		wikiShip, exists := wikiData.Ships[shipName]
		if !exists {
			continue
		}

		t.Run(shipName, func(t *testing.T) {
			got := CalculateShipBaseSpeed(unitType)
			if got != wikiShip.BaseSpeed {
				t.Logf("DISCREPANCY: CalculateShipBaseSpeed(%s): expected=%d, got=%d", shipName, wikiShip.BaseSpeed, got)
			}
		})
	}
}

func TestWiki_ShipStructuralIntegrity(t *testing.T) {
	wikiData := loadWikiShipData(t)

	shipMap := map[string]domain.UnitType{
		"small_cargo":     domain.UnitSmallCargo,
		"large_cargo":     domain.UnitLargeCargo,
		"light_fighter":   domain.UnitLightFighter,
		"heavy_fighter":   domain.UnitHeavyFighter,
		"cruiser":         domain.UnitCruiser,
		"battleship":      domain.UnitBattleship,
		"battlecruiser":   domain.UnitBattlecruiser,
		"bomber":          domain.UnitBomber,
		"destroyer":       domain.UnitDestroyer,
		"deathstar":       domain.UnitDeathstar,
		"recycler":        domain.UnitRecycler,
		"espionage_probe": domain.UnitEspionageProbe,
		"solar_satellite": domain.UnitSolarSatellite,
		"colony_ship":     domain.UnitColonyShip,
		"crawler":         domain.UnitCrawler,
		"reaper":          domain.UnitReaper,
		"pathfinder":      domain.UnitPathfinder,
	}

	for shipName, unitType := range shipMap {
		wikiShip, exists := wikiData.Ships[shipName]
		if !exists {
			continue
		}

		t.Run(shipName, func(t *testing.T) {
			got := CalculateShipStructuralIntegrity(unitType)
			if got != wikiShip.StructuralIntegrity {
				t.Logf("DISCREPANCY: CalculateShipStructuralIntegrity(%s): expected=%d, got=%d", shipName, wikiShip.StructuralIntegrity, got)
			}
		})
	}
}

func TestWiki_ShipShield(t *testing.T) {
	wikiData := loadWikiShipData(t)

	shipMap := map[string]domain.UnitType{
		"small_cargo":     domain.UnitSmallCargo,
		"large_cargo":     domain.UnitLargeCargo,
		"light_fighter":   domain.UnitLightFighter,
		"heavy_fighter":   domain.UnitHeavyFighter,
		"cruiser":         domain.UnitCruiser,
		"battleship":      domain.UnitBattleship,
		"battlecruiser":   domain.UnitBattlecruiser,
		"bomber":          domain.UnitBomber,
		"destroyer":       domain.UnitDestroyer,
		"deathstar":       domain.UnitDeathstar,
		"recycler":        domain.UnitRecycler,
		"espionage_probe": domain.UnitEspionageProbe,
		"solar_satellite": domain.UnitSolarSatellite,
		"colony_ship":     domain.UnitColonyShip,
		"crawler":         domain.UnitCrawler,
		"reaper":          domain.UnitReaper,
		"pathfinder":      domain.UnitPathfinder,
	}

	for shipName, unitType := range shipMap {
		wikiShip, exists := wikiData.Ships[shipName]
		if !exists {
			continue
		}

		t.Run(shipName, func(t *testing.T) {
			got := CalculateShipShield(unitType)
			if got != wikiShip.Shield {
				t.Logf("DISCREPANCY: CalculateShipShield(%s): expected=%d, got=%d", shipName, wikiShip.Shield, got)
			}
		})
	}
}

func TestWiki_ShipWeapon(t *testing.T) {
	wikiData := loadWikiShipData(t)

	shipMap := map[string]domain.UnitType{
		"small_cargo":     domain.UnitSmallCargo,
		"large_cargo":     domain.UnitLargeCargo,
		"light_fighter":   domain.UnitLightFighter,
		"heavy_fighter":   domain.UnitHeavyFighter,
		"cruiser":         domain.UnitCruiser,
		"battleship":      domain.UnitBattleship,
		"battlecruiser":   domain.UnitBattlecruiser,
		"bomber":          domain.UnitBomber,
		"destroyer":       domain.UnitDestroyer,
		"deathstar":       domain.UnitDeathstar,
		"recycler":        domain.UnitRecycler,
		"espionage_probe": domain.UnitEspionageProbe,
		"solar_satellite": domain.UnitSolarSatellite,
		"colony_ship":     domain.UnitColonyShip,
		"crawler":         domain.UnitCrawler,
		"reaper":          domain.UnitReaper,
		"pathfinder":      domain.UnitPathfinder,
	}

	for shipName, unitType := range shipMap {
		wikiShip, exists := wikiData.Ships[shipName]
		if !exists {
			continue
		}

		t.Run(shipName, func(t *testing.T) {
			got := CalculateShipWeapon(unitType)
			if got != wikiShip.Weapon {
				t.Logf("DISCREPANCY: CalculateShipWeapon(%s): expected=%d, got=%d", shipName, wikiShip.Weapon, got)
			}
		})
	}
}

func TestWiki_ShipSpeedWithUpgrades(t *testing.T) {
	tests := []struct {
		name             string
		unitType         domain.UnitType
		combustionLevel  int
		impulseLevel     int
		hyperspaceLevel  int
		expectedSpeed    int64
	}{
		{"Light Fighter base", domain.UnitLightFighter, 0, 0, 0, 12500},
		{"Light Fighter combustion 5", domain.UnitLightFighter, 5, 0, 0, 13750},
		{"Cruiser base", domain.UnitCruiser, 0, 0, 0, 15000},
		{"Cruiser impulse 5", domain.UnitCruiser, 0, 5, 0, 27000},
		{"Battleship base", domain.UnitBattleship, 0, 0, 0, 10000},
		{"Battleship hyperspace 5", domain.UnitBattleship, 0, 0, 5, 25000},
		{"Deathstar base", domain.UnitDeathstar, 0, 0, 0, 100},
		{"Deathstar hyperspace 10", domain.UnitDeathstar, 0, 0, 10, 400},
		{"Espionage Probe base", domain.UnitEspionageProbe, 0, 0, 0, 100000000},
		{"Espionage Probe combustion 10", domain.UnitEspionageProbe, 10, 0, 0, 200000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateShipSpeed(tt.unitType, tt.combustionLevel, tt.impulseLevel, tt.hyperspaceLevel)
			if got != tt.expectedSpeed {
				t.Logf("DISCREPANCY: CalculateShipSpeed(%s, combustion=%d, impulse=%d, hyperspace=%d): expected=%d, got=%d",
					tt.name, tt.combustionLevel, tt.impulseLevel, tt.hyperspaceLevel, tt.expectedSpeed, got)
			}
		})
	}
}
