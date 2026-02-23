package service

import (
	"context"
	"encoding/json"
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
	ErrBattleEngineFailed  = errors.New("battle engine failed")
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
	shipsJSON, _ := json.Marshal(params.Ships)
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
		Ships:             string(shipsJSON),
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
	case 15: // Expedition
		return s.processExpedition(ctx, mission)
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
	targetPlanet, err := s.planetRepo.GetByCoords(ctx, 0, mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition)
	if err != nil {
		return err
	}

	attackerShips := make(map[int16]int16)
	if err := json.Unmarshal([]byte(mission.Ships), &attackerShips); err != nil {
		return err
	}

	defenderFleets, err := s.fleetRepo.GetActive(ctx)
	if err != nil {
		return err
	}

	hasDefender := false
	defenderShips := make(map[int16]int16)
	for _, fleet := range defenderFleets {
		if fleet.TargetGalaxy == mission.TargetGalaxy &&
			fleet.TargetSystem == mission.TargetSystem &&
			fleet.TargetPosition == mission.TargetPosition &&
			fleet.MissionType == 1 {

			var fleetShips map[int16]int16
			if err := json.Unmarshal([]byte(fleet.Ships), &fleetShips); err != nil {
				continue
			}
			for shipID, amount := range fleetShips {
				defenderShips[shipID] += amount
			}
			hasDefender = true
		}
	}

	if targetPlanet != nil && targetPlanet.UserID != 0 {
		hasDefender = true
		shipsOnPlanet := map[int16]int16{
			202: int16(targetPlanet.SmallCargo),
			203: int16(targetPlanet.LargeCargo),
			204: int16(targetPlanet.LightFighter),
			205: int16(targetPlanet.HeavyFighter),
			206: int16(targetPlanet.Cruiser),
			207: int16(targetPlanet.Battleship),
			208: int16(targetPlanet.ColonyShip),
			209: int16(targetPlanet.Recycler),
			210: int16(targetPlanet.EspionageProbe),
			211: int16(targetPlanet.Bomber),
			213: int16(targetPlanet.Destroyer),
			214: int16(targetPlanet.Deathstar),
			215: int16(targetPlanet.Battlecruiser),
			218: int16(targetPlanet.Reaper),
			219: int16(targetPlanet.Pathfinder),
			401: int16(targetPlanet.RocketLauncher),
			402: int16(targetPlanet.LightLaser),
			403: int16(targetPlanet.HeavyLaser),
			404: int16(targetPlanet.IonCannon),
			405: int16(targetPlanet.GaussCannon),
			406: int16(targetPlanet.PlasmaTurret),
			407: int16(targetPlanet.ShieldDome),
			408: int16(targetPlanet.MissileInterceptor),
			409: int16(targetPlanet.MissileLauncher),
		}
		for shipID, amount := range shipsOnPlanet {
			defenderShips[shipID] += amount
		}
	}

	if !hasDefender {
		return s.processUndefendedAttack(ctx, mission, targetPlanet)
	}

	attackerTech, _ := s.techRepo.GetByUserID(ctx, mission.UserID)
	attackerStats := s.getFleetStats(attackerShips, attackerTech)

	defenderUserID := targetPlanet.UserID
	defenderTech, _ := s.techRepo.GetByUserID(ctx, defenderUserID)
	defenderStats := s.getFleetStats(defenderShips, defenderTech)

	battleResult := s.simulateBattle(attackerStats, defenderStats)

	attackerLosses := battleResult.AttackerLosses
	defenderLosses := battleResult.DefenderLosses
	totalDestroyed := s.calculateDestroyedResources(attackerLosses, defenderLosses)

	moonChance := s.calculateMoonChance(totalDestroyed)

	if battleResult.Winner == "attacker" {
		loot := s.calculateLoot(targetPlanet)
		_ = s.planetRepo.AddResources(ctx, mission.UserID, loot.Metal, loot.Crystal, loot.Deuterium)

		for shipID, count := range battleResult.AttackerRemaining {
			s.removeShipsFromPlanet(ctx, targetPlanet, int(shipID), int(count))
		}

		_ = s.updateFleetAfterBattle(ctx, mission, battleResult.AttackerRemaining)
		mission.Status = 2

		existingPlanet, _ := s.planetRepo.GetByCoordsAny(ctx, mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition)
		if existingPlanet == nil && moonChance > 0 {
			s.createMoon(ctx, mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition, moonChance)
		}
	} else {
		_ = s.updateFleetAfterBattle(ctx, mission, battleResult.DefenderRemaining)
		mission.Status = 2
	}

	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) calculateDestroyedResources(attackerLosses, defenderLosses map[int16]int16) int64 {
	costs := map[int16]int64{
		202: 4000, 203: 12000, 204: 18000, 205: 45000, 206: 40000,
		207: 90000, 208: 20000, 209: 18000, 210: 1000, 211: 125000,
		212: 2500, 213: 20000, 214: 250000, 215: 4000, 218: 8000,
		219: 40000, 401: 2000, 402: 2000, 403: 8000, 404: 8000,
		405: 35000, 406: 130000, 407: 20000, 408: 10000, 409: 20000,
	}

	var total int64 = 0
	for shipID, amount := range attackerLosses {
		if cost, ok := costs[shipID]; ok {
			total += cost * int64(amount)
		}
	}
	for shipID, amount := range defenderLosses {
		if cost, ok := costs[shipID]; ok {
			total += cost * int64(amount)
		}
	}

	return total
}

