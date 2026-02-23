package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) UserRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &userRepository{db: gdb}
}

func (r *userRepository) Create(ctx context.Context, user *schema.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*schema.User, error) {
	var user schema.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*schema.User, error) {
	var user schema.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*schema.User, error) {
	var user schema.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByAuthToken(ctx context.Context, token string) (*schema.User, error) {
	var user schema.User
	err := r.db.WithContext(ctx).Where("auth_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetNPCUsers(ctx context.Context) ([]*schema.User, error) {
	var users []*schema.User
	err := r.db.WithContext(ctx).Where("is_npc = ?", true).Find(&users).Error
	return users, err
}

func (r *userRepository) Update(ctx context.Context, user *schema.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*schema.User, error) {
	var users []*schema.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
