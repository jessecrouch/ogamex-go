package schema

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username        string         `gorm:"size:64;uniqueIndex" json:"username"`
	Email           string         `gorm:"size:255;uniqueIndex" json:"email"`
	Password        string         `gorm:"size:255" json:"-"`
	PlayerName      string         `gorm:"size:64" json:"player_name"`
	AuthToken       string         `gorm:"size:64;uniqueIndex" json:"auth_token,omitempty"`
	CharacterClass  int           `gorm:"default:0" json:"character_class"`
	DarkMatter      int64          `gorm:"default:0" json:"dark_matter"`
	CurrentPlanetID *uint          `json:"current_planet_id"`
	OnVacation      bool           `gorm:"default:false" json:"on_vacation"`
	VacationEndTime *time.Time     `json:"vacation_end_time"`
	PremiumEndsAt  *time.Time     `json:"premium_ends_at"`
	LastOnline      time.Time      `json:"last_online"`
	RegisteredAt    time.Time      `json:"registered_at"`
	IsNPC           bool           `gorm:"default:false" json:"is_npc"`
}

func (User) TableName() string {
	return "users"
}

type Planet struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserID       uint           `gorm:"index" json:"user_id"`
	Name         string         `gorm:"size:100" json:"name"`
	Galaxy       int            `gorm:"index" json:"galaxy"`
	System       int            `gorm:"index" json:"system"`
	Position     int            `gorm:"index" json:"position"`
	IsMoon       bool           `gorm:"default:false" json:"is_moon"`
	PlanetType   int            `gorm:"default:1" json:"planet_type"`

	Metal     int64 `json:"metal"`
	Crystal   int64 `json:"crystal"`
	Deuterium int64 `json:"deuterium"`

	MetalCapacity     int64 `json:"metal_capacity"`
	CrystalCapacity   int64 `json:"crystal_capacity"`
	DeuteriumCapacity int64 `json:"deuterium_capacity"`

	MetalProduction     int64 `json:"metal_production"`
	CrystalProduction   int64 `json:"crystal_production"`
	DeuteriumProduction int64 `json:"deuterium_production"`

	EnergyAvailable int64 `json:"energy_available"`
	EnergyMax       int64 `json:"energy_max"`
	EnergyUsed      int64 `json:"energy_used"`

	TempMin      int `json:"temp_min"`
	TempMax      int `json:"temp_max"`
	FieldsUsed   int `json:"fields_used"`
	FieldsMax    int `json:"fields_max"`

	MetalMinePercent           int `json:"metal_mine_percent"`
	CrystalMinePercent         int `json:"crystal_mine_percent"`
	DeuteriumSynthesizerPercent int `json:"deuterium_synthesizer_percent"`
	SolarPlantPercent          int `json:"solar_plant_percent"`
	FusionPlantPercent         int `json:"fusion_plant_percent"`

	MetalMine           int `json:"metal_mine"`
	CrystalMine         int `json:"crystal_mine"`
	DeuteriumSynthesizer int `json:"deuterium_synthesizer"`
	SolarPlant          int `json:"solar_plant"`
	FusionPlant         int `json:"fusion_plant"`
	MetalStorageBuilding        int `json:"metal_store"`
	CrystalStorageBuilding      int `json:"crystal_store"`
	DeuteriumStorageBuilding    int `json:"deuterium_store"`

	RobotFactory      int `json:"robot_factory"`
	Shipyard          int `json:"shipyard"`
	ResearchLab       int `json:"research_lab"`
	MissileSilo       int `json:"missile_silo"`
	NaniteFactory     int `json:"nano_factory"`
	Terraformer       int `json:"terraformer"`
	SpaceDock         int `json:"space_dock"`
	LunarBase         int            `json:"lunar_base"`
	SensorPhalanx     int            `json:"sensor_phalanx"`
	JumpGate          int            `json:"jump_gate"`

	SolarSatellite   int            `json:"solar_satellite"`
	Crawler          int            `json:"crawler"`

	SmallCargo       int            `json:"small_cargo"`
	LargeCargo       int            `json:"large_cargo"`
	LightFighter    int            `json:"light_fighter"`
	HeavyFighter    int            `json:"heavy_fighter"`
	Cruiser          int            `json:"cruiser"`
	Battleship       int            `json:"battleship"`
	ColonyShip       int            `json:"colony_ship"`
	Recycler         int            `json:"recycler"`
	EspionageProbe   int            `json:"espionage_probe"`
	Bomber           int            `json:"bomber"`
	Destroyer        int            `json:"destroyer"`
	Deathstar        int            `json:"deathstar"`
	Battlecruiser    int            `json:"battlecruiser"`
	Reaper           int            `json:"reaper"`
	Pathfinder       int            `json:"pathfinder"`

	RocketLauncher   int            `json:"rocket_launcher"`
	LightLaser      int            `json:"light_laser"`
	HeavyLaser      int            `json:"heavy_laser"`
	IonCannon       int            `json:"ion_cannon"`
	GaussCannon     int            `json:"gauss_cannon"`
	PlasmaTurret    int            `json:"plasma_turret"`
	ShieldDome      int            `json:"shield_dome"`
	MissileInterceptor int         `json:"missile_interceptor"`
	MissileLauncher int            `json:"missile_launcher"`

	DefenseActivated bool          `gorm:"default:false" json:"defense_activated"`

	RapidfireFrom   string         `json:"rapidfire_from"`
	RapidfireTo     string         `json:"rapidfire_to"`
}

