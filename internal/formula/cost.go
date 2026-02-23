package formula

import (
	"math"

	"ogamex-go/internal/domain"
)

func GetCostFactor(level int, baseFactor float64) float64 {
	if level <= 1 {
		return 1.0
	}
	return math.Pow(baseFactor, float64(level-1))
}

type BuildingBaseCost struct {
	Metal     int
	Crystal   int
	Deuterium int
	CostFactor float64
}

var buildingBaseCosts = map[domain.BuildingType]BuildingBaseCost{
	domain.BuildingMetalMine:            {Metal: 60, Crystal: 0, Deuterium: 0, CostFactor: 1.5},
	domain.BuildingCrystalMine:          {Metal: 48, Crystal: 24, Deuterium: 0, CostFactor: 1.6},
	domain.BuildingDeuteriumSynthesizer: {Metal: 225, Crystal: 0, Deuterium: 0, CostFactor: 1.5},
	domain.BuildingSolarPlant:           {Metal: 75, Crystal: 0, Deuterium: 0, CostFactor: 1.5},
	domain.BuildingFusionPlant:          {Metal: 900, Crystal: 360, Deuterium: 180, CostFactor: 1.8},
	domain.BuildingMetalStorage:         {Metal: 100, Crystal: 0, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingCrystalStorage:       {Metal: 100, Crystal: 50, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingDeuteriumStorage:     {Metal: 100, Crystal: 100, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingRoboticsFactory:      {Metal: 400, Crystal: 120, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingShipyard:             {Metal: 400, Crystal: 200, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingResearchLab:          {Metal: 200, Crystal: 400, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingNaniteFactory:       {Metal: 1000000, Crystal: 200000, Deuterium: 0, CostFactor: 2.0},
	domain.BuildingTerraformer:         {Metal: 50000, Crystal: 100000, Deuterium: 1000000, CostFactor: 2.0},
	domain.BuildingCrawler:              {Metal: 400, Crystal: 0, Deuterium: 100, CostFactor: 1.5},
	domain.BuildingSpaceDock:           {Metal: 400, Crystal: 200, Deuterium: 100, CostFactor: 2.0},
	domain.BuildingMissileSilo:          {Metal: 20000, Crystal: 20000, Deuterium: 1000, CostFactor: 2.0},
	domain.BuildingMetalMine2:          {Metal: 60, Crystal: 30, Deuterium: 0, CostFactor: 1.5},
	domain.BuildingCrystalMine2:        {Metal: 48, Crystal: 24, Deuterium: 0, CostFactor: 1.6},
	domain.BuildingDeuteriumSynthesizer2: {Metal: 225, Crystal: 0, Deuterium: 0, CostFactor: 1.5},
	domain.BuildingLunarBase:           {Metal: 20000, Crystal: 40000, Deuterium: 20000, CostFactor: 2.0},
	domain.BuildingSensorPhalanx:       {Metal: 20000, Crystal: 40000, Deuterium: 20000, CostFactor: 2.0},
	domain.BuildingJumpGate:            {Metal: 2000000, Crystal: 4000000, Deuterium: 2000000, CostFactor: 2.0},
}

func CalculateBuildingCost(buildingType domain.BuildingType, level int) (metal, crystal, deuterium int64) {
	baseCost, exists := buildingBaseCosts[buildingType]
	if !exists {
		return 0, 0, 0
	}

	if level <= 0 {
		return 0, 0, 0
	}

	factor := GetCostFactor(level, baseCost.CostFactor)

	metal = int64(math.Round(float64(baseCost.Metal) * factor))
	crystal = int64(math.Round(float64(baseCost.Crystal) * factor))
	deuterium = int64(math.Round(float64(baseCost.Deuterium) * factor))

	return metal, crystal, deuterium
}

func CalculateStorageCapacity(level int) int64 {
	if level <= 0 {
		return 10000
	}
	capacity := 10000.0 * math.Pow(2.0, float64(level))
	return int64(math.Round(capacity))
}
