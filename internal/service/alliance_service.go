package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type AllianceService struct {
	allianceRepo repository.AllianceRepository
	userRepo     repository.UserRepository
	planetRepo   repository.PlanetRepository
}

func NewAllianceService(allianceRepo repository.AllianceRepository, userRepo repository.UserRepository, planetRepo repository.PlanetRepository) *AllianceService {
	return &AllianceService{
		allianceRepo: allianceRepo,
		userRepo:     userRepo,
		planetRepo:   planetRepo,
	}
}

type CreateAllianceInput struct {
	Name        string
	Tag         string
	Description string
	Logo        string
	Website     string
}

func (s *AllianceService) CreateAlliance(ctx context.Context, userID uint, input CreateAllianceInput) (*schema.Alliance, error) {
	if input.Name == "" {
		return nil, errors.New("alliance name is required")
	}
	if input.Tag == "" {
		return nil, errors.New("alliance tag is required")
	}

	existing, _ := s.allianceRepo.GetByTag(ctx, input.Tag)
	if existing != nil {
		return nil, errors.New("alliance tag already taken")
	}

	existing, _ = s.allianceRepo.GetByName(ctx, input.Name)
	if existing != nil {
		return nil, errors.New("alliance name already taken")
	}

	alliance := &schema.Alliance{
		Name:        input.Name,
		Tag:         input.Tag,
		FounderID:   userID,
		Description: input.Description,
		Logo:        input.Logo,
		Website:     input.Website,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.allianceRepo.Create(ctx, alliance)
	if err != nil {
		return nil, err
	}

	member := &schema.AllianceMember{
		AllianceID: alliance.ID,
		UserID:     userID,
		Rank:       "Founder",
		JoinedAt:   time.Now(),
	}
	_ = s.allianceRepo.AddMember(ctx, member)

	return alliance, nil
}

func (s *AllianceService) GetAlliance(ctx context.Context, allianceID uint) (*schema.Alliance, error) {
	return s.allianceRepo.GetByID(ctx, allianceID)
}

func (s *AllianceService) GetAllianceByTag(ctx context.Context, tag string) (*schema.Alliance, error) {
	return s.allianceRepo.GetByTag(ctx, tag)
}

func (s *AllianceService) GetAllAlliances(ctx context.Context) ([]*schema.Alliance, error) {
	return s.allianceRepo.GetAll(ctx)
}

func (s *AllianceService) GetMembers(ctx context.Context, allianceID uint) ([]*AllianceMemberInfo, error) {
	members, err := s.allianceRepo.GetMembers(ctx, allianceID)
	if err != nil {
		return nil, err
	}

	var result []*AllianceMemberInfo
	for _, m := range members {
		user, _ := s.userRepo.GetByID(ctx, m.UserID)
		planets, _ := s.planetRepo.GetByUserID(ctx, m.UserID)

		result = append(result, &AllianceMemberInfo{
			ID:         m.ID,
			UserID:     m.UserID,
			Username:   user.Username,
			PlayerName: user.PlayerName,
			Rank:       m.Rank,
			Planets:    len(planets),
			JoinedAt:   m.JoinedAt,
		})
	}

	return result, nil
}

type AllianceMemberInfo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	PlayerName string    `json:"player_name"`
	Rank       string    `json:"rank"`
	Planets    int       `json:"planets"`
	JoinedAt   time.Time `json:"joined_at"`
}

func (s *AllianceService) GetUserAlliance(ctx context.Context, userID uint) (*schema.Alliance, error) {
	alliances, err := s.allianceRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, a := range alliances {
		member, err := s.allianceRepo.GetMember(ctx, a.ID, userID)
		if err == nil && member != nil {
			return a, nil
		}
	}

	return nil, errors.New("user not in alliance")
}

