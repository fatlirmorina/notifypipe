package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pocketbase/pocketbase/models"
)

// listNotifications returns all notification channels
func (r *Router) listNotifications(c *fiber.Ctx) error {
	records, err := r.db.App().Dao().FindRecordsByFilter(
		"notifications",
		"",
		"",
		0,
		0,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var result []fiber.Map
	for _, record := range records {
		result = append(result, fiber.Map{
			"id":      record.Id,
			"name":    record.GetString("name"),
			"type":    record.GetString("type"),
			"url":     record.GetString("url"),
			"enabled": record.GetBool("enabled"),
		})
	}

	return c.JSON(result)
}

// createNotification creates a new notification channel
func (r *Router) createNotification(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.Name == "" || body.Type == "" || body.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name, type, and URL are required"})
	}

	collection, err := r.db.App().Dao().FindCollectionByNameOrId("notifications")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	record := models.NewRecord(collection)
	record.Set("name", body.Name)
	record.Set("type", body.Type)
	record.Set("url", body.URL)
	record.Set("enabled", true)

	if err := r.db.App().Dao().SaveRecord(record); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification channel created",
		"id":      record.Id,
	})
}

// updateNotification updates a notification channel
func (r *Router) updateNotification(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		URL     string `json:"url"`
		Enabled bool   `json:"enabled"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	record, err := r.db.App().Dao().FindRecordById("notifications", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	if body.Name != "" {
		record.Set("name", body.Name)
	}
	if body.Type != "" {
		record.Set("type", body.Type)
	}
	if body.URL != "" {
		record.Set("url", body.URL)
	}
	record.Set("enabled", body.Enabled)

	if err := r.db.App().Dao().SaveRecord(record); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification channel updated",
	})
}

// deleteNotification deletes a notification channel
func (r *Router) deleteNotification(c *fiber.Ctx) error {
	id := c.Params("id")

	record, err := r.db.App().Dao().FindRecordById("notifications", id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	if err := r.db.App().Dao().DeleteRecord(record); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Notification channel deleted",
	})
}

// testNotification tests a notification channel
func (r *Router) testNotification(c *fiber.Ctx) error {
	var body struct {
		URL string `json:"url"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if body.URL == "" {
		return c.Status(400).JSON(fiber.Map{"error": "URL is required"})
	}

	if err := r.notifier.TestNotification(body.URL); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Test notification sent successfully",
	})
}
