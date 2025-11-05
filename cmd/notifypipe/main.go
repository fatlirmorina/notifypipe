package main

import (
	"log"
	"os"

	"github.com/fatlirmorina/notifypipe/internal/api"
	"github.com/fatlirmorina/notifypipe/internal/config"
	"github.com/fatlirmorina/notifypipe/internal/database"
	"github.com/fatlirmorina/notifypipe/internal/docker"
	"github.com/fatlirmorina/notifypipe/internal/notifications"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DataDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Docker client
	dockerClient, err := docker.NewClient(cfg.DockerSocket)
	if err != nil {
		log.Fatalf("Failed to initialize Docker client: %v", err)
	}
	defer dockerClient.Close()

	// Initialize notification manager
	notificationManager := notifications.NewManager(db)

	// Initialize event monitor
	eventMonitor := docker.NewEventMonitor(dockerClient, db, notificationManager)

	// Start monitoring Docker events in background
	go func() {
		if err := eventMonitor.Start(); err != nil {
			log.Printf("Error monitoring Docker events: %v", err)
		}
	}()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "NotifyPipe v1.0.2",
		ServerHeader: "NotifyPipe",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Serve static files
	app.Static("/", "./web/dist")

	// API routes
	apiRouter := api.NewRouter(app, db, dockerClient, notificationManager, cfg)
	apiRouter.Setup()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 NotifyPipe is running on http://localhost:%s", port)
	log.Printf("📊 Dashboard: http://localhost:%s", port)
	log.Printf("🔗 API: http://localhost:%s/api", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
