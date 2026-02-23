package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/formula"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type ResearchService struct {
	userRepo       repository.UserRepository
	planetRepo     repository.PlanetRepository
	researchQueue  repository.ResearchQueueRepository
	techRepo       repository.UserTechRepository
}

func NewResearchService(
	userRepo repository.UserRepository,
	planetRepo repository.PlanetRepository,
	researchQueue repository.ResearchQueueRepository,
	techRepo repository.UserTechRepository,
) *ResearchService {
	return &ResearchService{
		userRepo:      userRepo,
		planetRepo:    planetRepo,
		researchQueue: researchQueue,
		techRepo:      techRepo,
	}
}

type ResearchCost struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
}

func (s *ResearchService) GetResearchCost(researchID int, level int) ResearchCost {
	baseCost := s.getBaseCost(researchID)
	factor := formula.GetCostFactor(level, 1.8)

	return ResearchCost{
		Metal:     int64(float64(baseCost.Metal) * factor),
		Crystal:   int64(float64(baseCost.Crystal) * factor),
		Deuterium: int64(float64(baseCost.Deuterium) * factor),
	}
}

func (s *ResearchService) getBaseCost(researchID int) ResearchCost {
	costs := map[int]ResearchCost{
		1:  {Metal: 0, Crystal: 800, Deuterium: 400},    // Energy
		2:  {Metal: 200, Crystal: 600, Deuterium: 0},    // Laser
		3:  {Metal: 1000, Crystal: 300, Deuterium: 0},   // Ion
		4:  {Metal: 4000, Crystal: 2000, Deuterium: 1000}, // Hyperspace
		5:  {Metal: 2400, Crystal: 1200, Deuterium: 600}, // Plasma
		6:  {Metal: 9000, Crystal: 4000, Deuterium: 6000}, // Fusion
		7:  {Metal: 4000, Crystal: 2000, Deuterium: 600}, // Impulse
		8:  {Metal: 10000, Crystal: 6000, Deuterium: 4000}, // Hyperspace Drive
		9:  {Metal: 200, Crystal: 1000, Deuterium: 200},  // Espionage
		10: {Metal: 100, Crystal: 400, Deuterium: 200},   // Computer
		11: {Metal: 8000, Crystal: 4000, Deuterium: 2000}, // Astrophysics
		12: {Metal: 240000, Crystal: 160000, Deuterium: 80000}, // IGR
		13: {Metal: 100000, Crystal: 50000, Deuterium: 50000}, // Graviton
		14: {Metal: 800, Crystal: 200, Deuterium: 0},       // Weapons
		15: {Metal: 400, Crystal: 600, Deuterium: 0},       // Shielding
		16: {Metal: 400, Crystal: 200, Deuterium: 0},        // Armor
	}

	if cost, ok := costs[researchID]; ok {
		return cost
	}
	return ResearchCost{}
}

func (s *ResearchService) StartResearch(ctx context.Context, userID uint, researchID int) error {
	tech, err := s.techRepo.GetByUserID(ctx, userID)
	if err != nil {
		tech = &schema.UserTech{UserID: userID}
		s.techRepo.Create(ctx, tech)
	}

	currentLevel := s.getTechLevel(tech, researchID)
	cost := s.GetResearchCost(researchID, currentLevel+1)

	planets, err := s.planetRepo.GetByUserID(ctx, userID)
	if err != nil || len(planets) == 0 {
		return errors.New("no planets found")
	}

	var labPlanet *schema.Planet
	for _, p := range planets {
		if p.ResearchLab > 0 {
			labPlanet = p
			break
		}
	}

	if labPlanet == nil {
		return errors.New("no research lab found")
	}

	if labPlanet.Metal < cost.Metal || labPlanet.Crystal < cost.Crystal || labPlanet.Deuterium < cost.Deuterium {
		return ErrInsufficientFunds
	}

	currentQueue, err := s.researchQueue.GetCurrent(ctx, userID)
	if err == nil && currentQueue != nil {
		return ErrLabBusy
	}

	err = s.planetRepo.SubResources(ctx, labPlanet.ID, cost.Metal, cost.Crystal, cost.Deuterium)
	if err != nil {
		return err
	}

	buildTime := s.calculateResearchTime(researchID, currentLevel+1, labPlanet)

	now := time.Now()
	queueItem := &schema.ResearchQueue{
		UserID:     userID,
		ResearchID: researchID,
		Level:      currentLevel + 1,
		StartTime:  now,
		EndTime:    now.Add(buildTime),
	}

	return s.researchQueue.Create(ctx, queueItem)
}

