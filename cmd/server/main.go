package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"ogamex-go/internal/api"
	"ogamex-go/internal/api/middleware"
	"ogamex-go/internal/database"
	appLogger "ogamex-go/internal/logger"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/scheduler"
	"ogamex-go/internal/service"
)

func main() {
	production := viper.GetString("ENV") == "production"
	appLogger.Init(production)

	if err := godotenv.Load(); err != nil {
		appLogger.Warn().Msg("No .env file found")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		appLogger.Warn().Err(err).Msg("No config file found, using defaults")
	}

	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.name")

	db, err := database.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		appLogger.Error().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	if err := db.AutoMigrate(); err != nil {
		appLogger.Error().Err(err).Msg("Failed to migrate database")
		os.Exit(1)
	}
	appLogger.Info().Msg("Database connected and migrated")

	userRepo := repository.NewUserRepository(db)
	planetRepo := repository.NewPlanetRepository(db)
	techRepo := repository.NewUserTechRepository(db)
	buildingQueueRepo := repository.NewBuildingQueueRepository(db)
	researchQueueRepo := repository.NewResearchQueueRepository(db)
	unitQueueRepo := repository.NewUnitQueueRepository(db)
	fleetMissionRepo := repository.NewFleetMissionRepository(db)

	buildingService := service.NewBuildingService(planetRepo, buildingQueueRepo, techRepo)
	researchService := service.NewResearchService(userRepo, planetRepo, researchQueueRepo, techRepo)
	
	universeSpeed := viper.GetInt("universe.speed")
	if universeSpeed == 0 {
		universeSpeed = 1
	}
	fleetService := service.NewFleetService(planetRepo, fleetMissionRepo, userRepo, techRepo, universeSpeed)
	
	economySpeed := viper.GetInt("universe.economy_speed")
	if economySpeed == 0 {
		economySpeed = 1
	}
	productionService := service.NewProductionService(planetRepo, techRepo, economySpeed)
	authService := service.NewAuthService(userRepo, planetRepo)
	unitService := service.NewUnitService(planetRepo, unitQueueRepo, techRepo)
	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo, userRepo)
	planetService := service.NewPlanetService(planetRepo)

	sched := scheduler.NewScheduler(buildingService, researchService, fleetService, productionService, unitService, planetRepo, buildingQueueRepo, researchQueueRepo, unitQueueRepo)
	sched.Start()

	handlers := api.NewHandlers(buildingService, researchService, fleetService, planetRepo, authService, unitService, productionService, messageService, planetService)

	app := fiber.New(fiber.Config{
		AppName: "ogamex-go",
	})

	app.Use(recover.New())
	app.Use(middleware.TraceIDMiddleware())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${trace_id} | ${method} ${path} ${error}\n",
		CustomTags: map[string]logger.LogFunc{
			"trace_id": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				traceID := c.Get("X-Trace-ID")
				if traceID == "" {
					if id, ok := c.Locals("trace_id").(string); ok {
						traceID = id
					}
				}
				return output.WriteString(traceID)
			},
		},
	}))

	rateLimit := limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if c.Locals("user_id") != nil {
				return fmt.Sprintf("user:%d", c.Locals("user_id"))
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "rate limit exceeded",
				"retry_after": 60,
			})
		},
	})

	app.Use(rateLimit)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/health/ready", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"status": "unhealthy",
				"database": "disconnected",
			})
		}
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
		})
	})

	handlers.SetupRoutes(app)

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	go func() {
		addr := fmt.Sprintf(":%s", port)
		appLogger.Info().Str("address", addr).Msg("Starting ogamex-go server")
		if err := app.Listen(addr); err != nil {
			appLogger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info().Msg("Shutting down...")
	sched.Stop()

	if err := db.Close(); err != nil {
		appLogger.Error().Err(err).Msg("Error closing database")
	}

	appLogger.Info().Msg("Server stopped")
}
