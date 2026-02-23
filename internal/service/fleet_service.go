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
	ErrFleetNotFound       = errors.New("fleet not found")
	ErrInvalidFleetParams  = errors.New("invalid fleet parameters")
	ErrInsufficientFleet   = errors.New("insufficient ships")
	ErrDestinationOccupied = errors.New("destination occupied")
)

type FleetService struct {
	planetRepo    repository.PlanetRepository
	fleetRepo     repository.FleetMissionRepository
	userRepo      repository.UserRepository
	techRepo      repository.UserTechRepository
	universeSpeed int
}

func NewFleetService(
	planetRepo repository.PlanetRepository,
	fleetRepo repository.FleetMissionRepository,
	userRepo repository.UserRepository,
	techRepo repository.UserTechRepository,
	universeSpeed int,
) *FleetService {
	return &FleetService{
		planetRepo:    planetRepo,
		fleetRepo:     fleetRepo,
		userRepo:      userRepo,
		techRepo:      techRepo,
		universeSpeed: universeSpeed,
	}
}

type FleetParams struct {
	Origin      schema.Planet
	Target      schema.Planet
	MissionType int
	Resources   struct {
		Metal     int64
		Crystal   int64
		Deuterium int64
	}
	Ships map[int]int
}

func (s *FleetService) SendFleet(ctx context.Context, userID uint, params FleetParams) error {
	tech, err := s.techRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	origin := params.Origin
	if origin.UserID != userID {
		return ErrInvalidFleetParams
	}

	totalCargo := s.calculateCargoCapacity(params.Ships, tech)
	totalResources := params.Resources.Metal + params.Resources.Crystal + params.Resources.Deuterium

	if totalResources > totalCargo {
		return errors.New("insufficient cargo capacity")
	}

	distance := formula.CalculateDistance(origin.Galaxy, origin.System, origin.Position,
		params.Target.Galaxy, params.Target.System, params.Target.Position)

	speed := s.calculateFleetSpeed(params.Ships, tech)

	travelTime := formula.CalculateFlightTime(distance, speed, s.universeSpeed)
	if travelTime == 0 {
		return ErrInvalidFleetParams
	}

	consumption := s.calculateFuelConsumption(params.Ships, distance, speed, tech)
	if origin.Deuterium < consumption {
		return ErrInsufficientFunds
	}

	err = s.planetRepo.SubResources(ctx, origin.ID, 0, 0, consumption)
	if err != nil {
		return err
	}

	fuelCost := params.Resources.Metal + params.Resources.Crystal + params.Resources.Deuterium
	err = s.planetRepo.SubResources(ctx, origin.ID, params.Resources.Metal, params.Resources.Crystal, params.Resources.Deuterium)
	if err != nil {
		return err
	}
	_ = fuelCost

	now := time.Now()
	mission := &schema.FleetMission{
		UserID:            userID,
		MissionType:       params.MissionType,
		OriginGalaxy:      origin.Galaxy,
		OriginSystem:      origin.System,
		OriginPosition:    origin.Position,
		OriginPlanetType:  1,
		TargetGalaxy:      params.Target.Galaxy,
		TargetSystem:      params.Target.System,
		TargetPosition:    params.Target.Position,
		TargetPlanetType:  1,
		LaunchTime:        now,
		ArrivalTime:       now.Add(time.Duration(travelTime) * time.Second),
		ReturnTime:        now.Add(time.Duration(travelTime*2) * time.Second),
		Metal:             params.Resources.Metal,
		Crystal:           params.Resources.Crystal,
		Deuterium:         params.Resources.Deuterium,
		Status:            0,
	}

	return s.fleetRepo.Create(ctx, mission)
}

func (s *FleetService) calculateCargoCapacity(ships map[int]int, tech *schema.UserTech) int64 {
	baseCargo := map[int]int64{
		202: 5000,  // Small Cargo
		203: 25000, // Large Cargo
		208: 7500,  // Colony Ship
		209: 20000, // Recycler
		210: 0,     // Espionage Probe
		212: 0,     // Solar Satellite
		217: 0,     // Crawler
		219: 10000, // Pathfinder
	}

	total := int64(0)
	for shipID, count := range ships {
		if cargo, ok := baseCargo[shipID]; ok {
			total += cargo * int64(count)
		}
	}

	return total
}

