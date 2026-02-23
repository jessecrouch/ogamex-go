package service

import (
	"context"
	"encoding/json"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type EspionageService struct {
	espionageRepo repository.EspionageRepository
	planetRepo    repository.PlanetRepository
	userRepo     repository.UserRepository
	messageRepo   repository.MessageRepository
}

func NewEspionageService(
	espionageRepo repository.EspionageRepository,
	planetRepo repository.PlanetRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
) *EspionageService {
	return &EspionageService{
		espionageRepo: espionageRepo,
		planetRepo:    planetRepo,
		userRepo:     userRepo,
		messageRepo:  messageRepo,
	}
}

type SpyReportInput struct {
	OriginGalaxy   int
	OriginSystem   int
	OriginPosition int
	TargetGalaxy   int
	TargetSystem   int
	TargetPosition int
	Ships          map[int]int
}

func (s *EspionageService) ExecuteSpyMission(ctx context.Context, userID uint, input SpyReportInput) (*schema.EspionageReport, error) {
	probeCount := input.Ships[210]
	if probeCount <= 0 {
		return nil, nil
	}

	targetPlanet, err := s.planetRepo.GetByCoordsAny(ctx, input.TargetGalaxy, input.TargetSystem, input.TargetPosition)
	if err != nil {
		return nil, nil
	}

	reportType := "no_activity"
	if targetPlanet != nil && targetPlanet.UserID != 0 {
		reportType = s.calculateReportLevel(probeCount)

		report := &schema.EspionageReport{
			UserID:        userID,
			TargetUserID:  targetPlanet.UserID,
			TargetPlanetID: targetPlanet.ID,
			Galaxy:        input.TargetGalaxy,
			System:        input.TargetSystem,
			Position:      input.TargetPosition,
			ReportType:    reportType,
			Metal:         targetPlanet.Metal,
			Crystal:       targetPlanet.Crystal,
			Deuterium:    targetPlanet.Deuterium,
			Energy:        targetPlanet.EnergyAvailable,
			Ships:         s.formatShips(targetPlanet),
			Defense:       s.formatDefense(targetPlanet),
			Buildings:     s.formatBuildings(targetPlanet),
			Read:          false,
			CreatedAt:     time.Now(),
		}

		if targetPlanet.UserID != userID {
			user, _ := s.userRepo.GetByID(ctx, targetPlanet.UserID)
			username := "Unknown"
			if user != nil {
				username = user.Username
			}

			msg := &schema.Message{
				UserID:    targetPlanet.UserID,
				Type:      50,
				FromName:  username,
				Subject:   "Espionage Report",
				Body:      "Your planet was spied on!",
				CreatedAt: time.Now(),
			}
			_ = s.messageRepo.Create(ctx, msg)
		}

		err = s.espionageRepo.Create(ctx, report)
		if err != nil {
			return nil, err
		}

		return report, nil
	}

	report := &schema.EspionageReport{
		UserID:     userID,
		Galaxy:     input.TargetGalaxy,
		System:     input.TargetSystem,
		Position:   input.TargetPosition,
		ReportType: "no_activity",
		Read:       false,
		CreatedAt:  time.Now(),
	}

	_ = s.espionageRepo.Create(ctx, report)
	return report, nil
}

func (s *EspionageService) calculateReportLevel(probes int) string {
	if probes >= 100 {
		return "full"
	} else if probes >= 50 {
		return "partial"
	} else if probes >= 10 {
		return "minimal"
	} else if probes >= 1 {
		return "short"
	}
	return "no_activity"
}

func (s *EspionageService) formatShips(planet *schema.Planet) string {
	ships := map[string]int{
		"small_cargo":   planet.SmallCargo,
		"large_cargo":   planet.LargeCargo,
		"light_fighter": planet.LightFighter,
		"heavy_fighter": planet.HeavyFighter,
		"cruiser":       planet.Cruiser,
		"battleship":    planet.Battleship,
		"colony_ship":   planet.ColonyShip,
		"recycler":      planet.Recycler,
		"espionage_probe": planet.EspionageProbe,
		"bomber":        planet.Bomber,
		"destroyer":     planet.Destroyer,
		"deathstar":     planet.Deathstar,
		"battlecruiser": planet.Battlecruiser,
		"reaper":        planet.Reaper,
		"pathfinder":    planet.Pathfinder,
	}

	result := make(map[string]int)
	for k, v := range ships {
		if v > 0 {
			result[k] = v
		}
	}

	b, _ := json.Marshal(result)
	return string(b)
}

func (s *EspionageService) formatDefense(planet *schema.Planet) string {
	defense := map[string]int{
		"rocket_launcher":      planet.RocketLauncher,
		"light_laser":         planet.LightLaser,
		"heavy_laser":         planet.HeavyLaser,
		"ion_cannon":          planet.IonCannon,
		"gauss_cannon":        planet.GaussCannon,
		"plasma_turret":       planet.PlasmaTurret,
		"shield_dome":         planet.ShieldDome,
		"missile_interceptor": planet.MissileInterceptor,
		"missile_launcher":    planet.MissileLauncher,
	}

	result := make(map[string]int)
	for k, v := range defense {
		if v > 0 {
			result[k] = v
		}
	}

	b, _ := json.Marshal(result)
	return string(b)
}

func (s *EspionageService) formatBuildings(planet *schema.Planet) string {
	buildings := map[string]int{
		"metal_mine":            planet.MetalMine,
		"crystal_mine":         planet.CrystalMine,
		"deuterium_synthesizer": planet.DeuteriumSynthesizer,
		"solar_plant":          planet.SolarPlant,
		"fusion_plant":        planet.FusionPlant,
		"metal_storage":       planet.MetalStorageBuilding,
		"crystal_storage":     planet.CrystalStorageBuilding,
		"deuterium_storage":   planet.DeuteriumStorageBuilding,
		"robot_factory":       planet.RobotFactory,
		"shipyard":            planet.Shipyard,
		"research_lab":        planet.ResearchLab,
		"nanite_factory":      planet.NaniteFactory,
		"terraformer":          planet.Terraformer,
		"space_dock":           planet.SpaceDock,
	}

	b, _ := json.Marshal(buildings)
	return string(b)
}

func (s *EspionageService) GetReports(ctx context.Context, userID uint, limit, offset int) ([]*schema.EspionageReport, error) {
	return s.espionageRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *EspionageService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.espionageRepo.GetUnreadCount(ctx, userID)
}

func (s *EspionageService) MarkAsRead(ctx context.Context, reportID uint) error {
	return s.espionageRepo.MarkAsRead(ctx, reportID)
}

func (s *EspionageService) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.espionageRepo.MarkAllAsRead(ctx, userID)
}

func (s *EspionageService) DeleteReport(ctx context.Context, reportID uint) error {
	return s.espionageRepo.Delete(ctx, reportID)
}

func (s *EspionageService) DeleteAllReports(ctx context.Context, userID uint) error {
	return s.espionageRepo.DeleteAll(ctx, userID)
}
