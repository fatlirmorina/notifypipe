package notifications

import (
	"fmt"
	"log"

	"github.com/containrrr/shoutrrr"
	"github.com/fatlirmorina/notifypipe/internal/database"
)

// Manager manages notifications
type Manager struct {
	db *database.Database
}

// NewManager creates a new notification manager
func NewManager(db *database.Database) *Manager {
	return &Manager{db: db}
}

// Send sends a notification to all enabled channels
func (m *Manager) Send(message string) {
	records, err := m.db.App().Dao().FindRecordsByExpr("notifications", nil)
	if err != nil {
		log.Printf("Error fetching notifications: %v", err)
		return
	}

	if len(records) == 0 {
		log.Println("No notification channels configured")
		return
	}

	for _, record := range records {
		if !record.GetBool("enabled") {
			continue
		}

		url := record.GetString("url")
		name := record.GetString("name")

		if err := m.SendToURL(url, message); err != nil {
			log.Printf("Error sending notification to %s: %v", name, err)
		} else {
			log.Printf("✅ Notification sent to %s", name)
		}
	}
}

// SendToURL sends a notification to a specific URL
func (m *Manager) SendToURL(url, message string) error {
	sender, err := shoutrrr.CreateSender(url)
	if err != nil {
		return fmt.Errorf("failed to create sender: %w", err)
	}

	errs := sender.Send(message, nil)
	if len(errs) > 0 {
		return fmt.Errorf("send errors: %v", errs)
	}

	return nil
}

// TestNotification tests a notification URL
func (m *Manager) TestNotification(url string) error {
	testMessage := "🔔 Test notification from NotifyPipe! If you receive this, your notification channel is configured correctly."
	return m.SendToURL(url, testMessage)
}
