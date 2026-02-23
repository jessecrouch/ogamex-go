package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type debrisRepository struct {
	db *gorm.DB
}

func NewDebrisRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) DebrisRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &debrisRepository{db: gdb}
}

func (r *debrisRepository) Create(ctx context.Context, debris *schema.DebrisField) error {
	return r.db.WithContext(ctx).Create(debris).Error
}

func (r *debrisRepository) GetByID(ctx context.Context, id uint) (*schema.DebrisField, error) {
	var debris schema.DebrisField
	err := r.db.WithContext(ctx).First(&debris, id).Error
	if err != nil {
		return nil, err
	}
	return &debris, nil
}

func (r *debrisRepository) GetByCoords(ctx context.Context, galaxy, system, position int) (*schema.DebrisField, error) {
	var debris schema.DebrisField
	err := r.db.WithContext(ctx).
		Where("galaxy = ? AND system = ? AND position = ?", galaxy, system, position).
		First(&debris).Error
	if err != nil {
		return nil, err
	}
	return &debris, nil
}

func (r *debrisRepository) GetAll(ctx context.Context) ([]*schema.DebrisField, error) {
	var debris []*schema.DebrisField
	err := r.db.WithContext(ctx).Find(&debris).Error
	return debris, err
}

func (r *debrisRepository) GetActive(ctx context.Context) ([]*schema.DebrisField, error) {
	var debris []*schema.DebrisField
	err := r.db.WithContext(ctx).
		Where("expires_at > ?", time.Now()).
		Find(&debris).Error
	return debris, err
}

func (r *debrisRepository) Update(ctx context.Context, debris *schema.DebrisField) error {
	return r.db.WithContext(ctx).Save(debris).Error
}

func (r *debrisRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.DebrisField{}, id).Error
}

func (r *debrisRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&schema.DebrisField{}).Error
}

func (r *debrisRepository) AddResources(ctx context.Context, galaxy, system, position int, metal, crystal int64) error {
	debris, err := r.GetByCoords(ctx, galaxy, system, position)
	if err != nil {
		newDebris := &schema.DebrisField{
			Galaxy:    galaxy,
			System:    system,
			Position:  position,
			Metal:     metal,
			Crystal:   crystal,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(12 * time.Hour),
		}
		return r.Create(ctx, newDebris)
	}

	debris.Metal += metal
	debris.Crystal += crystal
	debris.ExpiresAt = time.Now().Add(12 * time.Hour)
	return r.Update(ctx, debris)
}