func (s *FleetService) calculateMoonChance(destroyedResources int64) int {
	if destroyedResources < 100000 {
		return 0
	}

	baseChance := int((destroyedResources - 100000) / 100000)
	if baseChance > 25 {
		baseChance = 25
	}

	return baseChance
}

func (s *FleetService) createMoon(ctx context.Context, galaxy, system, position int, chance int) bool {
	if chance <= 0 {
		return false
	}

	random := int64(0) % 100
	if random < int64(chance) {
		moon := &schema.Planet{
			UserID:      0,
			Name:        "Moon",
			Galaxy:      galaxy,
			System:      system,
			Position:    position,
			IsMoon:      true,
			PlanetType:  3,
			Metal:       0,
			Crystal:     0,
			Deuterium:   0,
			MetalCapacity:     0,
			CrystalCapacity:   0,
			DeuteriumCapacity: 0,
			MetalProduction:     0,
			CrystalProduction:   0,
			DeuteriumProduction: 0,
			EnergyAvailable:     0,
			EnergyMax:          0,
			EnergyUsed:          0,
			TempMin:      -50,
			TempMax:      50,
			FieldsUsed:   0,
			FieldsMax:    0,
			LunarBase:   1,
		}

		_ = s.planetRepo.Create(ctx, moon)
		return true
	}

	return false
}

type FleetStats struct {
	Units map[int16]struct {
		AttackPower  float32
		ShieldPoints float32
		HullPlating  float32
	}
}

