package service

import (
	"context"
	"crypto/rand"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type CounterEspionageService struct {
	espionageRepo repository.EspionageRepository
	planetRepo    repository.PlanetRepository
	userRepo      repository.UserRepository
}

func NewCounterEspionageService(
	espionageRepo repository.EspionageRepository,
	planetRepo repository.PlanetRepository,
	userRepo repository.UserRepository,
) *CounterEspionageService {
	return &CounterEspionageService{
		espionageRepo: espionageRepo,
		planetRepo:    planetRepo,
		userRepo:      userRepo,
	}
}

func (s *CounterEspionageService) DetectEspionage(ctx context.Context, targetPlanetID uint, attackerID uint) (bool, error) {
	planet, err := s.planetRepo.GetByID(ctx, targetPlanetID)
	if err != nil {
		return false, err
	}

	counterChance := s.calculateCounterChance(planet)
	
	bytes := make([]byte, 1)
	rand.Read(bytes)
	roll := float64(bytes[0]) / 255.0
	
	detected := roll < counterChance
	
	return detected, nil
}

func (s *CounterEspionageService) calculateCounterChance(planet *schema.Planet) float64 {
	baseChance := 0.0
	
	counterChance := baseChance
	
	counterChance += float64(planet.SensorPhalanx) * 0.02
	
	if planet.SensorPhalanx >= 3 {
		counterChance += 0.15
	}
	if planet.SensorPhalanx >= 5 {
		counterChance += 0.20
	}
	
	if counterChance > 1.0 {
		counterChance = 1.0
	}
	
	return counterChance
}

func (s *CounterEspionageService) CounterAttackerFleet(ctx context.Context, attackerID, targetPlanetID uint) error {
	planet, err := s.planetRepo.GetByID(ctx, targetPlanetID)
	if err != nil {
		return err
	}
	
	if planet.SensorPhalanx < 1 {
		return nil
	}
	
	detected, err := s.DetectEspionage(ctx, targetPlanetID, attackerID)
	if err != nil {
		return err
	}
	
	if detected {
		attacker, _ := s.userRepo.GetByID(ctx, attackerID)
		if attacker != nil {
			_ = s.espionageRepo.Create(ctx, &schema.EspionageReport{
				UserID:        planet.UserID,
				TargetUserID:  attackerID,
				Galaxy:        planet.Galaxy,
				System:        planet.System,
				Position:      planet.Position,
				ReportType:    "counter",
				Metal:         0,
				Crystal:       0,
				Deuterium:     0,
				Energy:        0,
				Ships:         "Counter-espionage detected",
				Defense:       "",
				Buildings:     "",
				CreatedAt:     time.Now(),
			})
		}
	}
	
	return nil
}

func (s *CounterEspionageService) GetCounterEspionageLevel(ctx context.Context, planetID uint) (int, error) {
	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return 0, err
	}
	return planet.SensorPhalanx, nil
}

func (s *CounterEspionageService) CalculateDetectionChance(planet *schema.Planet) float64 {
	return s.calculateCounterChance(planet)
}

func (s *CounterEspionageService) GetCounterReports(ctx context.Context, userID uint) ([]*schema.EspionageReport, error) {
	allReports, err := s.espionageRepo.GetByUserID(ctx, userID, 100, 0)
	if err != nil {
		return nil, err
	}

	var counterReports []*schema.EspionageReport
	for _, r := range allReports {
		if r.ReportType == "counter" {
			counterReports = append(counterReports, r)
		}
	}

	return counterReports, nil
}
