package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type PremiumService struct {
	userRepo  repository.UserRepository
	planetRepo repository.PlanetRepository
}

func NewPremiumService(
	userRepo repository.UserRepository,
	planetRepo repository.PlanetRepository,
) *PremiumService {
	return &PremiumService{
		userRepo:    userRepo,
		planetRepo: planetRepo,
	}
}

type MerchantOffer struct {
	ResourceType string
	Amount       int64
	Price        int64
}

func (s *PremiumService) BuyResource(ctx context.Context, userID uint, planetID uint, resourceType string, amount int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return err
	}

	if user.DarkMatter < amount {
		return errors.New("insufficient dark matter")
	}

	exchangeRate := int64(2)
	price := amount * exchangeRate

	if user.DarkMatter < price {
		return errors.New("insufficient dark matter for this purchase")
	}

	user.DarkMatter -= price

	switch resourceType {
	case "metal":
		planet.Metal += amount * 1000
	case "crystal":
		planet.Crystal += amount * 1000
	case "deuterium":
		planet.Deuterium += amount * 1000
	default:
		return errors.New("invalid resource type")
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	return s.planetRepo.Update(ctx, planet)
}

func (s *PremiumService) SellResource(ctx context.Context, userID uint, planetID uint, resourceType string, amount int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return err
	}

	exchangeRate := int64(2)

	switch resourceType {
	case "metal":
		if planet.Metal < amount*1000 {
			return errors.New("insufficient metal")
		}
		planet.Metal -= amount * 1000
		user.DarkMatter += amount * exchangeRate
	case "crystal":
		if planet.Crystal < amount*1000 {
			return errors.New("insufficient crystal")
		}
		planet.Crystal -= amount * 1000
		user.DarkMatter += amount * exchangeRate
	case "deuterium":
		if planet.Deuterium < amount*1000 {
			return errors.New("insufficient deuterium")
		}
		planet.Deuterium -= amount * 1000
		user.DarkMatter += amount * exchangeRate
	default:
		return errors.New("invalid resource type")
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	return s.planetRepo.Update(ctx, planet)
}

func (s *PremiumService) ActivatePremium(ctx context.Context, userID uint, days int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	pricePerDay := int64(500)
	totalPrice := pricePerDay * int64(days)

	if user.DarkMatter < totalPrice {
		return errors.New("insufficient dark matter")
	}

	user.DarkMatter -= totalPrice

	if user.PremiumEndsAt == nil {
		now := time.Now()
		user.PremiumEndsAt = &now
	}
	*user.PremiumEndsAt = user.PremiumEndsAt.Add(time.Duration(days) * 24 * time.Hour)

	return s.userRepo.Update(ctx, user)
}

func (s *PremiumService) GetPremiumStatus(ctx context.Context, userID uint) (bool, *time.Time, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, nil, err
	}

	if user.PremiumEndsAt == nil {
		return false, nil, nil
	}

	isActive := time.Now().Before(*user.PremiumEndsAt)
	return isActive, user.PremiumEndsAt, nil
}

func (s *PremiumService) GetResourcePrice(resourceType string) int64 {
	switch resourceType {
	case "metal":
		return 2
	case "crystal":
		return 2
	case "deuterium":
		return 2
	default:
		return 0
	}
}

func (s *PremiumService) CalculateLootValue(planet *schema.Planet) int64 {
	return planet.Metal + planet.Crystal + planet.Deuterium
}

func (s *PremiumService) ApplyPremiumBonus(production int64, isPremium bool) int64 {
	if !isPremium {
		return production
	}
	return int64(float64(production) * 1.25)
}
