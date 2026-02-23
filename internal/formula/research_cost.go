package formula

import (
	"ogamex-go/internal/domain"
)

type ResearchBaseCost struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
}

var researchBaseCosts = map[domain.ResearchType]ResearchBaseCost{
	domain.ResearchEnergyTechnology:         {Metal: 0, Crystal: 800, Deuterium: 400},
	domain.ResearchLaserTechnology:          {Metal: 200, Crystal: 600, Deuterium: 0},
	domain.ResearchIonTechnology:            {Metal: 1000, Crystal: 300, Deuterium: 0},
	domain.ResearchHyperspaceTechnology:    {Metal: 4000, Crystal: 2000, Deuterium: 1000},
	domain.ResearchPlasmaTechnology:        {Metal: 2400, Crystal: 1200, Deuterium: 600},
	domain.ResearchFusionDrive:              {Metal: 9000, Crystal: 4000, Deuterium: 6000},
	domain.ResearchImpulseDrive:            {Metal: 4000, Crystal: 2000, Deuterium: 600},
	domain.ResearchHyperspaceDrive:          {Metal: 10000, Crystal: 6000, Deuterium: 4000},
	domain.ResearchEspionageTechnology:      {Metal: 200, Crystal: 1000, Deuterium: 200},
	domain.ResearchComputerTechnology:       {Metal: 100, Crystal: 400, Deuterium: 200},
	domain.ResearchAstrophysics:             {Metal: 8000, Crystal: 4000, Deuterium: 2000},
	domain.ResearchIntergalacticResearchNetwork: {Metal: 240000, Crystal: 160000, Deuterium: 80000},
	domain.ResearchGravitonTechnology:       {Metal: 100000, Crystal: 50000, Deuterium: 50000},
	domain.ResearchWeaponsTechnology:        {Metal: 800, Crystal: 200, Deuterium: 0},
	domain.ResearchShieldingTechnology:     {Metal: 400, Crystal: 600, Deuterium: 0},
	domain.ResearchArmorTechnology:         {Metal: 400, Crystal: 200, Deuterium: 0},
}

func CalculateResearchCost(researchType domain.ResearchType, level int) (metal, crystal, deuterium int64) {
	baseCost, exists := researchBaseCosts[researchType]
	if !exists {
		return 0, 0, 0
	}

	if level <= 0 {
		return 0, 0, 0
	}

	factor := GetCostFactor(level, 1.8)

	metal = int64(float64(baseCost.Metal) * factor)
	crystal = int64(float64(baseCost.Crystal) * factor)
	deuterium = int64(float64(baseCost.Deuterium) * factor)

	return metal, crystal, deuterium
}
