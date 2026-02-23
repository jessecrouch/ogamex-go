package service

import (
	"context"
	"time"

	"ogamex-go/internal/domain"
	"ogamex-go/internal/formula"
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

func (s *UnitService) intToUnitType(unitID int) domain.UnitType {
	mapping := map[int]domain.UnitType{
		202: domain.UnitSmallCargo,
		203: domain.UnitLargeCargo,
		204: domain.UnitLightFighter,
		205: domain.UnitHeavyFighter,
		206: domain.UnitCruiser,
		207: domain.UnitBattleship,
		215: domain.UnitBattlecruiser,
		211: domain.UnitBomber,
		213: domain.UnitDestroyer,
		214: domain.UnitDeathstar,
		209: domain.UnitRecycler,
		210: domain.UnitEspionageProbe,
		212: domain.UnitSolarSatellite,
		217: domain.UnitCrawler,
		208: domain.UnitColonyShip,
		218: domain.UnitReaper,
		219: domain.UnitPathfinder,
		401: domain.UnitRocketLauncher,
		402: domain.UnitLightLaser,
		403: domain.UnitHeavyLaser,
		405: domain.UnitGaussCannon,
		404: domain.UnitIonCannon,
		406: domain.UnitPlasmaTurret,
		407: domain.UnitSmallShieldDome,
		408: domain.UnitAntiBallisticMissiles,
		409: domain.UnitInterplanetaryMissiles,
	}
	if ut, ok := mapping[unitID]; ok {
		return ut
	}
	return domain.UnitType(unitID)
}

func (s *UnitService) GetUnitCost(unitID int) (metal, crystal, deuterium int64, buildTime time.Duration) {
	metal, crystal, deuterium, seconds := formula.CalculateUnitCost(s.intToUnitType(unitID))
	return metal, crystal, deuterium, time.Duration(seconds) * time.Second
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
	unitType := s.intToUnitType(unitID)
	return unitType.IsShip()
}

func (s *UnitService) IsDefenseUnit(unitID int) bool {
	unitType := s.intToUnitType(unitID)
	return unitType.IsDefense()
}

func (s *UnitService) IsValidUnit(unitID int) bool {
	unitType := s.intToUnitType(unitID)
	return unitType.IsShip() || unitType.IsDefense()
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
