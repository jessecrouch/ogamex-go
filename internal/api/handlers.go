package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"ogamex-go/internal/dto"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
	"ogamex-go/internal/service"
)

type BuildingRequest struct {
	Level int `json:"level"`
}

type ResearchRequest struct {
	ResearchID int `json:"research_id"`
}

type FleetRequest struct {
	MissionType    int `json:"mission_type"`
	OriginGalaxy   int `json:"origin_galaxy"`
	OriginSystem   int `json:"origin_system"`
	OriginPosition int `json:"origin_position"`
	TargetGalaxy   int `json:"target_galaxy"`
	TargetSystem   int `json:"target_system"`
	TargetPosition int `json:"target_position"`
	Resources      struct {
		Metal     int64 `json:"metal"`
		Crystal   int64 `json:"crystal"`
		Deuterium int64 `json:"deuterium"`
	} `json:"resources"`
	Ships map[string]int `json:"ships"`
}

type BuildUnitRequest struct {
	UnitID int `json:"unit_id"`
	Amount int `json:"amount"`
}

type CreateNoteRequest struct {
	Galaxy   int    `json:"galaxy"`
	System   int    `json:"system"`
	Position int    `json:"position"`
	Type     int    `json:"type"`
	Subject  string `json:"subject"`
	Text     string `json:"text"`
}

type UpdateNoteRequest struct {
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type CreateAllianceRequest struct {
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Website     string `json:"website"`
}

type ApplyToAllianceRequest struct {
	Message string `json:"message"`
}

type SendBuddyRequest struct {
	ReceiverID uint   `json:"receiver_id"`
	Message    string `json:"message"`
}

type JumpGateRequest struct {
	TargetMoonID uint           `json:"target_moon_id"`
	Ships        map[string]int `json:"ships"`
}

type PhalanxScanRequest struct {
	Galaxy   int `json:"galaxy"`
	System   int `json:"system"`
	Position int `json:"position"`
}

type Handlers struct {
	buildingService       *service.BuildingService
	researchService       *service.ResearchService
	fleetService          *service.FleetService
	planetRepo            repository.PlanetRepository
	authService           *service.AuthService
	unitService           *service.UnitService
	productionService     *service.ProductionService
	messageService        *service.MessageService
	planetService         *service.PlanetService
	noteService           *service.NoteService
	allianceService       *service.AllianceService
	buddyService          *service.BuddyService
	moonService           *service.MoonService
	espionageService      *service.EspionageService
	debrisService         *service.DebrisService
	npcService            *service.NPCService
	acsService            *service.ACSService
	premiumService        *service.PremiumService
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
		buildingService:       buildingService,
		researchService:       researchService,
		fleetService:          fleetService,
		planetRepo:            planetRepo,
		authService:           authService,
		unitService:           unitService,
		productionService:     productionService,
		messageService:        messageService,
		planetService:         planetService,
		noteService:           noteService,
		allianceService:       allianceService,
		buddyService:          buddyService,
		moonService:           moonService,
		espionageService:      espionageService,
		debrisService:         debrisService,
		npcService:            npcService,
		acsService:            acsService,
		premiumService:        premiumService,
		characterClassService: characterClassService,
	}
}

func (h *Handlers) SetupRoutes(app *fiber.App) {
	// TEST ROUTE - should be completely public
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test works"})
	})

	api := app.Group("/api/v1")

	api.Get("/status", h.GetStatus)

	api.Post("/auth/register", h.Register)
	api.Post("/auth/login", h.Login)

	api.Post("/battle/simulate", h.SimulateBattle)

	// Reference data endpoints (public, no auth required)
	api.Get("/buildings", h.GetBuildings)
	api.Get("/buildings/:id", h.GetBuilding)
	api.Get("/ships", h.GetShips)
	api.Get("/ships/:id", h.GetShip)
	api.Get("/defense", h.GetDefense)
	api.Get("/defense/:id", h.GetDefenseUnit)
	api.Get("/research", h.GetResearch)
	api.Get("/research/:id", h.GetResearchType)
	api.Get("/missions", h.GetMissions)
	api.Get("/game", h.GetGame)

	// Protected routes - these require auth
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

	// Debug/test endpoint to fix production percentages - requires auth but not admin
	protectedDebug := protected.Group("")
	protectedDebug.Post("/debug/fix-production", h.FixProductionPercentages)

	// Admin group - requires admin middleware
	admin := protected.Group("")
	admin.Post("/admin/fix-planets", h.FixPlanets)

	protected.Post("/research/start", h.StartResearch)
	protected.Get("/research/queue", h.GetResearchQueue)
	protected.Delete("/research/queue/:id", h.CancelResearch)

	protected.Post("/fleets/send", h.SendFleet)
	protected.Get("/fleets", h.GetFleets)
	protected.Post("/fleets/:id/recall", h.RecallFleet)

	protected.Post("/planets/:id/units/build", h.BuildUnit)
	protected.Get("/planets/:id/units/queue", h.GetUnitQueue)
	protected.Delete("/units/queue/:id", h.CancelUnit)

	// NEW: Ship building endpoints (cleaner URLs)
	protected.Post("/planets/:id/ships/:ship_id", h.BuildUnit)
	protected.Get("/planets/:id/ships", h.GetPlanetShips)

	// NEW: Defense building endpoints
	protected.Post("/planets/:id/defense/:defense_id", h.BuildUnit)
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

func (h *Handlers) adminMiddleware(c *fiber.Ctx) error {
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

	if user.ID != 1 {
		return c.Status(403).JSON(fiber.Map{"error": "admin access required"})
	}

	c.Locals("user_id", user.ID)
	c.Locals("user", user)

	return c.Next()
}

// FixPlanets fixes planet data (admin only)
// @Summary Fix planets
// @Description Fix planet data - admin utility endpoint
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /admin/fix-planets [post]
func (h *Handlers) FixPlanets(c *fiber.Ctx) error {
	count, err := h.planetRepo.FixProductionPercentages(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"planets_fixed": count,
	})
}

// FixProductionPercentages fixes production percentages (debug only)
// @Summary Fix production
// @Description Fix production percentages - debug utility endpoint
// @Tags Debug
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /debug/fix-production [post]
func (h *Handlers) FixProductionPercentages(c *fiber.Ctx) error {
	count, err := h.planetRepo.FixProductionPercentages(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"planets_fixed": count,
	})
}

// GetStatus returns server status and version information
// @Summary Get server status
// @Description Get server status and version information
// @Tags Server
// @Produce json
// @Success 200
// @Router /status [get]
func (h *Handlers) GetStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"online":   true,
		"version":  "0.0.1",
		"universe": "ogamex-go",
		"api_base": "/api/v1",
		"time":     time.Now().Unix(),
	})
}

// GetUser returns the current user's information
// @Summary Get current user
// @Description Get current user information
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /user [get]
func (h *Handlers) GetUser(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}

	u := user.(*schema.User)

	return c.JSON(fiber.Map{
		"id":              u.ID,
		"username":        u.Username,
		"email":           u.Email,
		"player_name":     u.PlayerName,
		"dark_matter":     u.DarkMatter,
		"character_class": u.CharacterClass,
		"current_planet":  u.CurrentPlanetID,
	})
}

// GetUserStats returns statistics for the current user
// @Summary Get user stats
// @Description Get user statistics including total resources across all planets
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /user/stats [get]
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
		"user_id": userID,
		"planets": planetCount,
		"total_resources": fiber.Map{
			"metal":     totalMetal,
			"crystal":   totalCrystal,
			"deuterium": totalDeuterium,
		},
	})
}

// GetUserPlanets returns all planets owned by the current user
// @Summary Get user planets
// @Description Get all planets owned by the current user
// @Tags Planets
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets [get]
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

// GetPlanet returns basic information about a specific planet
// @Summary Get planet
// @Description Get basic planet information
// @Tags Planets
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Failure 404
// @Router /planets/{id} [get]
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

// GetPlanetDetails returns detailed information about a planet
// @Summary Get planet details
// @Description Get detailed planet information including buildings and resources
// @Tags Planets
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Failure 404
// @Router /planets/{id}/details [get]
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
		"id":                 planet.ID,
		"name":               planet.Name,
		"galaxy":             planet.Galaxy,
		"system":             planet.System,
		"position":           planet.Position,
		"is_moon":            planet.IsMoon,
		"planet_type":        planet.PlanetType,
		"metal":              planet.Metal,
		"crystal":            planet.Crystal,
		"deuterium":          planet.Deuterium,
		"metal_capacity":     planet.MetalCapacity,
		"crystal_capacity":   planet.CrystalCapacity,
		"deuterium_capacity": planet.DeuteriumCapacity,
		"energy_available":   planet.EnergyAvailable,
		"energy_max":         planet.EnergyMax,
		"energy_used":        planet.EnergyUsed,
		"fields_used":        planet.FieldsUsed,
		"fields_max":         planet.FieldsMax,
		"temp_min":           planet.TempMin,
		"temp_max":           planet.TempMax,
		"defense_activated":  planet.DefenseActivated,
	})
}

