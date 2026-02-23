package dto

type RegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	PlayerName string `json:"player_name"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success   bool   `json:"success"`
	UserID    uint   `json:"user_id"`
	AuthToken string `json:"auth_token"`
	Username  string `json:"username"`
}

type UserResponse struct {
	ID              uint    `json:"id"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	PlayerName      string  `json:"player_name"`
	DarkMatter      int64   `json:"dark_matter"`
	CharacterClass  int     `json:"character_class"`
	CurrentPlanet   *uint   `json:"current_planet_id,omitempty"`
}

type UserStatsResponse struct {
	UserID          uint   `json:"user_id"`
	Planets         int    `json:"planets"`
	TotalResources  ResourcesResponse `json:"total_resources"`
}

type ResourcesResponse struct {
	Metal      int64 `json:"metal"`
	Crystal    int64 `json:"crystal"`
	Deuterium  int64 `json:"deuterium"`
}

type PlanetResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Galaxy    int    `json:"galaxy"`
	System    int    `json:"system"`
	Position  int    `json:"position"`
	IsMoon    bool   `json:"is_moon"`
}

type PlanetOverviewResponse struct {
	PlanetID    uint              `json:"planet_id"`
	Name        string            `json:"name"`
	Resources   ResourcesResponse `json:"resources"`
	Production  ProductionResponse `json:"production"`
	Buildings   BuildingsResponse  `json:"buildings"`
	HasQueue    bool              `json:"has_queue"`
	FieldsUsed  int               `json:"fields_used"`
	FieldsMax   int               `json:"fields_max"`
}

type ProductionResponse struct {
	Metal     int64 `json:"metal"`
	Crystal   int64 `json:"crystal"`
	Deuterium int64 `json:"deuterium"`
}

type BuildingsResponse struct {
	MetalMine            int `json:"metal_mine"`
	CrystalMine          int `json:"crystal_mine"`
	DeuteriumSynthesizer int `json:"deuterium_synthesizer"`
	SolarPlant           int `json:"solar_plant"`
	FusionPlant          int `json:"fusion_plant"`
	RobotFactory         int `json:"robot_factory"`
	Shipyard             int `json:"shipyard"`
	ResearchLab          int `json:"research_lab"`
	NaniteFactory        int `json:"nanite_factory"`
	Terraformer          int `json:"terraformer"`
	SpaceDock            int `json:"space_dock"`
	MetalStorage         int `json:"metal_storage"`
	CrystalStorage       int `json:"crystal_storage"`
	DeuteriumStorage     int `json:"deuterium_storage"`
}

type StartBuildingRequest struct {
	Level int `json:"level"`
}

type StartBuildingResponse struct {
	Success     bool   `json:"success"`
	PlanetID    uint   `json:"planet_id"`
	BuildingID  int    `json:"building_id"`
	Level       int    `json:"level"`
}

type BuildingQueueItem struct {
	ID        uint   `json:"id"`
	Building  int    `json:"building_id"`
	Level     int    `json:"level"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type BuildingQueueResponse struct {
	PlanetID uint                `json:"planet_id"`
	Queue    []BuildingQueueItem `json:"queue"`
}

type CancelQueueResponse struct {
	Success  bool `json:"success"`
	QueueID  uint `json:"queue_id"`
}

type StartResearchRequest struct {
	ResearchID int `json:"research_id"`
}

type StartResearchResponse struct {
	Success     bool `json:"success"`
	UserID      uint `json:"user_id"`
	ResearchID  int  `json:"research_id"`
}

type ResearchQueueItem struct {
	ID         uint   `json:"id"`
	ResearchID int    `json:"research_id"`
	Level      int    `json:"level"`
	StartTime  int64  `json:"start_time"`
	EndTime    int64  `json:"end_time"`
}

type ResearchQueueResponse struct {
	UserID  uint                  `json:"user_id"`
	Queue   []ResearchQueueItem   `json:"queue"`
}

type FleetResources struct {
	Metal     int64 `json:"metal"`
	Crystal   int64 `json:"crystal"`
	Deuterium int64 `json:"deuterium"`
}

type SendFleetRequest struct {
	MissionType    int             `json:"mission_type"`
	OriginGalaxy   int             `json:"origin_galaxy"`
	OriginSystem   int             `json:"origin_system"`
	OriginPosition int             `json:"origin_position"`
	TargetGalaxy   int             `json:"target_galaxy"`
	TargetSystem   int             `json:"target_system"`
	TargetPosition int             `json:"target_position"`
	Resources      FleetResources  `json:"resources"`
	Ships          map[string]int  `json:"ships"`
}

type SendFleetResponse struct {
	Success bool `json:"success"`
}

type FleetInfo struct {
	ID          uint   `json:"id"`
	MissionType int    `json:"mission_type"`
	Origin      string `json:"origin"`
	Target      string `json:"target"`
	LaunchTime  int64  `json:"launch_time"`
	ArrivalTime int64  `json:"arrival_time"`
	Status      int    `json:"status"`
}

type FleetsResponse struct {
	UserID  uint       `json:"user_id"`
	Fleets  []FleetInfo `json:"fleets"`
}

type BuildUnitRequest struct {
	UnitID  int `json:"unit_id"`
	Amount  int `json:"amount"`
}

type BuildUnitResponse struct {
	Success   bool `json:"success"`
	PlanetID  uint `json:"planet_id"`
	UnitID    int  `json:"unit_id"`
	Amount    int  `json:"amount"`
}

type UnitQueueItem struct {
	ID        uint   `json:"id"`
	UnitID    int    `json:"unit_id"`
	Amount    int    `json:"amount"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}

