package service

import (
	"context"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type UnitService struct {
	planetRepo   repository.PlanetRepository
	unitQueueRepo repository.UnitQueueRepository
	techRepo     repository.UserTechRepository
}

func NewUnitService(
	planetRepo repository.PlanetRepository,
	unitQueueRepo repository.UnitQueueRepository,
	techRepo repository.UserTechRepository,
) *UnitService {
	return &UnitService{
		planetRepo:   planetRepo,
		unitQueueRepo: unitQueueRepo,
		techRepo:     techRepo,
	}
}

var shipyardUnits = map[int]struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
	BuildTime time.Duration
}{
	202: {2000, 2000, 0, 5 * time.Second},   // Small Cargo
	203: {6000, 6000, 0, 8 * time.Second},   // Large Cargo
	204: {10000, 6000, 2000, 20 * time.Second}, // Light Fighter
	205: {25000, 15000, 5000, 40 * time.Second}, // Heavy Fighter
	206: {10000, 20000, 10000, 10 * time.Second}, // Cruiser
	207: {50000, 25000, 15000, 80 * time.Second}, // Battleship
	208: {10000, 10000, 0, 50 * time.Second}, // Colony Ship
	209: {10000, 6000, 2000, 15 * time.Second}, // Recycler
	210: {0, 1000, 0, 30 * time.Second},       // Espionage Probe
	211: {50000, 50000, 25000, 200 * time.Second}, // Bomber
	212: {0, 2000, 500, 3 * time.Second},     // Solar Satellite
	213: {10000, 10000, 0, 30 * time.Second}, // Destroyer
	214: {100000, 100000, 50000, 400 * time.Second}, // Deathstar
	215: {3000, 1000, 0, 4 * time.Second},     // Battlecruiser
	217: {2000, 2000, 1000, 10 * time.Second}, // Crawler
	218: {8000, 0, 0, 20 * time.Second},      // Reaper
	219: {20000, 10000, 10000, 75 * time.Second}, // Pathfinder
}

var defenseUnits = map[int]struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
	BuildTime time.Duration
}{
	401: {2000, 0, 0, 10 * time.Second},      // Rocket Launcher
	402: {1500, 500, 0, 11 * time.Second},    // Light Laser
	403: {6000, 2000, 0, 22 * time.Second},   // Heavy Laser
	404: {2000, 6000, 0, 16 * time.Second},   // Ion Cannon
	405: {20000, 15000, 2000, 45 * time.Second}, // Gauss Cannon
	406: {50000, 50000, 30000, 90 * time.Second}, // Plasma Turret
	407: {10000, 10000, 0, 20 * time.Second}, // Shield Dome (Large Shield)
	408: {8000, 2000, 0, 15 * time.Second},  // Missile Interceptor
	409: {15000, 5000, 0, 20 * time.Second},  // Missile Launcher
}

func (s *UnitService) GetUnitCost(unitID int) (metal, crystal, deuterium int64, buildTime time.Duration) {
	if unit, ok := shipyardUnits[unitID]; ok {
		return unit.Metal, unit.Crystal, unit.Deuterium, unit.BuildTime
	}
	if unit, ok := defenseUnits[unitID]; ok {
		return unit.Metal, unit.Crystal, unit.Deuterium, unit.BuildTime
	}
	return 0, 0, 0, 0
}

func (s *UnitService) BuildUnit(ctx context.Context, planetID uint, unitID int, amount int) error {
	if amount <= 0 {
		amount = 1
	}

	metal, crystal, deuterium, buildTime := s.GetUnitCost(unitID)
	if metal == 0 && crystal == 0 && deuterium == 0 {
		return ErrInvalidUnit
	}

	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return err
	}

	totalMetal := metal * int64(amount)
	totalCrystal := crystal * int64(amount)
	totalDeuterium := deuterium * int64(amount)

	if planet.Metal < totalMetal || planet.Crystal < totalCrystal || planet.Deuterium < totalDeuterium {
		return ErrInsufficientFunds
	}

	currentQueue, _ := s.unitQueueRepo.GetCurrent(ctx, planetID)
	if currentQueue != nil {
		return ErrQueueFull
	}

	err = s.planetRepo.SubResources(ctx, planetID, totalMetal, totalCrystal, totalDeuterium)
	if err != nil {
		return err
	}

	tech, _ := s.techRepo.GetByUserID(ctx, planet.UserID)
	techBonus := 1.0
	if tech != nil && tech.AssemblyTechnology > 0 {
		techBonus = 1.0 - float64(tech.AssemblyTechnology)*0.02
		if techBonus < 0.1 {
			techBonus = 0.1
		}
	}

	robotFactoryBonus := 1.0
	if planet.RobotFactory > 0 {
		robotFactoryBonus = 1.0 / (1.0 + float64(planet.RobotFactory)*0.25)
	}

	naniteBonus := 1.0
	if planet.NaniteFactory > 0 {
		naniteBonus = 1.0 / (1.0 + float64(planet.NaniteFactory)*0.5)
	}

	actualBuildTime := time.Duration(float64(buildTime) * techBonus * robotFactoryBonus * naniteBonus)

	now := time.Now()
	queue := &schema.UnitQueue{
		PlanetID: planetID,
		UnitID:   unitID,
		Amount:   amount,
		StartTime: now,
		EndTime:  now.Add(actualBuildTime),
	}

	return s.unitQueueRepo.Create(ctx, queue)
}