// SetCurrentPlanet sets the specified planet as the user's current active planet
// @Summary Set current planet
// @Description Set a planet as the current active planet
// @Tags Planets
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Failure 404
// @Router /planets/{id}/set-current [put]
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
		userPlanets, _ := h.planetRepo.GetByUserID(c.Context(), userID)
		planetIDs := make([]uint, len(userPlanets))
		for i, p := range userPlanets {
			planetIDs[i] = p.ID
		}
		return c.Status(403).JSON(fiber.Map{
			"error":        "planet does not belong to user",
			"your_planets": planetIDs,
		})
	}

	err = h.authService.SetCurrentPlanet(c.Context(), userID, uint(planetID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true, "current_planet_id": planetID})
}

// GetPlanetOverview returns planet overview with resources, production, and buildings
// @Summary Get planet overview
// @Description Get detailed overview of a planet including resources, production rates, and buildings
// @Tags Planets
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Success 200
// @Failure 404
// @Security BearerAuth
// @Router /planets/{id}/overview [get]
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

	user := c.Locals("user").(*schema.User)
	tech, _ := h.authService.GetUserTech(c.Context(), user.ID)

	production := h.productionService.CalculateProduction(planet, tech)

	return c.JSON(fiber.Map{
		"planet_id": planet.ID,
		"name":      planet.Name,
		"resources": fiber.Map{
			"metal":     planet.Metal,
			"crystal":   planet.Crystal,
			"deuterium": planet.Deuterium,
			"energy":    planet.EnergyAvailable,
		},
		"production": fiber.Map{
			"metal":     production.Metal,
			"crystal":   production.Crystal,
			"deuterium": production.Deuterium,
		},
		"buildings": fiber.Map{
			"metal_mine":      planet.MetalMine,
			"crystal_mine":    planet.CrystalMine,
			"deuterium_synth": planet.DeuteriumSynthesizer,
			"solar_plant":     planet.SolarPlant,
			"fusion_plant":    planet.FusionPlant,
		},
		"production_percent": fiber.Map{
			"metal_mine":      planet.MetalMinePercent,
			"crystal_mine":    planet.CrystalMinePercent,
			"deuterium_synth": planet.DeuteriumSynthesizerPercent,
		},
		"has_queue":   hasQueue,
		"fields_used": planet.FieldsUsed,
		"fields_max":  planet.FieldsMax,
	})
}

// GetPlanetResources returns current resources on a planet
// @Summary Get planet resources
// @Description Get current resources on a planet
// @Tags Planets
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/resources [get]
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
		"planet_id":          planet.ID,
		"metal":              planet.Metal,
		"crystal":            planet.Crystal,
		"deuterium":          planet.Deuterium,
		"energy":             planet.EnergyAvailable,
		"metal_capacity":     planet.MetalCapacity,
		"crystal_capacity":   planet.CrystalCapacity,
		"deuterium_capacity": planet.DeuteriumCapacity,
	})
}

// GetPlanetBuildings returns all buildings on a planet with their levels
// @Summary Get planet buildings
// @Description Get all buildings on a planet with their levels
// @Tags Buildings
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/buildings [get]
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
			"crystal_mine":          planet.CrystalMine,
			"deuterium_synthesizer": planet.DeuteriumSynthesizer,
			"solar_plant":           planet.SolarPlant,
			"fusion_plant":          planet.FusionPlant,
			"metal_storage":         planet.MetalStorageBuilding,
			"crystal_storage":       planet.CrystalStorageBuilding,
			"deuterium_storage":     planet.DeuteriumStorageBuilding,
			"robot_factory":         planet.RobotFactory,
			"shipyard":              planet.Shipyard,
			"research_lab":          planet.ResearchLab,
			"nanite_factory":        planet.NaniteFactory,
			"terraformer":           planet.Terraformer,
			"space_dock":            planet.SpaceDock,
		},
		"fields_used": planet.FieldsUsed,
		"fields_max":  planet.FieldsMax,
	})
}

// StartBuilding starts constructing a building
// @Summary Start building construction
// @Description Start building a structure on a planet (mine, power plant, etc.)
// @Tags Buildings
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Param building_id path int true "Building Type ID"
// @Param body body BuildingRequest true "Building level to construct"
// @Success 200
// @Failure 400
// @Security BearerAuth
// @Router /planets/{id}/buildings/{building_id} [post]
func (h *Handlers) StartBuilding(c *fiber.Ctx) error {
	planetID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid planet id"})
	}

	buildingID, err := c.ParamsInt("building_id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid building id"})
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

// GetBuildingQueue returns the building queue for a planet
// @Summary Get building queue
// @Description Get building queue for a planet
// @Tags Buildings
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/queue [get]
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
		ID        uint  `json:"id"`
		Building  int   `json:"building_id"`
		Level     int   `json:"level"`
		StartTime int64 `json:"start_time"`
		EndTime   int64 `json:"end_time"`
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

// CancelBuilding cancels a building in the queue
// @Summary Cancel building
// @Description Cancel a building in the construction queue
// @Tags Buildings
// @Produce json
// @Param id path int true "Queue ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /queue/{id} [delete]
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

// StartResearch starts a research project
// @Summary Start research
// @Description Start a research project
// @Tags Research
// @Accept json
// @Produce json
// @Param body body ResearchRequest true "Research ID to start"
// @Security BearerAuth
// @Success 200
// @Failure 400
// @Failure 401
// @Router /research/start [post]
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
		"success":     true,
		"user_id":     userID,
		"research_id": req.ResearchID,
	})
}

// GetResearchQueue returns the research queue for the current user
// @Summary Get research queue
// @Description Get research queue for the current user
// @Tags Research
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /research/queue [get]
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
		ID         uint  `json:"id"`
		ResearchID int   `json:"research_id"`
		Level      int   `json:"level"`
		StartTime  int64 `json:"start_time"`
		EndTime    int64 `json:"end_time"`
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

// SendFleet sends a fleet mission
// @Summary Send fleet
// @Description Send a fleet of ships on a mission (attack, transport, colonize, recycle, etc.)
// @Tags Fleets
// @Accept json
// @Produce json
// @Param body body FleetRequest true "Fleet mission details"
// @Success 200
// @Failure 400
// @Security BearerAuth
// @Router /fleets/send [post]
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

// GetFleets returns all active fleet missions for the current user
// @Summary Get fleets
// @Description Get all active fleet missions for the current user
// @Tags Fleets
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /fleets [get]
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

