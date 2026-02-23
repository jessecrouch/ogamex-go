package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type planetRepository struct {
	db *gorm.DB
}

func NewPlanetRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) PlanetRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &planetRepository{db: gdb}
}

func (r *planetRepository) Create(ctx context.Context, planet *schema.Planet) error {
	return r.db.WithContext(ctx).Create(planet).Error
}

func (r *planetRepository) GetByID(ctx context.Context, id uint) (*schema.Planet, error) {
	var planet schema.Planet
	err := r.db.WithContext(ctx).First(&planet, id).Error
	if err != nil {
		return nil, err
	}
	return &planet, nil
}

func (r *planetRepository) GetByCoords(ctx context.Context, userID uint, galaxy, system, position int) (*schema.Planet, error) {
	var planet schema.Planet
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND galaxy = ? AND system = ? AND position = ?", userID, galaxy, system, position).
		First(&planet).Error
	if err != nil {
		return nil, err
	}
	return &planet, nil
}

func (r *planetRepository) GetByUserID(ctx context.Context, userID uint) ([]*schema.Planet, error) {
	var planets []*schema.Planet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&planets).Error
	return planets, err
}

func (r *planetRepository) GetAll(ctx context.Context) ([]*schema.Planet, error) {
	var planets []*schema.Planet
	err := r.db.WithContext(ctx).Find(&planets).Error
	return planets, err
}

func (r *planetRepository) Update(ctx context.Context, planet *schema.Planet) error {
	return r.db.WithContext(ctx).Save(planet).Error
}

func (r *planetRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.Planet{}, id).Error
}

func (r *planetRepository) AddResources(ctx context.Context, id uint, metal, crystal, deuterium int64) error {
	return r.db.WithContext(ctx).Model(&schema.Planet{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"metal":     gorm.Expr("metal + ?", metal),
			"crystal":   gorm.Expr("crystal + ?", crystal),
			"deuterium": gorm.Expr("deuterium + ?", deuterium),
		}).Error
}

func (r *planetRepository) SubResources(ctx context.Context, id uint, metal, crystal, deuterium int64) error {
	return r.db.WithContext(ctx).Model(&schema.Planet{}).
		Where("id = ? AND metal >= ? AND crystal >= ? AND deuterium >= ?", id, metal, crystal, deuterium).
		Updates(map[string]interface{}{
			"metal":     gorm.Expr("metal - ?", metal),
			"crystal":   gorm.Expr("crystal - ?", crystal),
			"deuterium": gorm.Expr("deuterium - ?", deuterium),
		}).Error
}
