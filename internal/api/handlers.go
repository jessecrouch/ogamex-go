package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
	"ogamex-go/internal/service"
)

type Handlers struct {
	buildingService   *service.BuildingService
	researchService   *service.ResearchService
	fleetService      *service.FleetService
	planetRepo        repository.PlanetRepository
	authService       *service.AuthService
	unitService       *service.UnitService
	productionService *service.ProductionService
	messageService    *service.MessageService
	planetService     *service.PlanetService
	noteService       *service.NoteService
	allianceService   *service.AllianceService
	buddyService     *service.BuddyService
	moonService      *service.MoonService
	espionageService *service.EspionageService
	debrisService    *service.DebrisService
	npcService       *service.NPCService
	acsService       *service.ACSService
	premiumService   *service.PremiumService
	characterClassService *service.CharacterClassService
}

func NewHandlers(
	buildingService *service.BuildingService,
	researchService *service.ResearchService,
	fleetService *service.FleetService,
	planetRepo repository.PlanetRepository,
	authService *service.AuthService,
	unitService *service.UnitService,
	productionService *service.ProductionService,
	messageService *service.MessageService,
	planetService *service.PlanetService,
	noteService *service.NoteService,
	allianceService *service.AllianceService,
	buddyService *service.BuddyService,
	moonService *service.MoonService,
	espionageService *service.EspionageService,
	debrisService *service.DebrisService,
	npcService *service.NPCService,
	acsService *service.ACSService,
	premiumService *service.PremiumService,
	characterClassService *service.CharacterClassService,
) *Handlers {
	return &Handlers{
		buildingService:   buildingService,
		researchService:   researchService,
		fleetService:      fleetService,
		planetRepo:        planetRepo,
		authService:       authService,
		unitService:       unitService,
		productionService: productionService,
		messageService:    messageService,
		planetService:     planetService,
		noteService:       noteService,
		allianceService:   allianceService,
		buddyService:     buddyService,
		moonService:      moonService,
		espionageService: espionageService,
		debrisService:    debrisService,
		npcService:       npcService,
		acsService:       acsService,
		premiumService:   premiumService,
		characterClassService: characterClassService,
	}
}

func (h *Handlers) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/status", h.GetStatus)
	
	api.Post("/auth/register", h.Register)
	api.Post("/auth/login", h.Login)

	protected := api.Group("", h.authMiddleware)

	protected.Get("/user", h.GetUser)
	protected.Get("/user/stats", h.GetUserStats)
	protected.Get("/planets", h.GetUserPlanets)
	protected.Get("/planets/:id", h.GetPlanet)
	protected.Get("/planets/:id/details", h.GetPlanetDetails)
	protected.Put("/planets/:id/set-current", h.SetCurrentPlanet)
	protected.Get("/planets/:id/resources", h.GetPlanetResources)
	protected.Get("/planets/:id/buildings", h.GetPlanetBuildings)
	protected.Post("/planets/:id/buildings/:building_id", h.StartBuilding)
	protected.Get("/planets/:id/queue", h.GetBuildingQueue)
	protected.Delete("/queue/:id", h.CancelBuilding)
	protected.Get("/planets/:id/overview", h.GetPlanetOverview)

	protected.Post("/research/start", h.StartResearch)
	protected.Get("/research/queue", h.GetResearchQueue)
	protected.Delete("/research/queue/:id", h.CancelResearch)

	protected.Post("/fleets/send", h.SendFleet)
	protected.Get("/fleets", h.GetFleets)
	protected.Post("/fleets/:id/recall", h.RecallFleet)

	protected.Post("/planets/:id/units/build", h.BuildUnit)
	protected.Get("/planets/:id/units/queue", h.GetUnitQueue)
	protected.Delete("/units/queue/:id", h.CancelUnit)
	protected.Get("/units/available", h.GetAvailableUnits)
	protected.Get("/planets/:id/units", h.GetPlanetShips)
	protected.Get("/planets/:id/defense", h.GetPlanetDefense)
	protected.Get("/planets/:id/production", h.GetPlanetProduction)

	protected.Get("/messages", h.GetMessages)
	protected.Get("/messages/unread", h.GetUnreadCount)
	protected.Post("/messages/:id/read", h.MarkMessageRead)
	protected.Delete("/messages/:id", h.DeleteMessage)

	protected.Get("/galaxy/:galaxy/:system", h.GetGalaxy)

	protected.Get("/highscore/:category", h.GetHighscore)

	protected.Get("/notes", h.GetNotes)
	protected.Post("/notes", h.CreateNote)
	protected.Put("/notes/:id", h.UpdateNote)
	protected.Delete("/notes/:id", h.DeleteNote)

	protected.Get("/alliances", h.GetAlliances)
	protected.Post("/alliances", h.CreateAlliance)
	protected.Get("/alliances/:id", h.GetAlliance)
	protected.Get("/alliances/:id/members", h.GetAllianceMembers)
	protected.Post("/alliances/:id/apply", h.ApplyToAlliance)
	protected.Get("/alliances/:id/applications", h.GetAllianceApplications)
	protected.Post("/alliances/applications/:id/accept", h.AcceptApplication)
	protected.Post("/alliances/applications/:id/reject", h.RejectApplication)
	protected.Post("/alliances/leave", h.LeaveAlliance)
	protected.Put("/alliances", h.UpdateAlliance)

	protected.Get("/buddy", h.GetBuddies)
	protected.Post("/buddy/request", h.SendBuddyRequest)
	protected.Get("/buddy/pending", h.GetPendingBuddyRequests)
	protected.Post("/buddy/:id/accept", h.AcceptBuddyRequest)
	protected.Post("/buddy/:id/reject", h.RejectBuddyRequest)
	protected.Delete("/buddy/:id", h.RemoveBuddy)

	protected.Get("/planets/:id/jump-gate/targets", h.GetJumpGateTargets)
	protected.Post("/planets/:id/jump-gate/execute", h.ExecuteJumpGate)
	protected.Post("/planets/:id/phalanx/scan", h.ScanWithPhalanx)

	protected.Get("/espionage", h.GetEspionageReports)
	protected.Get("/espionage/unread", h.GetEspionageUnreadCount)
	protected.Post("/espionage/:id/read", h.MarkEspionageReportRead)
	protected.Delete("/espionage/:id", h.DeleteEspionageReport)

	protected.Get("/user/vacation", h.GetVacationStatus)
	protected.Post("/user/vacation/enable", h.EnableVacationMode)
	protected.Post("/user/vacation/disable", h.DisableVacationMode)

	protected.Get("/debris", h.GetDebrisFields)
	protected.Get("/debris/:galaxy/:system/:position", h.GetDebrisField)
	protected.Post("/debris/:galaxy/:system/:position/collect", h.CollectDebris)

	protected.Get("/wrecks", h.GetWreckFields)
	protected.Get("/wrecks/:galaxy/:system/:position", h.GetWreckField)
	protected.Post("/wrecks/:galaxy/:system/:position/collect", h.CollectWreckField)

	protected.Get("/npc/planets", h.GetNPCPlanets)
	protected.Get("/npc/planets/:galaxy/:system/:position", h.GetNPCPlanet)
	protected.Post("/npc/planets", h.CreateNPCPlanet)
	protected.Post("/npc/fleets/expedition", h.GenerateExpeditionFleet)
	protected.Post("/npc/fleets/pirate", h.GeneratePirateRaid)

	protected.Post("/acs/create", h.CreateACS)
	protected.Post("/acs/:id/join", h.JoinACS)
	protected.Get("/acs/:id", h.GetACS)
	protected.Get("/acs/:id/fleets", h.GetACSFleets)

	protected.Get("/premium/status", h.GetPremiumStatus)
	protected.Post("/premium/activate", h.ActivatePremium)
	protected.Post("/merchant/buy", h.MerchantBuy)
	protected.Post("/merchant/sell", h.MerchantSell)

	protected.Get("/character-class", h.GetCharacterClass)
	protected.Post("/character-class/select", h.SelectCharacterClass)

	protected.Post("/planets/:id/move", h.MovePlanet)
}

