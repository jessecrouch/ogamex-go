package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ogamex-go/internal/database"
	"ogamex-go/internal/schema"
)

type buildingQueueRepository struct {
	db *gorm.DB
}

func NewBuildingQueueRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) BuildingQueueRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &buildingQueueRepository{db: gdb}
}

func (r *buildingQueueRepository) Create(ctx context.Context, queue *schema.BuildingQueue) error {
	return r.db.WithContext(ctx).Create(queue).Error
}

func (r *buildingQueueRepository) GetByID(ctx context.Context, id uint) (*schema.BuildingQueue, error) {
	var queue schema.BuildingQueue
	err := r.db.WithContext(ctx).First(&queue, id).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *buildingQueueRepository) GetByPlanetID(ctx context.Context, planetID uint) ([]*schema.BuildingQueue, error) {
	var queues []*schema.BuildingQueue
	err := r.db.WithContext(ctx).Where("planet_id = ?", planetID).Find(&queues).Error
	return queues, err
}

func (r *buildingQueueRepository) GetCurrent(ctx context.Context, planetID uint) (*schema.BuildingQueue, error) {
	var queue schema.BuildingQueue
	err := r.db.WithContext(ctx).
		Where("planet_id = ? AND is_cancelled = ? AND end_time > ?", planetID, false, time.Now()).
		Order("start_time ASC").
		First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *buildingQueueRepository) GetCompleted(ctx context.Context, planetID uint) ([]*schema.BuildingQueue, error) {
	var queues []*schema.BuildingQueue
	err := r.db.WithContext(ctx).
		Where("planet_id = ? AND is_cancelled = ? AND end_time <= ?", planetID, false, time.Now()).
		Order("start_time ASC").
		Find(&queues).Error
	return queues, err
}

func (r *buildingQueueRepository) GetAllWithActive(ctx context.Context) ([]*schema.BuildingQueue, error) {
	var queues []*schema.BuildingQueue
	err := r.db.WithContext(ctx).
		Where("is_cancelled = ?", false).
		Order("start_time ASC").
		Find(&queues).Error
	return queues, err
}

func (r *buildingQueueRepository) Update(ctx context.Context, queue *schema.BuildingQueue) error {
	return r.db.WithContext(ctx).Save(queue).Error
}

func (r *buildingQueueRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.BuildingQueue{}, id).Error
}

func (r *buildingQueueRepository) Cancel(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&schema.BuildingQueue{}).
		Where("id = ?", id).
		Update("is_cancelled", true).Error
}

type researchQueueRepository struct {
	db *gorm.DB
}

func NewResearchQueueRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) ResearchQueueRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &researchQueueRepository{db: gdb}
}

func (r *researchQueueRepository) Create(ctx context.Context, queue *schema.ResearchQueue) error {
	return r.db.WithContext(ctx).Create(queue).Error
}

func (r *researchQueueRepository) GetByID(ctx context.Context, id uint) (*schema.ResearchQueue, error) {
	var queue schema.ResearchQueue
	err := r.db.WithContext(ctx).First(&queue, id).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *researchQueueRepository) GetByUserID(ctx context.Context, userID uint) ([]*schema.ResearchQueue, error) {
	var queues []*schema.ResearchQueue
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&queues).Error
	return queues, err
}

func (r *researchQueueRepository) GetCurrent(ctx context.Context, userID uint) (*schema.ResearchQueue, error) {
	var queue schema.ResearchQueue
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND end_time > ?", userID, time.Now()).
		Order("start_time ASC").
		First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *researchQueueRepository) GetAllWithActive(ctx context.Context) ([]*schema.ResearchQueue, error) {
	var queues []*schema.ResearchQueue
	err := r.db.WithContext(ctx).
		Order("start_time ASC").
		Find(&queues).Error
	return queues, err
}

func (r *researchQueueRepository) Update(ctx context.Context, queue *schema.ResearchQueue) error {
	return r.db.WithContext(ctx).Save(queue).Error
}

func (r *researchQueueRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.ResearchQueue{}, id).Error
}

func (r *researchQueueRepository) Cancel(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.ResearchQueue{}, id).Error
}

type unitQueueRepository struct {
	db *gorm.DB
}

func NewUnitQueueRepository(db interface{ WithContext(ctx context.Context) *gorm.DB }) UnitQueueRepository {
	var gdb *gorm.DB
	switch d := db.(type) {
	case *gorm.DB:
		gdb = d
	case *database.Database:
		gdb = d.DB
	}
	return &unitQueueRepository{db: gdb}
}

func (r *unitQueueRepository) Create(ctx context.Context, queue *schema.UnitQueue) error {
	return r.db.WithContext(ctx).Create(queue).Error
}

func (r *unitQueueRepository) GetByID(ctx context.Context, id uint) (*schema.UnitQueue, error) {
	var queue schema.UnitQueue
	err := r.db.WithContext(ctx).First(&queue, id).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *unitQueueRepository) GetByPlanetID(ctx context.Context, planetID uint) ([]*schema.UnitQueue, error) {
	var queues []*schema.UnitQueue
	err := r.db.WithContext(ctx).Where("planet_id = ?", planetID).Find(&queues).Error
	return queues, err
}

func (r *unitQueueRepository) GetCurrent(ctx context.Context, planetID uint) (*schema.UnitQueue, error) {
	var queue schema.UnitQueue
	err := r.db.WithContext(ctx).
		Where("planet_id = ? AND end_time > ?", planetID, time.Now()).
		Order("start_time ASC").
		First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *unitQueueRepository) Update(ctx context.Context, queue *schema.UnitQueue) error {
	return r.db.WithContext(ctx).Save(queue).Error
}

func (r *unitQueueRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&schema.UnitQueue{}, id).Error
}
