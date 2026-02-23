package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type NPCService struct {
	userRepo    repository.UserRepository
	planetRepo  repository.PlanetRepository
	fleetRepo   repository.FleetMissionRepository
}

func NewNPCService(
	userRepo repository.UserRepository,
	planetRepo repository.PlanetRepository,
	fleetRepo repository.FleetMissionRepository,
) *NPCService {
	return &NPCService{
		userRepo:   userRepo,
		planetRepo: planetRepo,
		fleetRepo:  fleetRepo,
	}
}

type NPCPlanetConfig struct {
	Galaxy       int
	System       int
	Position     int
	Level        int
	Resources    int64
	DefenseLevel int
}

func (s *NPCService) CreateNPCPlanet(ctx context.Context, config NPCPlanetConfig) (*schema.Planet, error) {
	npcUser, err := s.getOrCreateNPCUser(ctx, config.Galaxy, config.System)
	if err != nil {
		return nil, err
	}

	existingPlanet, _ := s.planetRepo.GetByCoords(ctx, npcUser.ID, config.Galaxy, config.System, config.Position)
	if existingPlanet != nil {
		return existingPlanet, nil
	}

	metal := config.Resources + int64(rand.Intn(10000))
	crystal := config.Resources/2 + int64(rand.Intn(5000))
	deuterium := config.Resources/4 + int64(rand.Intn(2500))

	planet := &schema.Planet{
		UserID:   npcUser.ID,
		Name:     fmt.Sprintf("Planet-%d:%d:%d", config.Galaxy, config.System, config.Position),
		Galaxy:   config.Galaxy,
		System:   config.System,
		Position: config.Position,
		PlanetType: 1,

		Metal:     metal,
		Crystal:   crystal,
		Deuterium: deuterium,

		MetalCapacity:     100000,
		CrystalCapacity:   100000,
		DeuteriumCapacity: 100000,

		FieldsUsed: 0,
		FieldsMax:  10 + config.Level,

		TempMin: -50 + config.Level*5,
		TempMax:  0 + config.Level*5,
	}

	planet.MetalMine = config.Level
	planet.CrystalMine = config.Level
	planet.DeuteriumSynthesizer = config.Level
	planet.SolarPlant = config.Level
	planet.FusionPlant = config.Level / 2

	if config.DefenseLevel > 0 {
		planet.RocketLauncher = config.DefenseLevel * 10
		planet.LightLaser = config.DefenseLevel * 5
		planet.HeavyLaser = config.DefenseLevel * 2
	}

	err = s.planetRepo.Create(ctx, planet)
	if err != nil {
		return nil, err
	}

	return planet, nil
}

func (s *NPCService) getOrCreateNPCUser(ctx context.Context, galaxy, system int) (*schema.User, error) {
	username := fmt.Sprintf("npc_%d_%d", galaxy, system)
	
	npcUser, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil && npcUser != nil {
		return npcUser, nil
	}

	email := fmt.Sprintf("npc_%d_%d@ogamex.local", galaxy, system)
	password := make([]byte, 16)
	rand.Read(password)
	hashedPassword := hex.EncodeToString(password)

	npcUser = &schema.User{
		Username:       username,
		Email:          email,
		Password:       hashedPassword,
		PlayerName:     fmt.Sprintf("NPC %d:%d", galaxy, system),
		CharacterClass: 0,
		DarkMatter:     0,
		IsNPC:          true,
		RegisteredAt:   time.Now(),
		LastOnline:     time.Now(),
	}

	err = s.userRepo.Create(ctx, npcUser)
	if err != nil {
		return nil, err
	}

	planet := &schema.User{
		ID: npcUser.ID,
	}
	err = s.userRepo.Update(ctx, planet)
	if err != nil {
		return nil, err
	}

	return npcUser, nil
}

type NPCFleetConfig struct {
	OriginGalaxy   int
	OriginSystem   int
	OriginPosition int
	TargetGalaxy   int
	TargetSystem   int
	TargetPosition int
	MissionType   int
	Ships          map[int]int64
}

func (s *NPCService) CreateNPCFleet(ctx context.Context, config NPCFleetConfig) (*schema.FleetMission, error) {
	npcUser, err := s.getOrCreateNPCUser(ctx, config.OriginGalaxy, config.OriginSystem)
	if err != nil {
		return nil, err
	}

	fleet := &schema.FleetMission{
		UserID:            npcUser.ID,
		MissionType:       config.MissionType,
		OriginGalaxy:     config.OriginGalaxy,
		OriginSystem:     config.OriginSystem,
		OriginPosition:   config.OriginPosition,
		OriginPlanetType: 1,
		TargetGalaxy:     config.TargetGalaxy,
		TargetSystem:     config.TargetSystem,
		TargetPosition:   config.TargetPosition,
		TargetPlanetType: 1,
		LaunchTime:       time.Now(),
		ArrivalTime:      time.Now().Add(30 * time.Minute),
		ReturnTime:       time.Now().Add(60 * time.Minute),
		Metal:            0,
		Crystal:          0,
		Deuterium:        0,
		Status:           0,
		Ships:            s.shipsToJSON(config.Ships),
	}

	err = s.fleetRepo.Create(ctx, fleet)
	if err != nil {
		return nil, err
	}

	return fleet, nil
}

