package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type allianceRepository struct {
	db *gorm.DB
}

func NewAllianceRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) AllianceRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &allianceRepository{db: gdb}
}

func (r *allianceRepository) Create(ctx context.Context, alliance *schema.Alliance) error {
	return r.db.WithContext(ctx).Create(alliance).Error
}

func (r *allianceRepository) GetByID(ctx context.Context, id uint) (*schema.Alliance, error) {
	var alliance schema.Alliance
	err := r.db.WithContext(ctx).First(&alliance, id).Error
	if err != nil {
		return nil, err
	}
	return &alliance, nil
}

func (r *allianceRepository) GetByTag(ctx context.Context, tag string) (*schema.Alliance, error) {
	var alliance schema.Alliance
	err := r.db.WithContext(ctx).Where("tag = ?", tag).First(&alliance).Error
	if err != nil {
		return nil, err
	}
	return &alliance, nil
}

func (r *allianceRepository) GetByName(ctx context.Context, name string) (*schema.Alliance, error) {
	var alliance schema.Alliance
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&alliance).Error
	if err != nil {
		return nil, err
	}
	return &alliance, nil
}

func (r *allianceRepository) GetAll(ctx context.Context) ([]*schema.Alliance, error) {
	var alliances []*schema.Alliance
	err := r.db.WithContext(ctx).Find(&alliances).Error
	return alliances, err
}

func (r *allianceRepository) Update(ctx context.Context, alliance *schema.Alliance) error {
	return r.db.WithContext(ctx).Save(alliance).Error
}

func (r *allianceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.Alliance{}, id).Error
}

func (r *allianceRepository) AddMember(ctx context.Context, member *schema.AllianceMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *allianceRepository) GetMembers(ctx context.Context, allianceID uint) ([]*schema.AllianceMember, error) {
	var members []*schema.AllianceMember
	err := r.db.WithContext(ctx).Where("alliance_id = ?", allianceID).Find(&members).Error
	return members, err
}

func (r *allianceRepository) GetMember(ctx context.Context, allianceID, userID uint) (*schema.AllianceMember, error) {
	var member schema.AllianceMember
	err := r.db.WithContext(ctx).Where("alliance_id = ? AND user_id = ?", allianceID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *allianceRepository) RemoveMember(ctx context.Context, allianceID, userID uint) error {
	return r.db.WithContext(ctx).Where("alliance_id = ? AND user_id = ?", allianceID, userID).Delete(&schema.AllianceMember{}).Error
}

func (r *allianceRepository) CreateApplication(ctx context.Context, app *schema.AllianceApplication) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *allianceRepository) GetApplication(ctx context.Context, allianceID, userID uint) (*schema.AllianceApplication, error) {
	var app schema.AllianceApplication
	err := r.db.WithContext(ctx).Where("alliance_id = ? AND user_id = ? AND status = ?", allianceID, userID, "pending").First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *allianceRepository) GetApplications(ctx context.Context, allianceID uint) ([]*schema.AllianceApplication, error) {
	var apps []*schema.AllianceApplication
	err := r.db.WithContext(ctx).Where("alliance_id = ? AND status = ?", allianceID, "pending").Find(&apps).Error
	return apps, err
}

func (r *allianceRepository) UpdateApplication(ctx context.Context, app *schema.AllianceApplication) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *allianceRepository) DeleteApplication(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.AllianceApplication{}, id).Error
}
