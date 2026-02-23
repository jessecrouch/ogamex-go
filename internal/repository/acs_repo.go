package repository

import (
	"context"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type ACSRepository interface {
	Create(ctx context.Context, acs *schema.ACS) error
	GetByID(ctx context.Context, id uint) (*schema.ACS, error)
	GetByFleetID(ctx context.Context, fleetID uint) (*schema.ACS, error)
	AddFleet(ctx context.Context, acsID uint, fleetID uint) error
	GetFleets(ctx context.Context, acsID uint) ([]*schema.FleetMission, error)
	Delete(ctx context.Context, id uint) error
}

type acsRepository struct {
	db *gorm.DB
}

func NewACSRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) ACSRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &acsRepository{db: gdb}
}

func (r *acsRepository) Create(ctx context.Context, acs *schema.ACS) error {
	return r.db.WithContext(ctx).Create(acs).Error
}

func (r *acsRepository) GetByID(ctx context.Context, id uint) (*schema.ACS, error) {
	var acs schema.ACS
	err := r.db.WithContext(ctx).First(&acs, id).Error
	if err != nil {
		return nil, err
	}
	return &acs, nil
}

func (r *acsRepository) GetByFleetID(ctx context.Context, fleetID uint) (*schema.ACS, error) {
	var fleet schema.FleetMission
	err := r.db.WithContext(ctx).First(&fleet, fleetID).Error
	if err != nil {
		return nil, err
	}
	if fleet.ACSId == nil {
		return nil, nil
	}
	return r.GetByID(ctx, *fleet.ACSId)
}

func (r *acsRepository) AddFleet(ctx context.Context, acsID uint, fleetID uint) error {
	acs, err := r.GetByID(ctx, acsID)
	if err != nil {
		return err
	}
	acs.FleetIDs = append(acs.FleetIDs, fleetID)
	return r.db.WithContext(ctx).Save(acs).Error
}

func (r *acsRepository) GetFleets(ctx context.Context, acsID uint) ([]*schema.FleetMission, error) {
	var fleets []*schema.FleetMission
	err := r.db.WithContext(ctx).Where("acs_id = ?", acsID).Find(&fleets).Error
	return fleets, err
}

func (r *acsRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.ACS{}, id).Error
}
