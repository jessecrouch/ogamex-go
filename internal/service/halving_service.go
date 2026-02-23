package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/repository"
)

type HalvingService struct {
	userRepo repository.UserRepository
}

func NewHalvingService(userRepo repository.UserRepository) *HalvingService {
	return &HalvingService{
		userRepo: userRepo,
	}
}

func (s *HalvingService) IsHalvingActive(ctx context.Context, userID uint) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if user.DMHolding == nil || *user.DMHolding <= 0 {
		return false, nil
	}

	return true, nil
}

func (s *HalvingService) GetHalvingRemaining(ctx context.Context, userID uint) (int64, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}

	if user.DMHolding == nil {
		return 0, nil
	}

	return *user.DMHolding, nil
}

func (s *HalvingService) PurchaseHalving(ctx context.Context, userID uint, dmAmount int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.DarkMatter < dmAmount {
		return errors.New("insufficient dark matter")
	}

	user.DarkMatter -= dmAmount

	if user.DMHolding == nil {
		user.DMHolding = new(int64)
	}
	*user.DMHolding += dmAmount

	return s.userRepo.Update(ctx, user)
}

func (s *HalvingService) ConsumeHalving(ctx context.Context, userID uint, amount int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.DMHolding == nil || *user.DMHolding <= 0 {
		return errors.New("no halving available")
	}

	if *user.DMHolding < amount {
		return errors.New("insufficient halving amount")
	}

	*user.DMHolding -= amount

	return s.userRepo.Update(ctx, user)
}

func (s *HalvingService) ApplyHalvingToDuration(duration time.Duration, userID uint, isActive bool) time.Duration {
	if !isActive {
		return duration
	}
	return time.Duration(float64(duration) * 0.5)
}

func (s *HalvingService) CalculateHalvingCost(baseDuration time.Duration, targetDuration time.Duration) int64 {
	if targetDuration >= baseDuration {
		return 0
	}

	minutesReduced := float64(baseDuration-targetDuration) / float64(time.Minute)
	
	costPerMinute := int64(2)
	
	return int64(minutesReduced) * costPerMinute
}

func (s *HalvingService) UseHalvingForBuild(ctx context.Context, userID uint, buildingID int, level int, baseDuration time.Duration) (time.Duration, error) {
	active, err := s.IsHalvingActive(ctx, userID)
	if err != nil {
		return baseDuration, err
	}

	if !active {
		return baseDuration, nil
	}

	err = s.ConsumeHalving(ctx, userID, 1)
	if err != nil {
		return baseDuration, nil
	}

	return s.ApplyHalvingToDuration(baseDuration, userID, true), nil
}

func (s *HalvingService) UseHalvingForResearch(ctx context.Context, userID uint, researchID int, level int, baseDuration time.Duration) (time.Duration, error) {
	active, err := s.IsHalvingActive(ctx, userID)
	if err != nil {
		return baseDuration, err
	}

	if !active {
		return baseDuration, nil
	}

	err = s.ConsumeHalving(ctx, userID, 1)
	if err != nil {
		return baseDuration, nil
	}

	return s.ApplyHalvingToDuration(baseDuration, userID, true), nil
}

func (s *HalvingService) UseHalvingForUnit(ctx context.Context, userID uint, unitID int, amount int, baseDuration time.Duration) (time.Duration, error) {
	active, err := s.IsHalvingActive(ctx, userID)
	if err != nil {
		return baseDuration, err
	}

	if !active {
		return baseDuration, nil
	}

	err = s.ConsumeHalving(ctx, userID, 1)
	if err != nil {
		return baseDuration, nil
	}

	return s.ApplyHalvingToDuration(baseDuration, userID, true), nil
}