func (h *Handlers) authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "authorization required"})
	}

	parts := []string{}
	for _, p := range strings.Split(authHeader, " ") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
	}

	token := parts[1]
	user, err := h.authService.ValidateToken(c.Context(), token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}

	c.Locals("user_id", user.ID)
	c.Locals("user", user)

	return c.Next()
}

func (h *Handlers) GetStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"online":    true,
		"version":   "0.0.1",
		"universe":  "ogamex-go",
		"api_base":  "/api/v1",
		"time":      time.Now().Unix(),
	})
}

func (h *Handlers) GetUser(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	
	u := user.(*schema.User)
	
	return c.JSON(fiber.Map{
		"id":              u.ID,
		"username":       u.Username,
		"email":          u.Email,
		"player_name":    u.PlayerName,
		"dark_matter":    u.DarkMatter,
		"character_class": u.CharacterClass,
		"current_planet": u.CurrentPlanetID,
	})
}

func (h *Handlers) GetUserStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	planets, err := h.planetRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	totalMetal := int64(0)
	totalCrystal := int64(0)
	totalDeuterium := int64(0)
	planetCount := len(planets)

	for _, p := range planets {
		totalMetal += p.Metal
		totalCrystal += p.Crystal
		totalDeuterium += p.Deuterium
	}

	return c.JSON(fiber.Map{
		"user_id":           userID,
		"planets":           planetCount,
		"total_resources": fiber.Map{
			"metal":      totalMetal,
			"crystal":    totalCrystal,
			"deuterium":  totalDeuterium,
		},
	})
}

func (h *Handlers) GetUserPlanets(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	
	planets, err := h.planetRepo.GetByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	
	type PlanetInfo struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Galaxy   int    `json:"galaxy"`
		System   int    `json:"system"`
		Position int    `json:"position"`
		IsMoon   bool   `json:"is_moon"`
	}
	
	items := make([]PlanetInfo, len(planets))
	for i, p := range planets {
		items[i] = PlanetInfo{
			ID:       p.ID,
			Name:     p.Name,
			Galaxy:   p.Galaxy,
			System:   p.System,
			Position: p.Position,
			IsMoon:   p.IsMoon,
		}
	}
	
	return c.JSON(fiber.Map{
		"planets": items,
	})
}

func (h *Handlers) GetPlanet(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	return c.JSON(fiber.Map{
		"id":       planet.ID,
		"name":     planet.Name,
		"galaxy":   planet.Galaxy,
		"system":   planet.System,
		"position": planet.Position,
		"is_moon":  planet.IsMoon,
	})
}

func (h *Handlers) GetPlanetDetails(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	return c.JSON(fiber.Map{
		"id":              planet.ID,
		"name":            planet.Name,
		"galaxy":          planet.Galaxy,
		"system":          planet.System,
		"position":        planet.Position,
		"is_moon":         planet.IsMoon,
		"planet_type":     planet.PlanetType,
		"metal":           planet.Metal,
		"crystal":         planet.Crystal,
		"deuterium":       planet.Deuterium,
		"metal_capacity":  planet.MetalCapacity,
		"crystal_capacity": planet.CrystalCapacity,
		"deuterium_capacity": planet.DeuteriumCapacity,
		"energy_available": planet.EnergyAvailable,
		"energy_max":      planet.EnergyMax,
		"energy_used":     planet.EnergyUsed,
		"fields_used":     planet.FieldsUsed,
		"fields_max":      planet.FieldsMax,
		"temp_min":        planet.TempMin,
		"temp_max":        planet.TempMax,
		"defense_activated": planet.DefenseActivated,
	})
}

