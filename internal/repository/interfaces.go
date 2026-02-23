package repository

import (
	"context"

	"ogamex-go/internal/schema"
)

type UserRepository interface {
	Create(ctx context.Context, user *schema.User) error
	GetByID(ctx context.Context, id uint) (*schema.User, error)
	GetByUsername(ctx context.Context, username string) (*schema.User, error)
	GetByEmail(ctx context.Context, email string) (*schema.User, error)
	GetByAuthToken(ctx context.Context, token string) (*schema.User, error)
	Update(ctx context.Context, user *schema.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*schema.User, error)
}

type PlanetRepository interface {
	Create(ctx context.Context, planet *schema.Planet) error
	GetByID(ctx context.Context, id uint) (*schema.Planet, error)
	GetByCoords(ctx context.Context, userID uint, galaxy, system, position int) (*schema.Planet, error)
	GetByCoordsAny(ctx context.Context, galaxy, system, position int) (*schema.Planet, error)
	GetByUserID(ctx context.Context, userID uint) ([]*schema.Planet, error)
	GetAll(ctx context.Context) ([]*schema.Planet, error)
	Update(ctx context.Context, planet *schema.Planet) error
	Delete(ctx context.Context, id uint) error
	AddResources(ctx context.Context, id uint, metal, crystal, deuterium int64) error
	SubResources(ctx context.Context, id uint, metal, crystal, deuterium int64) error
}

type UserTechRepository interface {
	GetByUserID(ctx context.Context, userID uint) (*schema.UserTech, error)
	Create(ctx context.Context, tech *schema.UserTech) error
	Update(ctx context.Context, tech *schema.UserTech) error
}

type FleetMissionRepository interface {
	Create(ctx context.Context, mission *schema.FleetMission) error
	GetByID(ctx context.Context, id uint) (*schema.FleetMission, error)
	GetByUserID(ctx context.Context, userID uint) ([]*schema.FleetMission, error)
	GetActive(ctx context.Context) ([]*schema.FleetMission, error)
	GetArriving(ctx context.Context, before interface{}) ([]*schema.FleetMission, error)
	GetReturning(ctx context.Context, before interface{}) ([]*schema.FleetMission, error)
	Update(ctx context.Context, mission *schema.FleetMission) error
	Delete(ctx context.Context, id uint) error
}

type BuildingQueueRepository interface {
	Create(ctx context.Context, queue *schema.BuildingQueue) error
	GetByID(ctx context.Context, id uint) (*schema.BuildingQueue, error)
	GetByPlanetID(ctx context.Context, planetID uint) ([]*schema.BuildingQueue, error)
	GetCurrent(ctx context.Context, planetID uint) (*schema.BuildingQueue, error)
	GetCompleted(ctx context.Context, planetID uint) ([]*schema.BuildingQueue, error)
	GetAllWithActive(ctx context.Context) ([]*schema.BuildingQueue, error)
	Update(ctx context.Context, queue *schema.BuildingQueue) error
	Delete(ctx context.Context, id uint) error
	Cancel(ctx context.Context, id uint) error
}

type ResearchQueueRepository interface {
	Create(ctx context.Context, queue *schema.ResearchQueue) error
	GetByID(ctx context.Context, id uint) (*schema.ResearchQueue, error)
	GetByUserID(ctx context.Context, userID uint) ([]*schema.ResearchQueue, error)
	GetCurrent(ctx context.Context, userID uint) (*schema.ResearchQueue, error)
	GetAllWithActive(ctx context.Context) ([]*schema.ResearchQueue, error)
	Update(ctx context.Context, queue *schema.ResearchQueue) error
	Delete(ctx context.Context, id uint) error
	Cancel(ctx context.Context, id uint) error
}

type UnitQueueRepository interface {
	Create(ctx context.Context, queue *schema.UnitQueue) error
	GetByID(ctx context.Context, id uint) (*schema.UnitQueue, error)
	GetByPlanetID(ctx context.Context, planetID uint) ([]*schema.UnitQueue, error)
	GetCurrent(ctx context.Context, planetID uint) (*schema.UnitQueue, error)
	GetAllWithActive(ctx context.Context) ([]*schema.UnitQueue, error)
	Update(ctx context.Context, queue *schema.UnitQueue) error
	Delete(ctx context.Context, id uint) error
}
