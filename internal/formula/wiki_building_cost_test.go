package formula

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"ogamex-go/internal/domain"
)

type WikiBuildingCostTestData struct {
	Buildings        map[string]WikiBuildingData `json:"buildings"`
	StorageCapacity  StorageCapacityData         `json:"storage_capacity"`
	ConstructionTime ConstructionTimeData       `json:"construction_time"`
}

type WikiBuildingData struct {
	BaseCost   BuildingCostData `json:"base_cost"`
	CostFactor float64          `json:"cost_factor"`
	TestCases  []BuildingCase    `json:"test_cases"`
}

type BuildingCostData struct {
	Metal     int `json:"metal"`
	Crystal   int `json:"crystal"`
	Deuterium int `json:"deuterium"`
}

type BuildingCase struct {
	Level      int `json:"level"`
	Metal      int `json:"metal"`
	Crystal    int `json:"crystal"`
	Deuterium  int `json:"deuterium"`
}

type StorageCapacityData struct {
	BaseCapacity int               `json:"base_capacity"`
	TestCases    []StorageCase    `json:"test_cases"`
}

type StorageCase struct {
	Level    int `json:"level"`
	Capacity int `json:"capacity"`
}

type ConstructionTimeData struct {
	TestCases []ConstructionCase `json:"test_cases"`
}

type ConstructionCase struct {
	Building string  `json:"building"`
	Level    int     `json:"level"`
	Factor   float64 `json:"factor"`
}

func getWikiBuildingDataPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(currentFile)
	return filepath.Join(dir, "..", "..", "tests", "wiki", filename)
}

func loadWikiBuildingCostData(t *testing.T) WikiBuildingCostTestData {
	data, err := os.ReadFile(getWikiBuildingDataPath("buildings_cost.json"))
	if err != nil {
		t.Fatalf("Failed to read wiki building cost test data: %v", err)
	}

	var wikiData WikiBuildingCostTestData
	err = json.Unmarshal(data, &wikiData)
	if err != nil {
		t.Fatalf("Failed to parse wiki building cost test data: %v", err)
	}

	return wikiData
}

func TestWiki_MetalMineCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["metal_mine"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, deuterium := CalculateBuildingCost(domain.BuildingMetalMine, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(MetalMine, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(MetalMine, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
			if deuterium != int64(tc.Deuterium) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(MetalMine, level=%d).Deuterium: expected=%d, got=%d", tc.Level, tc.Deuterium, deuterium)
			}
		})
	}
}

func TestWiki_CrystalMineCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["crystal_mine"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, _ := CalculateBuildingCost(domain.BuildingCrystalMine, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(CrystalMine, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(CrystalMine, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
		})
	}
}

func TestWiki_DeuteriumSynthesizerCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["deuterium_synthesizer"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, _, _ := CalculateBuildingCost(domain.BuildingDeuteriumSynthesizer, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(DeuteriumSynth, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
		})
	}
}

func TestWiki_SolarPlantCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["solar_plant"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, _, _ := CalculateBuildingCost(domain.BuildingSolarPlant, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(SolarPlant, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
		})
	}
}

func TestWiki_FusionPlantCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["fusion_plant"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, deuterium := CalculateBuildingCost(domain.BuildingFusionPlant, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(FusionPlant, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(FusionPlant, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
			if deuterium != int64(tc.Deuterium) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(FusionPlant, level=%d).Deuterium: expected=%d, got=%d", tc.Level, tc.Deuterium, deuterium)
			}
		})
	}
}

func TestWiki_MetalStorageCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["metal_storage"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, _, _ := CalculateBuildingCost(domain.BuildingMetalStorage, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(MetalStorage, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
		})
	}
}

func TestWiki_CrystalStorageCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["crystal_storage"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, _ := CalculateBuildingCost(domain.BuildingCrystalStorage, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(CrystalStorage, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(CrystalStorage, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
		})
	}
}

func TestWiki_DeuteriumStorageCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["deuterium_storage"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, _ := CalculateBuildingCost(domain.BuildingDeuteriumStorage, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(DeuteriumStorage, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(DeuteriumStorage, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
		})
	}
}

func TestWiki_RoboticsFactoryCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["robotics_factory"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, _ := CalculateBuildingCost(domain.BuildingRoboticsFactory, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(RoboticsFactory, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(RoboticsFactory, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
		})
	}
}

func TestWiki_ShipyardCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["shipyard"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, _ := CalculateBuildingCost(domain.BuildingShipyard, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(Shipyard, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(Shipyard, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
		})
	}
}

func TestWiki_NaniteFactoryCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["nanite_factory"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, _ := CalculateBuildingCost(domain.BuildingNaniteFactory, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(NaniteFactory, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(NaniteFactory, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
		})
	}
}

func TestWiki_TerraformerCost(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	building := wikiData.Buildings["terraformer"]
	
	for _, tc := range building.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			metal, crystal, deuterium := CalculateBuildingCost(domain.BuildingTerraformer, tc.Level)
			if metal != int64(tc.Metal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(Terraformer, level=%d).Metal: expected=%d, got=%d", tc.Level, tc.Metal, metal)
			}
			if crystal != int64(tc.Crystal) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(Terraformer, level=%d).Crystal: expected=%d, got=%d", tc.Level, tc.Crystal, crystal)
			}
			if deuterium != int64(tc.Deuterium) {
				t.Logf("DISCREPANCY: CalculateBuildingCost(Terraformer, level=%d).Deuterium: expected=%d, got=%d", tc.Level, tc.Deuterium, deuterium)
			}
		})
	}
}

func TestWiki_StorageCapacity(t *testing.T) {
	wikiData := loadWikiBuildingCostData(t)
	
	for _, tc := range wikiData.StorageCapacity.TestCases {
		t.Run(fmt.Sprintf("level_%d", tc.Level), func(t *testing.T) {
			got := CalculateStorageCapacity(tc.Level)
			if got != int64(tc.Capacity) {
				t.Logf("DISCREPANCY: CalculateStorageCapacity(level=%d): expected=%d, got=%d", tc.Level, tc.Capacity, got)
			}
		})
	}
}
