package service

import (
	"context"
	"fmt"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type ACSService struct {
	acsRepo     repository.ACSRepository
	fleetRepo   repository.FleetMissionRepository
	userRepo    repository.UserRepository
	planetRepo  repository.PlanetRepository
}

func NewACSService(
	acsRepo repository.ACSRepository,
	fleetRepo repository.FleetMissionRepository,
	userRepo repository.UserRepository,
	planetRepo repository.PlanetRepository,
) *ACSService {
	return &ACSService{
		acsRepo:    acsRepo,
		fleetRepo:  fleetRepo,
		userRepo:   userRepo,
		planetRepo: planetRepo,
	}
}

type ACSCreateParams struct {
	TargetGalaxy   int
	TargetSystem   int
	TargetPosition int
	ArrivalTime    time.Time
	FleetID        uint
}

func (s *ACSService) CreateACS(ctx context.Context, params ACSCreateParams) (*schema.ACS, error) {
	acs := &schema.ACS{
		TargetGalaxy:   params.TargetGalaxy,
		TargetSystem:   params.TargetSystem,
		TargetPosition: params.TargetPosition,
		ArrivalTime:    params.ArrivalTime,
		FleetIDs:       []uint{params.FleetID},
		CreatedAt:      time.Now(),
	}

	err := s.acsRepo.Create(ctx, acs)
	if err != nil {
		return nil, err
	}

	fleet, err := s.fleetRepo.GetByID(ctx, params.FleetID)
	if err != nil {
		return nil, err
	}

	fleet.ACSId = &acs.ID
	err = s.fleetRepo.Update(ctx, fleet)
	if err != nil {
		return nil, err
	}

	return acs, nil
}

func (s *ACSService) JoinACS(ctx context.Context, acsID uint, fleetID uint) (*schema.ACS, error) {
	acs, err := s.acsRepo.GetByID(ctx, acsID)
	if err != nil {
		return nil, err
	}

	fleet, err := s.fleetRepo.GetByID(ctx, fleetID)
	if err != nil {
		return nil, err
	}

	if fleet.TargetGalaxy != acs.TargetGalaxy ||
		fleet.TargetSystem != acs.TargetSystem ||
		fleet.TargetPosition != acs.TargetPosition {
		return nil, fmt.Errorf("fleet target does not match ACS target")
	}

	timeDiff := fleet.ArrivalTime.Sub(acs.ArrivalTime)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	if timeDiff > 5*time.Minute {
		return nil, fmt.Errorf("fleet arrival time differs too much from ACS")
	}

	fleet.ACSId = &acsID
	err = s.fleetRepo.Update(ctx, fleet)
	if err != nil {
		return nil, err
	}

	acs.FleetIDs = append(acs.FleetIDs, fleetID)
	err = s.acsRepo.Create(ctx, acs)
	if err != nil {
		return nil, err
	}

	return acs, nil
}

func (s *ACSService) GetACS(ctx context.Context, acsID uint) (*schema.ACS, error) {
	return s.acsRepo.GetByID(ctx, acsID)
}

func (s *ACSService) GetACSFleets(ctx context.Context, acsID uint) ([]*schema.FleetMission, error) {
	return s.acsRepo.GetFleets(ctx, acsID)
}

func (s *ACSService) FindMatchingACS(ctx context.Context, fleetID uint) (*schema.ACS, error) {
	fleet, err := s.fleetRepo.GetByID(ctx, fleetID)
	if err != nil {
		return nil, err
	}

	if fleet.ACSId != nil {
		return s.acsRepo.GetByID(ctx, *fleet.ACSId)
	}

	arrivingFleets, err := s.fleetRepo.GetArriving(ctx, fleet.ArrivalTime.Add(5*time.Minute))
	if err != nil {
		return nil, err
	}

	for _, f := range arrivingFleets {
		if f.ID == fleetID {
			continue
		}
		if f.UserID == fleet.UserID {
			continue
		}
		if f.TargetGalaxy != fleet.TargetGalaxy ||
			f.TargetSystem != fleet.TargetSystem ||
			f.TargetPosition != fleet.TargetPosition {
			continue
		}

		timeDiff := f.ArrivalTime.Sub(fleet.ArrivalTime)
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}
		if timeDiff <= 5*time.Minute && f.ACSId != nil {
			return s.acsRepo.GetByID(ctx, *f.ACSId)
		}
	}

	return nil, nil
}

func (s *ACSService) GetACSShips(ctx context.Context, acsID uint) (map[int16]int16, error) {
	fleets, err := s.acsRepo.GetFleets(ctx, acsID)
	if err != nil {
		return nil, err
	}

	combinedShips := make(map[int16]int16)
	for _, fleet := range fleets {
		ships := s.parseShips(fleet.Ships)
		for shipID, count := range ships {
			combinedShips[shipID] += count
		}
	}

	return combinedShips, nil
}

func (s *ACSService) parseShips(shipsStr string) map[int16]int16 {
	ships := make(map[int16]int16)
	if shipsStr == "" {
		return ships
	}

	var currentID int16
	for _, part := range splitAndTrim(shipsStr, ";") {
		for i, pair := range splitAndTrim(part, ",") {
			if i == 0 {
				var id int
				fmt.Sscanf(pair, "%d", &id)
				currentID = int16(id)
			} else {
				var count int
				fmt.Sscanf(pair, "%d", &count)
				ships[currentID] += int16(count)
			}
		}
	}
	return ships
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range split(s, sep) {
		trimmed := trim(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func split(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func (s *ACSService) DeleteACS(ctx context.Context, acsID uint) error {
	fleets, err := s.acsRepo.GetFleets(ctx, acsID)
	if err != nil {
		return err
	}

	for _, fleet := range fleets {
		fleet.ACSId = nil
		_ = s.fleetRepo.Update(ctx, fleet)
	}

	return s.acsRepo.Delete(ctx, acsID)
}
