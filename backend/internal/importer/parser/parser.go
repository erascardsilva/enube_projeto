package parser

import (
	"time"
)

// Record represents a parsed billing record
type Record struct {
	ClientID    string
	ClientName  string
	Category    string
	Resource    string
	Amount      float64
	BillingDate time.Time
	Description string
}

// Parser handles parsing of CSV records
type Parser struct{}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse converts a CSV record into a structured Record
func (p *Parser) Parse(record []string) (*Record, error) {
	// Assuming CSV columns are in order:
	// client_id, client_name, category, resource, amount, billing_date, description
	if len(record) < 7 {
		return nil, ErrInvalidRecord
	}

	billingDate, err := time.Parse("2006-01-02", record[5])
	if err != nil {
		return nil, ErrInvalidDate
	}

	amount, err := parseAmount(record[4])
	if err != nil {
		return nil, ErrInvalidAmount
	}

	return &Record{
		ClientID:    record[0],
		ClientName:  record[1],
		Category:    record[2],
		Resource:    record[3],
		Amount:      amount,
		BillingDate: billingDate,
		Description: record[6],
	}, nil
}

func parseAmount(s string) (float64, error) {
	// Implement amount parsing logic
	// This is a placeholder - implement according to your data format
	return 0, nil
}