// SimulateBattle simulates a battle without sending fleets
// @Summary Simulate battle
// @Description Run a battle simulation (SpeedSim style) - supports multiple attackers, defense, ACS, and Monte Carlo simulations
// @Tags Battle
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Router /battle/simulate [post]
func (h *Handlers) SimulateBattle(c *fiber.Ctx) error {
	var req dto.BattleSimulateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if len(req.AttackerFleets) == 0 || len(req.DefenderFleets) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "attacker and defender fleets required"})
	}

	// Set defaults
	if req.SimulationCount < 1 {
		req.SimulationCount = 1
	}
	if req.SimulationCount > 1000 {
		req.SimulationCount = 1000 // Cap at 1000 for performance
	}

	// Run Monte Carlo simulations
	attackerWins := 0
	defenderWins := 0
	draws := 0

	var totalAttackerLeft map[string]int
	var totalDefenderLeft map[string]int
	var totalAttackerLost map[string]int
	var totalDefenderLost map[string]int
	var totalDebrisMetal int64
	var totalDebrisCrystal int64
	var totalMoonChance float64

	// Run simulations
	for sim := 0; sim < req.SimulationCount; sim++ {
		// Convert fleets to Rust format
		attackerFleets := convertFleetsToRust(req.AttackerFleets)
		defenderFleets := convertFleetsToRust(req.DefenderFleets)

		var attackerTech, defenderTech *schema.UserTech
		if req.AttackerTech != nil {
			attackerTech = &schema.UserTech{
				WeaponsTechnology:   req.AttackerTech.Weapons,
				ShieldingTechnology: req.AttackerTech.Shielding,
				ArmorTechnology:     req.AttackerTech.Armor,
			}
		}
		if req.DefenderTech != nil {
			defenderTech = &schema.UserTech{
				WeaponsTechnology:   req.DefenderTech.Weapons,
				ShieldingTechnology: req.DefenderTech.Shielding,
				ArmorTechnology:     req.DefenderTech.Armor,
			}
		}

		result, err := h.fleetService.SimulateBattleWithRust(attackerFleets, defenderFleets, attackerTech, defenderTech)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		// Aggregate results
		attackerLeft := convertRustResultToMap(result.AttackerRemaining)
		defenderLeft := convertRustResultToMap(result.DefenderRemaining)
		attackerLost := convertRustResultToMap(result.AttackerLosses)
		defenderLost := convertRustResultToMap(result.DefenderLosses)

		// Track wins/losses/draws
		if result.Winner == "attacker" {
			attackerWins++
		} else if result.Winner == "defender" {
			defenderWins++
		} else {
			draws++
		}

		// Accumulate for averaging
		if totalAttackerLeft == nil {
			totalAttackerLeft = make(map[string]int)
			totalDefenderLeft = make(map[string]int)
			totalAttackerLost = make(map[string]int)
			totalDefenderLost = make(map[string]int)
		}

		for k, v := range attackerLeft {
			totalAttackerLeft[k] += v
		}
		for k, v := range defenderLeft {
			totalDefenderLeft[k] += v
		}
		for k, v := range attackerLost {
			totalAttackerLost[k] += v
		}
		for k, v := range defenderLost {
			totalDefenderLost[k] += v
		}

		// Calculate debris for this simulation
		debrisMetal, debrisCrystal := calculateDebris(attackerLost, defenderLost)
		totalDebrisMetal += debrisMetal
		totalDebrisCrystal += debrisCrystal

		// Calculate moon chance
		debrisTotal := debrisMetal + debrisCrystal
		if debrisTotal > 100000 {
			moonChance := float64(debrisTotal) / 1000000.0
			if moonChance > 0.2 {
				moonChance = 0.2
			}
			totalMoonChance += moonChance
		}
	}

	// Average the results for Monte Carlo
	if req.SimulationCount > 1 {
		avgMap(totalAttackerLeft, req.SimulationCount)
		avgMap(totalDefenderLeft, req.SimulationCount)
		avgMap(totalAttackerLost, req.SimulationCount)
		avgMap(totalDefenderLost, req.SimulationCount)
		totalDebrisMetal /= int64(req.SimulationCount)
		totalDebrisCrystal /= int64(req.SimulationCount)
		totalMoonChance /= float64(req.SimulationCount)
	}

	// Calculate ruins (only if moon chance > 0)
	var ruinsMetal, ruinsCrystal int64
	if totalMoonChance > 0 {
		ruinsMetal = totalDebrisMetal / 2
		ruinsCrystal = totalDebrisCrystal / 2
	}

	// Calculate plunder
	plunder := calculatePlunder(req.TargetResources, totalDefenderLost, req.AttackerFleets)

	// Calculate fuel and flight time
	var fuelUsed int64
	var flightTime *dto.FlightResult
	if req.OriginCoords != nil && req.TargetCoords != nil {
		fuelUsed = calculateFuel(req.AttackerFleets, req.OriginCoords, req.TargetCoords, req.AttackerTech)
		flightTime = calculateFlightTime(req.AttackerFleets, req.OriginCoords, req.TargetCoords, req.AttackerTech)
	}

	// Determine winner
	winner := "draw"
	winChance := float64(draws) / float64(req.SimulationCount) * 100
	if attackerWins > defenderWins && attackerWins > draws {
		winner = "attacker"
		winChance = float64(attackerWins) / float64(req.SimulationCount) * 100
	} else if defenderWins > attackerWins && defenderWins > draws {
		winner = "defender"
		winChance = float64(defenderWins) / float64(req.SimulationCount) * 100
	}

	response := dto.BattleSimulateResponse{
		Simulations:  req.SimulationCount,
		Winner:       winner,
		WinChance:    winChance,
		AttackerWins: attackerWins,
		DefenderWins: defenderWins,
		Draws:        draws,

		AttackerLeft: totalAttackerLeft,
		DefenderLeft: totalDefenderLeft,
		AttackerLost: totalAttackerLost,
		DefenderLost: totalDefenderLost,

		Debris: dto.DebrisInfo{
			Metal:   totalDebrisMetal,
			Crystal: totalDebrisCrystal,
			Total:   totalDebrisMetal + totalDebrisCrystal,
		},
		Ruins: dto.RuinsInfo{
			Metal:   ruinsMetal,
			Crystal: ruinsCrystal,
			Total:   ruinsMetal + ruinsCrystal,
		},
		MoonChance: totalMoonChance,

		Plunder:    plunder,
		FuelUsed:   fuelUsed,
		FlightTime: flightTime,
	}

	return c.JSON(response)
}

func convertFleetsToRust(fleets []dto.FleetComposition) map[int16]int16 {
	result := make(map[int16]int16)
	for _, fleet := range fleets {
		for name, count := range fleet.Ships {
			shipID := shipNameToID(name)
			if shipID == 0 {
				continue
			}
			result[shipID] += int16(count)
		}
		// Add defense units (they're treated as ships in battle)
		for name, count := range fleet.Defense {
			defID := defenseNameToID(name)
			if defID == 0 {
				continue
			}
			result[defID] += int16(count)
		}
	}
	return result
}

func convertRustResultToMap(result map[int16]int16) map[string]int {
	m := make(map[string]int)
	for id, count := range result {
		m[shipIDToName(id)] = int(count)
	}
	return m
}

func avgMap(m map[string]int, divisor int) {
	for k := range m {
		m[k] = m[k] / divisor
	}
}

func calculateDebris(attackerLost, defenderLost map[string]int) (metal, crystal int64) {
	for name, count := range attackerLost {
		cost := getShipCost(shipNameToID(name))
		metal += cost.Metal * int64(count)
		crystal += cost.Crystal * int64(count)
	}
	for name, count := range defenderLost {
		cost := getShipCost(shipNameToID(name))
		// Defense goes 70% to debris by default
		if isDefense(name) {
			metal += cost.Metal * int64(count) * 70 / 100
			crystal += cost.Crystal * int64(count) * 70 / 100
		} else {
			metal += cost.Metal * int64(count) / 2
			crystal += cost.Crystal * int64(count) / 2
		}
	}
	return
}

func isDefense(name string) bool {
	defenseNames := map[string]bool{
		"rocket_launcher":     true,
		"light_laser":         true,
		"heavy_laser":         true,
		"ion_cannon":          true,
		"gauss_cannon":        true,
		"plasma_turret":       true,
		"small_shield_dome":   true,
		"large_shield_dome":   true,
		"missile_interceptor": true,
		"missile_launcher":    true,
	}
	return defenseNames[strings.ToLower(name)]
}

func calculatePlunder(targetResources *dto.TargetResources, defenderLost map[string]int, attackerFleets []dto.FleetComposition) dto.PlunderInfo {
	result := dto.PlunderInfo{
		Theoretical: dto.ResourcesResponse{},
		Actual:      dto.ResourcesResponse{},
	}

	if targetResources == nil {
		return result
	}

	// Theoretical: 50% of resources on planet (max 75% of fleet capacity)
	// For simulation, we assume full resources
	result.Theoretical.Metal = targetResources.Metal / 2
	result.Theoretical.Crystal = targetResources.Crystal / 2
	result.Theoretical.Deuterium = targetResources.Deuterium / 2

	// Calculate available cargo space
	var totalCargo int64
	for _, fleet := range attackerFleets {
		for name, count := range fleet.Ships {
			cargo := getShipCargo(shipNameToID(name))
			totalCargo += cargo * int64(count)
		}
	}

	result.CargoNeeded = result.Theoretical.Metal + result.Theoretical.Crystal + result.Theoretical.Deuterium

	// Actual plunder is limited by cargo capacity
	if result.CargoNeeded > totalCargo {
		result.CargoNeeded = totalCargo
		// Proportional distribution
		total := result.Theoretical.Metal + result.Theoretical.Crystal + result.Theoretical.Deuterium
		if total > 0 {
			ratio := float64(totalCargo) / float64(total)
			result.Actual.Metal = int64(float64(result.Theoretical.Metal) * ratio)
			result.Actual.Crystal = int64(float64(result.Theoretical.Crystal) * ratio)
			result.Actual.Deuterium = int64(float64(result.Theoretical.Deuterium) * ratio)
		}
	} else {
		result.Actual = result.Theoretical
	}

	return result
}

func getShipCargo(shipID int16) int64 {
	cargo := map[int16]int64{
		202: 5000,    // small_cargo
		203: 25000,   // large_cargo
		204: 50,      // light_fighter
		205: 100,     // heavy_fighter
		206: 800,     // cruiser
		207: 1500,    // battleship
		208: 7500,    // colony_ship
		209: 20000,   // recycler
		210: 0,       // espionage_probe
		211: 500,     // bomber
		213: 2000,    // destroyer
		214: 1000000, // deathstar
		215: 750,     // battlecruiser
		218: 70000,   // reaper
		219: 15000,   // pathfinder
	}
	if c, ok := cargo[shipID]; ok {
		return c
	}
	return 0
}

func calculateFuel(fleets []dto.FleetComposition, origin, target *dto.Coordinates, tech *dto.TechLevel) int64 {
	if origin == nil || target == nil {
		return 0
	}

	distance := calculateDistance(origin.Galaxy, origin.System, origin.Position,
		target.Galaxy, target.System, target.Position)

	speed := 10000 // Base speed
	if tech != nil {
		if tech.HyperspaceDrive > 0 {
			speed = 5000 + (tech.HyperspaceDrive * 1000)
		} else if tech.ImpulseDrive > 0 {
			speed = 2000 + (tech.ImpulseDrive * 500)
		} else if tech.CombustionDrive > 0 {
			speed = 500 + (tech.CombustionDrive * 100)
		}
	}

	// Base consumption formula: distance * 1.5 * (ships * base_consumption) / speed
	var totalConsumption int64
	for _, fleet := range fleets {
		for name, count := range fleet.Ships {
			baseFuel := getShipBaseFuel(shipNameToID(name))
			totalConsumption += int64(count) * baseFuel * int64(distance) * 2 / int64(speed)
		}
	}

	return totalConsumption
}

func getShipBaseFuel(shipID int16) int64 {
	fuel := map[int16]int64{
		202: 1,
		203: 5,
		204: 1,
		205: 3,
		206: 10,
		207: 25,
		208: 15,
		209: 20,
		210: 1,
		211: 50,
		213: 30,
		214: 1000,
		215: 15,
		218: 40,
		219: 25,
	}
	if f, ok := fuel[shipID]; ok {
		return f
	}
	return 1
}

