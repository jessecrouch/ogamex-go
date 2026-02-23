package formula

import (
	"math"

	"ogamex-go/internal/domain"
)

func CalculateMetalMineProduction(level int, energyFactor float64) int64 {
	if level <= 0 {
		return 0
	}
	base := 30 * float64(level) * math.Pow(1.1, float64(level))
	production := base * energyFactor
	return int64(math.Ceil(production))
}

func CalculateMetalMineEnergyConsumption(level int) int64 {
	if level <= 0 {
		return 0
	}
	consumption := 10 * float64(level) * math.Pow(1.1, float64(level))
	return -int64(math.Ceil(consumption))
}

func CalculateCrystalMineProduction(level int, energyFactor float64) int64 {
	if level <= 0 {
		return 0
	}
	base := 20 * float64(level) * math.Pow(1.1, float64(level))
	production := base * energyFactor
	return int64(math.Ceil(production))
}

func CalculateCrystalMineEnergyConsumption(level int) int64 {
	if level <= 0 {
		return 0
	}
	consumption := 10 * float64(level) * math.Pow(1.1, float64(level))
	return -int64(math.Ceil(consumption))
}

func CalculateDeuteriumProduction(level int, temperature int, energyFactor float64) int64 {
	if level <= 0 {
		return 0
	}
	tempFactor := 1.44 - 0.004*float64(temperature)
	base := 10 * float64(level) * math.Pow(1.1, float64(level)) * tempFactor
	production := base * energyFactor
	return int64(math.Ceil(production))
}

func CalculateDeuteriumSynthesizerEnergyConsumption(level int) int64 {
	if level <= 0 {
		return 0
	}
	consumption := 20 * float64(level) * math.Pow(1.1, float64(level))
	return -int64(math.Ceil(consumption))
}

func CalculateSolarPlantEnergyProduction(level int) int64 {
	if level <= 0 {
		return 0
	}
	production := 20 * float64(level) * math.Pow(1.1, float64(level))
	return int64(math.Ceil(production))
}

func CalculateFusionPlantEnergyProduction(level int, energyTechLevel int) int64 {
	if level <= 0 {
		return 0
	}
	factor := 1.05 + float64(energyTechLevel)*0.01
	production := 30 * float64(level) * math.Pow(factor, float64(level))
	return int64(math.Ceil(production))
}

func CalculateFusionPlantDeuteriumConsumption(level int) int64 {
	if level <= 0 {
		return 0
	}
	consumption := 10 * float64(level) * math.Pow(1.1, float64(level))
	return -int64(math.Ceil(consumption))
}

func CalculateSolarSatelliteEnergyProduction(temperature int, count int) int64 {
	if count <= 0 {
		return 0
	}
	production := math.Floor(float64(temperature+140)/6.0) * float64(count)
	return int64(production)
}

func GetProduction(building domain.BuildingType, level int, temp int, energyFactor float64, energyTechLevel int) domain.Resources {
	var resources domain.Resources

	switch building {
	case domain.BuildingMetalMine:
		resources.Metal = CalculateMetalMineProduction(level, energyFactor)
		resources.Energy = CalculateMetalMineEnergyConsumption(level)
	case domain.BuildingCrystalMine:
		resources.Crystal = CalculateCrystalMineProduction(level, energyFactor)
		resources.Energy = CalculateCrystalMineEnergyConsumption(level)
	case domain.BuildingDeuteriumSynthesizer:
		resources.Deuterium = CalculateDeuteriumProduction(level, temp, energyFactor)
		resources.Energy = CalculateDeuteriumSynthesizerEnergyConsumption(level)
	case domain.BuildingSolarPlant:
		resources.Energy = CalculateSolarPlantEnergyProduction(level)
	case domain.BuildingFusionPlant:
		resources.Energy = CalculateFusionPlantEnergyProduction(level, energyTechLevel)
		resources.Deuterium = CalculateFusionPlantDeuteriumConsumption(level)
	}

	return resources
}