func (Planet) TableName() string {
	return "planets"
}

type UserTech struct {
	ID                      uint           `gorm:"primaryKey" json:"id"`
	UserID                  uint           `gorm:"uniqueIndex" json:"user_id"`
	EnergyTechnology        int            `gorm:"default:0" json:"energy_technology"`
	LaserTechnology         int            `gorm:"default:0" json:"laser_technology"`
	IonTechnology           int            `gorm:"default:0" json:"ion_technology"`
	HyperspaceTechnology    int            `gorm:"default:0" json:"hyperspace_technology"`
	PlasmaTechnology        int            `gorm:"default:0" json:"plasma_technology"`
	CombatDrive             int            `gorm:"default:0" json:"combat_drive"`
	ImpulseDrive            int            `gorm:"default:0" json:"impulse_drive"`
	HyperspaceDrive         int            `gorm:"default:0" json:"hyperspace_drive"`
	EspionageTechnology     int            `gorm:"default:0" json:"espionage_technology"`
	ComputerTechnology      int            `gorm:"default:0" json:"computer_technology"`
	Astrophysics            int            `gorm:"default:0" json:"astrophysics"`
	IntergalacticResearch  int            `gorm:"default:0" json:"intergalactic_research"`
	GravitonTechnology     int            `gorm:"default:0" json:"graviton_technology"`
	WeaponsTechnology      int            `gorm:"default:0" json:"weapons_technology"`
	ShieldingTechnology    int            `gorm:"default:0" json:"shielding_technology"`
	ArmorTechnology        int            `gorm:"default:0" json:"armor_technology"`
	AssemblyTechnology     int            `gorm:"default:0" json:"assembly_technology"`
}

func (UserTech) TableName() string {
	return "users_tech"
}

type FleetMission struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"index" json:"user_id"`
	MissionType       int            `gorm:"index" json:"mission_type"`
	OriginGalaxy     int            `json:"origin_galaxy"`
	OriginSystem     int            `json:"origin_system"`
	OriginPosition   int            `json:"origin_position"`
	OriginPlanetType int            `json:"origin_planet_type"`
	TargetGalaxy     int            `json:"target_galaxy"`
	TargetSystem     int            `json:"target_system"`
	TargetPosition   int            `json:"target_position"`
	TargetPlanetType int            `json:"target_planet_type"`
	LaunchTime       time.Time      `json:"launch_time"`
	ArrivalTime      time.Time      `json:"arrival_time"`
	ReturnTime       time.Time      `json:"return_time"`
	Metal            int64          `json:"metal"`
	Crystal          int64          `json:"crystal"`
	Deuterium        int64          `json:"deuterium"`
	Status           int            `gorm:"default:0" json:"status"`
	ACSId            *uint          `json:"acs_id"`
	Ships            string         `gorm:"type:text" json:"ships"`
}

func (FleetMission) TableName() string {
	return "fleet_missions"
}

type ACS struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TargetGalaxy   int         `json:"target_galaxy"`
	TargetSystem   int         `json:"target_system"`
	TargetPosition int        `json:"target_position"`
	ArrivalTime    time.Time   `json:"arrival_time"`
	FleetIDs       []uint      `gorm:"-" json:"fleet_ids"`
	CreatedAt      time.Time   `json:"created_at"`
}

func (ACS) TableName() string {
	return "acs_fleets"
}

type BuildingQueue struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	PlanetID    uint           `gorm:"index" json:"planet_id"`
	BuildingID  int            `json:"building_id"`
	Level      int            `json:"level"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
	IsCancelled bool          `gorm:"default:false" json:"is_cancelled"`
}