func (h *Handlers) SetCurrentPlanet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	if planet.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "planet does not belong to user"})
	}

	err = h.authService.SetCurrentPlanet(c.Context(), userID, uint(planetID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "current_planet_id": planetID})
}

func (h *Handlers) GetPlanetOverview(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	queue, _ := h.buildingService.GetQueue(c.Context(), uint(planetID))
	
	hasQueue := len(queue) > 0

	return c.JSON(fiber.Map{
		"planet_id":     planet.ID,
		"name":          planet.Name,
		"resources": fiber.Map{
			"metal":      planet.Metal,
			"crystal":    planet.Crystal,
			"deuterium":  planet.Deuterium,
			"energy":     planet.EnergyAvailable,
		},
		"production": fiber.Map{
			"metal":      planet.MetalProduction,
			"crystal":    planet.CrystalProduction,
			"deuterium":  planet.DeuteriumProduction,
		},
		"buildings": fiber.Map{
			"metal_mine":     planet.MetalMine,
			"crystal_mine":  planet.CrystalMine,
			"deuterium_synth": planet.DeuteriumSynthesizer,
			"solar_plant":    planet.SolarPlant,
			"fusion_plant":   planet.FusionPlant,
		},
		"has_queue":     hasQueue,
		"fields_used":   planet.FieldsUsed,
		"fields_max":    planet.FieldsMax,
	})
}

func (h *Handlers) GetPlanetResources(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	return c.JSON(fiber.Map{
		"planet_id":  planet.ID,
		"metal":      planet.Metal,
		"crystal":    planet.Crystal,
		"deuterium":  planet.Deuterium,
		"energy":     planet.EnergyAvailable,
		"metal_capacity":      planet.MetalCapacity,
		"crystal_capacity":    planet.CrystalCapacity,
		"deuterium_capacity":  planet.DeuteriumCapacity,
	})
}

func (h *Handlers) GetPlanetBuildings(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	return c.JSON(fiber.Map{
		"planet_id": planet.ID,
		"buildings": fiber.Map{
			"metal_mine":            planet.MetalMine,
			"crystal_mine":         planet.CrystalMine,
			"deuterium_synthesizer": planet.DeuteriumSynthesizer,
			"solar_plant":           planet.SolarPlant,
			"fusion_plant":         planet.FusionPlant,
			"metal_storage":        planet.MetalStorageBuilding,
			"crystal_storage":      planet.CrystalStorageBuilding,
			"deuterium_storage":     planet.DeuteriumStorageBuilding,
			"robot_factory":        planet.RobotFactory,
			"shipyard":              planet.Shipyard,
			"research_lab":         planet.ResearchLab,
			"nanite_factory":        planet.NaniteFactory,
			"terraformer":           planet.Terraformer,
			"space_dock":            planet.SpaceDock,
		},
		"fields_used": planet.FieldsUsed,
		"fields_max":  planet.FieldsMax,
	})
}

func (h *Handlers) StartBuilding(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	buildingID, err := c.ParamsInt("building_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid building id"})
	}

	type BuildingRequest struct {
		Level int `json:"level"`
	}

	var req BuildingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.Level < 1 {
		req.Level = 1
	}

	err = h.buildingService.StartBuilding(c.Context(), uint(planetID), buildingID, req.Level)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"planet_id":   planetID,
		"building_id": buildingID,
		"level":       req.Level,
	})
}