func calculateDistance(g1, s1, p1, g2, s2, p2 int) int {
	galaxyDiff := abs(g1 - g2)
	if galaxyDiff == 0 {
		systemDiff := abs(s1 - s2)
		if systemDiff == 0 {
			return abs(p1-p2) * 5
		}
		return systemDiff*20 + 2700
	}
	return galaxyDiff*20000 + 2700
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func calculateFlightTime(fleets []dto.FleetComposition, origin, target *dto.Coordinates, tech *dto.TechLevel) *dto.FlightResult {
	if origin == nil || target == nil || len(fleets) == 0 {
		return nil
	}

	distance := calculateDistance(origin.Galaxy, origin.System, origin.Position,
		target.Galaxy, target.System, target.Position)

	// Find slowest ship in fleet
	slowestSpeed := 1000000
	for _, fleet := range fleets {
		for name := range fleet.Ships {
			speed := getShipSpeed(shipNameToID(name), tech)
			if speed < slowestSpeed {
				slowestSpeed = speed
			}
		}
	}

	if slowestSpeed == 0 {
		slowestSpeed = 1
	}

	// Flight time in seconds: distance / speed * 3600 (hours to seconds)
	seconds := (distance * 3600) / slowestSpeed

	result := &dto.FlightResult{
		Seconds:   seconds,
		Hours:     seconds / 3600,
		Minutes:   (seconds % 3600) / 60,
		Formatted: formatTime(seconds),
	}

	// Add return time for round trip
	returnTime := seconds * 2
	result.ReturnTime = &dto.FlightResult{
		Seconds:   returnTime,
		Hours:     returnTime / 3600,
		Minutes:   (returnTime % 3600) / 60,
		Formatted: formatTime(returnTime),
	}

	return result
}

func getShipSpeed(shipID int16, tech *dto.TechLevel) int {
	baseSpeed := map[int16]int{
		202: 5000,
		203: 3000,
		204: 12500,
		205: 10000,
		206: 15000,
		207: 10000,
		208: 2500,
		209: 2000,
		210: 100000000, // Very fast probes
		211: 4000,
		213: 5000,
		214: 200,
		215: 10000,
		218: 7000,
		219: 10000,
	}

	speed := baseSpeed[shipID]

	// Apply drive tech bonuses
	if tech != nil {
		combustionBonus := 1.0 + float64(tech.CombustionDrive)*0.1
		impulseBonus := 1.0 + float64(tech.ImpulseDrive)*0.2
		hyperspaceBonus := 1.0 + float64(tech.HyperspaceDrive)*0.3

		// Use highest applicable bonus
		if tech.HyperspaceDrive > 0 {
			speed = int(float64(speed) * hyperspaceBonus)
		} else if tech.ImpulseDrive > 0 {
			speed = int(float64(speed) * impulseBonus)
		} else if tech.CombustionDrive > 0 {
			speed = int(float64(speed) * combustionBonus)
		}
	}

	return speed
}

func formatTime(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func shipNameToID(name string) int16 {
	mapping := map[string]int16{
		"small_cargo":     202,
		"large_cargo":     203,
		"light_fighter":   204,
		"heavy_fighter":   205,
		"cruiser":         206,
		"battleship":      207,
		"colony_ship":     208,
		"colonizer":       208,
		"recycler":        209,
		"espionage_probe": 210,
		"bomber":          211,
		"destroyer":       213,
		"deathstar":       214,
		"battlecruiser":   215,
		"reaper":          218,
		"pathfinder":      219,
		"solar_satellite": 212,
		"crawler":         217,
	}
	if id, ok := mapping[strings.ToLower(name)]; ok {
		return id
	}
	return defenseNameToID(name)
}

func defenseNameToID(name string) int16 {
	mapping := map[string]int16{
		"rocket_launcher":     401,
		"light_laser":         402,
		"heavy_laser":         403,
		"ion_cannon":          404,
		"gauss_cannon":        405,
		"plasma_turret":       406,
		"small_shield_dome":   407,
		"small_shield":        407,
		"large_shield_dome":   408,
		"large_shield":        408,
		"missile_interceptor": 428,
		"missile_launcher":    429,
	}
	if id, ok := mapping[strings.ToLower(name)]; ok {
		return id
	}
	return 0
}

func shipIDToName(id int16) string {
	mapping := map[int16]string{
		202: "small_cargo",
		203: "large_cargo",
		204: "light_fighter",
		205: "heavy_fighter",
		206: "cruiser",
		207: "battleship",
		208: "colonizer",
		209: "recycler",
		210: "espionage_probe",
		211: "bomber",
		212: "solar_satellite",
		213: "destroyer",
		214: "deathstar",
		215: "battlecruiser",
		217: "crawler",
		218: "reaper",
		219: "pathfinder",
		401: "rocket_launcher",
		402: "light_laser",
		403: "heavy_laser",
		404: "ion_cannon",
		405: "gauss_cannon",
		406: "plasma_turret",
		407: "small_shield_dome",
		408: "large_shield_dome",
		428: "missile_interceptor",
		429: "missile_launcher",
	}
	if name, ok := mapping[id]; ok {
		return name
	}
	return fmt.Sprintf("unit_%d", id)
}

type shipCost struct {
	Metal   int64
	Crystal int64
}

func getShipCost(shipID int16) shipCost {
	costs := map[int16]shipCost{
		202: {Metal: 2000, Crystal: 0},
		203: {Metal: 6000, Crystal: 0},
		204: {Metal: 3000, Crystal: 1000},
		205: {Metal: 6000, Crystal: 4000},
		206: {Metal: 20000, Crystal: 7000},
		207: {Metal: 45000, Crystal: 15000},
		208: {Metal: 10000, Crystal: 20000},
		209: {Metal: 10000, Crystal: 6000},
		210: {Metal: 0, Crystal: 1000},
		211: {Metal: 50000, Crystal: 25000},
		212: {Metal: 0, Crystal: 2000}, // solar satellite
		213: {Metal: 60000, Crystal: 50000},
		214: {Metal: 5000000, Crystal: 4000000},
		215: {Metal: 30000, Crystal: 4000},
		217: {Metal: 10000, Crystal: 5000}, // crawler
		218: {Metal: 70000, Crystal: 40000},
		219: {Metal: 40000, Crystal: 20000},
		// Defense
		401: {Metal: 2000, Crystal: 0},
		402: {Metal: 1500, Crystal: 500},
		403: {Metal: 6000, Crystal: 2000},
		404: {Metal: 20000, Crystal: 15000},
		405: {Metal: 10000, Crystal: 20000},
		406: {Metal: 50000, Crystal: 50000},
		407: {Metal: 10000, Crystal: 0},
		408: {Metal: 50000, Crystal: 50000},
		428: {Metal: 8000, Crystal: 2000},  // missile_interceptor
		429: {Metal: 12500, Crystal: 2500}, // missile_launcher
	}
	if cost, ok := costs[shipID]; ok {
		return cost
	}
	return shipCost{Metal: 0, Crystal: 0}
}

// RecallFleet recalls a fleet mission
// @Summary Recall fleet
// @Description Recall a fleet mission that is currently in progress
// @Tags Fleets
// @Produce json
// @Param id path int true "Fleet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /fleets/{id}/recall [post]
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

// CancelResearch cancels a research in the queue
// @Summary Cancel research
// @Description Cancel a research in the queue
// @Tags Research
// @Produce json
// @Param id path int true "Queue ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /research/queue/{id} [delete]
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

// Register creates a new user account
// @Summary Register new user
// @Description Create a new user account and receive an auth token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register request"
// @Success 201
// @Failure 400
// @Router /auth/register [post]
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
		Username:   req.Username,
		Email:      req.Email,
		Password:   req.Password,
		PlayerName: req.PlayerName,
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"user_id":    user.ID,
		"auth_token": user.AuthToken,
		"username":   user.Username,
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login authenticates a user and returns an auth token
// @Summary Login user
// @Description Authenticate with username and password to receive an auth token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login request"
// @Success 200
// @Failure 401
// @Router /auth/login [post]
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
		"success":    true,
		"user_id":    user.ID,
		"auth_token": user.AuthToken,
		"username":   user.Username,
	})
}

// BuildUnit builds a ship or defense unit
// @Summary Build unit
// @Description Build a ship or defense unit on a planet
// @Tags Units
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Param body body BuildUnitRequest true "Unit ID and amount to build"
// @Security BearerAuth
// @Success 200
// @Failure 400
// @Failure 401
// @Router /planets/{id}/units/build [post]
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
		userPlanets, _ := h.planetRepo.GetByUserID(c.Context(), userID)
		planetIDs := make([]uint, len(userPlanets))
		for i, p := range userPlanets {
			planetIDs[i] = p.ID
		}
		return c.Status(403).JSON(fiber.Map{
			"error":        "planet does not belong to user",
			"your_planets": planetIDs,
		})
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
		"success":   true,
		"planet_id": planetID,
		"unit_id":   req.UnitID,
		"amount":    req.Amount,
	})
}

// GetUnitQueue returns the unit production queue for a planet
// @Summary Get unit queue
// @Description Get unit production queue for a planet
// @Tags Units
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/units/queue [get]
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
		ID        uint  `json:"id"`
		UnitID    int   `json:"unit_id"`
		Amount    int   `json:"amount"`
		StartTime int64 `json:"start_time"`
		EndTime   int64 `json:"end_time"`
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