func (s *FleetService) getFleetStats(ships map[int16]int16, tech *schema.UserTech) FleetStats {
	baseStats := map[int16]struct {
		AttackPower  float32
		ShieldPoints float32
		HullPlating  float32
	}{
		202: {AttackPower: 5, ShieldPoints: 10, HullPlating: 40},
		203: {AttackPower: 5, ShieldPoints: 25, HullPlating: 120},
		204: {AttackPower: 50, ShieldPoints: 10, HullPlating: 100},
		205: {AttackPower: 150, ShieldPoints: 25, HullPlating: 250},
		206: {AttackPower: 400, ShieldPoints: 50, HullPlating: 800},
		207: {AttackPower: 1000, ShieldPoints: 200, HullPlating: 2000},
		208: {AttackPower: 50, ShieldPoints: 100, HullPlating: 500},
		209: {AttackPower: 20, ShieldPoints: 10, HullPlating: 100},
		210: {AttackPower: 0.1, ShieldPoints: 0, HullPlating: 5},
		211: {AttackPower: 700, ShieldPoints: 150, HullPlating: 1700},
		212: {AttackPower: 1, ShieldPoints: 10, HullPlating: 20},
		213: {AttackPower: 1100, ShieldPoints: 300, HullPlating: 2200},
		214: {AttackPower: 200000, ShieldPoints: 50000, HullPlating: 700000},
		215: {AttackPower: 700, ShieldPoints: 100, HullPlating: 1500},
		218: {AttackPower: 1300, ShieldPoints: 200, HullPlating: 2800},
		219: {AttackPower: 300, ShieldPoints: 50, HullPlating: 600},
		401: {AttackPower: 80, ShieldPoints: 20, HullPlating: 200},
		402: {AttackPower: 100, ShieldPoints: 25, HullPlating: 200},
		403: {AttackPower: 250, ShieldPoints: 100, HullPlating: 500},
		404: {AttackPower: 150, ShieldPoints: 80, HullPlating: 400},
		405: {AttackPower: 1100, ShieldPoints: 300, HullPlating: 2200},
		406: {AttackPower: 3000, ShieldPoints: 500, HullPlating: 5000},
		407: {AttackPower: 1, ShieldPoints: 5000, HullPlating: 100},
		408: {AttackPower: 80, ShieldPoints: 20, HullPlating: 200},
		409: {AttackPower: 150, ShieldPoints: 40, HullPlating: 300},
	}

	stats := FleetStats{Units: make(map[int16]struct {
		AttackPower  float32
		ShieldPoints float32
		HullPlating  float32
	})}

	weaponsBonus := float32(1.0)
	shieldBonus := float32(1.0)
	armorBonus := float32(1.0)
	if tech != nil {
		weaponsBonus = float32(1.0 + float64(tech.WeaponsTechnology)*0.1)
		shieldBonus = float32(1.0 + float64(tech.ShieldingTechnology)*0.1)
		armorBonus = float32(1.0 + float64(tech.ArmorTechnology)*0.1)
	}

	for shipID, amount := range ships {
		if base, ok := baseStats[shipID]; ok {
			stats.Units[shipID] = struct {
				AttackPower  float32
				ShieldPoints float32
				HullPlating  float32
			}{
				AttackPower:  float32(float32(base.AttackPower) * weaponsBonus * float32(amount)),
				ShieldPoints: float32(float32(base.ShieldPoints) * shieldBonus * float32(amount)),
				HullPlating:  float32(float32(base.HullPlating) * armorBonus * float32(amount)),
			}
		}
	}

	return stats
}

type BattleResult struct {
	Winner            string         `json:"winner"`
	AttackerRemaining map[int16]int16 `json:"attacker_remaining"`
	DefenderRemaining map[int16]int16 `json:"defender_remaining"`
	AttackerLosses    map[int16]int16 `json:"attacker_losses"`
	DefenderLosses    map[int16]int16 `json:"defender_losses"`
}

func (s *FleetService) simulateBattle(attackerStats FleetStats, defenderStats FleetStats) BattleResult {
	var attackerShield, attackerHull, attackerPower float32
	var defenderShield, defenderHull, defenderPower float32

	for _, u := range attackerStats.Units {
		attackerShield += u.ShieldPoints
		attackerHull += u.HullPlating
		attackerPower += u.AttackPower
	}
	for _, u := range defenderStats.Units {
		defenderShield += u.ShieldPoints
		defenderHull += u.HullPlating
		defenderPower += u.AttackPower
	}

	attackerRemaining := make(map[int16]int16)
	defenderRemaining := make(map[int16]int16)
	for k := range attackerStats.Units {
		attackerRemaining[k] = 1
	}
	for k := range defenderStats.Units {
		defenderRemaining[k] = 1
	}

	rounds := 6
	for r := 0; r < rounds; r++ {
		attackerDamage := attackerPower
		defenderDamage := defenderPower

		defenderShield -= attackerDamage
		if defenderShield < 0 {
			defenderHull += defenderShield
			defenderShield = 0
		}

		attackerShield -= defenderDamage
		if attackerShield < 0 {
			attackerHull += attackerShield
			attackerShield = 0
		}

		if attackerHull <= 0 {
			return BattleResult{
				Winner:            "defender",
				AttackerRemaining: map[int16]int16{},
				DefenderRemaining: defenderRemaining,
			}
		}
		if defenderHull <= 0 {
			return BattleResult{
				Winner:            "attacker",
				AttackerRemaining: attackerRemaining,
				DefenderRemaining: map[int16]int16{},
			}
		}
	}

	lossFactor := attackerHull / (attackerHull + defenderHull)
	attackerRet := make(map[int16]int16)
	defenderRet := make(map[int16]int16)

	for k := range attackerStats.Units {
		attackerRet[k] = int16(float32(1) * (lossFactor + 0.5))
		if attackerRet[k] < 1 {
			attackerRet[k] = 1
		}
	}
	for k := range defenderStats.Units {
		defenderRet[k] = int16(float32(1) * (1 - lossFactor + 0.5))
		if defenderRet[k] < 1 {
			defenderRet[k] = 1
		}
	}

	winner := "draw"
	if attackerHull > defenderHull {
		winner = "attacker"
	} else if defenderHull > attackerHull {
		winner = "defender"
	}

	return BattleResult{
		Winner:            winner,
		AttackerRemaining: attackerRet,
		DefenderRemaining: defenderRet,
	}
}

