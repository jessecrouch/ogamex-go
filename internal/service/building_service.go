package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/formula"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

var (
	ErrBuildingNotFound    = errors.New("building not found")
	ErrInsufficientFunds   = errors.New("insufficient resources")
	ErrBuildingNotAvailable = errors.New("building not available on this planet")
	ErrQueueFull          = errors.New("building queue is full")
	ErrTechNotFound       = errors.New("technology not found")
	ErrLabBusy            = errors.New("research lab is busy")
)

type BuildingService struct {
	planetRepo     repository.PlanetRepository
	buildingQueue  repository.BuildingQueueRepository
	techRepo       repository.UserTechRepository
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
	baseCost := s.getBaseCost(buildingID)
	factor := s.getCostFactor(buildingID)
	
	multiplier := formula.GetCostFactor(level, factor)
	
	return BuildingCost{
		Metal:     int64(float64(baseCost.Metal) * multiplier),
		Crystal:   int64(float64(baseCost.Crystal) * multiplier),
		Deuterium: int64(float64(baseCost.Deuterium) * multiplier),
	}
}

func (s *BuildingService) getBaseCost(buildingID int) BuildingCost {
	costs := map[int]BuildingCost{
		1:  {Metal: 60, Crystal: 15, Deuterium: 0},    // Metal Mine
		2:  {Metal: 48, Crystal: 24, Deuterium: 0},   // Crystal Mine
		3:  {Metal: 225, Crystal: 75, Deuterium: 0},  // Deuterium Synthesizer
		4:  {Metal: 75, Crystal: 30, Deuterium: 0},   // Solar Plant
		12: {Metal: 900, Crystal: 360, Deuterium: 180}, // Fusion Plant
	}
	
	if cost, ok := costs[buildingID]; ok {
		return cost
	}
	return BuildingCost{}
}

func (s *BuildingService) getCostFactor(buildingID int) float64 {
	factors := map[int]float64{
		1:  1.5,
		2:  1.6,
		3:  1.5,
		4:  1.5,
		12: 1.8,
	}
	
	if factor, ok := factors[buildingID]; ok {
		return factor
	}
	return 1.5
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
	fields := map[int]int{
		planet.MetalMine:           1,
		planet.CrystalMine:         1,
		planet.DeuteriumSynthesizer: 1,
		planet.SolarPlant:          1,
		planet.FusionPlant:         1,
		planet.MetalStorageBuilding:   1,
		planet.CrystalStorageBuilding: 1,
		planet.DeuteriumStorageBuilding: 1,
		planet.RobotFactory:        1,
		planet.Shipyard:            1,
		planet.ResearchLab:         1,
		planet.MissileSilo:         1,
		planet.NaniteFactory:       1,
		planet.Terraformer:         1,
		planet.SpaceDock:           1,
		planet.LunarBase:           1,
		planet.SensorPhalanx:       1,
		planet.JumpGate:            1,
	}

	total := 0
	for _, f := range fields {
		total += f
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