func (h *Handlers) GetBuildingQueue(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	queue, err := h.buildingService.GetQueue(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	type QueueItem struct {
		ID        uint   `json:"id"`
		Building  int    `json:"building_id"`
		Level     int    `json:"level"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
	}

	items := make([]QueueItem, len(queue))
	for i, q := range queue {
		items[i] = QueueItem{
			ID:        q.ID,
			Building:  q.BuildingID,
			Level:     q.Level,
			StartTime: q.StartTime.Unix(),
			EndTime:   q.EndTime.Unix(),
		}
	}

	return c.JSON(fiber.Map{
		"planet_id": planetID,
		"queue":     items,
	})
}

func (h *Handlers) CancelBuilding(c *fiber.Ctx) error {
	queueID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid queue id"})
	}

	err = h.buildingService.CancelBuilding(c.Context(), uint(queueID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"queue_id": queueID,
	})
}

func (h *Handlers) StartResearch(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	type ResearchRequest struct {
		ResearchID int `json:"research_id"`
	}

	var req ResearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	err := h.researchService.StartResearch(c.Context(), userID, req.ResearchID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"user_id":      userID,
		"research_id":  req.ResearchID,
	})
}

func (h *Handlers) GetResearchQueue(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	queue, err := h.researchService.GetQueue(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	type QueueItem struct {
		ID         uint   `json:"id"`
		ResearchID int    `json:"research_id"`
		Level      int    `json:"level"`
		StartTime  int64  `json:"start_time"`
		EndTime    int64  `json:"end_time"`
	}

	items := make([]QueueItem, len(queue))
	for i, q := range queue {
		items[i] = QueueItem{
			ID:         q.ID,
			ResearchID: q.ResearchID,
			Level:      q.Level,
			StartTime:  q.StartTime.Unix(),
			EndTime:    q.EndTime.Unix(),
		}
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
		"queue":   items,
	})
}

type FleetRequest struct {
	MissionType    int    `json:"mission_type"`
	OriginGalaxy   int    `json:"origin_galaxy"`
	OriginSystem   int    `json:"origin_system"`
	OriginPosition int    `json:"origin_position"`
	TargetGalaxy   int    `json:"target_galaxy"`
	TargetSystem   int    `json:"target_system"`
	TargetPosition int    `json:"target_position"`
	Resources      struct {
		Metal     int64 `json:"metal"`
		Crystal   int64 `json:"crystal"`
		Deuterium int64 `json:"deuterium"`
	} `json:"resources"`
	Ships map[string]int `json:"ships"`
}

func (h *Handlers) SendFleet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req FleetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	origin, err := h.planetRepo.GetByCoords(c.Context(), userID, req.OriginGalaxy, req.OriginSystem, req.OriginPosition)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "origin planet not found"})
	}

	target, err := h.planetRepo.GetByCoords(c.Context(), 0, req.TargetGalaxy, req.TargetSystem, req.TargetPosition)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "target planet not found"})
	}

	ships := make(map[int]int)
	for shipType, count := range req.Ships {
		shipID, err := strconv.Atoi(shipType)
		if err == nil && count > 0 {
			ships[shipID] = count
		}
	}

	params := service.FleetParams{
		Origin:      *origin,
		Target:      *target,
		MissionType: req.MissionType,
	}
	params.Resources.Metal = req.Resources.Metal
	params.Resources.Crystal = req.Resources.Crystal
	params.Resources.Deuterium = req.Resources.Deuterium
	params.Ships = ships

	err = h.fleetService.SendFleet(c.Context(), userID, params)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (h *Handlers) GetFleets(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	fleets, err := h.fleetService.GetActiveMissions(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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

	items := make([]FleetInfo, len(fleets))
	for i, f := range fleets {
		items[i] = FleetInfo{
			ID:          f.ID,
			MissionType: f.MissionType,
			Origin:      fmt.Sprintf("%d:%d:%d", f.OriginGalaxy, f.OriginSystem, f.OriginPosition),
			Target:      fmt.Sprintf("%d:%d:%d", f.TargetGalaxy, f.TargetSystem, f.TargetPosition),
			LaunchTime:  f.LaunchTime.Unix(),
			ArrivalTime: f.ArrivalTime.Unix(),
			Status:      f.Status,
		}
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
		"fleets":  items,
	})
}

func (h *Handlers) RecallFleet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	fleetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid fleet id"})
	}

	err = h.fleetService.RecallFleet(c.Context(), uint(fleetID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"fleet_id": fleetID,
	})
}

func (h *Handlers) CancelResearch(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	queueID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid queue id"})
	}

	err = h.researchService.CancelResearch(c.Context(), uint(queueID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"queue_id": queueID,
	})
}

type RegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	PlayerName string `json:"player_name"`
}

func (h *Handlers) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username and password required"})
	}

	if req.PlayerName == "" {
		req.PlayerName = req.Username
	}

	user, err := h.authService.Register(c.Context(), service.RegisterInput{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		PlayerName: req.PlayerName,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"user_id":   user.ID,
		"auth_token": user.AuthToken,
		"username":  user.Username,
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handlers) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.authService.Login(c.Context(), service.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"user_id":   user.ID,
		"auth_token": user.AuthToken,
		"username":  user.Username,
	})
}

type BuildUnitRequest struct {
	UnitID int `json:"unit_id"`
	Amount int `json:"amount"`
}

func (h *Handlers) BuildUnit(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	if planet.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "planet does not belong to user"})
	}

	var req BuildUnitRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.UnitID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "unit_id is required"})
	}

	if req.Amount <= 0 {
		req.Amount = 1
	}

	err = h.unitService.BuildUnit(c.Context(), uint(planetID), req.UnitID, req.Amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"planet_id": planetID,
		"unit_id":  req.UnitID,
		"amount":   req.Amount,
	})
}

func (h *Handlers) GetUnitQueue(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	queue, err := h.unitService.GetQueue(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	type UnitQueueItem struct {
		ID        uint   `json:"id"`
		UnitID    int    `json:"unit_id"`
		Amount    int    `json:"amount"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
	}

	items := make([]UnitQueueItem, len(queue))
	for i, q := range queue {
		items[i] = UnitQueueItem{
			ID:        q.ID,
			UnitID:    q.UnitID,
			Amount:    q.Amount,
			StartTime: q.StartTime.Unix(),
			EndTime:   q.EndTime.Unix(),
		}
	}

	return c.JSON(fiber.Map{
		"planet_id": planetID,
		"queue":     items,
	})
}

func (h *Handlers) CancelUnit(c *fiber.Ctx) error {
	queueID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid queue id"})
	}

	err = h.unitService.CancelQueue(c.Context(), uint(queueID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":  true,
		"queue_id": queueID,
	})
}

