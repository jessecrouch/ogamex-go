package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"ogamex-go/internal/api"
	"ogamex-go/internal/database"
	"ogamex-go/internal/repository"
	"ogamex-go/internal/scheduler"
	"ogamex-go/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: No config file found: %v", err)
	}

	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbPassword := viper.GetString("database.password")
	dbName := viper.GetString("database.name")

	db, err := database.NewDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database connected and migrated")

	userRepo := repository.NewUserRepository(db)
	planetRepo := repository.NewPlanetRepository(db)
	techRepo := repository.NewUserTechRepository(db)
	buildingQueueRepo := repository.NewBuildingQueueRepository(db)
	researchQueueRepo := repository.NewResearchQueueRepository(db)
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

	sched := scheduler.NewScheduler(buildingService, researchService, fleetService, productionService)
	sched.Start()

	handlers := api.NewHandlers(buildingService, researchService, fleetService, planetRepo, authService)

	app := fiber.New(fiber.Config{
		AppName: "ogamex-go",
	})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	handlers.SetupRoutes(app)

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting ogamex-go on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	sched.Stop()

	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server stopped")
}
