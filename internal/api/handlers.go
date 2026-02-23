package api

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/service"
)

type Handlers struct {
	buildingService *service.BuildingService
	researchService *service.ResearchService
	fleetService    *service.FleetService
	planetRepo      repository.PlanetRepository
	authService     *service.AuthService
}

func NewHandlers(
	buildingService *service.BuildingService,
	researchService *service.ResearchService,
	fleetService *service.FleetService,
	planetRepo repository.PlanetRepository,
	authService *service.AuthService,
) *Handlers {
	return &Handlers{
		buildingService: buildingService,
		researchService: researchService,
		fleetService:    fleetService,
		planetRepo:      planetRepo,
		authService:     authService,
	}
}

func (h *Handlers) SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/status", h.GetStatus)
	
	api.Post("/auth/register", h.Register)
	api.Post("/auth/login", h.Login)

	api.Get("/planets/:id", h.GetPlanet)
	api.Get("/planets/:id/resources", h.GetPlanetResources)
	api.Get("/planets/:id/buildings", h.GetPlanetBuildings)
	api.Post("/planets/:id/buildings/:building_id", h.StartBuilding)
	api.Get("/planets/:id/queue", h.GetBuildingQueue)
	api.Delete("/queue/:id", h.CancelBuilding)

	api.Post("/research/start", h.StartResearch)
	api.Get("/research/queue", h.GetResearchQueue)

	api.Post("/fleets/send", h.SendFleet)
	api.Get("/fleets", h.GetFleets)
}

func (h *Handlers) GetStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"online":   true,
		"version": "0.0.1",
		"players":  0,
		"universe": "ogamex-go",
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
	type ResearchRequest struct {
		UserID     uint `json:"user_id"`
		ResearchID int  `json:"research_id"`
	}

	var req ResearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	err := h.researchService.StartResearch(c.Context(), req.UserID, req.ResearchID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"user_id":      req.UserID,
		"research_id":  req.ResearchID,
	})
}

func (h *Handlers) GetResearchQueue(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id required"})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	queue, err := h.researchService.GetQueue(c.Context(), uint(userID))
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
	UserID         uint   `json:"user_id"`
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
	var req FleetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	origin, err := h.planetRepo.GetByCoords(c.Context(), req.UserID, req.OriginGalaxy, req.OriginSystem, req.OriginPosition)
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

	err = h.fleetService.SendFleet(c.Context(), req.UserID, params)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
	})
}

func (h *Handlers) GetFleets(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id required"})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	fleets, err := h.fleetService.GetActiveMissions(c.Context(), uint(userID))
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
