package docker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/fatlirmorina/notifypipe/internal/database"
	"github.com/fatlirmorina/notifypipe/internal/notifications"
	"github.com/pocketbase/pocketbase/models"
)

// EventMonitor monitors Docker events
type EventMonitor struct {
	client      *Client
	db          *database.Database
	notifier    *notifications.Manager
	ctx         context.Context
	cancelFunc  context.CancelFunc
}

// NewEventMonitor creates a new event monitor
func NewEventMonitor(client *Client, db *database.Database, notifier *notifications.Manager) *EventMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventMonitor{
		client:     client,
		db:         db,
		notifier:   notifier,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Start starts monitoring Docker events
func (em *EventMonitor) Start() error {
	log.Println("🔍 Starting Docker event monitoring...")

	eventsChan, errChan := em.client.cli.Events(em.ctx, types.EventsOptions{})

	for {
		select {
		case event := <-eventsChan:
			em.handleEvent(event)
		case err := <-errChan:
			if err != nil {
				log.Printf("Error receiving Docker event: %v", err)
			}
		case <-em.ctx.Done():
			log.Println("Stopping Docker event monitoring...")
			return nil
		}
	}
}

// handleEvent handles a Docker event
func (em *EventMonitor) handleEvent(event events.Message) {
	if event.Type != events.ContainerEventType {
		return
	}

	containerID := event.Actor.ID
	containerName := event.Actor.Attributes["name"]
	action := event.Action

	log.Printf("📦 Container event: %s - %s (%s)", containerName, action, containerID[:12])

	// Handle different event types
	switch action {
	case "start":
		em.handleContainerStart(containerID, containerName)
	case "die":
		em.handleContainerDie(containerID, containerName, event.Actor.Attributes["exitCode"])
	case "create":
		em.handleContainerCreate(containerID, containerName)
	}
}

// handleContainerStart handles container start events
func (em *EventMonitor) handleContainerStart(containerID, containerName string) {
	// Wait a bit to ensure container is actually running
	time.Sleep(2 * time.Second)

	containerInfo, err := em.client.GetContainer(containerID)
	if err != nil {
		log.Printf("Error getting container info: %v", err)
		return
	}

	if !containerInfo.State.Running {
		return
	}

	// Log event
	em.logEvent(containerID, containerName, "start", "success", "Container started successfully")

	// Check if we should notify
	if em.shouldNotify(containerID, "success") {
		message := fmt.Sprintf("✅ Container '%s' deployed successfully", containerName)
		em.notifier.Send(message)
	}
}

// handleContainerDie handles container die events
func (em *EventMonitor) handleContainerDie(containerID, containerName, exitCode string) {
	status := "failure"
	message := fmt.Sprintf("Container stopped with exit code %s", exitCode)

	// If exit code is 0, it's a graceful shutdown
	if exitCode == "0" {
		status = "stopped"
		message = "Container stopped gracefully"
	}

	// Log event
	em.logEvent(containerID, containerName, "die", status, message)

	// Only notify on failures (non-zero exit codes)
	if exitCode != "0" && em.shouldNotify(containerID, "failure") {
		notifMessage := fmt.Sprintf("❌ Container '%s' failed to deploy. Exit code: %s", containerName, exitCode)
		em.notifier.Send(notifMessage)
	}
}

// handleContainerCreate handles container create events
func (em *EventMonitor) handleContainerCreate(containerID, containerName string) {
	containerInfo, err := em.client.GetContainer(containerID)
	if err != nil {
		log.Printf("Error getting container info: %v", err)
		return
	}

	// Store or update container in database
	em.upsertContainer(containerID, containerName, containerInfo.Config.Image)
}

// shouldNotify checks if we should send notification for this container
func (em *EventMonitor) shouldNotify(containerID, eventType string) bool {
	records, err := em.db.App().Dao().FindRecordsByFilter("containers", "", "", 0, 0)
	if err != nil {
		return false
	}

	for _, record := range records {
		if record.GetString("container_id") == containerID {
			if eventType == "success" {
				return record.GetBool("notify_on_success")
			} else if eventType == "failure" {
				return record.GetBool("notify_on_failure")
			}
		}
	}

	// Default: notify on failures
	return eventType == "failure"
}

// logEvent logs an event to the database
func (em *EventMonitor) logEvent(containerID, containerName, eventType, status, message string) {
	collection, err := em.db.App().Dao().FindCollectionByNameOrId("events_log")
	if err != nil {
		log.Printf("Error finding events_log collection: %v", err)
		return
	}

	record := models.NewRecord(collection)
	record.Set("container_id", containerID)
	record.Set("container_name", containerName)
	record.Set("event_type", eventType)
	record.Set("status", status)
	record.Set("message", message)
	record.Set("timestamp", time.Now())

	if err := em.db.App().Dao().SaveRecord(record); err != nil {
		log.Printf("Error saving event log: %v", err)
	}
}

// upsertContainer creates or updates a container in the database
func (em *EventMonitor) upsertContainer(containerID, containerName, image string) {
	collection, err := em.db.App().Dao().FindCollectionByNameOrId("containers")
	if err != nil {
		log.Printf("Error finding containers collection: %v", err)
		return
	}

	// Try to find existing record
	records, err := em.db.App().Dao().FindRecordsByFilter("containers", "", "", 0, 0)
	if err != nil {
		log.Printf("Error finding container records: %v", err)
		return
	}

	var existingRecord *models.Record
	for _, record := range records {
		if record.GetString("container_id") == containerID {
			existingRecord = record
			break
		}
	}

	if existingRecord != nil {
		// Update existing
		existingRecord.Set("name", containerName)
		existingRecord.Set("image", image)
		if err := em.db.App().Dao().SaveRecord(existingRecord); err != nil {
			log.Printf("Error updating container: %v", err)
		}
	} else {
		// Create new
		record := models.NewRecord(collection)
		record.Set("container_id", containerID)
		record.Set("name", strings.TrimPrefix(containerName, "/"))
		record.Set("image", image)
		record.Set("notify_on_success", false)
		record.Set("notify_on_failure", true) // Default: notify on failures
		record.Set("status", "created")

		if err := em.db.App().Dao().SaveRecord(record); err != nil {
			log.Printf("Error creating container record: %v", err)
		}
	}
}

// Stop stops the event monitor
func (em *EventMonitor) Stop() {
	em.cancelFunc()
}
