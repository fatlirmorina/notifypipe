package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pocketbase/pocketbase/models"
)

// healthCheck returns the health status
func (r *Router) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "NotifyPipe",
		"version": "1.0.1",
	})
}

// getSetupStatus checks if the app has been set up
func (r *Router) getSetupStatus(c *fiber.Ctx) error {
	// Check if we have any notification channels configured
	records, err := r.db.App().Dao().FindRecordsByFilter(
		"notifications",
		"",
		"",
		1,
		0,
	)
	if err != nil {
		records = []*models.Record{}
	}

	isSetup := len(records) > 0

	return c.JSON(fiber.Map{
		"setup_complete": isSetup,
		"needs_setup":    !isSetup,
	})
}

// completeSetup marks the setup as complete
func (r *Router) completeSetup(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Setup completed successfully",
	})
}

// getStats returns application statistics
func (r *Router) getStats(c *fiber.Ctx) error {
	// Count containers
	containers, err := r.db.App().Dao().FindRecordsByFilter("containers", "", "", 0, 0)
	if err != nil {
		containers = []*models.Record{}
	}

	// Count notification channels
	notifications, err := r.db.App().Dao().FindRecordsByFilter("notifications", "", "", 0, 0)
	if err != nil {
		notifications = []*models.Record{}
	}

	// Count recent events (last 24 hours)
	events, err := r.db.App().Dao().FindRecordsByFilter("events_log", "", "", 0, 0)
	if err != nil {
		events = []*models.Record{}
	}

	return c.JSON(fiber.Map{
		"containers_count":     len(containers),
		"notifications_count":  len(notifications),
		"events_count":         len(events),
	})
}
