package service

import (
	"context"

	"ogamex-go/internal/formula"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type ProductionService struct {
	planetRepo repository.PlanetRepository
	techRepo   repository.UserTechRepository
	economySpeed int
}

func NewProductionService(
	planetRepo repository.PlanetRepository,
	techRepo repository.UserTechRepository,
	economySpeed int,
) *ProductionService {
	return &ProductionService{
		planetRepo:   planetRepo,
		techRepo:     techRepo,
		economySpeed: economySpeed,
	}
}

type ProductionResult struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
	Energy    int64
}

func (s *ProductionService) CalculateProduction(planet *schema.Planet, tech *schema.UserTech) ProductionResult {
	if tech == nil {
		tech = &schema.UserTech{UserID: planet.UserID}
	}
	
	energyAvailable := s.calculateEnergyProduction(planet, tech)
	energyUsed := s.calculateEnergyUsage(planet)
	
	netEnergy := energyAvailable + energyUsed
	
	production := ProductionResult{
		Energy: netEnergy,
	}

	if netEnergy < 0 {
		deficitFactor := 1.0 + float64(netEnergy)/200.0
		if deficitFactor < 0.5 {
			deficitFactor = 0.5
		}
		production.Metal = int64(float64(s.calculateMetalProduction(planet, tech)) * deficitFactor)
		production.Crystal = int64(float64(s.calculateCrystalProduction(planet, tech)) * deficitFactor)
		production.Deuterium = int64(float64(s.calculateDeuteriumProduction(planet, tech)) * deficitFactor)
	} else {
		production.Metal = s.calculateMetalProduction(planet, tech)
		production.Crystal = s.calculateCrystalProduction(planet, tech)
		production.Deuterium = s.calculateDeuteriumProduction(planet, tech)
	}

	return production
}

func (s *ProductionService) calculateEnergyProduction(planet *schema.Planet, tech *schema.UserTech) int64 {
	solarPlantEnergy := formula.CalculateSolarPlantEnergyProduction(planet.SolarPlant)
	fusionPlantEnergy := formula.CalculateFusionPlantEnergyProduction(planet.FusionPlant, tech.EnergyTechnology)
	satelliteEnergy := formula.CalculateSolarSatelliteEnergyProduction(planet.TempMax, planet.SolarSatellite)
	
	return solarPlantEnergy + fusionPlantEnergy + satelliteEnergy
}

func (s *ProductionService) calculateEnergyUsage(planet *schema.Planet) int64 {
	metalMineUsage := formula.CalculateMetalMineEnergyConsumption(planet.MetalMine)
	crystalMineUsage := formula.CalculateCrystalMineEnergyConsumption(planet.CrystalMine)
	deuteriumSynthUsage := formula.CalculateDeuteriumSynthesizerEnergyConsumption(planet.DeuteriumSynthesizer)
	
	return metalMineUsage + crystalMineUsage + deuteriumSynthUsage
}

func (s *ProductionService) calculateMetalProduction(planet *schema.Planet, tech *schema.UserTech) int64 {
	energyFactor := float64(planet.MetalMinePercent) / 100.0
	
	baseProd := formula.CalculateMetalMineProduction(planet.MetalMine, energyFactor)
	
	return int64(float64(baseProd) * float64(s.economySpeed))
}

func (s *ProductionService) calculateCrystalProduction(planet *schema.Planet, tech *schema.UserTech) int64 {
	energyFactor := float64(planet.CrystalMinePercent) / 100.0
	
	baseProd := formula.CalculateCrystalMineProduction(planet.CrystalMine, energyFactor)
	
	return int64(float64(baseProd) * float64(s.economySpeed))
}

func (s *ProductionService) calculateDeuteriumProduction(planet *schema.Planet, tech *schema.UserTech) int64 {
	energyFactor := float64(planet.DeuteriumSynthesizerPercent) / 100.0
	
	baseProd := formula.CalculateDeuteriumProduction(planet.DeuteriumSynthesizer, planet.TempMax, energyFactor)
	
	return int64(float64(baseProd) * float64(s.economySpeed))
}

func (s *ProductionService) CalculateProductionForType(planet *schema.Planet, tech *schema.UserTech, resourceType string) int64 {
	result := s.CalculateProductionFull(planet, tech)
	switch resourceType {
	case "metal":
		return result.Metal
	case "crystal":
		return result.Crystal
	case "deuterium":
		return result.Deuterium
	}
	return 0
}

func (s *ProductionService) CalculateEnergy(planet *schema.Planet, tech *schema.UserTech) int64 {
	result := s.CalculateProductionFull(planet, tech)
	return result.Energy
}

func (s *ProductionService) CalculateProductionFull(planet *schema.Planet, tech *schema.UserTech) ProductionResult {
	energyAvailable := s.calculateEnergyProduction(planet, tech)
	energyUsed := s.calculateEnergyUsage(planet)

	netEnergy := energyAvailable + energyUsed

	production := ProductionResult{
		Energy: netEnergy,
	}

	if netEnergy < 0 {
		deficitFactor := 1.0 + float64(netEnergy)/200.0
		if deficitFactor < 0.5 {
			deficitFactor = 0.5
		}
		production.Metal = int64(float64(s.calculateMetalProduction(planet, tech)) * deficitFactor)
		production.Crystal = int64(float64(s.calculateCrystalProduction(planet, tech)) * deficitFactor)
		production.Deuterium = int64(float64(s.calculateDeuteriumProduction(planet, tech)) * deficitFactor)
	} else {
		production.Metal = s.calculateMetalProduction(planet, tech)
		production.Crystal = s.calculateCrystalProduction(planet, tech)
		production.Deuterium = s.calculateDeuteriumProduction(planet, tech)
	}

	return production
}

func (s *ProductionService) UpdatePlanetResources(ctx context.Context, planetID uint) error {
	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return err
	}

	tech, err := s.techRepo.GetByUserID(ctx, planet.UserID)
	if err != nil {
		tech = &schema.UserTech{UserID: planet.UserID}
	}

	production := s.CalculateProduction(planet, tech)
	
	planet.Metal += production.Metal
	planet.Crystal += production.Crystal
	planet.Deuterium += production.Deuterium
	planet.EnergyAvailable = production.Energy
	
	if planet.Metal > planet.MetalCapacity {
		planet.Metal = planet.MetalCapacity
	}
	if planet.Crystal > planet.CrystalCapacity {
		planet.Crystal = planet.CrystalCapacity
	}
	if planet.Deuterium > planet.DeuteriumCapacity {
		planet.Deuterium = planet.DeuteriumCapacity
	}

	return s.planetRepo.Update(ctx, planet)
}

func (s *ProductionService) UpdateAllPlanets(ctx context.Context) error {
	planets, err := s.planetRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	
	for _, planet := range planets {
		if err := s.UpdatePlanetResources(ctx, planet.ID); err != nil {
			continue
		}
	}
	
	return nil
}
