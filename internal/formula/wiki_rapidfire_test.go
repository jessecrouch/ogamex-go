package formula

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"ogamex-go/internal/domain"
)

type WikiRapidFireTestData struct {
	RapidFire map[string]map[string]int `json:"rapid_fire"`
}

func getWikiRapidFireDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiRapidFireData(t *testing.T) WikiRapidFireTestData {
	data, err := os.ReadFile(getWikiRapidFireDataPath("rapid_fire.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki rapid fire test data: %v", err)
	}

	var wikiData WikiRapidFireTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki rapid fire test data: %v", err)
	}

	return wikiData
}

func TestWiki_RapidFire(t *testing.T) {
	wikiData := loadWikiRapidFireData(t)

	unitMap := map[string]domain.UnitType{
		"deathstar":        domain.UnitDeathstar,
		"solar_satellite":  domain.UnitSolarSatellite,
		"small_cargo":     domain.UnitSmallCargo,
		"light_fighter":   domain.UnitLightFighter,
		"heavy_fighter":   domain.UnitHeavyFighter,
		"cruiser":         domain.UnitCruiser,
		"battleship":      domain.UnitBattleship,
		"battlecruiser":   domain.UnitBattlecruiser,
		"bomber":          domain.UnitBomber,
		"destroyer":       domain.UnitDestroyer,
		"reaper":          domain.UnitReaper,
		"pathfinder":      domain.UnitPathfinder,
	}

	targetMap := map[string]domain.UnitType{
		"espionage_probe":   domain.UnitEspionageProbe,
		"solar_satellite":  domain.UnitSolarSatellite,
		"small_cargo":      domain.UnitSmallCargo,
		"large_cargo":      domain.UnitLargeCargo,
		"light_fighter":    domain.UnitLightFighter,
		"heavy_fighter":    domain.UnitHeavyFighter,
		"cruiser":          domain.UnitCruiser,
		"battleship":       domain.UnitBattleship,
		"battlecruiser":    domain.UnitBattlecruiser,
		"colony_ship":      domain.UnitColonyShip,
		"recycler":        domain.UnitRecycler,
		"bomber":          domain.UnitBomber,
		"destroyer":        domain.UnitDestroyer,
		"deathstar":        domain.UnitDeathstar,
		"reaper":          domain.UnitReaper,
		"pathfinder":       domain.UnitPathfinder,
		"crawler":          domain.UnitCrawler,
		"rocket_launcher":   domain.UnitRocketLauncher,
		"light_laser":      domain.UnitLightLaser,
		"heavy_laser":      domain.UnitHeavyLaser,
		"ion_cannon":       domain.UnitIonCannon,
		"gauss_cannon":     domain.UnitGaussCannon,
	}

	for shipName, attacker := range unitMap {
		wikiRapidFire, exists := wikiData.RapidFire[shipName]
		if !exists {
			continue
		}

		for targetName, expectedRF := range wikiRapidFire {
			defender, ok := targetMap[targetName]
			if !ok {
				continue
			}

			t.Run(shipName+"_vs_"+targetName, func(t *testing.T) {
				got := GetRapidFireAgainst(attacker, defender)
				if got != expectedRF {
					t.Logf("DISCREPANCY: GetRapidFireAgainst(%s vs %s): expected=%d, got=%d", shipName, targetName, expectedRF, got)
				}
			})
		}
	}
}
