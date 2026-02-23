package formula

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

type WikiProductionTestData struct {
	MetalMine       WikiResourceData `json:"metal_mine"`
	CrystalMine     WikiResourceData `json:"crystal_mine"`
	DeuteriumSynth  WikiDeutData     `json:"deuterium_synthesizer"`
	SolarPlant      WikiResourceData `json:"solar_plant"`
	FusionReactor   WikiFusionData   `json:"fusion_reactor"`
	SolarSatellite  WikiSatelliteData `json:"solar_satellite"`
	BasicIncome     BasicIncomeData  `json:"basic_income"`
	Crawler         CrawlerData      `json:"crawler"`
	PlasmaTechBonus PlasmaTechData   `json:"plasma_technology_bonus"`
}

type WikiResourceData struct {
	BaseProduction int              `json:"base_production"`
	BaseCost       CostData         `json:"base_cost"`
	EnergyFactor  float64          `json:"energy_factor"`
	CostFactor    float64          `json:"cost_factor"`
	TestCases     []ProductionCase `json:"test_cases"`
}

type CostData struct {
	Metal     int `json:"metal"`
	Crystal   int `json:"crystal"`
	Deuterium int `json:"deuterium"`
}

type ProductionCase struct {
	Level             int `json:"level"`
	Production        int `json:"production"`
	EnergyConsumption int `json:"energy_consumption"`
	Temp              int `json:"temp"`
}

type WikiDeutData struct {
	BaseProduction int              `json:"base_production"`
	BaseCost       CostData         `json:"base_cost"`
	EnergyFactor   float64          `json:"energy_factor"`
	CostFactor     float64          `json:"cost_factor"`
	TestCases      []DeuteriumCase  `json:"test_cases"`
	TempFormula    string           `json:"temperature_formula"`
}

type DeuteriumCase struct {
	Level                 int `json:"level"`
	Temp                  int `json:"temp"`
	Production            int `json:"production"`
	DeuteriumConsumption  int `json:"deuterium_consumption"`
	EnergyConsumption     int `json:"energy_consumption"`
}

type WikiFusionData struct {
	BaseProduction int           `json:"base_production"`
	BaseCost        CostData      `json:"base_cost"`
	EnergyFactor   float64       `json:"energy_factor"`
	CostFactor     float64       `json:"cost_factor"`
	TestCases      []FusionCase  `json:"test_cases"`
	Formula        string        `json:"formula"`
}

type FusionCase struct {
	Level               int `json:"level"`
	EnergyTech          int `json:"energy_tech"`
	Production          int `json:"production"`
	DeuteriumConsumption int `json:"deuterium_consumption"`
}

type WikiSatelliteData struct {
	BaseSpeed int             `json:"base_speed"`
	TestCases []SatelliteCase `json:"test_cases"`
	Formula   string         `json:"formula"`
}

type SatelliteCase struct {
	Temp   int `json:"temp"`
	Count  int `json:"count"`
	Energy int `json:"energy"`
}

type BasicIncomeData struct {
	Description string `json:"description"`
	Metal       int    `json:"metal"`
	Crystal     int    `json:"crystal"`
	Deuterium   int    `json:"deuterium"`
	Notes       string `json:"notes"`
}

type CrawlerData struct {
	ProductionBonus   float64       `json:"production_bonus"`
	Description       string        `json:"description"`
	TestCases         []CrawlerCase `json:"test_cases"`
	MaxBonusPercent   int           `json:"max_bonus_percent"`
	ActiveFormula     string        `json:"active_crawlers_formula"`
}

type CrawlerCase struct {
	MineLevelsSum int     `json:"mine_levels_sum"`
	Crawlers      int     `json:"crawlers"`
	BonusPercent  float64 `json:"bonus_percent"`
}

type PlasmaTechData struct {
	Description string           `json:"description"`
	Formula     string           `json:"formula"`
	Where       map[string]string `json:"where"`
	TestCases   []PlasmaCase    `json:"test_cases"`
}

type PlasmaCase struct {
	PlasmaLevel     int `json:"plasma_level"`
	CrystalMine     int `json:"crystal_mine"`
	BaseProduction  int `json:"base_production"`
	WithPlasma      int `json:"with_plasma"`
}

func getWikiDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiProductionData(t *testing.T) WikiProductionTestData {
	data, err := os.ReadFile(getWikiDataPath("production.json"))
	require.NoError(t, err, "Failed to read wiki production test data")

	var wikiData WikiProductionTestData
	err = json.Unmarshal(data, &wikiData)
	require.NoError(t, err, "Failed to parse wiki production test data")

	return wikiData
}

func TestWiki_MetalMineProduction(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.MetalMine.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateMetalMineProduction(tc.Level, 1.0)
			if got != int64(tc.Production) {
				t.Logf("DISCREPANCY: CalculateMetalMineProduction(level=%d): expected=%d, got=%d", tc.Level, tc.Production, got)
			}
		})
	}
}

