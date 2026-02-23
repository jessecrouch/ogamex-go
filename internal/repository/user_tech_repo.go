package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type userTechRepository struct {
	db *gorm.DB
}

func NewUserTechRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) UserTechRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &userTechRepository{db: gdb}
}

func (r *userTechRepository) GetByUserID(ctx context.Context, userID uint) (*schema.UserTech, error) {
	var tech schema.UserTech
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&tech).Error
	if err != nil {
		return nil, err
	}
	return &tech, nil
}

func (r *userTechRepository) Create(ctx context.Context, tech *schema.UserTech) error {
	return r.db.WithContext(ctx).Create(tech).Error
}

func (r *userTechRepository) Update(ctx context.Context, tech *schema.UserTech) error {
	return r.db.WithContext(ctx).Save(tech).Error
}