func (h *Handlers) GetAvailableUnits(c *fiber.Ctx) error {
	type UnitInfo struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Metal       int64  `json:"metal"`
		Crystal     int64  `json:"crystal"`
		Deuterium   int64  `json:"deuterium"`
		BuildTime   int64  `json:"build_time_seconds"`
		Category    string `json:"category"`
	}

	units := []UnitInfo{
		{202, "Small Cargo", 2000, 2000, 0, 5, "ship"},
		{203, "Large Cargo", 6000, 6000, 0, 8, "ship"},
		{204, "Light Fighter", 10000, 6000, 2000, 20, "ship"},
		{205, "Heavy Fighter", 25000, 15000, 5000, 40, "ship"},
		{206, "Cruiser", 10000, 20000, 10000, 10, "ship"},
		{207, "Battleship", 50000, 25000, 15000, 80, "ship"},
		{208, "Colony Ship", 10000, 10000, 0, 50, "ship"},
		{209, "Recycler", 10000, 6000, 2000, 15, "ship"},
		{210, "Espionage Probe", 0, 1000, 0, 30, "ship"},
		{211, "Bomber", 50000, 50000, 25000, 200, "ship"},
		{212, "Solar Satellite", 0, 2000, 500, 3, "ship"},
		{213, "Destroyer", 10000, 10000, 0, 30, "ship"},
		{214, "Deathstar", 100000, 100000, 50000, 400, "ship"},
		{215, "Battlecruiser", 3000, 1000, 0, 4, "ship"},
		{217, "Crawler", 2000, 2000, 1000, 10, "ship"},
		{218, "Reaper", 8000, 0, 0, 20, "ship"},
		{219, "Pathfinder", 20000, 10000, 10000, 75, "ship"},
		{401, "Rocket Launcher", 2000, 0, 0, 10, "defense"},
		{402, "Light Laser", 1500, 500, 0, 11, "defense"},
		{403, "Heavy Laser", 6000, 2000, 0, 22, "defense"},
		{404, "Ion Cannon", 2000, 6000, 0, 16, "defense"},
		{405, "Gauss Cannon", 20000, 15000, 2000, 45, "defense"},
		{406, "Plasma Turret", 50000, 50000, 30000, 90, "defense"},
		{407, "Shield Dome", 10000, 10000, 0, 20, "defense"},
		{408, "Missile Interceptor", 8000, 2000, 0, 15, "defense"},
		{409, "Missile Launcher", 15000, 5000, 0, 20, "defense"},
	}

	return c.JSON(fiber.Map{
		"units": units,
	})
}

func (h *Handlers) GetPlanetShips(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	return c.JSON(fiber.Map{
		"planet_id": planet.ID,
		"ships": fiber.Map{
			"small_cargo":       planet.SmallCargo,
			"large_cargo":       planet.LargeCargo,
			"light_fighter":    planet.LightFighter,
			"heavy_fighter":    planet.HeavyFighter,
			"cruiser":          planet.Cruiser,
			"battleship":       planet.Battleship,
			"colony_ship":      planet.ColonyShip,
			"recycler":         planet.Recycler,
			"espionage_probe":  planet.EspionageProbe,
			"bomber":           planet.Bomber,
			"destroyer":        planet.Destroyer,
			"deathstar":        planet.Deathstar,
			"battlecruiser":    planet.Battlecruiser,
			"reaper":           planet.Reaper,
			"pathfinder":       planet.Pathfinder,
			"solar_satellite":  planet.SolarSatellite,
			"crawler":          planet.Crawler,
		},
	})
}

func (h *Handlers) GetPlanetDefense(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	return c.JSON(fiber.Map{
		"planet_id": planet.ID,
		"defense": fiber.Map{
			"rocket_launcher":      planet.RocketLauncher,
			"light_laser":         planet.LightLaser,
			"heavy_laser":         planet.HeavyLaser,
			"ion_cannon":          planet.IonCannon,
			"gauss_cannon":        planet.GaussCannon,
			"plasma_turret":       planet.PlasmaTurret,
			"shield_dome":         planet.ShieldDome,
			"missile_interceptor": planet.MissileInterceptor,
			"missile_launcher":    planet.MissileLauncher,
		},
	})
}

func (h *Handlers) GetPlanetProduction(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	planet, err := h.planetRepo.GetByID(c.Context(), uint(planetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "planet not found"})
	}

	tech, _ := h.researchService.GetTech(c.Context(), planet.UserID)

	metalProd := h.productionService.CalculateProductionForType(planet, tech, "metal")
	crystalProd := h.productionService.CalculateProductionForType(planet, tech, "crystal")
	deuteriumProd := h.productionService.CalculateProductionForType(planet, tech, "deuterium")
	energyProd := h.productionService.CalculateEnergy(planet, tech)

	return c.JSON(fiber.Map{
		"planet_id": planet.ID,
		"production": fiber.Map{
			"metal":      metalProd,
			"crystal":    crystalProd,
			"deuterium":  deuteriumProd,
			"energy":     energyProd,
		},
		"consumption": fiber.Map{
			"energy": planet.EnergyUsed,
		},
	})
}

func (h *Handlers) GetMessages(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	messages, err := h.messageService.GetMessages(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"messages": messages,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *Handlers) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	count, err := h.messageService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"unread_count": count,
	})
}

func (h *Handlers) MarkMessageRead(c *fiber.Ctx) error {
	messageID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid message id"})
	}

	err = h.messageService.MarkAsRead(c.Context(), uint(messageID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) DeleteMessage(c *fiber.Ctx) error {
	messageID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid message id"})
	}

	err = h.messageService.DeleteMessage(c.Context(), uint(messageID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetGalaxy(c *fiber.Ctx) error {
	galaxy, err := c.ParamsInt("galaxy")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}

	system, err := c.ParamsInt("system")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}

	if galaxy < 1 || galaxy > 10 {
		return c.Status(400).JSON(fiber.Map{"error": "galaxy must be between 1 and 10"})
	}

	if system < 1 || system > 499 {
		return c.Status(400).JSON(fiber.Map{"error": "system must be between 1 and 499"})
	}

	positions, err := h.planetService.GetGalaxy(c.Context(), galaxy, system)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"galaxy":   galaxy,
		"system":   system,
		"positions": positions,
	})
}

func (h *Handlers) GetHighscore(c *fiber.Ctx) error {
	category := c.Params("category")
	if category == "" {
		category = "points"
	}

	validCategories := map[string]bool{
		"points":   true,
		"military": true,
		"defense":  true,
		"research": true,
		"fleets":   true,
	}
	if !validCategories[category] {
		return c.Status(400).JSON(fiber.Map{"error": "invalid category"})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "100"))

	entries, err := h.planetService.GetHighscore(c.Context(), category, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"category": category,
		"entries":  entries,
	})
}

