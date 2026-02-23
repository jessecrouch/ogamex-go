package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type espionageRepository struct {
	db *gorm.DB
}

func NewEspionageRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) EspionageRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &espionageRepository{db: gdb}
}

func (r *espionageRepository) Create(ctx context.Context, report *schema.EspionageReport) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *espionageRepository) GetByID(ctx context.Context, id uint) (*schema.EspionageReport, error) {
	var report schema.EspionageReport
	err := r.db.WithContext(ctx).First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *espionageRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*schema.EspionageReport, error) {
	var reports []*schema.EspionageReport
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&reports).Error
	return reports, err
}

func (r *espionageRepository) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&schema.EspionageReport{}).
		Where("user_id = ? AND read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *espionageRepository) MarkAsRead(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&schema.EspionageReport{}).Where("id = ?", id).Update("read", true).Error
}

func (r *espionageRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&schema.EspionageReport{}).Where("user_id = ?", userID).Update("read", true).Error
}

func (r *espionageRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.EspionageReport{}, id).Error
}

func (r *espionageRepository) DeleteAll(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&schema.EspionageReport{}).Error
}
