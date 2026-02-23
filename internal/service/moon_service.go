package service

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type MoonService struct {
	planetRepo    repository.PlanetRepository
	fleetRepo     repository.FleetMissionRepository
	userRepo      repository.UserRepository
}

func NewMoonService(planetRepo repository.PlanetRepository, fleetRepo repository.FleetMissionRepository, userRepo repository.UserRepository) *MoonService {
	return &MoonService{
		planetRepo:    planetRepo,
		fleetRepo:     fleetRepo,
		userRepo:      userRepo,
	}
}

type JumpGateTarget struct {
	PlanetID   uint   `json:"planet_id"`
	PlanetName string `json:"planet_name"`
	Galaxy     int    `json:"galaxy"`
	System     int    `json:"system"`
	Position   int    `json:"position"`
}

func (s *MoonService) GetJumpGateTargets(ctx context.Context, moonID uint, userID uint) ([]*JumpGateTarget, error) {
	moon, err := s.planetRepo.GetByID(ctx, moonID)
	if err != nil {
		return nil, errors.New("moon not found")
	}

	if moon.UserID != userID {
		return nil, errors.New("moon does not belong to user")
	}

	if moon.JumpGate <= 0 {
		return nil, errors.New("jump gate not built")
	}

	moons, err := s.planetRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var targets []*JumpGateTarget
	for _, p := range moons {
		if p.IsMoon && p.UserID == userID && p.ID != moonID && p.JumpGate > 0 {
			targets = append(targets, &JumpGateTarget{
				PlanetID:   p.ID,
				PlanetName: p.Name,
				Galaxy:     p.Galaxy,
				System:     p.System,
				Position:   p.Position,
			})
		}
	}

	return targets, nil
}

type JumpFleetInput struct {
	OriginMoonID  uint              `json:"origin_moon_id"`
	TargetMoonID  uint              `json:"target_moon_id"`
	Ships         map[string]int    `json:"ships"`
}

func (s *MoonService) ExecuteJumpGate(ctx context.Context, userID uint, input JumpFleetInput) error {
	originMoon, err := s.planetRepo.GetByID(ctx, input.OriginMoonID)
	if err != nil {
		return errors.New("origin moon not found")
	}

	if originMoon.UserID != userID {
		return errors.New("origin moon does not belong to user")
	}

	if originMoon.JumpGate <= 0 {
		return errors.New("jump gate not built")
	}

	targetMoon, err := s.planetRepo.GetByID(ctx, input.TargetMoonID)
	if err != nil {
		return errors.New("target moon not found")
	}

	if targetMoon.UserID != userID {
		return errors.New("target moon does not belong to user")
	}

	if targetMoon.JumpGate <= 0 {
		return errors.New("target moon has no jump gate")
	}

	jumpFuel := s.calculateJumpFuel(input.Ships)
	if originMoon.Deuterium < jumpFuel {
		return errors.New("insufficient deuterium for jump")
	}

	originMoon.Deuterium -= jumpFuel
	err = s.planetRepo.Update(ctx, originMoon)
	if err != nil {
		return err
	}

	s.addShipsToPlanet(ctx, targetMoon, input.Ships)

	return nil
}

func (s *MoonService) calculateJumpFuel(ships map[string]int) int64 {
	fuelPerShip := map[string]int{
		"202": 50,
		"203": 125,
		"204": 75,
		"205": 150,
		"206": 400,
		"207": 1250,
		"208": 3000,
		"209": 225,
		"210": 1,
		"211": 1500,
		"213": 2000,
		"214": 10000,
		"215": 375,
		"218": 750,
		"219": 1000,
	}

	var totalFuel int
	for shipType, count := range ships {
		if fuel, ok := fuelPerShip[shipType]; ok {
			totalFuel += fuel * count
		}
	}

	return int64(totalFuel)
}