func (s *ResearchService) getTechLevel(tech *schema.UserTech, researchID int) int {
	switch researchID {
	case 1:
		return tech.EnergyTechnology
	case 2:
		return tech.LaserTechnology
	case 3:
		return tech.IonTechnology
	case 4:
		return tech.HyperspaceTechnology
	case 5:
		return tech.PlasmaTechnology
	case 6:
		return tech.CombatDrive
	case 7:
		return tech.ImpulseDrive
	case 8:
		return tech.HyperspaceDrive
	case 9:
		return tech.EspionageTechnology
	case 10:
		return tech.ComputerTechnology
	case 11:
		return tech.Astrophysics
	case 12:
		return tech.IntergalacticResearch
	case 13:
		return tech.GravitonTechnology
	case 14:
		return tech.WeaponsTechnology
	case 15:
		return tech.ShieldingTechnology
	case 16:
		return tech.ArmorTechnology
	}
	return 0
}

func (s *ResearchService) calculateResearchTime(researchID int, level int, planet *schema.Planet) time.Duration {
	baseTime := time.Duration(60*level) * time.Second

	labBonus := 1.0
	if planet.ResearchLab > 0 {
		labBonus = 1.0 / (1.0 + float64(planet.ResearchLab)*0.25)
	}

	naniteBonus := 1.0
	if planet.NaniteFactory > 0 {
		naniteBonus = 1.0 / (1.0 + float64(planet.NaniteFactory)*0.5)
	}

	return time.Duration(float64(baseTime) * labBonus * naniteBonus)
}

func (s *ResearchService) CompleteResearch(ctx context.Context, queueID uint) error {
	queue, err := s.researchQueue.GetByID(ctx, queueID)
	if err != nil {
		return err
	}

	if time.Now().Before(queue.EndTime) {
		return errors.New("research not yet complete")
	}

	tech, err := s.techRepo.GetByUserID(ctx, queue.UserID)
	if err != nil {
		tech = &schema.UserTech{UserID: queue.UserID}
	}

	switch queue.ResearchID {
	case 1:
		tech.EnergyTechnology = queue.Level
	case 2:
		tech.LaserTechnology = queue.Level
	case 3:
		tech.IonTechnology = queue.Level
	case 4:
		tech.HyperspaceTechnology = queue.Level
	case 5:
		tech.PlasmaTechnology = queue.Level
	case 6:
		tech.CombatDrive = queue.Level
	case 7:
		tech.ImpulseDrive = queue.Level
	case 8:
		tech.HyperspaceDrive = queue.Level
	case 9:
		tech.EspionageTechnology = queue.Level
	case 10:
		tech.ComputerTechnology = queue.Level
	case 11:
		tech.Astrophysics = queue.Level
	case 12:
		tech.IntergalacticResearch = queue.Level
	case 13:
		tech.GravitonTechnology = queue.Level
	case 14:
		tech.WeaponsTechnology = queue.Level
	case 15:
		tech.ShieldingTechnology = queue.Level
	case 16:
		tech.ArmorTechnology = queue.Level
	}

	err = s.techRepo.Update(ctx, tech)
	if err != nil {
		return err
	}

	return s.researchQueue.Delete(ctx, queueID)
}

func (s *ResearchService) GetQueue(ctx context.Context, userID uint) ([]*schema.ResearchQueue, error) {
	return s.researchQueue.GetByUserID(ctx, userID)
}
