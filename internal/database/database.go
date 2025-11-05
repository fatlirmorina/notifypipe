package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// Database wraps PocketBase
type Database struct {
	app *pocketbase.PocketBase
}

// New creates a new database instance
func New(dataDir string) (*Database, error) {
	// Ensure data directory exists
	dbPath := filepath.Join(dataDir, "pb_data")
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, err
	}
	
	// Set environment variable for PocketBase data directory
	os.Setenv("PB_DATA_DIR", dbPath)
	
	app := pocketbase.New()
	db := &Database{app: app}

	// Bootstrap the application
	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	// Setup collections immediately after bootstrap
	if err := db.setupCollections(); err != nil {
		log.Printf("Warning: Error setting up collections: %v", err)
	}

	return db, nil
}

// setupCollections creates the required collections
func (db *Database) setupCollections() error {
	// Check if collections already exist
	collections, err := db.app.Dao().FindCollectionsByType("base")
	if err != nil {
		log.Printf("Error checking existing collections: %v", err)
		// Continue anyway to try creating them
	}

	// If we already have collections, skip setup
	if len(collections) >= 4 {
		log.Println("Collections already exist, skipping setup")
		return nil
	}

	log.Println("Setting up database collections...")

	// Create notifications collection
	notificationsCollection := &models.Collection{}
	notificationsCollection.Name = "notifications"
	notificationsCollection.Type = models.CollectionTypeBase
	notificationsCollection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name:     "name",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name:     "type",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name:     "url",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name: "enabled",
			Type: schema.FieldTypeBool,
		},
	)

	if err := db.app.Dao().SaveCollection(notificationsCollection); err != nil {
		log.Printf("Error creating notifications collection: %v", err)
	} else {
		log.Println("✅ Created notifications collection")
	}

	// Create containers collection
	containersCollection := &models.Collection{}
	containersCollection.Name = "containers"
	containersCollection.Type = models.CollectionTypeBase
	containersCollection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name:     "container_id",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name:     "name",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name: "notify_on_success",
			Type: schema.FieldTypeBool,
		},
		&schema.SchemaField{
			Name: "notify_on_failure",
			Type: schema.FieldTypeBool,
		},
		&schema.SchemaField{
			Name: "image",
			Type: schema.FieldTypeText,
		},
		&schema.SchemaField{
			Name: "status",
			Type: schema.FieldTypeText,
		},
	)

	if err := db.app.Dao().SaveCollection(containersCollection); err != nil {
		log.Printf("Error creating containers collection: %v", err)
	} else {
		log.Println("✅ Created containers collection")
	}

	// Create events_log collection
	eventsCollection := &models.Collection{}
	eventsCollection.Name = "events_log"
	eventsCollection.Type = models.CollectionTypeBase
	eventsCollection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name:     "container_id",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name: "container_name",
			Type: schema.FieldTypeText,
		},
		&schema.SchemaField{
			Name:     "event_type",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name:     "status",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name: "message",
			Type: schema.FieldTypeText,
		},
		&schema.SchemaField{
			Name: "timestamp",
			Type: schema.FieldTypeDate,
		},
	)

	if err := db.app.Dao().SaveCollection(eventsCollection); err != nil {
		log.Printf("Error creating events_log collection: %v", err)
	} else {
		log.Println("✅ Created events_log collection")
	}

	// Create settings collection
	settingsCollection := &models.Collection{}
	settingsCollection.Name = "settings"
	settingsCollection.Type = models.CollectionTypeBase
	settingsCollection.Schema = schema.NewSchema(
		&schema.SchemaField{
			Name:     "key",
			Type:     schema.FieldTypeText,
			Required: true,
		},
		&schema.SchemaField{
			Name: "value",
			Type: schema.FieldTypeText,
		},
	)

	if err := db.app.Dao().SaveCollection(settingsCollection); err != nil {
		log.Printf("Error creating settings collection: %v", err)
	} else {
		log.Println("✅ Created settings collection")
	}

	log.Println("✅ Database setup completed")
	return nil
}

// App returns the PocketBase app instance
func (db *Database) App() *pocketbase.PocketBase {
	return db.app
}

// Close closes the database connection
func (db *Database) Close() error {
	// PocketBase handles cleanup internally
	return nil
}
