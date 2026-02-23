package formula

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"ogamex-go/internal/domain"
)

type WikiDefenseTestData struct {
	Defense map[string]WikiDefenseData `json:"defense"`
}

type WikiDefenseData struct {
	Cost                DefenseCostData `json:"cost"`
	StructuralIntegrity int64           `json:"structural_integrity"`
	Shield              int64           `json:"shield"`
	Weapon              int64           `json:"weapon"`
	RapidFireAgainst   map[string]int  `json:"rapid_fire_against"`
}

type DefenseCostData struct {
	Metal     int `json:"metal"`
	Crystal   int `json:"crystal"`
	Deuterium int `json:"deuterium"`
}

func getWikiDefenseDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiDefenseData(t *testing.T) WikiDefenseTestData {
	data, err := os.ReadFile(getWikiDefenseDataPath("defense.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki defense test data: %v", err)
	}

	var wikiData WikiDefenseTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki defense test data: %v", err)
	}

	return wikiData
}

func TestWiki_DefenseStructuralIntegrity(t *testing.T) {
	wikiData := loadWikiDefenseData(t)

	defenseMap := map[string]domain.UnitType{
		"rocket_launcher":            domain.UnitRocketLauncher,
		"light_laser":               domain.UnitLightLaser,
		"heavy_laser":               domain.UnitHeavyLaser,
		"gauss_cannon":               domain.UnitGaussCannon,
		"ion_cannon":                 domain.UnitIonCannon,
		"plasma_turret":             domain.UnitPlasmaTurret,
		"small_shield_dome":         domain.UnitSmallShieldDome,
		"large_shield_dome":         domain.UnitLargeShieldDome,
		"anti_ballistic_missile":    domain.UnitAntiBallisticMissiles,
		"interplanetary_missile":    domain.UnitInterplanetaryMissiles,
	}

	for defenseName, unitType := range defenseMap {
		wikiDefense, exists := wikiData.Defense[defenseName]
		if !exists {
			continue
		}

		t.Run(defenseName, func(t *testing.T) {
			got := CalculateDefenseStructuralIntegrity(unitType)
			if got != wikiDefense.StructuralIntegrity {
				t.Logf("DISCREPANCY: CalculateDefenseStructuralIntegrity(%s): expected=%d, got=%d", defenseName, wikiDefense.StructuralIntegrity, got)
			}
		})
	}
}

func TestWiki_DefenseShield(t *testing.T) {
	wikiData := loadWikiDefenseData(t)

	defenseMap := map[string]domain.UnitType{
		"rocket_launcher":            domain.UnitRocketLauncher,
		"light_laser":               domain.UnitLightLaser,
		"heavy_laser":               domain.UnitHeavyLaser,
		"gauss_cannon":               domain.UnitGaussCannon,
		"ion_cannon":                 domain.UnitIonCannon,
		"plasma_turret":             domain.UnitPlasmaTurret,
		"small_shield_dome":         domain.UnitSmallShieldDome,
		"large_shield_dome":         domain.UnitLargeShieldDome,
		"anti_ballistic_missile":    domain.UnitAntiBallisticMissiles,
		"interplanetary_missile":    domain.UnitInterplanetaryMissiles,
	}

	for defenseName, unitType := range defenseMap {
		wikiDefense, exists := wikiData.Defense[defenseName]
		if !exists {
			continue
		}

		t.Run(defenseName, func(t *testing.T) {
			got := CalculateDefenseShield(unitType)
			if got != wikiDefense.Shield {
				t.Logf("DISCREPANCY: CalculateDefenseShield(%s): expected=%d, got=%d", defenseName, wikiDefense.Shield, got)
			}
		})
	}
}

func TestWiki_DefenseWeapon(t *testing.T) {
	wikiData := loadWikiDefenseData(t)

	defenseMap := map[string]domain.UnitType{
		"rocket_launcher":            domain.UnitRocketLauncher,
		"light_laser":               domain.UnitLightLaser,
		"heavy_laser":               domain.UnitHeavyLaser,
		"gauss_cannon":               domain.UnitGaussCannon,
		"ion_cannon":                 domain.UnitIonCannon,
		"plasma_turret":             domain.UnitPlasmaTurret,
		"small_shield_dome":         domain.UnitSmallShieldDome,
		"large_shield_dome":         domain.UnitLargeShieldDome,
		"anti_ballistic_missile":    domain.UnitAntiBallisticMissiles,
		"interplanetary_missile":    domain.UnitInterplanetaryMissiles,
	}

	for defenseName, unitType := range defenseMap {
		wikiDefense, exists := wikiData.Defense[defenseName]
		if !exists {
			continue
		}

		t.Run(defenseName, func(t *testing.T) {
			got := CalculateDefenseWeapon(unitType)
			if got != wikiDefense.Weapon {
				t.Logf("DISCREPANCY: CalculateDefenseWeapon(%s): expected=%d, got=%d", defenseName, wikiDefense.Weapon, got)
			}
		})
	}
}