func (s *FleetService) calculateFleetSpeed(ships map[int]int, tech *schema.UserTech) int {
	baseSpeeds := map[int]int{
		202: 5000,  // Small Cargo (5000/10000 with impulse)
		203: 7500,  // Large Cargo
		208: 2500,  // Colony Ship
		209: 2000,  // Recycler (4000/6000 with upgrades)
		210: 100000000, // Espionage Probe
		217: 4000,  // Crawler
		219: 12000, // Pathfinder
	}

	slowestSpeed := 0
	for shipID, count := range ships {
		if count > 0 {
			if speed, ok := baseSpeeds[shipID]; ok {
				if slowestSpeed == 0 || speed < slowestSpeed {
					slowestSpeed = speed
				}
			}
		}
	}

	if slowestSpeed == 0 {
		return 0
	}

	techBonus := 1.0 + float64(tech.ComputerTechnology)*0.05
	return int(float64(slowestSpeed) * techBonus)
}

func (s *FleetService) calculateFuelConsumption(ships map[int]int, distance int, speed int, tech *schema.UserTech) int64 {
	baseConsumption := map[int]int{
		202: 1,
		203: 10,
		209: 5,
		208: 10,
		210: 1,
		212: 1,
		217: 2,
		219: 8,
	}

	total := 0
	for shipID, count := range ships {
		if consumption, ok := baseConsumption[shipID]; ok {
			total += consumption * count
		}
	}

	consumptionFactor := float64(distance) / float64(speed) / 10.0
	return int64(float64(total) * consumptionFactor * float64(s.universeSpeed))
}

func (s *FleetService) GetActiveMissions(ctx context.Context, userID uint) ([]*schema.FleetMission, error) {
	return s.fleetRepo.GetByUserID(ctx, userID)
}

func (s *FleetService) ProcessArrivingMissions(ctx context.Context) error {
	now := time.Now()
	missions, err := s.fleetRepo.GetArriving(ctx, now)
	if err != nil {
		return err
	}

	for _, mission := range missions {
		s.processMission(ctx, mission)
	}

	return nil
}

func (s *FleetService) processMission(ctx context.Context, mission *schema.FleetMission) error {
	switch mission.MissionType {
	case 1: // Attack
		return s.processAttack(ctx, mission)
	case 3: // Transport
		return s.processTransport(ctx, mission)
	case 4: // Deploy
		return s.processDeploy(ctx, mission)
	case 7: // Harvest
		return s.processHarvest(ctx, mission)
	case 8: // Colonize
		return s.processColonize(ctx, mission)
	case 9: // Recycle
		return s.processRecycle(ctx, mission)
	}

	mission.Status = 1
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) processTransport(ctx context.Context, mission *schema.FleetMission) error {
	target, err := s.planetRepo.GetByCoords(ctx, 0, mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition)
	if err != nil {
		return err
	}

	err = s.planetRepo.AddResources(ctx, target.ID, mission.Metal, mission.Crystal, mission.Deuterium)
	if err != nil {
		return err
	}

	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) processDeploy(ctx context.Context, mission *schema.FleetMission) error {
	target, err := s.planetRepo.GetByCoords(ctx, 0, mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition)
	if err != nil {
		return err
	}

	err = s.planetRepo.AddResources(ctx, target.ID, mission.Metal, mission.Crystal, mission.Deuterium)
	if err != nil {
		return err
	}

	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) processAttack(ctx context.Context, mission *schema.FleetMission) error {
	return nil
}

func (s *FleetService) processHarvest(ctx context.Context, mission *schema.FleetMission) error {
	return nil
}

func (s *FleetService) processColonize(ctx context.Context, mission *schema.FleetMission) error {
	return nil
}

func (s *FleetService) processRecycle(ctx context.Context, mission *schema.FleetMission) error {
	return nil
}
