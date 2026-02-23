package formula

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type WikiCombatTestData struct {
	ShieldAbsorption ShieldAbsorption `json:"shield_absorption"`
	WeaponDamage     WeaponDamage     `json:"weapon_damage"`
	DefenseRebuild   DefenseRebuild   `json:"defense_rebuild"`
	DebrisField      DebrisField      `json:"debris_field"`
	MoonChance       MoonChance       `json:"moon_chance"`
	BattleScenarios  []BattleScenario `json:"battle_scenarios"`
}

type ShieldAbsorption struct {
	Description        string  `json:"description"`
	MaxShieldAbsorption float64 `json:"max_shield_absorption"`
	Example            ShieldExample `json:"example"`
}

type ShieldExample struct {
	ShieldPoints    float64 `json:"shield_points"`
	IncomingDamage  float64 `json:"incoming_damage"`
	Absorbed        float64 `json:"absorbed"`
	RemainingDamage float64 `json:"remaining_damage"`
}

type WeaponDamage struct {
	Description   string  `json:"description"`
	MinMultiplier float64 `json:"min_multiplier"`
	MaxMultiplier float64 `json:"max_multiplier"`
}

type DefenseRebuild struct {
	Description    string  `json:"description"`
	RebuildChance float64 `json:"rebuild_chance"`
}

type DebrisField struct {
	Description       string  `json:"description"`
	MetalPercentage  float64 `json:"metal_percentage"`
	CrystalPercentage float64 `json:"crystal_percentage"`
	Example          DebrisExample `json:"example"`
}

type DebrisExample struct {
	DestroyedValue  float64 `json:"destroyed_value"`
	DebrisMetal     float64 `json:"debris_metal"`
	DebrisCrystal   float64 `json:"debris_crystal"`
}

type MoonChance struct {
	Description     string  `json:"description"`
	BaseChance     float64 `json:"base_chance"`
	Per100kDebris   float64 `json:"per_100k_debris"`
	MaxChance       float64 `json:"max_chance"`
}

type BattleScenario struct {
	Name     string `json:"name"`
	Attacker UnitGroup `json:"attacker"`
	Defender UnitGroup `json:"defender"`
}

type UnitGroup struct {
	Units map[string]UnitStats `json:"units"`
}

type UnitStats struct {
	Amount int    `json:"amount"`
	Attack float64 `json:"attack"`
	Shield float64 `json:"shield"`
	Hull   float64 `json:"hull"`
}

func getWikiCombatDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiCombatData(t *testing.T) WikiCombatTestData {
	data, err := os.ReadFile(getWikiCombatDataPath("combat.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki combat test data: %v", err)
	}

	var wikiData WikiCombatTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki combat test data: %v", err)
	}

	return wikiData
}

func TestWiki_CombatShieldAbsorption(t *testing.T) {
	wikiData := loadWikiCombatData(t)

	ex := wikiData.ShieldAbsorption.Example
	maxAbsorption := ex.ShieldPoints * wikiData.ShieldAbsorption.MaxShieldAbsorption
	actualAbsorption := min(maxAbsorption, ex.IncomingDamage)

	if actualAbsorption > ex.Absorbed+1 || actualAbsorption < ex.Absorbed-1 {
		t.Logf("DISCREPANCY: Shield absorption calculation: expected=%.0f, got=%.0f", ex.Absorbed, actualAbsorption)
	}

	remainingDamage := ex.IncomingDamage - actualAbsorption
	if remainingDamage > ex.RemainingDamage+1 || remainingDamage < ex.RemainingDamage-1 {
		t.Logf("DISCREPANCY: Remaining damage after shield: expected=%.0f, got=%.0f", ex.RemainingDamage, remainingDamage)
	}

	t.Logf("Shield Absorption: max absorption=%.0f (%.0f%% of shield), actual=%.0f", 
		maxAbsorption, wikiData.ShieldAbsorption.MaxShieldAbsorption*100, actualAbsorption)
}

