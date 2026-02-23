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
}

func NewAuthService(userRepo repository.UserRepository, planetRepo repository.PlanetRepository) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		planetRepo: planetRepo,
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

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
