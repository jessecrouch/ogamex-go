package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) MessageRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &messageRepository{db: gdb}
}

func (r *messageRepository) Create(ctx context.Context, message *schema.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) GetByID(ctx context.Context, id uint) (*schema.Message, error) {
	var message schema.Message
	err := r.db.WithContext(ctx).First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*schema.Message, error) {
	var messages []*schema.Message
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&schema.Message{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *messageRepository) MarkAsRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&schema.Message{}).Where("id = ?", id).Update("read", true).Error
}

func (r *messageRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&schema.Message{}).Where("user_id = ?", userID).Update("read", true).Error
}

func (r *messageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.Message{}, id).Error
}

func (r *messageRepository) DeleteAll(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&schema.Message{}).Error
}