func (s *FleetService) processUndefendedAttack(ctx context.Context, mission *schema.FleetMission, targetPlanet *schema.Planet) error {
	if targetPlanet == nil {
		mission.Status = 2
		return s.fleetRepo.Update(ctx, mission)
	}

	loot := s.calculateLoot(targetPlanet)
	_ = s.planetRepo.AddResources(ctx, mission.UserID, loot.Metal, loot.Crystal, loot.Deuterium)

	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) calculateLoot(planet *schema.Planet) struct{ Metal, Crystal, Deuterium int64 } {
	maxLoot := (planet.Metal + planet.Crystal + planet.Deuterium) * 50 / 100
	if maxLoot <= 0 {
		return struct{ Metal, Crystal, Deuterium int64 }{0, 0, 0}
	}

	metalRatio := float64(planet.Metal) / float64(planet.Metal+planet.Crystal+planet.Deuterium)
	crystalRatio := float64(planet.Crystal) / float64(planet.Metal+planet.Crystal+planet.Deuterium)
	deutRatio := float64(planet.Deuterium) / float64(planet.Metal+planet.Crystal+planet.Deuterium)

	return struct{ Metal, Crystal, Deuterium int64 }{
		Metal:     int64(float64(maxLoot) * metalRatio),
		Crystal:   int64(float64(maxLoot) * crystalRatio),
		Deuterium: int64(float64(maxLoot) * deutRatio),
	}
}

func (s *FleetService) removeShipsFromPlanet(ctx context.Context, planet *schema.Planet, shipID, amount int) {
	switch shipID {
	case 202:
		planet.SmallCargo -= amount
	case 203:
		planet.LargeCargo -= amount
	case 204:
		planet.LightFighter -= amount
	case 205:
		planet.HeavyFighter -= amount
	case 206:
		planet.Cruiser -= amount
	case 207:
		planet.Battleship -= amount
	case 209:
		planet.Recycler -= amount
	case 211:
		planet.Bomber -= amount
	case 213:
		planet.Destroyer -= amount
	case 215:
		planet.Battlecruiser -= amount
	case 218:
		planet.Reaper -= amount
	case 219:
		planet.Pathfinder -= amount
	case 401:
		planet.RocketLauncher -= amount
	case 402:
		planet.LightLaser -= amount
	case 403:
		planet.HeavyLaser -= amount
	case 404:
		planet.IonCannon -= amount
	case 405:
		planet.GaussCannon -= amount
	case 406:
		planet.PlasmaTurret -= amount
	case 408:
		planet.MissileInterceptor -= amount
	case 409:
		planet.MissileLauncher -= amount
	}
	_ = s.planetRepo.Update(ctx, planet)
}

