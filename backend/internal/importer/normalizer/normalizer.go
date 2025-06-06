package normalizer

import (
	"strings"
	"time"

	"backend/internal/importer/parser"
)

// NormalizedRecord represents a record ready for database insertion
type NormalizedRecord struct {
	ClientID    string
	ClientName  string
	Category    string
	Resource    string
	Amount      float64
	BillingDate time.Time
	Description string
	Year        int
	Month       int
	Day         int
}

// Normalizer handles data normalization
type Normalizer struct{}

// NewNormalizer creates a new normalizer instance
func NewNormalizer() *Normalizer {
	return &Normalizer{}
}

// Normalize processes a parsed record and returns a normalized version
func (n *Normalizer) Normalize(record *parser.Record) (*NormalizedRecord, error) {
	if record == nil {
		return nil, ErrInvalidRecord
	}

	// Normalize strings
	clientName := strings.TrimSpace(record.ClientName)
	category := strings.TrimSpace(record.Category)
	resource := strings.TrimSpace(record.Resource)
	description := strings.TrimSpace(record.Description)

	// Extract date components
	year := record.BillingDate.Year()
	month := int(record.BillingDate.Month())
	day := record.BillingDate.Day()

	return &NormalizedRecord{
		ClientID:    record.ClientID,
		ClientName:  clientName,
		Category:    category,
		Resource:    resource,
		Amount:      record.Amount,
		BillingDate: record.BillingDate,
		Description: description,
		Year:        year,
		Month:       month,
		Day:         day,
	}, nil
}
