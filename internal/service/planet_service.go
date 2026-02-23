package service

import (
	"context"
	"errors"
	"fmt"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type PlanetService struct {
	planetRepo repository.PlanetRepository
	userRepo   repository.UserRepository
}

func NewPlanetService(planetRepo repository.PlanetRepository, userRepo repository.UserRepository) *PlanetService {
	return &PlanetService{
		planetRepo: planetRepo,
		userRepo:   userRepo,
	}
}

type GalaxyPosition struct {
	Galaxy   int    `json:"galaxy"`
	System   int    `json:"system"`
	Position int    `json:"position"`
	Planet   *struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		UserID     *uint  `json:"user_id"`
		UserName   string `json:"user_name,omitempty"`
		IsMoon     bool   `json:"is_moon"`
		PlanetType int    `json:"planet_type"`
	} `json:"planet,omitempty"`
	Debris struct {
		Metal     int64 `json:"metal"`
		Crystal   int64 `json:"crystal"`
	} `json:"debris,omitempty"`
}

func (s *PlanetService) GetGalaxy(ctx context.Context, galaxy, system int) ([]*GalaxyPosition, error) {
	positions := make([]*GalaxyPosition, 15)

	for i := 1; i <= 15; i++ {
		pos := &GalaxyPosition{
			Galaxy:   galaxy,
			System:   system,
			Position: i,
		}

		planet, err := s.planetRepo.GetByCoordsAny(ctx, galaxy, system, i)
		if err == nil && planet != nil {
			planetInfo := struct {
				ID         uint   `json:"id"`
				Name       string `json:"name"`
				UserID     *uint  `json:"user_id"`
				UserName   string `json:"user_name,omitempty"`
				IsMoon     bool   `json:"is_moon"`
				PlanetType int    `json:"planet_type"`
			}{
				ID:         planet.ID,
				Name:       planet.Name,
				UserID:     &planet.UserID,
				IsMoon:     planet.IsMoon,
				PlanetType: planet.PlanetType,
			}
			pos.Planet = &planetInfo
		}

		positions[i-1] = pos
	}

	return positions, nil
}

func (s *PlanetService) GetSystemInfo(ctx context.Context, galaxy, system int) (map[string]interface{}, error) {
	planets, err := s.planetRepo.GetBySystem(ctx, galaxy, system)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"galaxy":   galaxy,
		"system":   system,
		"planets":  planets,
		"max_positions": 15,
	}

	return result, nil
}

func (s *PlanetService) SetDefenseActivation(ctx context.Context, planetID uint, userID uint, activate bool) error {
	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return err
	}

	if planet.UserID != userID {
		return errors.New("planet does not belong to user")
	}

	planet.DefenseActivated = activate
	return s.planetRepo.Update(ctx, planet)
}

type ColonizeInput struct {
	Galaxy   int
	System   int
	Position int
	Name     string
}

func (s *PlanetService) ColonizePlanet(ctx context.Context, userID uint, input ColonizeInput) (*schema.Planet, error) {
	existing, _ := s.planetRepo.GetByCoordsAny(ctx, input.Galaxy, input.System, input.Position)
	if existing != nil {
		return nil, errors.New("position already occupied")
	}

	planet := &schema.Planet{
		UserID:       userID,
		Name:         input.Name,
		Galaxy:       input.Galaxy,
		System:       input.System,
		Position:     input.Position,
		IsMoon:       false,
		PlanetType:   1,
		Metal:        500,
		Crystal:      500,
		Deuterium:    0,
		MetalCapacity:     100000,
		CrystalCapacity:   100000,
		DeuteriumCapacity: 100000,
		MetalProduction:     0,
		CrystalProduction:   0,
		DeuteriumProduction: 0,
		EnergyAvailable: 0,
		EnergyMax:       0,
		EnergyUsed:      0,
		TempMin:      0,
		TempMax:      0,
		FieldsUsed:   0,
		FieldsMax:    163,
	}

	err := s.planetRepo.Create(ctx, planet)
	if err != nil {
		return nil, err
	}

	return planet, nil
}

type HighscoreEntry struct {
	Rank       int             `json:"rank"`
	UserID     uint            `json:"user_id"`
	Username   string          `json:"username"`
	PlayerName string          `json:"player_name"`
	Points     int64           `json:"points"`
	Planets    int             `json:"planets"`
	Fleets     int             `json:"fleets"`
	Research   int             `json:"research"`
	Defense    int64           `json:"defense"`
	Military   int64           `json:"military"`
}

