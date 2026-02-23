package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"ogamex-go/internal/repository"
	"ogamex-go/internal/schema"
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

	protected := api.Group("", h.authMiddleware)

	protected.Get("/user", h.GetUser)
	protected.Get("/planets", h.GetUserPlanets)
	protected.Get("/planets/:id", h.GetPlanet)
	protected.Get("/planets/:id/resources", h.GetPlanetResources)
	protected.Get("/planets/:id/buildings", h.GetPlanetBuildings)
	protected.Post("/planets/:id/buildings/:building_id", h.StartBuilding)
	protected.Get("/planets/:id/queue", h.GetBuildingQueue)
	protected.Delete("/queue/:id", h.CancelBuilding)

	protected.Post("/research/start", h.StartResearch)
	protected.Get("/research/queue", h.GetResearchQueue)

	protected.Post("/fleets/send", h.SendFleet)
	protected.Get("/fleets", h.GetFleets)
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
		"online":   true,
		"version": "0.0.1",
		"players":  0,
		"universe": "ogamex-go",
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
