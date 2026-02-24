package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidToken     = errors.New("invalid token")
)

type AuthService struct {
	userRepo    repository.UserRepository
	planetRepo  repository.PlanetRepository
	techRepo    repository.UserTechRepository
}

func NewAuthService(userRepo repository.UserRepository, planetRepo repository.PlanetRepository, techRepo repository.UserTechRepository) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		planetRepo: planetRepo,
		techRepo:   techRepo,
	}
}

type RegisterInput struct {
	Username  string
	Email     string
	Password  string
	PlayerName string
}

type LoginInput struct {
	Username string
	Password string
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*schema.User, error) {
	existing, _ := s.userRepo.GetByUsername(ctx, input.Username)
	if existing != nil {
		return nil, ErrUserExists
	}

	existing, _ = s.userRepo.GetByEmail(ctx, input.Email)
	if existing != nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	authToken := generateToken()

	user := &schema.User{
		Username:       input.Username,
		Email:          input.Email,
		Password:       string(hashedPassword),
		PlayerName:     input.PlayerName,
		AuthToken:      authToken,
		RegisteredAt:   time.Now(),
		LastOnline:     time.Now(),
		CharacterClass: 0,
		DarkMatter:     1000,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	homePlanet := &schema.Planet{
		UserID:   user.ID,
		Name:     "Homeworld",
		Galaxy:   1,
		System:   1,
		Position: 1,
		IsMoon:   false,
		PlanetType: 1,

		Metal:     500,
		Crystal:   250,
		Deuterium: 0,

		MetalCapacity:     1000,
		CrystalCapacity:   1000,
		DeuteriumCapacity: 1000,

		TempMin:      -50,
		TempMax:      50,
		FieldsUsed:   0,
		FieldsMax:    10,

		MetalMinePercent:           100,
		CrystalMinePercent:         100,
		DeuteriumSynthesizerPercent: 100,
		SolarPlantPercent:          100,
		FusionPlantPercent:         100,
	}

	err = s.planetRepo.Create(ctx, homePlanet)
	if err != nil {
		return nil, err
	}

	user.CurrentPlanetID = &homePlanet.ID
	s.userRepo.Update(ctx, user)

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*schema.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	user.LastOnline = time.Now()
	s.userRepo.Update(ctx, user)

	return user, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*schema.User, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.GetByAuthToken(ctx, token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return user, nil
}

func (s *AuthService) GetUserTech(ctx context.Context, userID uint) (*schema.UserTech, error) {
	tech, err := s.techRepo.GetByUserID(ctx, userID)
	if err != nil {
		return &schema.UserTech{
			UserID: userID,
		}, nil
	}
	return tech, nil
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *AuthService) SetVacationMode(ctx context.Context, userID uint, enabled bool, endTime *time.Time) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if enabled {
		if user.OnVacation {
			return errors.New("already on vacation")
		}

		planets, _ := s.planetRepo.GetByUserID(ctx, userID)
		for _, p := range planets {
			if p.DefenseActivated {
				p.DefenseActivated = false
				_ = s.planetRepo.Update(ctx, p)
			}
		}

		vacationEnd := time.Now().Add(7 * 24 * time.Hour)
		user.OnVacation = true
		user.VacationEndTime = &vacationEnd
	} else {
		if !user.OnVacation {
			return errors.New("not on vacation")
		}

		user.OnVacation = false
		user.VacationEndTime = nil
	}

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) GetVacationStatus(ctx context.Context, userID uint) (bool, *time.Time, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, nil, err
	}

	if user.OnVacation && user.VacationEndTime != nil {
		if time.Now().After(*user.VacationEndTime) {
			user.OnVacation = false
			user.VacationEndTime = nil
			_ = s.userRepo.Update(ctx, user)
			return false, nil, nil
		}
	}

	return user.OnVacation, user.VacationEndTime, nil
}

func (s *AuthService) SetCurrentPlanet(ctx context.Context, userID, planetID uint) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	planet, err := s.planetRepo.GetByID(ctx, planetID)
	if err != nil {
		return errors.New("planet not found")
	}

	if planet.UserID != userID {
		return errors.New("planet does not belong to user")
	}

	user.CurrentPlanetID = &planetID
	return s.userRepo.Update(ctx, user)
}
