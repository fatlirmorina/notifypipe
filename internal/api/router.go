package api

import (
	"github.com/fatlirmorina/notifypipe/internal/config"
	"github.com/fatlirmorina/notifypipe/internal/database"
	"github.com/fatlirmorina/notifypipe/internal/docker"
	"github.com/fatlirmorina/notifypipe/internal/notifications"
	"github.com/gofiber/fiber/v2"
)

// Router handles API routing
type Router struct {
	app      *fiber.App
	db       *database.Database
	docker   *docker.Client
	notifier *notifications.Manager
	config   *config.Config
}

// NewRouter creates a new API router
func NewRouter(
	app *fiber.App,
	db *database.Database,
	dockerClient *docker.Client,
	notifier *notifications.Manager,
	cfg *config.Config,
) *Router {
	return &Router{
		app:      app,
		db:       db,
		docker:   dockerClient,
		notifier: notifier,
		config:   cfg,
	}
}

// Setup sets up all API routes
func (r *Router) Setup() {
	api := r.app.Group("/api")

	// Health check
	api.Get("/health", r.healthCheck)

	// Setup
	api.Get("/setup/status", r.getSetupStatus)
	api.Post("/setup/complete", r.completeSetup)

	// Containers
	api.Get("/containers", r.listContainers)
	api.Get("/containers/:id", r.getContainer)
	api.Put("/containers/:id", r.updateContainer)

	// Notifications
	api.Get("/notifications", r.listNotifications)
	api.Post("/notifications", r.createNotification)
	api.Put("/notifications/:id", r.updateNotification)
	api.Delete("/notifications/:id", r.deleteNotification)
	api.Post("/notifications/test", r.testNotification)

	// Events
	api.Get("/events", r.listEvents)
	api.Get("/events/:containerId", r.getContainerEvents)

	// Statistics
	api.Get("/stats", r.getStats)
}
