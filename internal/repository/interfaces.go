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
	GetNPCUsers(ctx context.Context) ([]*schema.User, error)
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
	GetBySystem(ctx context.Context, galaxy, system int) ([]*schema.Planet, error)
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

type MessageRepository interface {
	Create(ctx context.Context, message *schema.Message) error
	GetByID(ctx context.Context, id uint) (*schema.Message, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*schema.Message, error)
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	Delete(ctx context.Context, id uint) error
	DeleteAll(ctx context.Context, userID uint) error
}

type NoteRepository interface {
	Create(ctx context.Context, note *schema.Note) error
	GetByID(ctx context.Context, id uint) (*schema.Note, error)
	GetByUserID(ctx context.Context, userID uint) ([]*schema.Note, error)
	Update(ctx context.Context, note *schema.Note) error
	Delete(ctx context.Context, id uint) error
}

type AllianceRepository interface {
	Create(ctx context.Context, alliance *schema.Alliance) error
	GetByID(ctx context.Context, id uint) (*schema.Alliance, error)
	GetByTag(ctx context.Context, tag string) (*schema.Alliance, error)
	GetByName(ctx context.Context, name string) (*schema.Alliance, error)
	GetAll(ctx context.Context) ([]*schema.Alliance, error)
	Update(ctx context.Context, alliance *schema.Alliance) error
	Delete(ctx context.Context, id uint) error
	AddMember(ctx context.Context, member *schema.AllianceMember) error
	GetMembers(ctx context.Context, allianceID uint) ([]*schema.AllianceMember, error)
	GetMember(ctx context.Context, allianceID, userID uint) (*schema.AllianceMember, error)
	RemoveMember(ctx context.Context, allianceID, userID uint) error
	CreateApplication(ctx context.Context, app *schema.AllianceApplication) error
	GetApplication(ctx context.Context, allianceID, userID uint) (*schema.AllianceApplication, error)
	GetApplications(ctx context.Context, allianceID uint) ([]*schema.AllianceApplication, error)
	UpdateApplication(ctx context.Context, app *schema.AllianceApplication) error
	DeleteApplication(ctx context.Context, id uint) error
}

type BuddyRepository interface {
	Create(ctx context.Context, buddy *schema.Buddy) error
	GetByID(ctx context.Context, id uint) (*schema.Buddy, error)
	GetByUserID(ctx context.Context, userID uint) ([]*schema.Buddy, error)
	GetPendingReceived(ctx context.Context, userID uint) ([]*schema.Buddy, error)
	GetPendingSent(ctx context.Context, userID uint) ([]*schema.Buddy, error)
	GetPendingRequest(ctx context.Context, senderID, receiverID uint) (*schema.Buddy, error)
	Update(ctx context.Context, buddy *schema.Buddy) error
	Delete(ctx context.Context, id uint) error
	DeleteByUserIDs(ctx context.Context, userID1, userID2 uint) error
}

type EspionageRepository interface {
	Create(ctx context.Context, report *schema.EspionageReport) error
	GetByID(ctx context.Context, id uint) (*schema.EspionageReport, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*schema.EspionageReport, error)
	GetUnreadCount(ctx context.Context, userID uint) (int64, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	Delete(ctx context.Context, id uint) error
	DeleteAll(ctx context.Context, userID uint) error
}

type DebrisRepository interface {
	Create(ctx context.Context, debris *schema.DebrisField) error
	GetByID(ctx context.Context, id uint) (*schema.DebrisField, error)
	GetByCoords(ctx context.Context, galaxy, system, position int) (*schema.DebrisField, error)
	GetAll(ctx context.Context) ([]*schema.DebrisField, error)
	GetActive(ctx context.Context) ([]*schema.DebrisField, error)
	Update(ctx context.Context, debris *schema.DebrisField) error
	Delete(ctx context.Context, id uint) error
	DeleteExpired(ctx context.Context) error
	AddResources(ctx context.Context, galaxy, system, position int, metal, crystal int64) error
}

type WreckRepository interface {
	Create(ctx context.Context, wreck *schema.WreckField) error
	GetByID(ctx context.Context, id uint) (*schema.WreckField, error)
	GetByCoords(ctx context.Context, galaxy, system, position int) (*schema.WreckField, error)
	GetAll(ctx context.Context) ([]*schema.WreckField, error)
	GetActive(ctx context.Context) ([]*schema.WreckField, error)
	Update(ctx context.Context, wreck *schema.WreckField) error
	Delete(ctx context.Context, id uint) error
	DeleteExpired(ctx context.Context) error
	CreateOrUpdate(ctx context.Context, galaxy, system, position int, metal, crystal, deuterium int64) error
}
