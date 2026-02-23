package service

import (
	"context"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type DebrisService struct {
	debrisRepo repository.DebrisRepository
	wreckRepo  repository.WreckRepository
	planetRepo repository.PlanetRepository
}

func NewDebrisService(
	debrisRepo repository.DebrisRepository,
	wreckRepo repository.WreckRepository,
	planetRepo repository.PlanetRepository,
) *DebrisService {
	return &DebrisService{
		debrisRepo: debrisRepo,
		wreckRepo:  wreckRepo,
		planetRepo: planetRepo,
	}
}

func (s *DebrisService) GetDebrisFields(ctx context.Context) ([]*schema.DebrisField, error) {
	return s.debrisRepo.GetActive(ctx)
}

func (s *DebrisService) GetDebrisField(ctx context.Context, galaxy, system, position int) (*schema.DebrisField, error) {
	return s.debrisRepo.GetByCoords(ctx, galaxy, system, position)
}

func (s *DebrisService) GetWreckFields(ctx context.Context) ([]*schema.WreckField, error) {
	return s.wreckRepo.GetActive(ctx)
}

func (s *DebrisService) GetWreckField(ctx context.Context, galaxy, system, position int) (*schema.WreckField, error) {
	return s.wreckRepo.GetByCoords(ctx, galaxy, system, position)
}

func (s *DebrisService) CreateDebrisFromBattle(ctx context.Context, galaxy, system, position int, attackerLosses, defenderLosses map[int16]int16) {
	metal := int64(0)
	crystal := int64(0)

	costs := map[int16]int64{
		202: 4000, 203: 12000, 204: 18000, 205: 45000, 206: 40000,
		207: 90000, 208: 20000, 209: 18000, 210: 1000, 211: 125000,
		212: 2500, 213: 20000, 214: 250000, 215: 4000, 218: 8000,
		219: 40000, 401: 2000, 402: 2000, 403: 8000, 404: 8000,
		405: 35000, 406: 130000, 407: 20000, 408: 10000, 409: 20000,
	}

	for shipID, amount := range attackerLosses {
		if cost, ok := costs[shipID]; ok {
			metal += (cost * int64(amount)) / 2
			crystal += (cost * int64(amount)) / 4
		}
	}

	for shipID, amount := range defenderLosses {
		if cost, ok := costs[shipID]; ok {
			metal += (cost * int64(amount)) / 2
			crystal += (cost * int64(amount)) / 4
		}
	}

	if metal > 0 || crystal > 0 {
		_ = s.debrisRepo.AddResources(ctx, galaxy, system, position, metal, crystal)
	}
}

func (s *DebrisService) CreateWreckFieldFromBattle(ctx context.Context, galaxy, system, position int, totalDestroyed int64) {
	if totalDestroyed < 1000000 {
		return
	}

	scrapMetal := totalDestroyed / 10
	scrapCrystal := totalDestroyed / 20
	scrapDeuterium := totalDestroyed / 40

	_ = s.wreckRepo.CreateOrUpdate(ctx, galaxy, system, position, scrapMetal, scrapCrystal, scrapDeuterium)
}

func (s *DebrisService) CollectDebris(ctx context.Context, userID uint, galaxy, system, position int, recyclerCapacity int64) (metal, crystal int64, err error) {
	debris, err := s.debrisRepo.GetByCoords(ctx, galaxy, system, position)
	if err != nil {
		return 0, 0, nil
	}

	recyclableMetal := debris.Metal
	recyclableCrystal := debris.Crystal

	if recyclerCapacity > 0 {
		totalResources := recyclableMetal + recyclableCrystal
		if totalResources > recyclerCapacity {
			ratio := float64(recyclerCapacity) / float64(totalResources)
			recyclableMetal = int64(float64(recyclableMetal) * ratio)
			recyclableCrystal = int64(float64(recyclableCrystal) * ratio)
		}
	}

	planet, _ := s.planetRepo.GetByCoords(ctx, userID, galaxy, system, position)
	if planet != nil {
		_ = s.planetRepo.AddResources(ctx, planet.ID, recyclableMetal, recyclableCrystal, 0)
	}

	debris.Metal -= recyclableMetal
	debris.Crystal -= recyclableCrystal

	if debris.Metal <= 0 && debris.Crystal <= 0 {
		_ = s.debrisRepo.Delete(ctx, debris.ID)
	} else {
		_ = s.debrisRepo.Update(ctx, debris)
	}

	return recyclableMetal, recyclableCrystal, nil
}

func (s *DebrisService) CollectWreckField(ctx context.Context, userID uint, galaxy, system, position int, cargoCapacity int64) (metal, crystal, deuterium int64, err error) {
	wreck, err := s.wreckRepo.GetByCoords(ctx, galaxy, system, position)
	if err != nil {
		return 0, 0, 0, nil
	}

	collectMetal := wreck.Metal
	collectCrystal := wreck.Crystal
	collectDeuterium := wreck.Deuterium

	totalResources := collectMetal + collectCrystal + collectDeuterium
	if cargoCapacity > 0 && totalResources > cargoCapacity {
		ratio := float64(cargoCapacity) / float64(totalResources)
		collectMetal = int64(float64(collectMetal) * ratio)
		collectCrystal = int64(float64(collectCrystal) * ratio)
		collectDeuterium = int64(float64(collectDeuterium) * ratio)
	}

	planet, _ := s.planetRepo.GetByCoords(ctx, userID, galaxy, system, position)
	if planet != nil {
		_ = s.planetRepo.AddResources(ctx, planet.ID, collectMetal, collectCrystal, collectDeuterium)
	}

	wreck.Metal -= collectMetal
	wreck.Crystal -= collectCrystal
	wreck.Deuterium -= collectDeuterium

	if wreck.Metal <= 0 && wreck.Crystal <= 0 && wreck.Deuterium <= 0 {
		_ = s.wreckRepo.Delete(ctx, wreck.ID)
	} else {
		_ = s.wreckRepo.Update(ctx, wreck)
	}

	return collectMetal, collectCrystal, collectDeuterium, nil
}

func (s *DebrisService) CleanupExpired(ctx context.Context) {
	_ = s.debrisRepo.DeleteExpired(ctx)
	_ = s.wreckRepo.DeleteExpired(ctx)
}

func (s *DebrisService) GetDebrisInSystem(ctx context.Context, galaxy, system int) ([]*schema.DebrisField, error) {
	all, err := s.debrisRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	var result []*schema.DebrisField
	for _, d := range all {
		if d.Galaxy == galaxy && d.System == system {
			result = append(result, d)
		}
	}
	return result, nil
}

func (s *DebrisService) GetWrecksInSystem(ctx context.Context, galaxy, system int) ([]*schema.WreckField, error) {
	all, err := s.wreckRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	var result []*schema.WreckField
	for _, w := range all {
		if w.Galaxy == galaxy && w.System == system {
			result = append(result, w)
		}
	}
	return result, nil
}
