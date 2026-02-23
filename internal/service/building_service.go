package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/domain"
	"ogamex-go/internal/formula"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

var (
	ErrBuildingNotFound     = errors.New("building not found")
	ErrInsufficientFunds    = errors.New("insufficient resources")
	ErrBuildingNotAvailable = errors.New("building not available on this planet")
	ErrQueueFull           = errors.New("building queue is full")
	ErrTechNotFound        = errors.New("technology not found")
	ErrLabBusy             = errors.New("research lab is busy")
)

type BuildingService struct {
	planetRepo    repository.PlanetRepository
	buildingQueue repository.BuildingQueueRepository
	techRepo      repository.UserTechRepository
}

func NewBuildingService(
	planetRepo repository.PlanetRepository,
	buildingQueue repository.BuildingQueueRepository,
	techRepo repository.UserTechRepository,
) *BuildingService {
	return &BuildingService{
		planetRepo:    planetRepo,
		buildingQueue: buildingQueue,
		techRepo:      techRepo,
	}
}

type BuildingCost struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
	Energy    int64
}

func (s *BuildingService) GetBuildingCost(buildingID int, level int) BuildingCost {
	buildingType := s.intToBuildingType(buildingID)
	metal, crystal, deuterium := formula.CalculateBuildingCost(buildingType, level)
	return BuildingCost{
		Metal:     metal,
		Crystal:   crystal,
		Deuterium: deuterium,
	}
}

func (s *BuildingService) intToBuildingType(buildingID int) domain.BuildingType {
	mapping := map[int]domain.BuildingType{
		1:  domain.BuildingMetalMine,
		2:  domain.BuildingCrystalMine,
		3:  domain.BuildingDeuteriumSynthesizer,
		4:  domain.BuildingSolarPlant,
		5:  domain.BuildingFusionPlant,
		6:  domain.BuildingMetalStorage,
		7:  domain.BuildingCrystalStorage,
		8:  domain.BuildingDeuteriumStorage,
		14: domain.BuildingMetalMine2,
		15: domain.BuildingCrystalMine2,
		16: domain.BuildingDeuteriumSynthesizer2,
		17: domain.BuildingSolarSatellite,
		18: domain.BuildingCrawler,
		19: domain.BuildingSpaceDock,
		20: domain.BuildingNaniteFactory,
		21: domain.BuildingTerraformer,
		22: domain.BuildingMissileSilo,
		23: domain.BuildingRoboticsFactory,
		24: domain.BuildingShipyard,
		25: domain.BuildingResearchLab,
		31: domain.BuildingResearchLab,
		33: domain.BuildingNaniteFactory,
		36: domain.BuildingSpaceDock,
		41: domain.BuildingLunarBase,
		42: domain.BuildingSensorPhalanx,
		43: domain.BuildingTerraformer,
		44: domain.BuildingMissileSilo,
		45: domain.BuildingJumpGate,
	}
	if bt, ok := mapping[buildingID]; ok {
		return bt
	}
	return domain.BuildingType(buildingID)
}

func (s *BuildingService) StartBuilding(ctx context.Context, planetID uint, buildingID int, level int) error {
	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return ErrBuildingNotFound
	}

	cost := s.GetBuildingCost(buildingID, level)

	if planet.Metal < cost.Metal || planet.Crystal < cost.Crystal || planet.Deuterium < cost.Deuterium {
		return ErrInsufficientFunds
	}

	currentQueue, err := s.buildingQueue.GetCurrent(ctx, planetID)
	if err == nil && currentQueue != nil {
		return ErrQueueFull
	}

	buildTime := s.calculateBuildTime(buildingID, level, planet)

	err = s.planetRepo.SubResources(ctx, planetID, cost.Metal, cost.Crystal, cost.Deuterium)
	if err != nil {
		return err
	}

	now := time.Now()
	queueItem := &schema.BuildingQueue{
		PlanetID:   planetID,
		BuildingID: buildingID,
		Level:      level,
		StartTime:  now,
		EndTime:    now.Add(buildTime),
	}

	return s.buildingQueue.Create(ctx, queueItem)
}