func (h *Handlers) GetNotes(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	notes, err := h.noteService.GetNotes(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"notes": notes})
}

type CreateNoteRequest struct {
	Galaxy   int    `json:"galaxy"`
	System   int    `json:"system"`
	Position int    `json:"position"`
	Type     int    `json:"type"`
	Subject  string `json:"subject"`
	Text     string `json:"text"`
}

func (h *Handlers) CreateNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req CreateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	if req.Subject == "" {
		return c.Status(400).JSON(fiber.Map{"error": "subject is required"})
	}

	note, err := h.noteService.CreateNote(c.Context(), userID, service.CreateNoteInput{
		Galaxy:   req.Galaxy,
		System:   req.System,
		Position: req.Position,
		Type:     req.Type,
		Subject:  req.Subject,
		Text:     req.Text,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"note":    note,
	})
}

type UpdateNoteRequest struct {
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

func (h *Handlers) UpdateNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	noteID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid note id"})
	}

	var req UpdateNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	note, err := h.noteService.UpdateNote(c.Context(), uint(noteID), req.Subject, req.Text)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"note":    note,
	})
}

func (h *Handlers) DeleteNote(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	noteID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid note id"})
	}

	err = h.noteService.DeleteNote(c.Context(), uint(noteID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetAlliances(c *fiber.Ctx) error {
	alliances, err := h.allianceService.GetAllAlliances(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"alliances": alliances})
}

type CreateAllianceRequest struct {
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Website     string `json:"website"`
}

func (h *Handlers) CreateAlliance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	var req CreateAllianceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	alliance, err := h.allianceService.CreateAlliance(c.Context(), userID, service.CreateAllianceInput{
		Name: req.Name, Tag: req.Tag, Description: req.Description, Logo: req.Logo, Website: req.Website,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "alliance": alliance})
}

func (h *Handlers) GetAlliance(c *fiber.Ctx) error {
	allianceID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid alliance id"})
	}
	alliance, err := h.allianceService.GetAlliance(c.Context(), uint(allianceID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "alliance not found"})
	}
	return c.JSON(fiber.Map{"alliance": alliance})
}

func (h *Handlers) GetAllianceMembers(c *fiber.Ctx) error {
	allianceID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid alliance id"})
	}
	members, err := h.allianceService.GetMembers(c.Context(), uint(allianceID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"members": members})
}

type ApplyToAllianceRequest struct {
	Message string `json:"message"`
}

func (h *Handlers) ApplyToAlliance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	allianceID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid alliance id"})
	}
	var req ApplyToAllianceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	err = h.allianceService.ApplyToAlliance(c.Context(), userID, uint(allianceID), req.Message)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetAllianceApplications(c *fiber.Ctx) error {
	allianceID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid alliance id"})
	}
	apps, err := h.allianceService.GetApplications(c.Context(), uint(allianceID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"applications": apps})
}

func (h *Handlers) AcceptApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	appID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid application id"})
	}
	err = h.allianceService.AcceptApplication(c.Context(), userID, uint(appID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) RejectApplication(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	appID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid application id"})
	}
	err = h.allianceService.RejectApplication(c.Context(), userID, uint(appID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) LeaveAlliance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	err := h.allianceService.LeaveAlliance(c.Context(), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) UpdateAlliance(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	var req CreateAllianceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	err := h.allianceService.UpdateAlliance(c.Context(), userID, service.CreateAllianceInput{
		Name: req.Name, Tag: req.Tag, Description: req.Description, Logo: req.Logo, Website: req.Website,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetBuddies(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	buddies, err := h.buddyService.GetBuddies(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"buddies": buddies})
}

type SendBuddyRequest struct {
	ReceiverID uint   `json:"receiver_id"`
	Message    string `json:"message"`
}

func (h *Handlers) SendBuddyRequest(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	var req SendBuddyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	err := h.buddyService.SendRequest(c.Context(), userID, req.ReceiverID, req.Message)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetPendingBuddyRequests(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	requests, err := h.buddyService.GetPendingRequests(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"requests": requests})
}

func (h *Handlers) AcceptBuddyRequest(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	requestID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request id"})
	}
	err = h.buddyService.AcceptRequest(c.Context(), userID, uint(requestID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) RejectBuddyRequest(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	requestID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request id"})
	}
	err = h.buddyService.RejectRequest(c.Context(), userID, uint(requestID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) RemoveBuddy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	buddyID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid buddy id"})
	}
	err = h.buddyService.RemoveBuddy(c.Context(), userID, uint(buddyID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetJumpGateTargets(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	moonID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid moon id"})
	}
	targets, err := h.moonService.GetJumpGateTargets(c.Context(), uint(moonID), userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"targets": targets})
}

type JumpGateRequest struct {
	TargetMoonID uint              `json:"target_moon_id"`
	Ships        map[string]int    `json:"ships"`
}

func (h *Handlers) ExecuteJumpGate(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	moonID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid moon id"})
	}
	var req JumpGateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	err = h.moonService.ExecuteJumpGate(c.Context(), userID, service.JumpFleetInput{
		OriginMoonID: uint(moonID), TargetMoonID: req.TargetMoonID, Ships: req.Ships,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

type PhalanxScanRequest struct {
	Galaxy   int `json:"galaxy"`
	System   int `json:"system"`
	Position int `json:"position"`
}

func (h *Handlers) ScanWithPhalanx(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	moonID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid moon id"})
	}
	var req PhalanxScanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	result, err := h.moonService.ScanWithPhalanx(c.Context(), uint(moonID), userID, req.Galaxy, req.System, req.Position)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"result": result})
}

func (h *Handlers) GetEspionageReports(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	reports, err := h.espionageService.GetReports(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reports": reports, "limit": limit, "offset": offset})
}

func (h *Handlers) GetEspionageUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	count, err := h.espionageService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"unread_count": count})
}

func (h *Handlers) MarkEspionageReportRead(c *fiber.Ctx) error {
	reportID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid report id"})
	}
	err = h.espionageService.MarkAsRead(c.Context(), uint(reportID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) DeleteEspionageReport(c *fiber.Ctx) error {
	reportID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid report id"})
	}
	err = h.espionageService.DeleteReport(c.Context(), uint(reportID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true})
}

func (h *Handlers) GetVacationStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	onVacation, endTime, err := h.authService.GetVacationStatus(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"on_vacation": onVacation, "vacation_end_time": endTime})
}

