package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type buddyRepository struct {
	db *gorm.DB
}

func NewBuddyRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) BuddyRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &buddyRepository{db: gdb}
}

func (r *buddyRepository) Create(ctx context.Context, buddy *schema.Buddy) error {
	return r.db.WithContext(ctx).Create(buddy).Error
}

func (r *buddyRepository) GetByID(ctx context.Context, id uint) (*schema.Buddy, error) {
	var buddy schema.Buddy
	err := r.db.WithContext(ctx).First(&buddy, id).Error
	if err != nil {
		return nil, err
	}
	return &buddy, nil
}

func (r *buddyRepository) GetByUserID(ctx context.Context, userID uint) ([]*schema.Buddy, error) {
	var buddies []*schema.Buddy
	err := r.db.WithContext(ctx).
		Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "accepted").
		Find(&buddies).Error
	return buddies, err
}

func (r *buddyRepository) GetPendingReceived(ctx context.Context, userID uint) ([]*schema.Buddy, error) {
	var buddies []*schema.Buddy
	err := r.db.WithContext(ctx).
		Where("receiver_id = ? AND status = ?", userID, "pending").
		Find(&buddies).Error
	return buddies, err
}

func (r *buddyRepository) GetPendingSent(ctx context.Context, userID uint) ([]*schema.Buddy, error) {
	var buddies []*schema.Buddy
	err := r.db.WithContext(ctx).
		Where("sender_id = ? AND status = ?", userID, "pending").
		Find(&buddies).Error
	return buddies, err
}

func (r *buddyRepository) GetPendingRequest(ctx context.Context, senderID, receiverID uint) (*schema.Buddy, error) {
	var buddy schema.Buddy
	err := r.db.WithContext(ctx).
		Where("((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)) AND status = ?",
			senderID, receiverID, receiverID, senderID, "pending").
		First(&buddy).Error
	if err != nil {
		return nil, err
	}
	return &buddy, nil
}

func (r *buddyRepository) Update(ctx context.Context, buddy *schema.Buddy) error {
	return r.db.WithContext(ctx).Save(buddy).Error
}

func (r *buddyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.Buddy{}, id).Error
}

func (r *buddyRepository) DeleteByUserIDs(ctx context.Context, userID1, userID2 uint) error {
	return r.db.WithContext(ctx).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1).
		Delete(&schema.Buddy{}).Error
}