// CancelUnit cancels a unit in the production queue
// @Summary Cancel unit
// @Description Cancel a unit in the production queue
// @Tags Units
// @Produce json
// @Param id path int true "Queue ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /units/queue/{id} [delete]
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

// GetAvailableUnits returns all units that can be built
// @Summary Get available units
// @Description Get all units (ships and defense) that can be built
// @Tags Units
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /units/available [get]
func (h *Handlers) GetAvailableUnits(c *fiber.Ctx) error {
	type UnitInfo struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Metal     int64  `json:"metal"`
		Crystal   int64  `json:"crystal"`
		Deuterium int64  `json:"deuterium"`
		BuildTime int64  `json:"build_time_seconds"`
		Category  string `json:"category"`
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

// GetPlanetShips returns all ships on a planet
// @Summary Get planet ships
// @Description Get all ships on a planet
// @Tags Units
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/units [get]
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
			"small_cargo":     planet.SmallCargo,
			"large_cargo":     planet.LargeCargo,
			"light_fighter":   planet.LightFighter,
			"heavy_fighter":   planet.HeavyFighter,
			"cruiser":         planet.Cruiser,
			"battleship":      planet.Battleship,
			"colony_ship":     planet.ColonyShip,
			"recycler":        planet.Recycler,
			"espionage_probe": planet.EspionageProbe,
			"bomber":          planet.Bomber,
			"destroyer":       planet.Destroyer,
			"deathstar":       planet.Deathstar,
			"battlecruiser":   planet.Battlecruiser,
			"reaper":          planet.Reaper,
			"pathfinder":      planet.Pathfinder,
			"solar_satellite": planet.SolarSatellite,
			"crawler":         planet.Crawler,
		},
	})
}

// GetPlanetDefense returns all defense units on a planet
// @Summary Get planet defense
// @Description Get all defense units on a planet
// @Tags Defense
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/defense [get]
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
			"rocket_launcher":     planet.RocketLauncher,
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

// GetPlanetProduction returns production rates for a planet
// @Summary Get planet production
// @Description Get production rates for a planet
// @Tags Planets
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/production [get]
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
			"metal":     metalProd,
			"crystal":   crystalProd,
			"deuterium": deuteriumProd,
			"energy":    energyProd,
		},
		"consumption": fiber.Map{
			"energy": planet.EnergyUsed,
		},
	})
}

// GetMessages returns all messages for the current user
// @Summary Get messages
// @Description Get all messages for the current user
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /messages [get]
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

// GetUnreadCount returns the count of unread messages
// @Summary Get unread count
// @Description Get count of unread messages
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /messages/unread [get]
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

// MarkMessageRead marks a message as read
// @Summary Mark message read
// @Description Mark a message as read
// @Tags Messages
// @Produce json
// @Param id path int true "Message ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /messages/{id}/read [post]
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

// DeleteMessage deletes a message
// @Summary Delete message
// @Description Delete a message
// @Tags Messages
// @Produce json
// @Param id path int true "Message ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /messages/{id} [delete]
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

// GetGalaxy returns galaxy information for a system
// @Summary Get galaxy
// @Description Get galaxy information for a system
// @Tags Galaxy
// @Produce json
// @Param galaxy path int true "Galaxy"
// @Param system path int true "System"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /galaxy/{galaxy}/{system} [get]
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
		"galaxy":    galaxy,
		"system":    system,
		"positions": positions,
	})
}

// GetHighscore returns highscore rankings
// @Summary Get highscore
// @Description Get highscore rankings
// @Tags Highscore
// @Produce json
// @Param category path int true "Category (0=Total, 1=Economy, 2=Research, 3=Military, 4=Military Built, 5=Military Destroyed, 6=Military Lost, 7=Honor)"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /highscore/{category} [get]
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

// GetNotes returns all notes for the current user
// @Summary Get notes
// @Description Get all notes for the current user
// @Tags Notes
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /notes [get]
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

// CreateNote creates a new note
// @Summary Create note
// @Description Create a new note
// @Tags Notes
// @Accept json
// @Produce json
// @Param body body CreateNoteRequest true "Note subject and text"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /notes [post]
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

// UpdateNote updates an existing note
// @Summary Update note
// @Description Update an existing note
// @Tags Notes
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param body body UpdateNoteRequest true "Updated note subject and text"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /notes/{id} [put]
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

// DeleteNote deletes a note
// @Summary Delete note
// @Description Delete a note
// @Tags Notes
// @Produce json
// @Param id path int true "Note ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /notes/{id} [delete]
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

// GetAlliances returns list of alliances
// @Summary Get alliances
// @Description Get list of alliances
// @Tags Alliances
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances [get]
func (h *Handlers) GetAlliances(c *fiber.Ctx) error {
	alliances, err := h.allianceService.GetAllAlliances(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"alliances": alliances})
}

// CreateAlliance creates a new alliance
// @Summary Create alliance
// @Description Create a new alliance
// @Tags Alliances
// @Accept json
// @Produce json
// @Param body body CreateAllianceRequest true "Alliance details"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances [post]
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

// GetAlliance returns information about a specific alliance
// @Summary Get alliance
// @Description Get information about a specific alliance
// @Tags Alliances
// @Produce json
// @Param id path int true "Alliance ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/{id} [get]
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

// GetAllianceMembers returns members of an alliance
// @Summary Get alliance members
// @Description Get members of an alliance
// @Tags Alliances
// @Produce json
// @Param id path int true "Alliance ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/{id}/members [get]
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

// ApplyToAlliance applies to join an alliance
// @Summary Apply to alliance
// @Description Apply to join an alliance
// @Tags Alliances
// @Accept json
// @Produce json
// @Param id path int true "Alliance ID"
// @Param body body ApplyToAllianceRequest true "Application message"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/{id}/apply [post]
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

// GetAllianceApplications returns applications to join an alliance
// @Summary Get alliance applications
// @Description Get applications to join an alliance
// @Tags Alliances
// @Produce json
// @Param id path int true "Alliance ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/{id}/applications [get]
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

// AcceptApplication accepts a player into the alliance
// @Summary Accept application
// @Description Accept a player into the alliance
// @Tags Alliances
// @Produce json
// @Param id path int true "Application ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/applications/{id}/accept [post]
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

// RejectApplication rejects a player from joining the alliance
// @Summary Reject application
// @Description Reject a player from joining the alliance
// @Tags Alliances
// @Produce json
// @Param id path int true "Application ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/applications/{id}/reject [post]
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

// LeaveAlliance leaves the current alliance
// @Summary Leave alliance
// @Description Leave the current alliance
// @Tags Alliances
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances/leave [post]
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

// UpdateAlliance updates alliance settings
// @Summary Update alliance
// @Description Update alliance settings
// @Tags Alliances
// @Accept json
// @Produce json
// @Param body body CreateAllianceRequest true "Updated alliance details"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /alliances [put]
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

// GetBuddies returns buddy list for the current user
// @Summary Get buddies
// @Description Get buddy list for the current user
// @Tags Buddy
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /buddy [get]
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

// SendBuddyRequest sends a buddy request
// @Summary Send buddy request
// @Description Send a buddy request to another player
// @Tags Buddy
// @Accept json
// @Produce json
// @Param body body SendBuddyRequest true "Player ID and message"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /buddy/request [post]
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

// GetPendingBuddyRequests returns pending buddy requests
// @Summary Get pending buddy requests
// @Description Get pending buddy requests
// @Tags Buddy
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /buddy/pending [get]
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

// AcceptBuddyRequest accepts a buddy request
// @Summary Accept buddy request
// @Description Accept a buddy request
// @Tags Buddy
// @Produce json
// @Param id path int true "Request ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /buddy/{id}/accept [post]
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

// RejectBuddyRequest rejects a buddy request
// @Summary Reject buddy request
// @Description Reject a buddy request
// @Tags Buddy
// @Produce json
// @Param id path int true "Request ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /buddy/{id}/reject [post]
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

// RemoveBuddy removes a buddy from the list
// @Summary Remove buddy
// @Description Remove a buddy from the list
// @Tags Buddy
// @Produce json
// @Param id path int true "Buddy ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /buddy/{id} [delete]
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

// GetJumpGateTargets returns available jump gate targets
// @Summary Get jump gate targets
// @Description Get available jump gate targets
// @Tags Planets
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/jump-gate/targets [get]
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

// ExecuteJumpGate uses the jump gate to transport fleets
// @Summary Execute jump gate
// @Description Execute jump gate to transport fleets to another planet
// @Tags Planets
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Param body body JumpGateRequest true "Target moon ID and ships to transport"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/jump-gate/execute [post]
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

// ScanWithPhalanx scans a galaxy position using phalanx
// @Summary Scan with phalanx
// @Description Scan a galaxy position using phalanx
// @Tags Planets
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Param body body PhalanxScanRequest true "Galaxy coordinates to scan"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /planets/{id}/phalanx/scan [post]
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

// GetEspionageReports returns espionage reports
// @Summary Get espionage reports
// @Description Get espionage reports for the current user
// @Tags Espionage
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /espionage [get]
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

// GetEspionageUnreadCount returns count of unread espionage reports
// @Summary Get espionage unread count
// @Description Get count of unread espionage reports
// @Tags Espionage
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /espionage/unread [get]
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

// MarkEspionageReportRead marks an espionage report as read
// @Summary Mark espionage report read
// @Description Mark an espionage report as read
// @Tags Espionage
// @Produce json
// @Param id path int true "Report ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /espionage/{id}/read [post]
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

