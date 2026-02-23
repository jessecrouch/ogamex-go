package formula

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type WikiResearchTestData struct {
	Research map[string]WikiResearchData `json:"research"`
}

type WikiResearchData struct {
	LabLevelRequired int            `json:"lab_level_required"`
	Prerequisites    map[string]int `json:"prerequisites"`
	BaseCost        ResearchCostData `json:"base_cost"`
}

type ResearchCostData struct {
	Metal     int `json:"metal"`
	Crystal   int `json:"crystal"`
	Deuterium int `json:"deuterium"`
}

func getWikiResearchDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiResearchData(t *testing.T) WikiResearchTestData {
	data, err := os.ReadFile(getWikiResearchDataPath("research.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki research test data: %v", err)
	}

	var wikiData WikiResearchTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki research test data: %v", err)
	}

	return wikiData
}

func TestWiki_ResearchRequirements(t *testing.T) {
	wikiData := loadWikiResearchData(t)

	researchMap := map[string]struct {
		labLevel   int
		prereqs   map[string]int
	}{
		"espionage_technology":         {3, map[string]int{}},
		"computer_technology":           {1, map[string]int{}},
		"weapons_technology":           {4, map[string]int{}},
		"shielding_technology":         {6, map[string]int{"energy_technology": 3}},
		"armor_technology":             {2, map[string]int{}},
		"energy_technology":            {1, map[string]int{}},
		"hyperspace_technology":       {7, map[string]int{"energy_technology": 5, "shielding_technology": 5}},
		"combustion_drive":             {1, map[string]int{"energy_technology": 1}},
		"impulse_drive":               {2, map[string]int{"energy_technology": 1}},
		"hyperspace_drive":            {7, map[string]int{"hyperspace_technology": 3}},
		"laser_technology":            {1, map[string]int{"energy_technology": 2}},
		"ion_technology":              {4, map[string]int{"laser_technology": 5, "energy_technology": 4}},
		"plasma_technology":           {4, map[string]int{"energy_technology": 8, "laser_technology": 10, "ion_technology": 5}},
		"intergalactic_research_network": {10, map[string]int{"computer_technology": 8, "hyperspace_technology": 8}},
		"graviton_technology":         {12, map[string]int{}},
		"astrophysics":                {3, map[string]int{"espionage_technology": 4, "impulse_drive": 3}},
	}

	for researchName, expected := range researchMap {
		t.Run(researchName, func(t *testing.T) {
			wikiResearch, exists := wikiData.Research[researchName]
			if !exists {
				t.Logf("DISCREPANCY: Research %s not found in test data", researchName)
				return
			}

			if wikiResearch.LabLevelRequired != expected.labLevel {
				t.Logf("DISCREPANCY: Research %s lab level: expected=%d, got=%d", researchName, expected.labLevel, wikiResearch.LabLevelRequired)
			}

			for prereq, level := range expected.prereqs {
				if wikiResearch.Prerequisites[prereq] != level {
					t.Logf("DISCREPANCY: Research %s prerequisite %s: expected=%d, got=%d", researchName, prereq, level, wikiResearch.Prerequisites[prereq])
				}
			}
		})
	}
}

func TestWiki_ResearchBaseCosts(t *testing.T) {
	wikiData := loadWikiResearchData(t)

	researchMap := map[string]struct {
		metal    int
		crystal  int
		deuterium int
	}{
		"espionage_technology":         {200, 1000, 200},
		"computer_technology":          {100, 400, 200},
		"weapons_technology":          {800, 200, 0},
		"shielding_technology":        {400, 600, 0},
		"armor_technology":            {400, 200, 0},
		"energy_technology":            {0, 800, 400},
		"hyperspace_technology":       {4000, 2000, 1000},
		"combustion_drive":            {400, 0, 0},
		"impulse_drive":              {4000, 2000, 600},
		"hyperspace_drive":           {10000, 6000, 4000},
		"laser_technology":            {200, 600, 0},
		"ion_technology":              {1000, 300, 0},
		"plasma_technology":           {2400, 1200, 600},
		"intergalactic_research_network": {240000, 160000, 80000},
		"graviton_technology":         {100000, 50000, 50000},
		"astrophysics":               {8000, 4000, 2000},
	}

	for researchName, expected := range researchMap {
		t.Run(researchName, func(t *testing.T) {
			wikiResearch, exists := wikiData.Research[researchName]
			if !exists {
				t.Logf("DISCREPANCY: Research %s not found in test data", researchName)
				return
			}

			if wikiResearch.BaseCost.Metal != expected.metal {
				t.Logf("DISCREPANCY: Research %s metal cost: expected=%d, got=%d", researchName, expected.metal, wikiResearch.BaseCost.Metal)
			}
			if wikiResearch.BaseCost.Crystal != expected.crystal {
				t.Logf("DISCREPANCY: Research %s crystal cost: expected=%d, got=%d", researchName, expected.crystal, wikiResearch.BaseCost.Crystal)
			}
			if wikiResearch.BaseCost.Deuterium != expected.deuterium {
				t.Logf("DISCREPANCY: Research %s deuterium cost: expected=%d, got=%d", researchName, expected.deuterium, wikiResearch.BaseCost.Deuterium)
			}
		})
	}
}