func (h *Handlers) EnableVacationMode(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	err := h.authService.SetVacationMode(c.Context(), userID, true, nil)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "on_vacation": true})
}

func (h *Handlers) DisableVacationMode(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	err := h.authService.SetVacationMode(c.Context(), userID, false, nil)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "on_vacation": false})
}

func (h *Handlers) GetDebrisFields(c *fiber.Ctx) error {
	debris, err := h.debrisService.GetDebrisFields(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": debris})
}

func (h *Handlers) GetDebrisField(c *fiber.Ctx) error {
	galaxy, err := strconv.Atoi(c.Params("galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}
	system, err := strconv.Atoi(c.Params("system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}
	position, err := strconv.Atoi(c.Params("position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid position"})
	}

	debris, err := h.debrisService.GetDebrisField(c.Context(), galaxy, system, position)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "debris field not found"})
	}
	return c.JSON(fiber.Map{"data": debris})
}

func (h *Handlers) CollectDebris(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	galaxy, err := strconv.Atoi(c.Params("galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}
	system, err := strconv.Atoi(c.Params("system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}
	position, err := strconv.Atoi(c.Params("position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid position"})
	}

	recyclerCapacity := int64(10000)
	if capStr := c.FormValue("recycler_capacity"); capStr != "" {
		if cap, err := strconv.ParseInt(capStr, 10, 64); err == nil {
			recyclerCapacity = cap
		}
	}

	metal, crystal, err := h.debrisService.CollectDebris(c.Context(), userID, galaxy, system, position, recyclerCapacity)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"metal":       metal,
		"crystal":     crystal,
		"recycled":    metal + crystal,
	})
}

func (h *Handlers) GetWreckFields(c *fiber.Ctx) error {
	wrecks, err := h.debrisService.GetWreckFields(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": wrecks})
}

func (h *Handlers) GetWreckField(c *fiber.Ctx) error {
	galaxy, err := strconv.Atoi(c.Params("galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}
	system, err := strconv.Atoi(c.Params("system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}
	position, err := strconv.Atoi(c.Params("position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid position"})
	}

	wreck, err := h.debrisService.GetWreckField(c.Context(), galaxy, system, position)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "wreck field not found"})
	}
	return c.JSON(fiber.Map{"data": wreck})
}

func (h *Handlers) CollectWreckField(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	galaxy, err := strconv.Atoi(c.Params("galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}
	system, err := strconv.Atoi(c.Params("system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}
	position, err := strconv.Atoi(c.Params("position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid position"})
	}

	cargoCapacity := int64(10000)
	if cargoStr := c.FormValue("cargo_capacity"); cargoStr != "" {
		if cap, err := strconv.ParseInt(cargoStr, 10, 64); err == nil {
			cargoCapacity = cap
		}
	}

	metal, crystal, deuterium, err := h.debrisService.CollectWreckField(c.Context(), userID, galaxy, system, position, cargoCapacity)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":     true,
		"metal":       metal,
		"crystal":     crystal,
		"deuterium":   deuterium,
		"recycled":    metal + crystal + deuterium,
	})
}

func (h *Handlers) GetNPCPlanets(c *fiber.Ctx) error {
	planets, err := h.npcService.GetNPCPlanets(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": planets})
}

func (h *Handlers) GetNPCPlanet(c *fiber.Ctx) error {
	galaxy, err := strconv.Atoi(c.Params("galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}
	system, err := strconv.Atoi(c.Params("system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}
	position, err := strconv.Atoi(c.Params("position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid position"})
	}

	planet, err := h.npcService.GetNPCPlanet(c.Context(), galaxy, system, position)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": planet})
}

func (h *Handlers) CreateNPCPlanet(c *fiber.Ctx) error {
	galaxy, err := strconv.Atoi(c.FormValue("galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid galaxy"})
	}
	system, err := strconv.Atoi(c.FormValue("system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid system"})
	}
	position, err := strconv.Atoi(c.FormValue("position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid position"})
	}

	level := 1
	if levelStr := c.FormValue("level"); levelStr != "" {
		if l, err := strconv.Atoi(levelStr); err == nil {
			level = l
		}
	}

	defenseLevel := 0
	if defStr := c.FormValue("defense_level"); defStr != "" {
		if d, err := strconv.Atoi(defStr); err == nil {
			defenseLevel = d
		}
	}

	resources := int64(10000)
	if resStr := c.FormValue("resources"); resStr != "" {
		if r, err := strconv.ParseInt(resStr, 10, 64); err == nil {
			resources = r
		}
	}

	config := service.NPCPlanetConfig{
		Galaxy:       galaxy,
		System:       system,
		Position:     position,
		Level:        level,
		Resources:    resources,
		DefenseLevel: defenseLevel,
	}

	planet, err := h.npcService.CreateNPCPlanet(c.Context(), config)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": planet})
}

func (h *Handlers) GenerateExpeditionFleet(c *fiber.Ctx) error {
	galaxy, err := strconv.Atoi(c.FormValue("galaxy"))
	if err != nil {
		galaxy = 1
	}
	system, err := strconv.Atoi(c.FormValue("system"))
	if err != nil {
		system = 250
	}
	position, err := strconv.Atoi(c.FormValue("position"))
	if err != nil {
		position = 8
	}

	fleet, err := h.npcService.GenerateExpeditionFleet(c.Context(), galaxy, system, position)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": fleet})
}

func (h *Handlers) GeneratePirateRaid(c *fiber.Ctx) error {
	galaxy, err := strconv.Atoi(c.FormValue("target_galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid target_galaxy"})
	}
	system, err := strconv.Atoi(c.FormValue("target_system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid target_system"})
	}
	position, err := strconv.Atoi(c.FormValue("target_position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid target_position"})
	}

	fleet, err := h.npcService.GeneratePirateRaid(c.Context(), galaxy, system, position)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": fleet})
}

func (h *Handlers) CreateACS(c *fiber.Ctx) error {
	fleetID, err := strconv.ParseUint(c.FormValue("fleet_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid fleet_id"})
	}

	targetGalaxy, _ := strconv.Atoi(c.FormValue("target_galaxy"))
	targetSystem, _ := strconv.Atoi(c.FormValue("target_system"))
	targetPosition, _ := strconv.Atoi(c.FormValue("target_position"))

	fleet, err := h.fleetService.GetFleetByID(c.Context(), uint(fleetID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "fleet not found"})
	}

	params := service.ACSCreateParams{
		TargetGalaxy:   targetGalaxy,
		TargetSystem:   targetSystem,
		TargetPosition: targetPosition,
		ArrivalTime:    fleet.ArrivalTime,
		FleetID:        uint(fleetID),
	}

	acs, err := h.acsService.CreateACS(c.Context(), params)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": acs})
}