func (s *AllianceService) ApplyToAlliance(ctx context.Context, userID uint, allianceID uint, message string) error {
	_, err := s.allianceRepo.GetByID(ctx, allianceID)
	if err != nil {
		return errors.New("alliance not found")
	}

	member, _ := s.allianceRepo.GetMember(ctx, allianceID, userID)
	if member != nil {
		return errors.New("already a member")
	}

	existingApp, _ := s.allianceRepo.GetApplication(ctx, allianceID, userID)
	if existingApp != nil {
		return errors.New("application already pending")
	}

	app := &schema.AllianceApplication{
		AllianceID: allianceID,
		UserID:     userID,
		Message:    message,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	return s.allianceRepo.CreateApplication(ctx, app)
}

func (s *AllianceService) GetApplications(ctx context.Context, allianceID uint) ([]*AllianceApplicationInfo, error) {
	apps, err := s.allianceRepo.GetApplications(ctx, allianceID)
	if err != nil {
		return nil, err
	}

	var result []*AllianceApplicationInfo
	for _, a := range apps {
		user, _ := s.userRepo.GetByID(ctx, a.UserID)
		result = append(result, &AllianceApplicationInfo{
			ID:         a.ID,
			UserID:     a.UserID,
			Username:   user.Username,
			PlayerName: user.PlayerName,
			Message:    a.Message,
			CreatedAt:  a.CreatedAt,
		})
	}

	return result, nil
}

type AllianceApplicationInfo struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	PlayerName string    `json:"player_name"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *AllianceService) AcceptApplication(ctx context.Context, adminUserID, applicationID uint) error {
	apps, err := s.allianceRepo.GetApplications(ctx, 0)
	if err != nil {
		return err
	}

	var app *schema.AllianceApplication
	for _, a := range apps {
		if a.ID == applicationID {
			app = a
			break
		}
	}
	if app == nil {
		return errors.New("application not found")
	}

	adminAlliance, err := s.GetUserAlliance(ctx, adminUserID)
	if err != nil {
		return errors.New("not an alliance admin")
	}

	if adminAlliance.ID != app.AllianceID {
		return errors.New("not authorized")
	}

	member := &schema.AllianceMember{
		AllianceID: app.AllianceID,
		UserID:     app.UserID,
		Rank:       "Member",
		JoinedAt:   time.Now(),
	}
	_ = s.allianceRepo.AddMember(ctx, member)

	app.Status = "accepted"
	return s.allianceRepo.UpdateApplication(ctx, app)
}

func (s *AllianceService) RejectApplication(ctx context.Context, adminUserID, applicationID uint) error {
	apps, err := s.allianceRepo.GetApplications(ctx, 0)
	if err != nil {
		return err
	}

	var app *schema.AllianceApplication
	for _, a := range apps {
		if a.ID == applicationID {
			app = a
			break
		}
	}
	if app == nil {
		return errors.New("application not found")
	}

	adminAlliance, err := s.GetUserAlliance(ctx, adminUserID)
	if err != nil {
		return errors.New("not an alliance admin")
	}

	if adminAlliance.ID != app.AllianceID {
		return errors.New("not authorized")
	}

	app.Status = "rejected"
	return s.allianceRepo.UpdateApplication(ctx, app)
}

func (s *AllianceService) LeaveAlliance(ctx context.Context, userID uint) error {
	alliance, err := s.GetUserAlliance(ctx, userID)
	if err != nil {
		return errors.New("not in alliance")
	}

	member, _ := s.allianceRepo.GetMember(ctx, alliance.ID, userID)
	if member != nil && member.Rank == "Founder" {
		return errors.New("founder cannot leave, transfer ownership first")
	}

	return s.allianceRepo.RemoveMember(ctx, alliance.ID, userID)
}

func (s *AllianceService) UpdateAlliance(ctx context.Context, userID uint, input CreateAllianceInput) error {
	alliance, err := s.GetUserAlliance(ctx, userID)
	if err != nil {
		return errors.New("not in alliance")
	}

	member, _ := s.allianceRepo.GetMember(ctx, alliance.ID, userID)
	if member == nil || (member.Rank != "Founder" && member.Rank != "Leader") {
		return errors.New("not authorized")
	}

	if input.Name != "" {
		alliance.Name = input.Name
	}
	if input.Description != "" {
		alliance.Description = input.Description
	}
	if input.Logo != "" {
		alliance.Logo = input.Logo
	}
	if input.Website != "" {
		alliance.Website = input.Website
	}

	alliance.UpdatedAt = time.Now()
	return s.allianceRepo.Update(ctx, alliance)
}