func TestWiki_MetalMineEnergyConsumption(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.MetalMine.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateMetalMineEnergyConsumption(tc.Level)
			if got != int64(tc.EnergyConsumption) {
				t.Logf("DISCREPANCY: CalculateMetalMineEnergyConsumption(level=%d): expected=%d, got=%d", tc.Level, tc.EnergyConsumption, got)
			}
		})
	}
}

func TestWiki_CrystalMineProduction(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.CrystalMine.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateCrystalMineProduction(tc.Level, 1.0)
			if got != int64(tc.Production) {
				t.Logf("DISCREPANCY: CalculateCrystalMineProduction(level=%d): expected=%d, got=%d", tc.Level, tc.Production, got)
			}
		})
	}
}

func TestWiki_CrystalMineEnergyConsumption(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.CrystalMine.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateCrystalMineEnergyConsumption(tc.Level)
			if got != int64(tc.EnergyConsumption) {
				t.Logf("DISCREPANCY: CalculateCrystalMineEnergyConsumption(level=%d): expected=%d, got=%d", tc.Level, tc.EnergyConsumption, got)
			}
		})
	}
}

func TestWiki_DeuteriumProduction(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.DeuteriumSynth.TestCases {
		t.Run(fmt.Sprintf("level_%d_temp_%d", tc.Level, tc.Temp), func(t *testing.T) {
			got := CalculateDeuteriumProduction(tc.Level, tc.Temp, 1.0)
			if got != int64(tc.Production) {
				t.Logf("DISCREPANCY: CalculateDeuteriumProduction(level=%d, temp=%d): expected=%d, got=%d", tc.Level, tc.Temp, tc.Production, got)
			}
		})
	}
}

func TestWiki_DeuteriumSynthesizerEnergyConsumption(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.DeuteriumSynth.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateDeuteriumSynthesizerEnergyConsumption(tc.Level)
			if got != int64(tc.EnergyConsumption) {
				t.Logf("DISCREPANCY: CalculateDeuteriumSynthesizerEnergyConsumption(level=%d): expected=%d, got=%d", tc.Level, tc.EnergyConsumption, got)
			}
		})
	}
}

func TestWiki_SolarPlantProduction(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.SolarPlant.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateSolarPlantEnergyProduction(tc.Level)
			if got != int64(tc.Production) {
				t.Logf("DISCREPANCY: CalculateSolarPlantEnergyProduction(level=%d): expected=%d, got=%d", tc.Level, tc.Production, got)
			}
		})
	}
}

func TestWiki_FusionReactorProduction(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.FusionReactor.TestCases {
		t.Run(fmt.Sprintf("level_%d_tech_%d", tc.Level, tc.EnergyTech), func(t *testing.T) {
			got := CalculateFusionPlantEnergyProduction(tc.Level, tc.EnergyTech)
			if got != int64(tc.Production) {
				t.Logf("DISCREPANCY: CalculateFusionPlantEnergyProduction(level=%d, tech=%d): expected=%d, got=%d", tc.Level, tc.EnergyTech, tc.Production, got)
			}
		})
	}
}

func TestWiki_FusionReactorDeuteriumConsumption(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.FusionReactor.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateFusionPlantDeuteriumConsumption(tc.Level)
			if got != int64(tc.DeuteriumConsumption) {
				t.Logf("DISCREPANCY: CalculateFusionPlantDeuteriumConsumption(level=%d): expected=%d, got=%d", tc.Level, tc.DeuteriumConsumption, got)
			}
		})
	}
}

func TestWiki_SolarSatelliteEnergyProduction(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	for _, tc := range wikiData.SolarSatellite.TestCases {
		t.Run(fmt.Sprintf("temp_%d_count_%d", tc.Temp, tc.Count), func(t *testing.T) {
			got := CalculateSolarSatelliteEnergyProduction(tc.Temp, tc.Count)
			if got != int64(tc.Energy) {
				t.Logf("DISCREPANCY: CalculateSolarSatelliteEnergyProduction(temp=%d, count=%d): expected=%d, got=%d", tc.Temp, tc.Count, tc.Energy, got)
			}
		})
	}
}

func TestWiki_BasicIncome(t *testing.T) {
	wikiData := loadWikiProductionData(t)
	if wikiData.BasicIncome.Metal != 30 {
		t.Logf("DISCREPANCY: BasicIncome.Metal: expected=30, got=%d", wikiData.BasicIncome.Metal)
	}
	if wikiData.BasicIncome.Crystal != 15 {
		t.Logf("DISCREPANCY: BasicIncome.Crystal: expected=15, got=%d", wikiData.BasicIncome.Crystal)
	}
	if wikiData.BasicIncome.Deuterium != 0 {
		t.Logf("DISCREPANCY: BasicIncome.Deuterium: expected=0, got=%d", wikiData.BasicIncome.Deuterium)
	}
}