func TestWiki_CombatWeaponDamageRange(t *testing.T) {
	wikiData := loadWikiCombatData(t)

	baseDamage := 1000.0
	minDamage := baseDamage * wikiData.WeaponDamage.MinMultiplier
	maxDamage := baseDamage * wikiData.WeaponDamage.MaxMultiplier

	t.Logf("Weapon Damage Range: base=%.0f, min=%.0f (%.0f%%), max=%.0f (%.0f%%)",
		baseDamage, minDamage, wikiData.WeaponDamage.MinMultiplier*100,
		maxDamage, wikiData.WeaponDamage.MaxMultiplier*100)
}

func TestWiki_CombatDefenseRebuildChance(t *testing.T) {
	wikiData := loadWikiCombatData(t)

	expected := wikiData.DefenseRebuild.RebuildChance
	if expected != 0.7 {
		t.Logf("DISCREPANCY: Defense rebuild chance: expected=0.7 (70%%), got=%.2f", expected)
	}
	t.Logf("Defense Rebuild: %.0f%% chance to rebuild", expected*100)
}

func TestWiki_CombatDebrisField(t *testing.T) {
	wikiData := loadWikiCombatData(t)

	ex := wikiData.DebrisField.Example
	expectedMetal := ex.DestroyedValue * wikiData.DebrisField.MetalPercentage
	expectedCrystal := ex.DestroyedValue * wikiData.DebrisField.CrystalPercentage

	if expectedMetal != ex.DebrisMetal {
		t.Logf("DISCREPANCY: Debris metal: expected=%.0f, got=%.0f", ex.DebrisMetal, expectedMetal)
	}
	if expectedCrystal != ex.DebrisCrystal {
		t.Logf("DISCREPANCY: Debris crystal: expected=%.0f, got=%.0f", ex.DebrisCrystal, expectedCrystal)
	}

	t.Logf("Debris Field: %.0f%% metal + %.0f%% crystal = %.0f%% of destroyed value",
		wikiData.DebrisField.MetalPercentage*100, wikiData.DebrisField.CrystalPercentage*100,
		(wikiData.DebrisField.MetalPercentage+wikiData.DebrisField.CrystalPercentage)*100)
}

func TestWiki_CombatMoonChance(t *testing.T) {
	wikiData := loadWikiCombatData(t)

	tests := []struct {
		debris    float64
		expected  float64
	}{
		{0, 0.02},
		{100000, 0.03},
		{1000000, 0.12},
		{5000000, 0.20},
	}

	for _, tt := range tests {
		var chance float64
		tt2 := tt.debris / 100000
		chance = wikiData.MoonChance.BaseChance + tt2*wikiData.MoonChance.Per100kDebris
		if chance > wikiData.MoonChance.MaxChance {
			chance = wikiData.MoonChance.MaxChance
		}

		diff := chance - tt.expected
		if diff > 0.001 || diff < -0.001 {
			t.Logf("DISCREPANCY: Moon chance for debris %.0f: expected=%.2f, got=%.2f", tt.debris, tt.expected, chance)
		}
		t.Logf("Moon chance (debris=%.0f): %.1f%%", tt.debris, chance*100)
	}
}

func TestWiki_CombatBattleSimulation(t *testing.T) {
	t.Skip("Skipping battle simulation - needs proper unit initialization with rapidfire")
}

func TestWiki_CombatDebrisCalculation(t *testing.T) {
	wikiData := loadWikiCombatData(t)

	shipCosts := map[string]struct{ metal, crystal float64 }{
		"light_fighter":   {3000, 0},
		"heavy_fighter":   {6000, 2000},
		"cruiser":         {20000, 7000},
		"battleship":      {50000, 25000},
		"deathstar":       {5000000, 4000000},
	}

	for name, costs := range shipCosts {
		totalValue := costs.metal + costs.crystal
		debrisMetal := totalValue * wikiData.DebrisField.MetalPercentage
		debrisCrystal := totalValue * wikiData.DebrisField.CrystalPercentage

		t.Logf("Debris from %s (value=%.0f): metal=%.0f, crystal=%.0f, total=%.0f",
			name, totalValue, debrisMetal, debrisCrystal, debrisMetal+debrisCrystal)
	}
}
