package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type fleetMissionRepository struct {
	db *gorm.DB
}

func NewFleetMissionRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) FleetMissionRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &fleetMissionRepository{db: gdb}
}

func (r *fleetMissionRepository) Create(ctx context.Context, mission *schema.FleetMission) error {
	return r.db.WithContext(ctx).Create(mission).Error
}

func (r *fleetMissionRepository) GetByID(ctx context.Context, id uint) (*schema.FleetMission, error) {
	var mission schema.FleetMission
	err := r.db.WithContext(ctx).First(&mission, id).Error
	if err != nil {
		return nil, err
	}
	return &mission, nil
}

func (r *fleetMissionRepository) GetByUserID(ctx context.Context, userID uint) ([]*schema.FleetMission, error) {
	var missions []*schema.FleetMission
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&missions).Error
	return missions, err
}

func (r *fleetMissionRepository) GetActive(ctx context.Context) ([]*schema.FleetMission, error) {
	var missions []*schema.FleetMission
	err := r.db.WithContext(ctx).Where("status IN ?", []int{1, 2}).Find(&missions).Error
	return missions, err
}

func (r *fleetMissionRepository) GetArriving(ctx context.Context, before interface{}) ([]*schema.FleetMission, error) {
	var missions []*schema.FleetMission
	query := r.db.WithContext(ctx).Where("status = ?", 1)
	
	switch t := before.(type) {
	case time.Time:
		query = query.Where("arrival_time <= ?", t)
	case int64:
		query = query.Where("arrival_time <= ?", time.Unix(t, 0))
	}
	
	err := query.Find(&missions).Error
	return missions, err
}

func (r *fleetMissionRepository) Update(ctx context.Context, mission *schema.FleetMission) error {
	return r.db.WithContext(ctx).Save(mission).Error
}

func (r *fleetMissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.FleetMission{}, id).Error
}