// DeleteEspionageReport deletes an espionage report
// @Summary Delete espionage report
// @Description Delete an espionage report
// @Tags Espionage
// @Produce json
// @Param id path int true "Report ID"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /espionage/{id} [delete]
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

// GetVacationStatus returns the vacation mode status for the current user
// @Summary Get vacation status
// @Description Get vacation mode status for current user
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /user/vacation [get]
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

// EnableVacationMode enables vacation mode for the current user
// @Summary Enable vacation mode
// @Description Enable vacation mode - protects planets while player is away
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /user/vacation/enable [post]
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

// DisableVacationMode disables vacation mode for the current user
// @Summary Disable vacation mode
// @Description Disable vacation mode
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /user/vacation/disable [post]
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

// GetDebrisFields returns all debris fields in the universe
// @Summary Get debris fields
// @Description Get all debris fields in the universe
// @Tags Debris
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /debris [get]
func (h *Handlers) GetDebrisFields(c *fiber.Ctx) error {
	debris, err := h.debrisService.GetDebrisFields(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": debris})
}

// GetDebrisField returns debris field at specific coordinates
// @Summary Get debris field
// @Description Get debris field at specific coordinates
// @Tags Debris
// @Produce json
// @Param galaxy path int true "Galaxy"
// @Param system path int true "System"
// @Param position path int true "Position"
// @Security BearerAuth
// @Success 200
// @Router /debris/{galaxy}/{system}/{position} [get]
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

// CollectDebris sends a recycler fleet to collect debris
// @Summary Collect debris
// @Description Send a recycler fleet to collect debris at specific coordinates
// @Tags Debris
// @Accept json
// @Produce json
// @Param galaxy path int true "Galaxy"
// @Param system path int true "System"
// @Param position path int true "Position"
// @Security BearerAuth
// @Success 200
// @Router /debris/{galaxy}/{system}/{position}/collect [post]
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
		"success":  true,
		"metal":    metal,
		"crystal":  crystal,
		"recycled": metal + crystal,
	})
}

// GetWreckFields returns all wreck fields in the universe
// @Summary Get wreck fields
// @Description Get all wreck fields in the universe
// @Tags Wrecks
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /wrecks [get]
func (h *Handlers) GetWreckFields(c *fiber.Ctx) error {
	wrecks, err := h.debrisService.GetWreckFields(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": wrecks})
}

// GetWreckField returns wreck field at specific coordinates
// @Summary Get wreck field
// @Description Get wreck field at specific coordinates
// @Tags Wrecks
// @Produce json
// @Param galaxy path int true "Galaxy"
// @Param system path int true "System"
// @Param position path int true "Position"
// @Security BearerAuth
// @Success 200
// @Router /wrecks/{galaxy}/{system}/{position} [get]
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

// CollectWreckField sends a fleet to collect wreck field
// @Summary Collect wreck field
// @Description Send a fleet to collect wreck field at specific coordinates
// @Tags Wrecks
// @Accept json
// @Produce json
// @Param galaxy path int true "Galaxy"
// @Param system path int true "System"
// @Param position path int true "Position"
// @Security BearerAuth
// @Success 200
// @Router /wrecks/{galaxy}/{system}/{position}/collect [post]
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
		"success":   true,
		"metal":     metal,
		"crystal":   crystal,
		"deuterium": deuterium,
		"recycled":  metal + crystal + deuterium,
	})
}

// GetNPCPlanets returns all NPC planets
// @Summary Get NPC planets
// @Description Get all NPC planets in the universe
// @Tags NPC
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /npc/planets [get]
func (h *Handlers) GetNPCPlanets(c *fiber.Ctx) error {
	planets, err := h.npcService.GetNPCPlanets(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": planets})
}

// GetNPCPlanet returns a specific NPC planet
// @Summary Get NPC planet
// @Description Get a specific NPC planet
// @Tags NPC
// @Produce json
// @Param galaxy path int true "Galaxy"
// @Param system path int true "System"
// @Param position path int true "Position"
// @Security BearerAuth
// @Success 200
// @Router /npc/planets/{galaxy}/{system}/{position} [get]
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

// CreateNPCPlanet creates a new NPC planet (admin only)
// @Summary Create NPC planet
// @Description Create a new NPC planet
// @Tags NPC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /npc/planets [post]
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

// GenerateExpeditionFleet generates a random expedition fleet
// @Summary Generate expedition fleet
// @Description Generate a random expedition fleet for testing
// @Tags NPC
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /npc/fleets/expedition [post]
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

// GeneratePirateRaid generates a pirate raid fleet
// @Summary Generate pirate raid
// @Description Generate a pirate raid fleet for testing
// @Tags NPC
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /npc/fleets/pirate [post]
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

// CreateACS creates an Attack Coordinate System (ACS) alliance
// @Summary Create ACS
// @Description Create an Attack Coordinate System (ACS) alliance for coordinated attacks
// @Tags ACS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /acs/create [post]
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

// JoinACS joins an existing ACS
// @Summary Join ACS
// @Description Join an existing Attack Coordinate System (ACS)
// @Tags ACS
// @Accept json
// @Produce json
// @Param id path int true "ACS ID"
// @Security BearerAuth
// @Success 200
// @Router /acs/{id}/join [post]
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

// GetACS returns information about an ACS
// @Summary Get ACS
// @Description Get information about an Attack Coordinate System (ACS)
// @Tags ACS
// @Produce json
// @Param id path int true "ACS ID"
// @Security BearerAuth
// @Success 200
// @Router /acs/{id} [get]
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

// GetACSFleets returns fleets in an ACS
// @Summary Get ACS fleets
// @Description Get all fleets in an Attack Coordinate System (ACS)
// @Tags ACS
// @Produce json
// @Param id path int true "ACS ID"
// @Security BearerAuth
// @Success 200
// @Router /acs/{id}/fleets [get]
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

// GetPremiumStatus returns premium status for the current user
// @Summary Get premium status
// @Description Get premium status for the current user
// @Tags Premium
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /premium [get]
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

// ActivatePremium activates premium features
// @Summary Activate premium
// @Description Activate premium features using dark matter
// @Tags Premium
// @Accept x-www-form-urlencoded
// @Produce json
// @Param days formData int true "Number of days to activate"
// @Security BearerAuth
// @Success 200
// @Router /premium/activate [post]
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

// MerchantBuy buys resources from the merchant
// @Summary Merchant buy
// @Description Buy resources from the merchant using dark matter
// @Tags Merchant
// @Accept x-www-form-urlencoded
// @Produce json
// @Param planet_id formData int true "Planet ID"
// @Param resource_type formData string true "Resource type to buy (metal, crystal, deuterium)"
// @Param amount formData int true "Amount to buy"
// @Security BearerAuth
// @Success 200
// @Router /merchant/buy [post]
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

// MerchantSell sells resources to the merchant
// @Summary Merchant sell
// @Description Sell resources to the merchant for dark matter
// @Tags Merchant
// @Accept x-www-form-urlencoded
// @Produce json
// @Param planet_id formData int true "Planet ID"
// @Param resource_type formData string true "Resource type to sell (metal, crystal, deuterium)"
// @Param amount formData int true "Amount to sell"
// @Security BearerAuth
// @Success 200
// @Router /merchant/sell [post]
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

// MovePlanet moves a planet to new coordinates
// @Summary Move planet
// @Description Move a planet to new coordinates (requires Dark Matter)
// @Tags Planets
// @Accept json
// @Produce json
// @Param id path int true "Planet ID"
// @Security BearerAuth
// @Success 200
// @Router /planets/{id}/move [post]
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

// GetCharacterClass returns the current user's character class
// @Summary Get character class
// @Description Get available character classes and current selection
// @Tags Character
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /character-class [get]
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
		"character_class":   classID,
		"class_info":        currentClass,
		"available_classes": classes,
	})
}

// SelectCharacterClass allows the user to select a character class
// @Summary Select character class
// @Description Select a character class (Collector, General, Engineer)
// @Tags Character
// @Accept x-www-form-urlencoded
// @Produce json
// @Param class_id formData int true "Class ID (1=Collector, 2=General, 3=Engineer)"
// @Security BearerAuth
// @Success 200
// @Failure 401
// @Router /character-class/select [post]
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

// ============ REFERENCE DATA ENDPOINTS ============

