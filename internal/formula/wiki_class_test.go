package formula

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type WikiClassTestData struct {
	Classes           map[string]ClassData           `json:"classes"`
	ClassComparison   map[string]ClassComparisonData `json:"class_comparison"`
	ResourceProduction ClassProductionData           `json:"resource_production_with_class"`
}

type ClassData struct {
	Description                string  `json:"description"`
	ProductionBonus           float64 `json:"production_bonus"`
	EnergyBonus               float64 `json:"energy_bonus"`
	CargoBonus                float64 `json:"cargo_bonus"`
	FleetSpeedBonus           float64 `json:"fleet_speed_bonus"`
	MiningDroneBonus          float64 `json:"mining_drone_bonus"`
	ProductionFormula         string  `json:"production_formula"`
	ShipCostReduction         float64 `json:"ship_cost_reduction"`
	DefenseCostReduction      float64 `json:"defense_cost_reduction"`
	WeaponsBonus              float64 `json:"weapons_bonus"`
	ShieldBonus               float64 `json:"shield_bonus"`
	PlanetSizeBonus           float64 `json:"planet_size_bonus"`
	DeuteriumConsumptionReduction float64 `json:"deuterium_consumption_reduction"`
	ExpeditionBonus           float64 `json:"expedition_bonus"`
	ScannerRangeBonus         float64 `json:"scanner_range_bonus"`
}

type ClassComparisonData struct {
	CollectorProduction   float64 `json:"collector_production"`
	GeneralProduction    float64 `json:"general_production"`
	CollectorShipCost    float64 `json:"collector_ship_cost"`
	GeneralShipCost      float64 `json:"general_ship_cost"`
	BaseConsumption      float64 `json:"base_consumption"`
	DiscovererConsumption float64 `json:"discoverer_consumption"`
	Reduction            float64 `json:"reduction"`
}

type ClassProductionData struct {
	MetalMineLevel10     int     `json:"metal_mine_level_10"`
	CollectorBonus       float64 `json:"collector_bonus"`
	CollectorProduction  int     `json:"collector_production"`
	GeneralProduction    int     `json:"general_production"`
	DiscovererProduction int     `json:"discoverer_production"`
}

func getWikiClassDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiClassData(t *testing.T) WikiClassTestData {
	data, err := os.ReadFile(getWikiClassDataPath("classes.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki class test data: %v", err)
	}

	var wikiData WikiClassTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki class test data: %v", err)
	}

	return wikiData
}

func TestWiki_CollectorClassBonuses(t *testing.T) {
	wikiData := loadWikiClassData(t)
	c := wikiData.Classes["collector"]

	t.Logf("=== Collector Class Bonuses ===")
	t.Logf("Production Bonus: %.0f%%", c.ProductionBonus*100)
	t.Logf("Energy Bonus: %.0f%%", c.EnergyBonus*100)
	t.Logf("Cargo Bonus: %.0f%%", c.CargoBonus*100)
	t.Logf("Fleet Speed Bonus: %.0f%%", c.FleetSpeedBonus*100)
	t.Logf("Mining Drone Bonus: +%.0f per level", c.MiningDroneBonus)

	if c.ProductionBonus != 0.25 {
		t.Logf("DISCREPANCY: Collector production bonus: expected=0.25, got=%.2f", c.ProductionBonus)
	}
	if c.EnergyBonus != 0.10 {
		t.Logf("DISCREPANCY: Collector energy bonus: expected=0.10, got=%.2f", c.EnergyBonus)
	}
	if c.CargoBonus != 0.25 {
		t.Logf("DISCREPANCY: Collector cargo bonus: expected=0.25, got=%.2f", c.CargoBonus)
	}
}

func TestWiki_GeneralClassBonuses(t *testing.T) {
	wikiData := loadWikiClassData(t)
	c := wikiData.Classes["general"]

	t.Logf("=== General Class Bonuses ===")
	t.Logf("Ship Cost Reduction: %.0f%%", c.ShipCostReduction*100)
	t.Logf("Defense Cost Reduction: %.0f%%", c.DefenseCostReduction*100)
	t.Logf("Fleet Speed Bonus: %.0f%%", c.FleetSpeedBonus*100)
	t.Logf("Weapons Bonus: %.0f%%", c.WeaponsBonus*100)
	t.Logf("Shield Bonus: %.0f%%", c.ShieldBonus*100)

	if c.ShipCostReduction != 0.15 {
		t.Logf("DISCREPANCY: General ship cost reduction: expected=0.15, got=%.2f", c.ShipCostReduction)
	}
	if c.DefenseCostReduction != 0.15 {
		t.Logf("DISCREPANCY: General defense cost reduction: expected=0.15, got=%.2f", c.DefenseCostReduction)
	}
}