func (s *BuildingService) calculateBuildTime(buildingID int, level int, planet *schema.Planet) time.Duration {
	baseTime := 60 * time.Second
	robotFactoryBonus := 1.0
	if planet.RobotFactory > 0 {
		robotFactoryBonus = 1.0 / (1.0 + float64(planet.RobotFactory)*0.25)
	}
	
	naniteBonus := 1.0
	if planet.NaniteFactory > 0 {
		naniteBonus = 1.0 / (1.0 + float64(planet.NaniteFactory)*0.5)
	}

	factor := robotFactoryBonus * naniteBonus
	baseMultiplier := formula.GetCostFactor(level, 1.0)

	return time.Duration(float64(baseTime) * baseMultiplier * factor)
}

func (s *BuildingService) CompleteBuilding(ctx context.Context, queueID uint) error {
	queue, err := s.buildingQueue.GetByID(ctx, queueID)
	if err != nil {
		return err
	}

	if time.Now().Before(queue.EndTime) {
		return errors.New("building not yet complete")
	}

	planet, err := s.planetRepo.GetByID(ctx, queue.PlanetID)
	if err != nil {
		return err
	}

	switch queue.BuildingID {
	case 1:
		planet.MetalMine = queue.Level
	case 2:
		planet.CrystalMine = queue.Level
	case 3:
		planet.DeuteriumSynthesizer = queue.Level
	case 4:
		planet.SolarPlant = queue.Level
	case 12:
		planet.FusionPlant = queue.Level
	case 22:
		planet.MetalStorageBuilding = queue.Level
	case 23:
		planet.CrystalStorageBuilding = queue.Level
	case 24:
		planet.DeuteriumStorageBuilding = queue.Level
	case 14:
		planet.RobotFactory = queue.Level
	case 21:
		planet.Shipyard = queue.Level
	case 31:
		planet.ResearchLab = queue.Level
	case 44:
		planet.MissileSilo = queue.Level
	case 33:
		planet.NaniteFactory = queue.Level
	case 43:
		planet.Terraformer = queue.Level
	case 36:
		planet.SpaceDock = queue.Level
	case 41:
		planet.LunarBase = queue.Level
	case 42:
		planet.SensorPhalanx = queue.Level
	case 45:
		planet.JumpGate = queue.Level
	}

	planet.FieldsUsed = s.calculateFieldsUsed(planet)

	err = s.planetRepo.Update(ctx, planet)
	if err != nil {
		return err
	}

	return s.buildingQueue.Delete(ctx, queueID)
}

func (s *BuildingService) calculateFieldsUsed(planet *schema.Planet) int {
	total := 0
	if planet.MetalMine > 0 {
		total++
	}
	if planet.CrystalMine > 0 {
		total++
	}
	if planet.DeuteriumSynthesizer > 0 {
		total++
	}
	if planet.SolarPlant > 0 {
		total++
	}
	if planet.FusionPlant > 0 {
		total++
	}
	if planet.MetalStorageBuilding > 0 {
		total++
	}
	if planet.CrystalStorageBuilding > 0 {
		total++
	}
	if planet.DeuteriumStorageBuilding > 0 {
		total++
	}
	if planet.RobotFactory > 0 {
		total++
	}
	if planet.Shipyard > 0 {
		total++
	}
	if planet.ResearchLab > 0 {
		total++
	}
	if planet.MissileSilo > 0 {
		total++
	}
	if planet.NaniteFactory > 0 {
		total++
	}
	if planet.Terraformer > 0 {
		total++
	}
	if planet.SpaceDock > 0 {
		total++
	}
	if planet.LunarBase > 0 {
		total++
	}
	if planet.SensorPhalanx > 0 {
		total++
	}
	if planet.JumpGate > 0 {
		total++
	}
	return total
}

func (s *BuildingService) GetQueue(ctx context.Context, planetID uint) ([]*schema.BuildingQueue, error) {
	return s.buildingQueue.GetByPlanetID(ctx, planetID)
}

func (s *BuildingService) CancelBuilding(ctx context.Context, queueID uint) error {
	queue, err := s.buildingQueue.GetByID(ctx, queueID)
	if err != nil {
		return err
	}

	refund := s.GetBuildingCost(queue.BuildingID, queue.Level)
	refund.Metal /= 2
	refund.Crystal /= 2
	refund.Deuterium /= 2

	_ = s.planetRepo.AddResources(ctx, queue.PlanetID, refund.Metal, refund.Crystal, refund.Deuterium)

	return s.buildingQueue.Cancel(ctx, queueID)
}