// GetBuildings returns all buildings with their costs and info
// @Summary Get all buildings
// @Description Get list of all buildings with IDs, names, base costs, and cost factors
// @Tags Reference
// @Produce json
// @Success 200
// @Router /buildings [get]
func (h *Handlers) GetBuildings(c *fiber.Ctx) error {
	type Building struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Metal      int     `json:"metal"`
		Crystal    int     `json:"crystal"`
		Deuterium  int     `json:"deuterium"`
		CostFactor float64 `json:"cost_factor"`
		Category   string  `json:"category"`
	}
	buildings := []Building{
		{1, "Metal Mine", 60, 0, 0, 1.5, "mine"},
		{2, "Crystal Mine", 48, 24, 0, 1.6, "mine"},
		{3, "Deuterium Synthesizer", 225, 0, 0, 1.5, "mine"},
		{4, "Solar Plant", 75, 0, 0, 1.5, "energy"},
		{5, "Fusion Reactor", 900, 360, 180, 1.8, "energy"},
		{6, "Metal Storage", 100, 0, 0, 2.0, "storage"},
		{7, "Crystal Storage", 100, 50, 0, 2.0, "storage"},
		{8, "Deuterium Storage", 100, 100, 0, 2.0, "storage"},
		{9, "Robot Factory", 400, 120, 0, 2.0, "facility"},
		{10, "Shipyard", 400, 200, 0, 2.0, "facility"},
		{11, "Research Lab", 200, 400, 0, 2.0, "facility"},
		{12, "Nanite Factory", 1000000, 200000, 0, 2.0, "facility"},
		{13, "Terraformer", 50000, 100000, 1000000, 2.0, "facility"},
		{14, "Space Dock", 400, 200, 100, 2.0, "facility"},
		{15, "Metal Silo", 100, 0, 0, 2.0, "storage"},
		{16, "Crystal Silo", 100, 50, 0, 2.0, "storage"},
		{17, "Deuterium Silo", 100, 100, 0, 2.0, "storage"},
		{18, "Lunar Base", 20000, 40000, 20000, 2.0, "moon"},
		{19, "Sensor Phalanx", 20000, 40000, 20000, 2.0, "moon"},
		{20, "Jump Gate", 2000000, 4000000, 2000000, 2.0, "moon"},
		{21, "Missile Silo", 20000, 20000, 1000, 2.0, "facility"},
	}
	return c.JSON(fiber.Map{"buildings": buildings})
}

// GetBuilding returns a single building by ID
// @Summary Get building by ID
// @Description Get detailed information about a specific building
// @Tags Reference
// @Produce json
// @Param id path int true "Building ID"
// @Success 200
// @Router /buildings/{id} [get]
func (h *Handlers) GetBuilding(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid building id"})
	}
	type Building struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Metal      int     `json:"metal"`
		Crystal    int     `json:"crystal"`
		Deuterium  int     `json:"deuterium"`
		CostFactor float64 `json:"cost_factor"`
		Category   string  `json:"category"`
	}
	buildings := map[int]Building{
		1:  {1, "Metal Mine", 60, 0, 0, 1.5, "mine"},
		2:  {2, "Crystal Mine", 48, 24, 0, 1.6, "mine"},
		3:  {3, "Deuterium Synthesizer", 225, 0, 0, 1.5, "mine"},
		4:  {4, "Solar Plant", 75, 0, 0, 1.5, "energy"},
		5:  {5, "Fusion Reactor", 900, 360, 180, 1.8, "energy"},
		6:  {6, "Metal Storage", 100, 0, 0, 2.0, "storage"},
		7:  {7, "Crystal Storage", 100, 50, 0, 2.0, "storage"},
		8:  {8, "Deuterium Storage", 100, 100, 0, 2.0, "storage"},
		9:  {9, "Robot Factory", 400, 120, 0, 2.0, "facility"},
		10: {10, "Shipyard", 400, 200, 0, 2.0, "facility"},
		11: {11, "Research Lab", 200, 400, 0, 2.0, "facility"},
		12: {12, "Nanite Factory", 1000000, 200000, 0, 2.0, "facility"},
		13: {13, "Terraformer", 50000, 100000, 1000000, 2.0, "facility"},
		14: {14, "Space Dock", 400, 200, 100, 2.0, "facility"},
		15: {15, "Metal Silo", 100, 0, 0, 2.0, "storage"},
		16: {16, "Crystal Silo", 100, 50, 0, 2.0, "storage"},
		17: {17, "Deuterium Silo", 100, 100, 0, 2.0, "storage"},
		18: {18, "Lunar Base", 20000, 40000, 20000, 2.0, "moon"},
		19: {19, "Sensor Phalanx", 20000, 40000, 20000, 2.0, "moon"},
		20: {20, "Jump Gate", 2000000, 4000000, 2000000, 2.0, "moon"},
		21: {21, "Missile Silo", 20000, 20000, 1000, 2.0, "facility"},
	}
	building, ok := buildings[id]
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "building not found"})
	}
	return c.JSON(building)
}

// GetShips returns all ships with full stats
// @Summary Get all ships
// @Description Get list of all ships with IDs, names, costs, speed, cargo, armor, weapons, shields
// @Tags Reference
// @Produce json
// @Success 200
// @Router /ships [get]
func (h *Handlers) GetShips(c *fiber.Ctx) error {
	type Ship struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Metal         int64  `json:"metal"`
		Crystal       int64  `json:"crystal"`
		Deuterium     int64  `json:"deuterium"`
		BuildTime     int    `json:"build_time_seconds"`
		CargoCapacity int64  `json:"cargo_capacity"`
		BaseSpeed     int64  `json:"base_speed"`
		StructuralInt int64  `json:"structural_integrity"`
		Shield        int64  `json:"shield"`
		Weapon        int64  `json:"weapon"`
		Engine        string `json:"engine_type"`
	}
	ships := []Ship{
		{202, "Small Cargo", 2000, 2000, 0, 5, 5000, 5000, 2000, 10, 5, "combustion"},
		{203, "Large Cargo", 6000, 6000, 0, 8, 25000, 7500, 6000, 25, 12, "combustion"},
		{204, "Light Fighter", 10000, 6000, 2000, 20, 50, 12500, 4000, 10, 50, "combustion"},
		{205, "Heavy Fighter", 25000, 15000, 5000, 40, 100, 10000, 10000, 25, 150, "impulse"},
		{206, "Cruiser", 10000, 20000, 10000, 10, 800, 15000, 27000, 50, 400, "impulse"},
		{207, "Battleship", 50000, 25000, 15000, 80, 1500, 10000, 60000, 200, 1000, "hyperspace"},
		{208, "Colony Ship", 10000, 10000, 0, 50, 7500, 2500, 30000, 100, 150, "impulse"},
		{209, "Recycler", 10000, 6000, 2000, 15, 20000, 6000, 16000, 10, 100, "hyperspace"},
		{210, "Espionage Probe", 0, 1000, 0, 30, 5, 100000000, 1000, 1, 0, "combustion"},
		{211, "Bomber", 50000, 50000, 25000, 200, 500, 4000, 75000, 500, 1000, "impulse"},
		{212, "Solar Satellite", 0, 2000, 500, 3, 0, 0, 2000, 1, 1, "none"},
		{213, "Destroyer", 10000, 10000, 0, 30, 2000, 5000, 110000, 500, 2000, "hyperspace"},
		{214, "Deathstar", 100000, 100000, 50000, 400, 1000000, 100, 9000000, 50000, 200000, "hyperspace"},
		{215, "Battlecruiser", 3000, 1000, 0, 4, 750, 10000, 70000, 400, 700, "hyperspace"},
		{217, "Crawler", 2000, 2000, 1000, 10, 0, 4000, 4000, 2, 8, "combustion"},
		{218, "Reaper", 8000, 0, 0, 20, 700, 10000, 140000, 700, 2800, "hyperspace"},
		{219, "Pathfinder", 20000, 10000, 10000, 75, 500, 12000, 23000, 100, 200, "hyperspace"},
	}
	return c.JSON(fiber.Map{"ships": ships})
}

// GetShip returns a single ship by ID
// @Summary Get ship by ID
// @Description Get detailed information about a specific ship
// @Tags Reference
// @Produce json
// @Param id path int true "Ship ID"
// @Success 200
// @Router /ships/{id} [get]
func (h *Handlers) GetShip(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ship id"})
	}
	type Ship struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Metal         int64  `json:"metal"`
		Crystal       int64  `json:"crystal"`
		Deuterium     int64  `json:"deuterium"`
		BuildTime     int    `json:"build_time_seconds"`
		CargoCapacity int64  `json:"cargo_capacity"`
		BaseSpeed     int64  `json:"base_speed"`
		StructuralInt int64  `json:"structural_integrity"`
		Shield        int64  `json:"shield"`
		Weapon        int64  `json:"weapon"`
		Engine        string `json:"engine_type"`
	}
	ships := map[int]Ship{
		202: {202, "Small Cargo", 2000, 2000, 0, 5, 5000, 5000, 2000, 10, 5, "combustion"},
		203: {203, "Large Cargo", 6000, 6000, 0, 8, 25000, 7500, 6000, 25, 12, "combustion"},
		204: {204, "Light Fighter", 10000, 6000, 2000, 20, 50, 12500, 4000, 10, 50, "combustion"},
		205: {205, "Heavy Fighter", 25000, 15000, 5000, 40, 100, 10000, 10000, 25, 150, "impulse"},
		206: {206, "Cruiser", 10000, 20000, 10000, 10, 800, 15000, 27000, 50, 400, "impulse"},
		207: {207, "Battleship", 50000, 25000, 15000, 80, 1500, 10000, 60000, 200, 1000, "hyperspace"},
		208: {208, "Colony Ship", 10000, 10000, 0, 50, 7500, 2500, 30000, 100, 150, "impulse"},
		209: {209, "Recycler", 10000, 6000, 2000, 15, 20000, 6000, 16000, 10, 100, "hyperspace"},
		210: {210, "Espionage Probe", 0, 1000, 0, 30, 5, 100000000, 1000, 1, 0, "combustion"},
		211: {211, "Bomber", 50000, 50000, 25000, 200, 500, 4000, 75000, 500, 1000, "impulse"},
		212: {212, "Solar Satellite", 0, 2000, 500, 3, 0, 0, 2000, 1, 1, "none"},
		213: {213, "Destroyer", 10000, 10000, 0, 30, 2000, 5000, 110000, 500, 2000, "hyperspace"},
		214: {214, "Deathstar", 100000, 100000, 50000, 400, 1000000, 100, 9000000, 50000, 200000, "hyperspace"},
		215: {215, "Battlecruiser", 3000, 1000, 0, 4, 750, 10000, 70000, 400, 700, "hyperspace"},
		217: {217, "Crawler", 2000, 2000, 1000, 10, 0, 4000, 4000, 2, 8, "combustion"},
		218: {218, "Reaper", 8000, 0, 0, 20, 700, 140000, 10000, 700, 2800, "hyperspace"},
		219: {219, "Pathfinder", 20000, 10000, 10000, 75, 500, 12000, 23000, 100, 200, "hyperspace"},
	}
	ship, ok := ships[id]
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "ship not found"})
	}
	return c.JSON(ship)
}

