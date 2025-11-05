package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pocketbase/dbx"
)

// listEvents returns recent events
func (r *Router) listEvents(c *fiber.Ctx) error {
	records, err := r.db.App().Dao().FindRecordsByFilter(
		"events_log",
		"",
		"-timestamp", // Sort by timestamp descending
		100,          // Limit
		0,            // Offset
	)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []fiber.Map
	for _, record := range records {
		result = append(result, fiber.Map{
			"id":             record.Id,
			"container_id":   record.GetString("container_id"),
			"container_name": record.GetString("container_name"),
			"event_type":     record.GetString("event_type"),
			"status":         record.GetString("status"),
			"message":        record.GetString("message"),
			"timestamp":      record.GetDateTime("timestamp"),
		})
	}

	return c.JSON(result)
}

// getContainerEvents returns events for a specific container
func (r *Router) getContainerEvents(c *fiber.Ctx) error {
	containerID := c.Params("containerId")

	records, err := r.db.App().Dao().FindRecordsByFilter(
		"events_log",
		"container_id = {:containerId}",
		"-timestamp",
		100,
		0,
		dbx.Params{"containerId": containerID},
	)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []fiber.Map
	for _, record := range records {
		result = append(result, fiber.Map{
			"id":             record.Id,
			"container_id":   record.GetString("container_id"),
			"container_name": record.GetString("container_name"),
			"event_type":     record.GetString("event_type"),
			"status":         record.GetString("status"),
			"message":        record.GetString("message"),
			"timestamp":      record.GetDateTime("timestamp"),
		})
	}

	return c.JSON(result)
}
