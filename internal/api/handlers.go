package api

import (
	"github.com/gofiber/fiber/v2"
)

// healthCheck returns the health status
func (r *Router) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "NotifyPipe",
		"version": "1.0.0",
	})
}

// getSetupStatus checks if the app has been set up
func (r *Router) getSetupStatus(c *fiber.Ctx) error {
	// Check if we have any notification channels configured
	records, err := r.db.App().Dao().FindRecordsByExpr("notifications", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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
	containers, err := r.db.App().Dao().FindRecordsByExpr("containers", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Count notification channels
	notifications, err := r.db.App().Dao().FindRecordsByExpr("notifications", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Count recent events (last 24 hours)
	events, err := r.db.App().Dao().FindRecordsByExpr("events_log", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"containers_count":     len(containers),
		"notifications_count":  len(notifications),
		"events_count":         len(events),
	})
}
