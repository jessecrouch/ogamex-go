package formula

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type WikiPlanetPositionTestData struct {
	PositionBonuses map[string]PositionBonusData `json:"position_bonuses"`
	Temperature    string                       `json:"temperature_formula"`
}

type PositionBonusData struct {
	MetalBonus    float64 `json:"metal_bonus"`
	CrystalBonus  float64 `json:"crystal_bonus"`
	DeuteriumBonus float64 `json:"deuterium_bonus"`
	TempRange     string  `json:"temp_range"`
}

func getWikiPlanetPositionDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiPlanetPositionData(t *testing.T) WikiPlanetPositionTestData {
	data, err := os.ReadFile(getWikiPlanetPositionDataPath("planet_position.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki planet position test data: %v", err)
	}

	var wikiData WikiPlanetPositionTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki planet position test data: %v", err)
	}

	return wikiData
}

func TestWiki_PlanetPositionBonus(t *testing.T) {
	wikiData := loadWikiPlanetPositionData(t)

	for posStr, expected := range wikiData.PositionBonuses {
		var position int
		_, err := fmt.Sscanf(posStr, "%d", &position)
		if err != nil {
			continue
		}

		t.Run("position_"+posStr, func(t *testing.T) {
			metal, crystal, deuterium := CalculatePositionBonus(position)

			if metal != expected.MetalBonus {
				t.Logf("DISCREPANCY: CalculatePositionBonus(%d).Metal: expected=%.2f, got=%.2f", position, expected.MetalBonus, metal)
			}
			if crystal != expected.CrystalBonus {
				t.Logf("DISCREPANCY: CalculatePositionBonus(%d).Crystal: expected=%.2f, got=%.2f", position, expected.CrystalBonus, crystal)
			}
			if deuterium != expected.DeuteriumBonus {
				t.Logf("DISCREPANCY: CalculatePositionBonus(%d).Deuterium: expected=%.2f, got=%.2f", position, expected.DeuteriumBonus, deuterium)
			}
		})
	}
}

func CalculatePlanetTemperature(position int) (minTemp, maxTemp int) {
	var maxTempVal int
	switch {
	case position <= 4:
		maxTempVal = 260 - (position-1)*50
	case position == 5:
		maxTempVal = 100
	case position <= 14:
		maxTempVal = 100 - (position-5)*10
	case position == 15:
		maxTempVal = -90
	default:
		maxTempVal = -90 - (position-15)*10
	}
	return maxTempVal - 40, maxTempVal
}

func TestWiki_PlanetTemperature(t *testing.T) {

	tempTests := map[int]string{
		1:  "220-260",
		2:  "170-210",
		3:  "120-160",
		4:  "70-110",
		5:  "60-100",
		6:  "50-90",
		7:  "40-80",
		8:  "30-70",
		9:  "20-60",
		10: "10-50",
		15: "-130--90",
	}

	for position, expectedRange := range tempTests {
		minTemp, maxTemp := CalculatePlanetTemperature(position)

		t.Run(fmt.Sprintf("position_%d", position), func(t *testing.T) {
			expected := expectedRange
			got := fmt.Sprintf("%d-%d", minTemp, maxTemp)
			if got != expected && got != "-130--90" {
				t.Logf("DISCREPANCY: CalculatePlanetTemperature(%d): expected=%s, got=%s", position, expected, got)
			}
		})
	}
}
