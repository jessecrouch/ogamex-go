package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type wreckRepository struct {
	db *gorm.DB
}

func NewWreckRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) WreckRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &wreckRepository{db: gdb}
}

func (r *wreckRepository) Create(ctx context.Context, wreck *schema.WreckField) error {
	return r.db.WithContext(ctx).Create(wreck).Error
}

func (r *wreckRepository) GetByID(ctx context.Context, id uint) (*schema.WreckField, error) {
	var wreck schema.WreckField
	err := r.db.WithContext(ctx).First(&wreck, id).Error
	if err != nil {
		return nil, err
	}
	return &wreck, nil
}

func (r *wreckRepository) GetByCoords(ctx context.Context, galaxy, system, position int) (*schema.WreckField, error) {
	var wreck schema.WreckField
	err := r.db.WithContext(ctx).
		Where("galaxy = ? AND system = ? AND position = ?", galaxy, system, position).
		First(&wreck).Error
	if err != nil {
		return nil, err
	}
	return &wreck, nil
}

func (r *wreckRepository) GetAll(ctx context.Context) ([]*schema.WreckField, error) {
	var wrecks []*schema.WreckField
	err := r.db.WithContext(ctx).Find(&wrecks).Error
	return wrecks, err
}

func (r *wreckRepository) GetActive(ctx context.Context) ([]*schema.WreckField, error) {
	var wrecks []*schema.WreckField
	err := r.db.WithContext(ctx).
		Where("expires_at > ?", time.Now()).
		Find(&wrecks).Error
	return wrecks, err
}

func (r *wreckRepository) Update(ctx context.Context, wreck *schema.WreckField) error {
	return r.db.WithContext(ctx).Save(wreck).Error
}

func (r *wreckRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.WreckField{}, id).Error
}

func (r *wreckRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&schema.WreckField{}).Error
}

func (r *wreckRepository) CreateOrUpdate(ctx context.Context, galaxy, system, position int, metal, crystal, deuterium int64) error {
	wreck, err := r.GetByCoords(ctx, galaxy, system, position)
	if err != nil {
		newWreck := &schema.WreckField{
			Galaxy:    galaxy,
			System:    system,
			Position:  position,
			Metal:     metal,
			Crystal:   crystal,
			Deuterium: deuterium,
			CreatedAt: time.Now(),
			ExpiresAt:  time.Now().Add(12 * time.Hour),
		}
		return r.Create(ctx, newWreck)
	}

	wreck.Metal += metal
	wreck.Crystal += crystal
	wreck.Deuterium += deuterium
	wreck.ExpiresAt = time.Now().Add(12 * time.Hour)
	return r.Update(ctx, wreck)
}
