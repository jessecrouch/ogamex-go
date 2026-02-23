package service

import (
	"context"
	"errors"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type BuddyService struct {
	buddyRepo  repository.BuddyRepository
	userRepo   repository.UserRepository
}

func NewBuddyService(buddyRepo repository.BuddyRepository, userRepo repository.UserRepository) *BuddyService {
	return &BuddyService{
		buddyRepo: buddyRepo,
		userRepo:  userRepo,
	}
}

type BuddyInfo struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	PlayerName  string    `json:"player_name"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *BuddyService) GetBuddies(ctx context.Context, userID uint) ([]*BuddyInfo, error) {
	buddies, err := s.buddyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []*BuddyInfo
	for _, b := range buddies {
		otherUserID := b.ReceiverID
		if b.SenderID == userID {
			otherUserID = b.ReceiverID
		}

		user, _ := s.userRepo.GetByID(ctx, otherUserID)
		if user != nil {
			result = append(result, &BuddyInfo{
				ID:          b.ID,
				UserID:      otherUserID,
				Username:    user.Username,
				PlayerName:  user.PlayerName,
				Message:     b.Message,
				Status:      b.Status,
				CreatedAt:   b.CreatedAt,
			})
		}
	}

	return result, nil
}

func (s *BuddyService) GetPendingRequests(ctx context.Context, userID uint) ([]*BuddyInfo, error) {
	buddies, err := s.buddyRepo.GetPendingReceived(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []*BuddyInfo
	for _, b := range buddies {
		user, _ := s.userRepo.GetByID(ctx, b.SenderID)
		result = append(result, &BuddyInfo{
			ID:          b.ID,
			UserID:      b.SenderID,
			Username:    user.Username,
			PlayerName:  user.PlayerName,
			Message:     b.Message,
			Status:      b.Status,
			CreatedAt:   b.CreatedAt,
		})
	}

	return result, nil
}

func (s *BuddyService) SendRequest(ctx context.Context, senderID uint, receiverID uint, message string) error {
	if senderID == receiverID {
		return errors.New("cannot send buddy request to yourself")
	}

	receiver, err := s.userRepo.GetByID(ctx, receiverID)
	if err != nil {
		return errors.New("receiver not found")
	}
	_ = receiver

	existing, _ := s.buddyRepo.GetPendingRequest(ctx, senderID, receiverID)
	if existing != nil {
		return errors.New("buddy request already pending")
	}

	buddy := &schema.Buddy{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	return s.buddyRepo.Create(ctx, buddy)
}

func (s *BuddyService) AcceptRequest(ctx context.Context, userID, requestID uint) error {
	buddy, err := s.buddyRepo.GetByID(ctx, requestID)
	if err != nil {
		return errors.New("buddy request not found")
	}

	if buddy.ReceiverID != userID {
		return errors.New("not authorized")
	}

	if buddy.Status != "pending" {
		return errors.New("request already processed")
	}

	buddy.Status = "accepted"
	return s.buddyRepo.Update(ctx, buddy)
}

func (s *BuddyService) RejectRequest(ctx context.Context, userID, requestID uint) error {
	buddy, err := s.buddyRepo.GetByID(ctx, requestID)
	if err != nil {
		return errors.New("buddy request not found")
	}

	if buddy.ReceiverID != userID {
		return errors.New("not authorized")
	}

	if buddy.Status != "pending" {
		return errors.New("request already processed")
	}

	buddy.Status = "rejected"
	return s.buddyRepo.Update(ctx, buddy)
}

func (s *BuddyService) RemoveBuddy(ctx context.Context, userID, buddyID uint) error {
	buddy, err := s.buddyRepo.GetByID(ctx, buddyID)
	if err != nil {
		return errors.New("buddy not found")
	}

	if buddy.SenderID != userID && buddy.ReceiverID != userID {
		return errors.New("not authorized")
	}

	return s.buddyRepo.Delete(ctx, buddyID)
}