// GetDefense returns all defense units with full stats
// @Summary Get all defense
// @Description Get list of all defense units with IDs, names, costs, armor, weapons, shields
// @Tags Reference
// @Produce json
// @Success 200
// @Router /defense [get]
func (h *Handlers) GetDefense(c *fiber.Ctx) error {
	type Defense struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Metal         int64  `json:"metal"`
		Crystal       int64  `json:"crystal"`
		Deuterium     int64  `json:"deuterium"`
		BuildTime     int    `json:"build_time_seconds"`
		StructuralInt int64  `json:"structural_integrity"`
		Shield        int64  `json:"shield"`
		Weapon        int64  `json:"weapon"`
	}
	defense := []Defense{
		{401, "Rocket Launcher", 2000, 0, 0, 10, 2000, 20, 80},
		{402, "Light Laser", 1500, 500, 0, 11, 2000, 25, 100},
		{403, "Heavy Laser", 6000, 2000, 0, 22, 8000, 100, 250},
		{404, "Ion Cannon", 2000, 6000, 0, 16, 8000, 500, 150},
		{405, "Gauss Cannon", 20000, 15000, 2000, 45, 35000, 200, 1100},
		{406, "Plasma Turret", 50000, 50000, 30000, 90, 100000, 300, 3000},
		{407, "Small Shield Dome", 10000, 10000, 0, 20, 20000, 2000, 1},
		{408, "Large Shield Dome", 50000, 50000, 0, 60, 100000, 10000, 1},
		{409, "Anti-Ballistic Missile", 8000, 0, 0, 1, 8000, 1, 1},
		{410, "Interplanetary Missile", 15000, 0, 0, 1, 15000, 1, 12000},
	}
	return c.JSON(fiber.Map{"defense": defense})
}

// GetDefenseUnit returns a single defense by ID
// @Summary Get defense by ID
// @Description Get detailed information about a specific defense unit
// @Tags Reference
// @Produce json
// @Param id path int true "Defense ID"
// @Success 200
// @Router /defense/{id} [get]
func (h *Handlers) GetDefenseUnit(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid defense id"})
	}
	type Defense struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Metal         int64  `json:"metal"`
		Crystal       int64  `json:"crystal"`
		Deuterium     int64  `json:"deuterium"`
		BuildTime     int    `json:"build_time_seconds"`
		StructuralInt int64  `json:"structural_integrity"`
		Shield        int64  `json:"shield"`
		Weapon        int64  `json:"weapon"`
	}
	defense := map[int]Defense{
		401: {401, "Rocket Launcher", 2000, 0, 0, 10, 2000, 20, 80},
		402: {402, "Light Laser", 1500, 500, 0, 11, 2000, 25, 100},
		403: {403, "Heavy Laser", 6000, 2000, 0, 22, 8000, 100, 250},
		404: {404, "Ion Cannon", 2000, 6000, 0, 16, 8000, 500, 150},
		405: {405, "Gauss Cannon", 20000, 15000, 2000, 45, 35000, 200, 1100},
		406: {406, "Plasma Turret", 50000, 50000, 30000, 90, 100000, 300, 3000},
		407: {407, "Small Shield Dome", 10000, 10000, 0, 20, 20000, 2000, 1},
		408: {408, "Large Shield Dome", 50000, 50000, 0, 60, 100000, 10000, 1},
		409: {409, "Anti-Ballistic Missile", 8000, 0, 0, 1, 8000, 1, 1},
		410: {410, "Interplanetary Missile", 15000, 0, 0, 1, 15000, 1, 12000},
	}
	def, ok := defense[id]
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "defense not found"})
	}
	return c.JSON(def)
}

// GetResearch returns all research types with costs
// @Summary Get all research
// @Description Get list of all research types with IDs, names, and base costs
// @Tags Reference
// @Produce json
// @Success 200
// @Router /research [get]
func (h *Handlers) GetResearch(c *fiber.Ctx) error {
	type Research struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Metal     int64  `json:"metal"`
		Crystal   int64  `json:"crystal"`
		Deuterium int64  `json:"deuterium"`
	}
	research := []Research{
		{113, "Energy Technology", 0, 800, 400},
		{120, "Laser Technology", 200, 600, 0},
		{121, "Ion Technology", 1000, 300, 0},
		{114, "Hyperspace Technology", 4000, 2000, 1000},
		{122, "Plasma Technology", 2400, 1200, 600},
		{115, "Combustion Drive", 4000, 2000, 600},
		{117, "Impulse Drive", 4000, 2000, 600},
		{118, "Hyperspace Drive", 10000, 6000, 4000},
		{106, "Espionage Technology", 200, 1000, 200},
		{108, "Computer Technology", 100, 400, 200},
		{124, "Astrophysics", 8000, 4000, 2000},
		{123, "Graviton Technology", 100000, 50000, 50000},
		{109, "Weapons Technology", 800, 200, 0},
		{110, "Shielding Technology", 400, 600, 0},
		{111, "Armour Technology", 400, 200, 0},
	}
	return c.JSON(fiber.Map{"research": research})
}

// GetResearchType returns a single research by ID
// @Summary Get research by ID
// @Description Get detailed information about a specific research
// @Tags Reference
// @Produce json
// @Param id path int true "Research ID"
// @Success 200
// @Router /research/{id} [get]
func (h *Handlers) GetResearchType(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid research id"})
	}
	type Research struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Metal     int64  `json:"metal"`
		Crystal   int64  `json:"crystal"`
		Deuterium int64  `json:"deuterium"`
	}
	research := map[int]Research{
		113: {113, "Energy Technology", 0, 800, 400},
		120: {120, "Laser Technology", 200, 600, 0},
		121: {121, "Ion Technology", 1000, 300, 0},
		114: {114, "Hyperspace Technology", 4000, 2000, 1000},
		122: {122, "Plasma Technology", 2400, 1200, 600},
		115: {115, "Combustion Drive", 4000, 2000, 600},
		117: {117, "Impulse Drive", 4000, 2000, 600},
		118: {118, "Hyperspace Drive", 10000, 6000, 4000},
		106: {106, "Espionage Technology", 200, 1000, 200},
		108: {108, "Computer Technology", 100, 400, 200},
		124: {124, "Astrophysics", 8000, 4000, 2000},
		123: {123, "Graviton Technology", 100000, 50000, 50000},
		109: {109, "Weapons Technology", 800, 200, 0},
		110: {110, "Shielding Technology", 400, 600, 0},
		111: {111, "Armour Technology", 400, 200, 0},
	}
	res, ok := research[id]
	if !ok {
		return c.Status(404).JSON(fiber.Map{"error": "research not found"})
	}
	return c.JSON(res)
}

// GetMissions returns all mission types
// @Summary Get all missions
// @Description Get list of all fleet mission types
// @Tags Reference
// @Produce json
// @Success 200
// @Router /missions [get]
func (h *Handlers) GetMissions(c *fiber.Ctx) error {
	type Mission struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	missions := []Mission{
		{1, "Attack", "Attack enemy fleets and planets"},
		{2, "Transport", "Transport resources between planets"},
		{3, "Colonize", "Colonize a new planet"},
		{4, "Recycle", "Harvest debris fields"},
		{5, "Expedition", "Explore for resources and artifacts"},
		{6, "ACS", "Attack Coordinate System - combine with allies"},
		{7, "Hold", "Hold position at target"},
		{8, "Deploy", "Deploy resources to another planet"},
	}
	return c.JSON(fiber.Map{"missions": missions})
}

// GetGame returns game meta information
// @Summary Get game info
// @Description Get highscore categories and other game meta information
// @Tags Reference
// @Produce json
// @Success 200
// @Router /game [get]
func (h *Handlers) GetGame(c *fiber.Ctx) error {
	type HighscoreCategory struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	type CharacterClass struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	return c.JSON(fiber.Map{
		"highscore_categories": []HighscoreCategory{
			{0, "Total Points"},
			{1, "Economy"},
			{2, "Research"},
			{3, "Military"},
			{4, "Military Built"},
			{5, "Military Destroyed"},
			{6, "Military Lost"},
			{7, "Honor"},
		},
		"character_classes": []CharacterClass{
			{0, "None", "No special class"},
			{1, "Collector", "Production bonuses"},
			{2, "General", "Military bonuses"},
			{3, "Discoverer", "Expedition bonuses"},
		},
	})
}

// End of reference data handlers