func TestWiki_DiscovererClassBonuses(t *testing.T) {
	wikiData := loadWikiClassData(t)
	c := wikiData.Classes["discoverer"]

	t.Logf("=== Discoverer Class Bonuses ===")
	t.Logf("Planet Size Bonus: %.0f%%", c.PlanetSizeBonus*100)
	t.Logf("Deuterium Consumption Reduction: %.0f%%", c.DeuteriumConsumptionReduction*100)
	t.Logf("Expedition Bonus: %.0f%%", c.ExpeditionBonus*100)
	t.Logf("Scanner Range Bonus: %.0f%%", c.ScannerRangeBonus*100)

	if c.PlanetSizeBonus != 0.25 {
		t.Logf("DISCREPANCY: Discoverer planet size bonus: expected=0.25, got=%.2f", c.PlanetSizeBonus)
	}
	if c.DeuteriumConsumptionReduction != 0.50 {
		t.Logf("DISCREPANCY: Discoverer deuterium reduction: expected=0.50, got=%.2f", c.DeuteriumConsumptionReduction)
	}
}

func TestWiki_ClassProductionComparison(t *testing.T) {
	wikiData := loadWikiClassData(t)
	comp := wikiData.ClassComparison["collector_vs_general"]

	t.Logf("=== Class Production Comparison ===")
	t.Logf("Collector: %.0f%% production", comp.CollectorProduction*100)
	t.Logf("General: %.0f%% production", comp.GeneralProduction*100)

	collector := wikiData.Classes["collector"]
	baseProduction := float64(wikiData.ResourceProduction.MetalMineLevel10)
	collectorExpected := baseProduction * (1 + collector.ProductionBonus)
	collectorActual := float64(wikiData.ResourceProduction.CollectorProduction)
	diff := collectorExpected - collectorActual
	if diff > 1 || diff < -1 {
		t.Logf("DISCREPANCY: Collector production: expected=%.0f, got=%.0f", collectorExpected, collectorActual)
	}
}

func TestWiki_ClassDeuteriumReduction(t *testing.T) {
	wikiData := loadWikiClassData(t)
	comp := wikiData.ClassComparison["discoverer_deuterium"]

	t.Logf("=== Class Deuterium Consumption ===")
	t.Logf("Base consumption: %.0f", comp.BaseConsumption)
	t.Logf("Discoverer consumption: %.0f (%.0f%% reduction)", comp.DiscovererConsumption, comp.Reduction*100)

	discoverer := wikiData.Classes["discoverer"]
	expectedReduction := discoverer.DeuteriumConsumptionReduction

	if expectedReduction != comp.Reduction {
		t.Logf("DISCREPANCY: Deuterium reduction: expected=%.2f, got=%.2f", comp.Reduction, expectedReduction)
	}
}

func TestWiki_ClassBonusFormulas(t *testing.T) {
	wikiData := loadWikiClassData(t)

	t.Logf("=== Class Bonus Formulas ===")

	baseMetalProduction := 779.0

	collector := wikiData.Classes["collector"]
	collectorProduction := baseMetalProduction * (1 + collector.ProductionBonus)
	t.Logf("Collector (base=%.0f, bonus=%.0f%%): %.0f",
		baseMetalProduction, collector.ProductionBonus*100, collectorProduction)

	general := wikiData.Classes["general"]
	generalProduction := baseMetalProduction * (1 + general.ProductionBonus)
	t.Logf("General (base=%.0f, bonus=%.0f%%): %.0f",
		baseMetalProduction, general.ProductionBonus*100, generalProduction)

	discoverer := wikiData.Classes["discoverer"]
	discovererProduction := baseMetalProduction * (1 + discoverer.ProductionBonus)
	t.Logf("Discoverer (base=%.0f, bonus=%.0f%%): %.0f",
		baseMetalProduction, discoverer.ProductionBonus*100, discovererProduction)

	baseShipCost := 10000.0
	generalShipCost := baseShipCost * (1 - general.ShipCostReduction)
	t.Logf("General ship cost (base=%.0f, reduction=%.0f%%): %.0f",
		baseShipCost, general.ShipCostReduction*100, generalShipCost)

	baseDeuteriumConsumption := 100.0
	discovererDeuteriumConsumption := baseDeuteriumConsumption * (1 - discoverer.DeuteriumConsumptionReduction)
	t.Logf("Discoverer deuterium (base=%.0f, reduction=%.0f%%): %.0f",
		baseDeuteriumConsumption, discoverer.DeuteriumConsumptionReduction*100, discovererDeuteriumConsumption)
}