func (s *MoonService) addShipsToPlanet(ctx context.Context, planet *schema.Planet, ships map[string]int) {
	for shipType, count := range ships {
		var shipID int
		json.Unmarshal([]byte(shipType), &shipID)

		switch shipID {
		case 202:
			planet.SmallCargo += count
		case 203:
			planet.LargeCargo += count
		case 204:
			planet.LightFighter += count
		case 205:
			planet.HeavyFighter += count
		case 206:
			planet.Cruiser += count
		case 207:
			planet.Battleship += count
		case 208:
			planet.ColonyShip += count
		case 209:
			planet.Recycler += count
		case 210:
			planet.EspionageProbe += count
		case 211:
			planet.Bomber += count
		case 213:
			planet.Destroyer += count
		case 214:
			planet.Deathstar += count
		case 215:
			planet.Battlecruiser += count
		case 217:
			planet.Crawler += count
		case 218:
			planet.Reaper += count
		case 219:
			planet.Pathfinder += count
		}
	}
	_ = s.planetRepo.Update(ctx, planet)
}

type PhalanxScanResult struct {
	DetectedFleets []DetectedFleet `json:"detected_fleets"`
	ScannedAt      time.Time       `json:"scanned_at"`
}

type DetectedFleet struct {
	Origin      string    `json:"origin"`
	Target      string    `json:"target"`
	ArrivalTime time.Time `json:"arrival_time"`
	Count       int       `json:"ship_count"`
}

func (s *MoonService) ScanWithPhalanx(ctx context.Context, moonID uint, userID uint, targetGalaxy, targetSystem, targetPosition int) (*PhalanxScanResult, error) {
	moon, err := s.planetRepo.GetByID(ctx, moonID)
	if err != nil {
		return nil, errors.New("moon not found")
	}

	if moon.UserID != userID {
		return nil, errors.New("moon does not belong to user")
	}

	if moon.SensorPhalanx <= 0 {
		return nil, errors.New("sensor phalanx not built")
	}

	scanRange := s.calculatePhalanxRange(moon.SensorPhalanx)
	originDistance := calculateDistance(moon.Galaxy, moon.System, moon.Position, moon.Galaxy, moon.System, moon.Position)
	if originDistance > int64(scanRange) {
		return nil, errors.New("target outside phalanx range")
	}

	fleets, err := s.fleetRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	var detected []DetectedFleet
	for _, fleet := range fleets {
		if fleet.TargetGalaxy == targetGalaxy &&
			fleet.TargetSystem == targetSystem &&
			fleet.TargetPosition == targetPosition {

			shipCount := 0
			var ships map[string]int
			json.Unmarshal([]byte(fleet.Ships), &ships)
			for _, count := range ships {
				shipCount += count
			}

			detected = append(detected, DetectedFleet{
				Origin:      formatCoords(fleet.OriginGalaxy, fleet.OriginSystem, fleet.OriginPosition),
				Target:      formatCoords(fleet.TargetGalaxy, fleet.TargetSystem, fleet.TargetPosition),
				ArrivalTime: fleet.ArrivalTime,
				Count:       shipCount,
			})
		}
	}

	return &PhalanxScanResult{
		DetectedFleets: detected,
		ScannedAt:      time.Now(),
	}, nil
}

func (s *MoonService) calculatePhalanxRange(phalanxLevel int) int {
	baseRange := 1
	additionalRange := (phalanxLevel - 1) * 1
	return baseRange + additionalRange
}

func calculateDistance(galaxy1, system1, pos1, galaxy2, system2, pos2 int) int64 {
	galaxyDiff := int64(abs(galaxy1 - galaxy2))
	systemDiff := int64(abs(system1 - system2)) * 2700
	positionDiff := int64(abs(pos1 - pos2))

	return galaxyDiff + systemDiff + positionDiff
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func formatCoords(galaxy, system, position int) string {
	return "[{galaxy}:{system}:{position}]"
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