func (BuildingQueue) TableName() string {
	return "building_queue"
}

type ResearchQueue struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	ResearchID int            `json:"research_id"`
	Level      int            `json:"level"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
}

func (ResearchQueue) TableName() string {
	return "research_queue"
}

type UnitQueue struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PlanetID   uint           `gorm:"index" json:"planet_id"`
	UnitID    int            `json:"unit_id"`
	Amount    int            `json:"amount"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
}

func (UnitQueue) TableName() string {
	return "unit_queue"
}

type Message struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	Type      int            `gorm:"index" json:"type"`
	FromUserID *uint         `json:"from_user_id,omitempty"`
	FromName  string         `json:"from_name"`
	Subject   string         `json:"subject"`
	Body      string         `gorm:"type:text" json:"body"`
	Read      bool           `gorm:"default:false" json:"read"`
	CreatedAt time.Time      `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

type Note struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	Galaxy     int            `json:"galaxy"`
	System     int            `json:"system"`
	Position   int            `json:"position"`
	Type       int            `json:"type"`
	Subject    string         `gorm:"size:100" json:"subject"`
	Text       string         `gorm:"type:text" json:"text"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (Note) TableName() string {
	return "notes"
}

type Alliance struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Name        string         `gorm:"size:50;uniqueIndex" json:"name"`
	Tag         string         `gorm:"size:10;uniqueIndex" json:"tag"`
	FounderID   uint           `json:"founder_id"`
	Description string         `gorm:"type:text" json:"description"`
	Logo        string         `gorm:"size:50" json:"logo"`
	Website     string         `gorm:"size:100" json:"website"`
	ApplicationText string     `gorm:"type:text" json:"application_text"`
}

func (Alliance) TableName() string {
	return "alliances"
}

type AllianceMember struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AllianceID uint           `gorm:"index" json:"alliance_id"`
	UserID     uint           `gorm:"uniqueIndex" json:"user_id"`
	Rank       string         `gorm:"size:30;default:'Member'" json:"rank"`
	JoinedAt   time.Time      `json:"joined_at"`
}

func (AllianceMember) TableName() string {
	return "alliance_members"
}

type AllianceApplication struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	AllianceID uint           `gorm:"index" json:"alliance_id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	Message    string         `gorm:"type:text" json:"message"`
	Status     string         `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (AllianceApplication) TableName() string {
	return "alliance_applications"
}

type Buddy struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SenderID   uint           `gorm:"index" json:"sender_id"`
	ReceiverID uint           `gorm:"index" json:"receiver_id"`
	Message    string         `gorm:"type:text" json:"message"`
	Status     string         `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (Buddy) TableName() string {
	return "buddy_requests"
}

type EspionageReport struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index" json:"user_id"`
	TargetUserID  uint           `gorm:"index" json:"target_user_id"`
	TargetPlanetID uint          `json:"target_planet_id"`
	Galaxy        int            `json:"galaxy"`
	System        int            `json:"system"`
	Position      int            `json:"position"`
	ReportType    string         `gorm:"size:20" json:"report_type"`
	Metal         int64          `json:"metal"`
	Crystal       int64          `json:"crystal"`
	Deuterium    int64          `json:"deuterium"`
	Energy        int64          `json:"energy"`
	Ships         string         `gorm:"type:text" json:"ships"`
	Defense       string         `gorm:"type:text" json:"defense"`
	Buildings     string         `gorm:"type:text" json:"buildings"`
	Read          bool           `gorm:"default:false" json:"read"`
	CreatedAt     time.Time      `json:"created_at"`
}

func (EspionageReport) TableName() string {
	return "espionage_reports"
}

type DebrisField struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Galaxy    int            `gorm:"index" json:"galaxy"`
	System    int            `gorm:"index" json:"system"`
	Position  int            `gorm:"index" json:"position"`
	Metal     int64          `json:"metal"`
	Crystal   int64          `json:"crystal"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
}

func (DebrisField) TableName() string {
	return "debris_fields"
}

type WreckField struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Galaxy    int            `gorm:"index" json:"galaxy"`
	System    int            `gorm:"index" json:"system"`
	Position  int            `gorm:"index" json:"position"`
	Metal     int64          `json:"metal"`
	Crystal   int64          `json:"crystal"`
	Deuterium int64          `json:"deuterium"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
}

func (WreckField) TableName() string {
	return "wreck_fields"
}