func (s *PlanetService) GetHighscore(ctx context.Context, category string, limit int) ([]*HighscoreEntry, error) {
	if limit <= 0 {
		limit = 100
	}

	planets, err := s.planetRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	userStats := make(map[uint]*HighscoreEntry)
	planetCount := make(map[uint]int)

	for _, p := range planets {
		planetCount[p.UserID]++

		if _, ok := userStats[p.UserID]; !ok {
			user := &schema.User{}
			u, err := s.userRepo.GetByID(ctx, p.UserID)
			if err == nil {
				user = u
			}

			userStats[p.UserID] = &HighscoreEntry{
				UserID:     p.UserID,
				Username:   user.Username,
				PlayerName: user.PlayerName,
				Points:     0,
				Planets:    0,
			}
		}

		var defensePoints int64
		defensePoints += int64(p.RocketLauncher) * 500
		defensePoints += int64(p.LightLaser) * 1000
		defensePoints += int64(p.HeavyLaser) * 2000
		defensePoints += int64(p.IonCannon) * 2000
		defensePoints += int64(p.GaussCannon) * 7000
		defensePoints += int64(p.PlasmaTurret) * 20000
		defensePoints += int64(p.ShieldDome) * 20000
		defensePoints += int64(p.MissileInterceptor) * 2000
		defensePoints += int64(p.MissileLauncher) * 15000

		var shipPoints int64
		shipPoints += int64(p.SmallCargo) * 2000
		shipPoints += int64(p.LargeCargo) * 6000
		shipPoints += int64(p.LightFighter) * 10000
		shipPoints += int64(p.HeavyFighter) * 25000
		shipPoints += int64(p.Cruiser) * 40000
		shipPoints += int64(p.Battleship) * 100000
		shipPoints += int64(p.ColonyShip) * 40000
		shipPoints += int64(p.Recycler) * 30000
		shipPoints += int64(p.EspionageProbe) * 1000
		shipPoints += int64(p.Bomber) * 75000
		shipPoints += int64(p.Destroyer) * 120000
		shipPoints += int64(p.Deathstar) * 2000000
		shipPoints += int64(p.Battlecruiser) * 150000
		shipPoints += int64(p.Reaper) * 140000
		shipPoints += int64(p.Pathfinder) * 65000

		userStats[p.UserID].Defense += defensePoints
		userStats[p.UserID].Military += shipPoints

		userStats[p.UserID].Points += defensePoints + shipPoints
	}

	entries := make([]*HighscoreEntry, 0, len(userStats))
	for _, e := range userStats {
		e.Planets = planetCount[e.UserID]
		e.Fleets = 0
		e.Research = 0
		entries = append(entries, e)
	}

	switch category {
	case "military":
		quickSort(entries, func(a, b *HighscoreEntry) bool { return a.Military > b.Military })
	case "defense":
		quickSort(entries, func(a, b *HighscoreEntry) bool { return a.Defense > b.Defense })
	case "research":
		quickSort(entries, func(a, b *HighscoreEntry) bool { return a.Research > b.Research })
	case "fleets":
		quickSort(entries, func(a, b *HighscoreEntry) bool { return a.Fleets > b.Fleets })
	default:
		quickSort(entries, func(a, b *HighscoreEntry) bool { return a.Points > b.Points })
	}

	if len(entries) > limit {
		entries = entries[:limit]
	}

	for i, e := range entries {
		e.Rank = i + 1
	}

	return entries, nil
}

func quickSort(arr []*HighscoreEntry, less func(a, b *HighscoreEntry) bool) {
	if len(arr) <= 1 {
		return
	}
	pivot := arr[len(arr)/2]
	i, j := 0, len(arr)-1
	for i <= j {
		for i <= j && less(arr[i], pivot) {
			i++
		}
		for i <= j && less(pivot, arr[j]) {
			j--
		}
		if i <= j {
			arr[i], arr[j] = arr[j], arr[i]
			i++
			j--
		}
	}
	quickSort(arr[:j+1], less)
	quickSort(arr[i:], less)
}

type MovePlanetInput struct {
	PlanetID    uint
	TargetGalaxy   int
	TargetSystem   int
	TargetPosition int
}

func (s *PlanetService) MovePlanet(ctx context.Context, userID uint, input MovePlanetInput) (*schema.Planet, error) {
	planet, err := s.planetRepo.GetByID(ctx, input.PlanetID)
	if err != nil {
		return nil, errors.New("planet not found")
	}

	if planet.UserID != userID {
		return nil, errors.New("planet does not belong to user")
	}

	if input.TargetGalaxy < 1 || input.TargetGalaxy > 10 {
		return nil, errors.New("invalid target galaxy (1-10)")
	}

	if input.TargetSystem < 1 || input.TargetSystem > 499 {
		return nil, errors.New("invalid target system (1-499)")
	}

	if input.TargetPosition < 1 || input.TargetPosition > 15 {
		return nil, errors.New("invalid target position (1-15)")
	}

	existing, _ := s.planetRepo.GetByCoordsAny(ctx, input.TargetGalaxy, input.TargetSystem, input.TargetPosition)
	if existing != nil && existing.ID != planet.ID {
		return nil, errors.New("target position already occupied")
	}

	currentPos := fmt.Sprintf("%d:%d:%d", planet.Galaxy, planet.System, planet.Position)
	targetPos := fmt.Sprintf("%d:%d:%d", input.TargetGalaxy, input.TargetSystem, input.TargetPosition)

	if currentPos == targetPos {
		return nil, errors.New("planet already at target position")
	}

	planet.Galaxy = input.TargetGalaxy
	planet.System = input.TargetSystem
	planet.Position = input.TargetPosition

	err = s.planetRepo.Update(ctx, planet)
	if err != nil {
		return nil, err
	}

	return planet, nil
}
