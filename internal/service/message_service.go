package service

import (
	"context"
	"time"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
)

type MessageService struct {
	messageRepo repository.MessageRepository
	userRepo    repository.UserRepository
}

func NewMessageService(
	messageRepo repository.MessageRepository,
	userRepo repository.UserRepository,
) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

const (
	MessageTypeGeneric     = 0
	MessageTypeBattle      = 1
	MessageTypeAlliance    = 2
	MessageTypeCommerce    = 3
	MessageTypeEspionage   = 4
	MessageTypeExpedition  = 5
	MessageTypeColony      = 6
	MessageTypeTransport  = 7
	MessageTypeSystem     = 8
)

func (s *MessageService) SendMessage(ctx context.Context, userID uint, msgType int, subject, body string, fromUserID *uint, fromName string) error {
	message := &schema.Message{
		UserID:    userID,
		Type:      msgType,
		Subject:   subject,
		Body:      body,
		FromUserID: fromUserID,
		FromName:  fromName,
		Read:      false,
		CreatedAt: time.Now(),
	}

	return s.messageRepo.Create(ctx, message)
}

func (s *MessageService) GetMessages(ctx context.Context, userID uint, limit, offset int) ([]*schema.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return s.messageRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *MessageService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.messageRepo.GetUnreadCount(ctx, userID)
}

func (s *MessageService) MarkAsRead(ctx context.Context, messageID uint) error {
	return s.messageRepo.MarkAsRead(ctx, messageID)
}

func (s *MessageService) MarkAllAsRead(ctx context.Context, userID uint) error {
	return s.messageRepo.MarkAllAsRead(ctx, userID)
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID uint) error {
	return s.messageRepo.Delete(ctx, messageID)
}

func (s *MessageService) DeleteAllMessages(ctx context.Context, userID uint) error {
	return s.messageRepo.DeleteAll(ctx, userID)
}

func (s *MessageService) SendBattleReport(ctx context.Context, userID uint, subject, body string) error {
	return s.SendMessage(ctx, userID, MessageTypeBattle, subject, body, nil, "System")
}

func (s *MessageService) SendSystemMessage(ctx context.Context, userID uint, subject, body string) error {
	return s.SendMessage(ctx, userID, MessageTypeSystem, subject, body, nil, "System")
}

func (s *MessageService) SendEspionageReport(ctx context.Context, userID uint, subject, body string) error {
	return s.SendMessage(ctx, userID, MessageTypeEspionage, subject, body, nil, "System")
}
