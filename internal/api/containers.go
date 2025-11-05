package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pocketbase/pocketbase/models"
)

// listContainers returns all containers
func (r *Router) listContainers(c *fiber.Ctx) error {
	// Get containers from Docker
	dockerContainers, err := r.docker.ListContainers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Get container settings from database
	dbRecords, err := r.db.App().Dao().FindRecordsByExpr("containers", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Create a map for quick lookup
	settingsMap := make(map[string]*models.Record)
	for _, record := range dbRecords {
		settingsMap[record.GetString("container_id")] = record
	}

	// Combine data
	var result []fiber.Map
	for _, container := range dockerContainers {
		settings := settingsMap[container.ID]
		
		name := ""
		if len(container.Names) > 0 {
			name = strings.TrimPrefix(container.Names[0], "/")
		}

		item := fiber.Map{
			"id":                 container.ID,
			"name":               name,
			"image":              container.Image,
			"state":              container.State,
			"status":             container.Status,
			"created":            container.Created,
			"notify_on_success":  false,
			"notify_on_failure":  true,
		}

		if settings != nil {
			item["notify_on_success"] = settings.GetBool("notify_on_success")
			item["notify_on_failure"] = settings.GetBool("notify_on_failure")
		}

		result = append(result, item)
	}

	return c.JSON(result)
}

// getContainer returns a specific container
func (r *Router) getContainer(c *fiber.Ctx) error {
	id := c.Params("id")

	containerInfo, err := r.docker.GetContainer(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Container not found"})
	}

	// Get settings from database
	records, err := r.db.App().Dao().FindRecordsByExpr("containers", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var settings *models.Record
	for _, record := range records {
		if record.GetString("container_id") == id {
			settings = record
			break
		}
	}

	result := fiber.Map{
		"id":                 containerInfo.ID,
		"name":               strings.TrimPrefix(containerInfo.Name, "/"),
		"image":              containerInfo.Config.Image,
		"state":              containerInfo.State.Status,
		"created":            containerInfo.Created,
		"notify_on_success":  false,
		"notify_on_failure":  true,
	}

	if settings != nil {
		result["notify_on_success"] = settings.GetBool("notify_on_success")
		result["notify_on_failure"] = settings.GetBool("notify_on_failure")
	}

	return c.JSON(result)
}

// updateContainer updates container notification settings
func (r *Router) updateContainer(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		NotifyOnSuccess bool `json:"notify_on_success"`
		NotifyOnFailure bool `json:"notify_on_failure"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Find or create container record
	records, err := r.db.App().Dao().FindRecordsByExpr("containers", nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var existingRecord *models.Record
	for _, record := range records {
		if record.GetString("container_id") == id {
			existingRecord = record
			break
		}
	}

	if existingRecord != nil {
		// Update existing
		existingRecord.Set("notify_on_success", body.NotifyOnSuccess)
		existingRecord.Set("notify_on_failure", body.NotifyOnFailure)
		if err := r.db.App().Dao().SaveRecord(existingRecord); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		// Create new
		containerInfo, err := r.docker.GetContainer(id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Container not found"})
		}

		collection, err := r.db.App().Dao().FindCollectionByNameOrId("containers")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		record := models.NewRecord(collection)
		record.Set("container_id", id)
		record.Set("name", strings.TrimPrefix(containerInfo.Name, "/"))
		record.Set("image", containerInfo.Config.Image)
		record.Set("notify_on_success", body.NotifyOnSuccess)
		record.Set("notify_on_failure", body.NotifyOnFailure)

		if err := r.db.App().Dao().SaveRecord(record); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Container settings updated",
	})
}
