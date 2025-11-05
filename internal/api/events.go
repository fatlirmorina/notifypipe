package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pocketbase/pocketbase/daos"
)

// listEvents returns recent events
func (r *Router) listEvents(c *fiber.Ctx) error {
	records, err := r.db.App().Dao().FindRecordsByExpr("events_log", nil, daos.SortFunc(func(sortData []daos.SortDataItem) {
		// Sort by timestamp descending
		sortData = append(sortData, daos.SortDataItem{
			Field:  "timestamp",
			IsSorted: true,
			IsDesc: true,
		})
	}))
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Limit to last 100 events
	limit := 100
	if len(records) > limit {
		records = records[:limit]
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

	records, err := r.db.App().Dao().FindRecordsByExpr("events_log", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []fiber.Map
	for _, record := range records {
		if record.GetString("container_id") == containerID {
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
	}

	return c.JSON(result)
}
