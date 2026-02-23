package service

import (
	"context"

	"ogamex-go/internal/repository"
)

type PlanetService struct {
	planetRepo repository.PlanetRepository
}

func NewPlanetService(planetRepo repository.PlanetRepository) *PlanetService {
	return &PlanetService{
		planetRepo: planetRepo,
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