func (s *NPCService) GenerateExpeditionFleet(ctx context.Context, galaxy, system, position int) (*schema.FleetMission, error) {
	shipOptions := []map[int]int64{
		{210: 1, 219: 2},
		{202: 5, 203: 3},
		{204: 1, 205: 2, 206: 1},
		{207: 1, 215: 3},
		{211: 1},
		{213: 1, 214: 1},
	}

	randIdx := rand.Intn(len(shipOptions))
	ships := shipOptions[randIdx]

	coords := []struct{ galaxy, system, position int }{
		{galaxy + rand.Intn(9) - 4, system + rand.Intn(9) - 4, rand.Intn(15) + 1},
		{galaxy + rand.Intn(9) - 4, system + rand.Intn(9) - 4, rand.Intn(15) + 1},
		{galaxy + rand.Intn(9) - 4, system + rand.Intn(9) - 4, rand.Intn(15) + 1},
	}
	target := coords[rand.Intn(len(coords))]

	return s.CreateNPCFleet(ctx, NPCFleetConfig{
		OriginGalaxy:   galaxy,
		OriginSystem:   system,
		OriginPosition: position,
		TargetGalaxy:   target.galaxy,
		TargetSystem:   target.system,
		TargetPosition: target.position,
		MissionType:    7,
		Ships:          ships,
	})
}

func (s *NPCService) GenerateDefenseWave(ctx context.Context, galaxy, system, position int, strength int) (*schema.FleetMission, error) {
	shipCount := int64(strength)
	
	ships := map[int]int64{
		202: shipCount,
		203: shipCount / 2,
		205: shipCount / 4,
	}

	return s.CreateNPCFleet(ctx, NPCFleetConfig{
		OriginGalaxy:   galaxy,
		OriginSystem:   system,
		OriginPosition: position,
		TargetGalaxy:   galaxy,
		TargetSystem:   system,
		TargetPosition: position,
		MissionType:    4,
		Ships:          ships,
	})
}

func (s *NPCService) GeneratePirateRaid(ctx context.Context, targetGalaxy, targetSystem, targetPosition int) (*schema.FleetMission, error) {
	pirateGalaxy := 1
	pirateSystem := rand.Intn(499) + 1

	ships := map[int]int64{
		202: int64(50 + rand.Intn(100)),
		203: int64(20 + rand.Intn(50)),
		204: int64(5 + rand.Intn(20)),
		205: int64(2 + rand.Intn(10)),
	}

	return s.CreateNPCFleet(ctx, NPCFleetConfig{
		OriginGalaxy:   pirateGalaxy,
		OriginSystem:   pirateSystem,
		OriginPosition: rand.Intn(15) + 1,
		TargetGalaxy:   targetGalaxy,
		TargetSystem:   targetSystem,
		TargetPosition: targetPosition,
		MissionType:    1,
		Ships:          ships,
	})
}

func (s *NPCService) GetNPCPlanets(ctx context.Context) ([]*schema.Planet, error) {
	users, err := s.userRepo.GetNPCUsers(ctx)
	if err != nil {
		return nil, err
	}

	var planets []*schema.Planet
	for _, user := range users {
		userPlanets, err := s.planetRepo.GetByUserID(ctx, user.ID)
		if err != nil {
			continue
		}
		planets = append(planets, userPlanets...)
	}

	return planets, nil
}

func (s *NPCService) GetNPCPlanet(ctx context.Context, galaxy, system, position int) (*schema.Planet, error) {
	users, err := s.userRepo.GetNPCUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		planet, err := s.planetRepo.GetByCoords(ctx, user.ID, galaxy, system, position)
		if err == nil && planet != nil {
			return planet, nil
		}
	}

	return nil, fmt.Errorf("npc planet not found at %d:%d:%d", galaxy, system, position)
}

func (s *NPCService) shipsToJSON(ships map[int]int64) string {
	result := ""
	for id, count := range ships {
		if result != "" {
			result += ";"
		}
		result += fmt.Sprintf("%d:%d", id, count)
	}
	return result
}

func (s *NPCService) LootResources(ctx context.Context, planet *schema.Planet, lootPercentage float64) (metal, crystal, deuterium int64) {
	metal = int64(float64(planet.Metal) * lootPercentage)
	crystal = int64(float64(planet.Crystal) * lootPercentage)
	deuterium = int64(float64(planet.Deuterium) * lootPercentage)

	planet.Metal -= metal
	planet.Crystal -= crystal
	planet.Deuterium -= deuterium

	s.planetRepo.Update(ctx, planet)

	return metal, crystal, deuterium
}
