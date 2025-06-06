package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"backend/internal/importer/normalizer"
)

// Repository handles database operations
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Save persists a normalized record to the database
func (r *Repository) Save(ctx context.Context, record *normalizer.NormalizedRecord) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert or update client
		client := &Client{
			ClientID: record.ClientID,
			Name:     record.ClientName,
		}
		if err := tx.Where("client_id = ?", record.ClientID).FirstOrCreate(client).Error; err != nil {
			return err
		}

		// Insert or update category
		category := &Category{
			Name: record.Category,
		}
		if err := tx.Where("name = ?", record.Category).FirstOrCreate(category).Error; err != nil {
			return err
		}

		// Insert or update resource
		resource := &Resource{
			Name: record.Resource,
		}
		if err := tx.Where("name = ?", record.Resource).FirstOrCreate(resource).Error; err != nil {
			return err
		}

		// Insert billing record
		billing := &Billing{
			ClientID:    client.ID,
			CategoryID:  category.ID,
			ResourceID:  resource.ID,
			Amount:      record.Amount,
			BillingDate: record.BillingDate,
			Description: record.Description,
			Year:        record.Year,
			Month:       record.Month,
			Day:         record.Day,
		}
		return tx.Create(billing).Error
	})
}

// Client represents a client in the database
type Client struct {
	gorm.Model
	ClientID string `gorm:"uniqueIndex"`
	Name     string
}

// Category represents a category in the database
type Category struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

// Resource represents a resource in the database
type Resource struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

// Billing represents a billing record in the database
type Billing struct {
	gorm.Model
	ClientID    uint
	CategoryID  uint
	ResourceID  uint
	Amount      float64
	BillingDate time.Time
	Description string
	Year        int
	Month       int
	Day         int
}
