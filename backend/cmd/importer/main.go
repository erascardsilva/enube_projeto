package main

import (
	"context"
	"flag"
	"log"
	"time"

	"backend/internal/db"
	"backend/internal/importer/service"
)

func main() {
	// Parse command line flags
	filePath := flag.String("file", "", "Path to the data file to import")
	batchSize := flag.Int("batch", 1000, "Number of records to process in each batch")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("File path is required")
	}

	// Initialize database connection
	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create importer service
	importer := service.NewImporter(db)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Start import process
	start := time.Now()
	stats, err := importer.Import(ctx, *filePath, *batchSize)
	if err != nil {
		log.Fatalf("Import failed: %v", err)
	}

	// Log results
	log.Printf("Import completed in %v", time.Since(start))
	log.Printf("Records processed: %d", stats.TotalRecords)
	log.Printf("Records imported: %d", stats.ImportedRecords)
	log.Printf("Records failed: %d", stats.FailedRecords)
	log.Printf("Categories found: %d", stats.Categories)
	log.Printf("Clients found: %d", stats.Clients)
	log.Printf("Resources found: %d", stats.Resources)
}