func (h *Handlers) JoinACS(c *fiber.Ctx) error {
	acsID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid acs_id"})
	}

	fleetID, err := strconv.ParseUint(c.FormValue("fleet_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid fleet_id"})
	}

	acs, err := h.acsService.JoinACS(c.Context(), uint(acsID), uint(fleetID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": acs})
}

func (h *Handlers) GetACS(c *fiber.Ctx) error {
	acsID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid acs_id"})
	}

	acs, err := h.acsService.GetACS(c.Context(), uint(acsID))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "acs not found"})
	}

	return c.JSON(fiber.Map{"data": acs})
}

func (h *Handlers) GetACSFleets(c *fiber.Ctx) error {
	acsID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid acs_id"})
	}

	fleets, err := h.acsService.GetACSFleets(c.Context(), uint(acsID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": fleets})
}

func (h *Handlers) GetPremiumStatus(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	isActive, endsAt, err := h.premiumService.GetPremiumStatus(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"is_active":       isActive,
		"premium_ends_at": endsAt,
	})
}

func (h *Handlers) ActivatePremium(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	days, err := strconv.Atoi(c.FormValue("days"))
	if err != nil || days <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid days"})
	}

	err = h.premiumService.ActivatePremium(c.Context(), userID, days)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "premium activated"})
}

func (h *Handlers) MerchantBuy(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	planetID, err := strconv.ParseUint(c.FormValue("planet_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet_id"})
	}

	resourceType := c.FormValue("resource_type")
	amount, err := strconv.ParseInt(c.FormValue("amount"), 10, 64)
	if err != nil || amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid amount"})
	}

	err = h.premiumService.BuyResource(c.Context(), userID, uint(planetID), resourceType, amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "resource purchased"})
}

func (h *Handlers) MerchantSell(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	planetID, err := strconv.ParseUint(c.FormValue("planet_id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet_id"})
	}

	resourceType := c.FormValue("resource_type")
	amount, err := strconv.ParseInt(c.FormValue("amount"), 10, 64)
	if err != nil || amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "invalid amount"})
	}

	err = h.premiumService.SellResource(c.Context(), userID, uint(planetID), resourceType, amount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "message": "resource sold"})
}

func (h *Handlers) MovePlanet(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	planetID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet_id"})
	}

	targetGalaxy, err := strconv.Atoi(c.FormValue("target_galaxy"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid target_galaxy"})
	}

	targetSystem, err := strconv.Atoi(c.FormValue("target_system"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid target_system"})
	}

	targetPosition, err := strconv.Atoi(c.FormValue("target_position"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid target_position"})
	}

	input := service.MovePlanetInput{
		PlanetID:       uint(planetID),
		TargetGalaxy:   targetGalaxy,
		TargetSystem:   targetSystem,
		TargetPosition: targetPosition,
	}

	planet, err := h.planetService.MovePlanet(c.Context(), userID, input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "data": planet})
}

func (h *Handlers) GetCharacterClass(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	classID, err := h.characterClassService.GetCharacterClass(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	classes := service.GetAllCharacterClasses()
	var currentClass service.CharacterClassInfo
	for _, c := range classes {
		if c.ID == classID {
			currentClass = c
			break
		}
	}

	return c.JSON(fiber.Map{
		"character_class": classID,
		"class_info":      currentClass,
		"available_classes": classes,
	})
}

func (h *Handlers) SelectCharacterClass(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	if userID == 0 {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	classID, err := strconv.Atoi(c.FormValue("class_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid class_id"})
	}

	err = h.characterClassService.SetCharacterClass(c.Context(), userID, classID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "character_class": classID})
}