func (s *FleetService) updateFleetAfterBattle(ctx context.Context, mission *schema.FleetMission, remaining map[int16]int16) error {
	if len(remaining) == 0 {
		return s.fleetRepo.Delete(ctx, mission.ID)
	}
	mission.Ships = func() string {
		b, _ := json.Marshal(remaining)
		return string(b)
	}()
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) processHarvest(ctx context.Context, mission *schema.FleetMission) error {
	_ = s.updateFleetAfterBattle(ctx, mission, map[int16]int16{})
	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) processColonize(ctx context.Context, mission *schema.FleetMission) error {
	targetPlanet, err := s.planetRepo.GetByCoords(ctx, 0, mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition)
	if err == nil && targetPlanet != nil {
		mission.Status = 2
		return s.fleetRepo.Update(ctx, mission)
	}

	ships := make(map[int16]int16)
	if err := json.Unmarshal([]byte(mission.Ships), &ships); err != nil {
		return err
	}

	colonyShipCount := ships[208]
	if colonyShipCount == 0 {
		mission.Status = 2
		return s.fleetRepo.Update(ctx, mission)
	}

	planetName := colonizeName(mission.TargetGalaxy, mission.TargetSystem, mission.TargetPosition)

	newPlanet := &schema.Planet{
		UserID:      mission.UserID,
		Name:        planetName,
		Galaxy:      mission.TargetGalaxy,
		System:      mission.TargetSystem,
		Position:    mission.TargetPosition,
		IsMoon:      false,
		PlanetType:  1,
		Metal:       500,
		Crystal:     250,
		Deuterium:   100,
		MetalCapacity:     10000,
		CrystalCapacity:   10000,
		DeuteriumCapacity: 10000,
		MetalProduction:     0,
		CrystalProduction:   0,
		DeuteriumProduction: 0,
		EnergyAvailable:     0,
		EnergyMax:          0,
		EnergyUsed:          0,
		TempMin:      -50,
		TempMax:      50,
		FieldsUsed:   0,
		FieldsMax:    163,
		MetalMine:           1,
		CrystalMine:         1,
		DeuteriumSynthesizer: 1,
		SolarPlant:          1,
		MetalStorageBuilding:        1,
		CrystalStorageBuilding:      1,
		DeuteriumStorageBuilding:    1,
		RobotFactory:      1,
		Shipyard:          1,
		SmallCargo:        5,
	}

	_ = s.planetRepo.Create(ctx, newPlanet)

	_ = s.updateFleetAfterBattle(ctx, mission, map[int16]int16{})
	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

func colonizeName(galaxy, system, position int) string {
	return "Planet"
}

func (s *FleetService) processRecycle(ctx context.Context, mission *schema.FleetMission) error {
	_ = s.updateFleetAfterBattle(ctx, mission, map[int16]int16{})
	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

func (s *FleetService) processExpedition(ctx context.Context, mission *schema.FleetMission) error {
	expeditionResult := s.calculateExpeditionResult()

	ships := make(map[int16]int16)
	if err := json.Unmarshal([]byte(mission.Ships), &ships); err != nil {
		return err
	}

	remainingShips := make(map[int16]int16)
	for shipID, amount := range ships {
		if amount > 0 {
			remainingShips[shipID] = amount
		}
	}

	switch expeditionResult.outcome {
	case "resources":
		_ = s.planetRepo.AddResources(ctx, mission.UserID, expeditionResult.metal, expeditionResult.crystal, expeditionResult.deuterium)
	case "dark_matter":
		user, err := s.userRepo.GetByID(ctx, mission.UserID)
		if err == nil {
			user.DarkMatter += int64(expeditionResult.darkMatter)
			_ = s.userRepo.Update(ctx, user)
		}
	case "nothing":
	case "lost":
		remainingShips = map[int16]int16{}
	}

	if len(remainingShips) > 0 {
		_ = s.updateFleetAfterBattle(ctx, mission, remainingShips)
	} else {
		_ = s.updateFleetAfterBattle(ctx, mission, map[int16]int16{})
	}

	mission.Status = 2
	return s.fleetRepo.Update(ctx, mission)
}

type expeditionOutcome struct {
	outcome     string
	metal       int64
	crystal     int64
	deuterium  int64
	darkMatter  int
}

func (s *FleetService) calculateExpeditionResult() expeditionOutcome {
	roll := int64(0) % 100

	if roll < 30 {
		return expeditionOutcome{outcome: "nothing"}
	} else if roll < 70 {
		metal := int64((int64(0) % 5000) + 500)
		crystal := int64((int64(0) % 3000) + 300)
		deuterium := int64((int64(0) % 1000) + 100)
		return expeditionOutcome{
			outcome:    "resources",
			metal:      metal,
			crystal:    crystal,
			deuterium: deuterium,
		}
	} else if roll < 90 {
		darkMatter := (int64(0) % 200) + 50
		return expeditionOutcome{
			outcome:    "dark_matter",
			darkMatter: int(darkMatter),
		}
	} else {
		return expeditionOutcome{outcome: "lost"}
	}
}