type UnitQueueResponse struct {
	PlanetID uint             `json:"planet_id"`
	Queue    []UnitQueueItem  `json:"queue"`
}

type UnitInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Metal     int64  `json:"metal"`
	Crystal   int64  `json:"crystal"`
	Deuterium int64  `json:"deuterium"`
	BuildTime int64  `json:"build_time_seconds"`
	Category  string `json:"category"`
}

type AvailableUnitsResponse struct {
	Units []UnitInfo `json:"units"`
}

type PlanetShipsResponse struct {
	PlanetID uint              `json:"planet_id"`
	Ships    map[string]int   `json:"ships"`
}

type PlanetDefenseResponse struct {
	PlanetID uint              `json:"planet_id"`
	Defense  map[string]int   `json:"defense"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	Code       string `json:"code,omitempty"`
	RetryAfter int    `json:"retry_after,omitempty"`
}

type StatusResponse struct {
	Online    bool   `json:"online"`
	Version   string `json:"version"`
	Universe  string `json:"universe"`
	APIBase   string `json:"api_base"`
	Time      int64  `json:"time"`
}

type BattleSimulateRequest struct {
	AttackerFleets []FleetComposition `json:"attacker_fleets"`
	DefenderFleets []FleetComposition `json:"defender_fleets"`
	AttackerTech   *TechLevel          `json:"attacker_tech,omitempty"`
	DefenderTech   *TechLevel          `json:"defender_tech,omitempty"`

	SimulationCount int  `json:"simulation_count,omitempty"` // 1 for deterministic, N for Monte Carlo
	RapidFire       bool `json:"rapid_fire,omitempty"`       // Default true
	DefToDebris    int  `json:"def_to_debris_percent,omitempty"` // 0-100, default 70

	OriginCoords *Coordinates `json:"origin_coords,omitempty"`
	TargetCoords *Coordinates `json:"target_coords,omitempty"`
	TargetResources *TargetResources `json:"target_resources,omitempty"`

	IPM *IPMRequest `json:"ipm,omitempty"`
}

type FleetComposition struct {
	Slot     int                    `json:"slot,omitempty"` // ACS slot number
	Ships    map[string]int         `json:"ships"`
	Defense  map[string]int        `json:"defense,omitempty"`
	Resources *FleetResources      `json:"resources,omitempty"`
}

type Coordinates struct {
	Galaxy   int `json:"galaxy"`
	System   int `json:"system"`
	Position int `json:"position"`
}

type TargetResources struct {
	Metal     int64 `json:"metal"`
	Crystal   int64 `json:"crystal"`
	Deuterium int64 `json:"deuterium"`
}

type IPMRequest struct {
	Count       int `json:"count"`
	TargetID    int `json:"target_id"` // 401=rocket_launcher, etc.
	Technologies *TechLevel `json:"technologies,omitempty"`
}

type TechLevel struct {
	Weapons         int `json:"weapons"`
	Shielding       int `json:"shielding"`
	Armor           int `json:"armor"`
	CombustionDrive int `json:"combustion_drive,omitempty"`
	ImpulseDrive    int `json:"impulse_drive,omitempty"`
	HyperspaceDrive int `json:"hyperspace_drive,omitempty"`
}

type BattleSimulateResponse struct {
	Simulations int `json:"simulations"`

	Winner          string                    `json:"winner"`
	WinChance       float64                   `json:"win_chance"` // percentage for Monte Carlo
	AttackerWins    int                       `json:"attacker_wins"`
	DefenderWins    int                       `json:"defender_wins"`
	Draws           int                       `json:"draws"`

	AttackerLeft   map[string]int            `json:"attacker_left"`
	DefenderLeft   map[string]int            `json:"defender_left"`
	AttackerLost   map[string]int            `json:"attacker_lost"`
	DefenderLost   map[string]int            `json:"defender_lost"`

	Debris         DebrisInfo                `json:"debris"`
	Ruins          RuinsInfo                 `json:"ruins"`
	MoonChance     float64                   `json:"moon_chance"`

	Plunder    PlunderInfo  `json:"plunder"`
	FuelUsed   int64        `json:"fuel_used"`
	FlightTime *FlightResult `json:"flight_time,omitempty"`

	IPMResult *IPMResult `json:"ipm_result,omitempty"`
}

type PlunderInfo struct {
	Theoretical  ResourcesResponse `json:"theoretical"`
	Actual       ResourcesResponse `json:"actual"`
	CargoNeeded int64             `json:"cargo_needed"`
}

type FlightResult struct {
	Seconds     int    `json:"seconds"`
	Hours       int    `json:"hours"`
	Minutes     int    `json:"minutes"`
	Formatted   string `json:"formatted"`
	ReturnTime  *FlightResult `json:"return_time,omitempty"`
}

type IPMResult struct {
	Destroyed     map[string]int `json:"destroyed"`
	RemainingDef map[string]int `json:"remaining_def"`
	IPMLosses    int            `json:"ipm_losses"`
}

type DebrisInfo struct {
	Metal     int64 `json:"metal"`
	Crystal   int64 `json:"crystal"`
	Total     int64 `json:"total"`
}

type RuinsInfo struct {
	Metal     int64 `json:"metal"`
	Crystal   int64 `json:"crystal"`
	Total     int64 `json:"total"`
}
