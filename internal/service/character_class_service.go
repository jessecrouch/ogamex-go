package service

import (
	"context"
	"errors"

	"ogamex-go/internal/repository"
)

type CharacterClassService struct {
	userRepo repository.UserRepository
}

func NewCharacterClassService(userRepo repository.UserRepository) *CharacterClassService {
	return &CharacterClassService{
		userRepo: userRepo,
	}
}

const (
	ClassNone    = 0
	ClassGeneral = 1
	ClassAdmiral = 2
	ClassEngineer = 3
	ClassGeologist = 4
)

type CharacterClassInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	NameKey  string `json:"name_key"`
	Bonus    string `json:"bonus"`
	BonusKey string `json:"bonus_key"`
}

func GetAllCharacterClasses() []CharacterClassInfo {
	return []CharacterClassInfo{
		{
			ID:       ClassNone,
			Name:     "No Class",
			NameKey:  "NoClass",
			Bonus:    "None",
			BonusKey: "NoClassBonus",
		},
		{
			ID:       ClassGeneral,
			Name:     "General",
			NameKey:  "class_general",
			Bonus:    "Military Bonus: +10% to all military capacities",
			BonusKey: "class_general_bonus",
		},
		{
			ID:       ClassAdmiral,
			Name:     "Admiral",
			NameKey:  "class_admiral",
			Bonus:    "Fleet Capacity: +25% cargo capacity",
			BonusKey: "class_admiral_bonus",
		},
		{
			ID:       ClassEngineer,
			Name:     "Engineer",
			NameKey:  "class_engineer",
			Bonus:    "Defense Bonus: -25% shield damage, +20% energy",
			BonusKey: "class_engineer_bonus",
		},
		{
			ID:       ClassGeologist,
			Name:     "Geologist",
			NameKey:  "class_geologist",
			Bonus:    "Resource Bonus: +10% metal, crystal, deuterium production",
			BonusKey: "class_geologist_bonus",
		},
	}
}

func (s *CharacterClassService) GetCharacterClass(ctx context.Context, userID uint) (int, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.CharacterClass, nil
}

func (s *CharacterClassService) SetCharacterClass(ctx context.Context, userID uint, classID int) error {
	if classID < 0 || classID > 4 {
		return errors.New("invalid character class")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.CharacterClass != 0 && classID != 0 {
		return errors.New("character class already set, cannot change")
	}

	user.CharacterClass = classID
	return s.userRepo.Update(ctx, user)
}

func (s *CharacterClassService) GetProductionBonus(ctx context.Context, userID uint) float64 {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 1.0
	}

	switch user.CharacterClass {
	case ClassGeologist:
		return 1.10
	default:
		return 1.0
	}
}

func (s *CharacterClassService) GetCargoBonus(ctx context.Context, userID uint) float64 {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 1.0
	}

	switch user.CharacterClass {
	case ClassAdmiral:
		return 1.25
	default:
		return 1.0
	}
}

func (s *CharacterClassService) GetMilitaryBonus(ctx context.Context, userID uint) float64 {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 1.0
	}

	switch user.CharacterClass {
	case ClassGeneral:
		return 1.10
	default:
		return 1.0
	}
}

func (s *CharacterClassService) GetShieldReduction(ctx context.Context, userID uint) float64 {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 1.0
	}

	switch user.CharacterClass {
	case ClassEngineer:
		return 0.75
	default:
		return 1.0
	}
}

func (s *CharacterClassService) GetEnergyBonus(ctx context.Context, userID uint) float64 {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 1.0
	}

	switch user.CharacterClass {
	case ClassEngineer:
		return 1.20
	default:
		return 1.0
	}
}

func (s *CharacterClassService) ApplyProductionBonus(baseProduction int64) int64 {
	return int64(float64(baseProduction) * 1.10)
}

func (s *CharacterClassService) ApplyCargoBonus(baseCapacity int64) int64 {
	return int64(float64(baseCapacity) * 1.25)
}