func (s *UnitService) CompleteUnit(ctx context.Context, queueID uint) error {
	queue, err := s.unitQueueRepo.GetByID(ctx, queueID)
	if err != nil {
		return err
	}

	if time.Now().Before(queue.EndTime) {
		return ErrNotReady
	}

	planet, err := s.planetRepo.GetByID(ctx, queue.PlanetID)
	if err != nil {
		return err
	}

	switch queue.UnitID {
	case 202:
		planet.SmallCargo += queue.Amount
	case 203:
		planet.LargeCargo += queue.Amount
	case 204:
		planet.LightFighter += queue.Amount
	case 205:
		planet.HeavyFighter += queue.Amount
	case 206:
		planet.Cruiser += queue.Amount
	case 207:
		planet.Battleship += queue.Amount
	case 208:
		planet.ColonyShip += queue.Amount
	case 209:
		planet.Recycler += queue.Amount
	case 210:
		planet.EspionageProbe += queue.Amount
	case 211:
		planet.Bomber += queue.Amount
	case 212:
		planet.SolarSatellite += queue.Amount
	case 213:
		planet.Destroyer += queue.Amount
	case 214:
		planet.Deathstar += queue.Amount
	case 215:
		planet.Battlecruiser += queue.Amount
	case 217:
		planet.Crawler += queue.Amount
	case 218:
		planet.Reaper += queue.Amount
	case 219:
		planet.Pathfinder += queue.Amount
	case 401:
		planet.RocketLauncher += queue.Amount
	case 402:
		planet.LightLaser += queue.Amount
	case 403:
		planet.HeavyLaser += queue.Amount
	case 404:
		planet.IonCannon += queue.Amount
	case 405:
		planet.GaussCannon += queue.Amount
	case 406:
		planet.PlasmaTurret += queue.Amount
	case 407:
		planet.ShieldDome += queue.Amount
	case 408:
		planet.MissileInterceptor += queue.Amount
	case 409:
		planet.MissileLauncher += queue.Amount
	}

	err = s.planetRepo.Update(ctx, planet)
	if err != nil {
		return err
	}

	return s.unitQueueRepo.Delete(ctx, queueID)
}

func (s *UnitService) GetQueue(ctx context.Context, planetID uint) ([]*schema.UnitQueue, error) {
	return s.unitQueueRepo.GetByPlanetID(ctx, planetID)
}

func (s *UnitService) CancelQueue(ctx context.Context, queueID uint) error {
	queue, err := s.unitQueueRepo.GetByID(ctx, queueID)
	if err != nil {
		return err
	}

	metal, crystal, deuterium, _ := s.GetUnitCost(queue.UnitID)
	refund := int64(queue.Amount)

	_ = s.planetRepo.AddResources(ctx, queue.PlanetID, metal*refund/2, crystal*refund/2, deuterium*refund/2)

	return s.unitQueueRepo.Delete(ctx, queueID)
}

func (s *UnitService) IsShipUnit(unitID int) bool {
	_, ok := shipyardUnits[unitID]
	return ok
}

func (s *UnitService) IsDefenseUnit(unitID int) bool {
	_, ok := defenseUnits[unitID]
	return ok
}

func (s *UnitService) IsValidUnit(unitID int) bool {
	return s.IsShipUnit(unitID) || s.IsDefenseUnit(unitID)
}

var (
	ErrInvalidUnit = &ServiceError{"invalid unit type"}
	ErrNotReady   = &ServiceError{"unit not ready yet"}
)

type ServiceError struct {
	msg string
}

func (e *ServiceError) Error() string {
	return e.msg
}
