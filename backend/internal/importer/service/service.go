// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack
//
// TODO: Implementar cache para melhorar performance
// TODO: Adicionar timeout configurável
// TODO: Melhorar tratamento de erros
// TODO: Implementar retry mechanism
// TODO: Adicionar mais métricas de performance
// TODO: Considerar usar context para cancelamento
// TODO: Implementar rate limiting
// TODO: Adicionar validação de arquivo antes do processamento
// TODO: Implementar progress tracking
// TODO: Adicionar suporte para diferentes formatos de arquivo

package service

import (
	"context"
	"sync"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"backend/internal/importer/normalizer"
	"backend/internal/importer/parser"
	"backend/internal/importer/repository"
)

// ImportStats holds statistics about the import process
type ImportStats struct {
	TotalRecords    int64
	ImportedRecords int64
	FailedRecords   int64
	Categories      int64
	Clients         int64
	Resources       int64
}

// Importer handles the data import process
type Importer struct {
	db         *gorm.DB
	repo       *repository.Repository
	parser     *parser.Parser
	normalizer *normalizer.Normalizer
}

// NewImporter creates a new importer instance
func NewImporter(db *gorm.DB) *Importer {
	return &Importer{
		db:         db,
		repo:       repository.NewRepository(db),
		parser:     parser.NewParser(),
		normalizer: normalizer.NewNormalizer(),
	}
}

// Import processes the Excel file and imports it into the database
func (i *Importer) Import(ctx context.Context, filePath string, batchSize int) (*ImportStats, error) {
	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, ErrNoSheets
	}

	// Get all rows from the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	// Skip header row
	if len(rows) < 2 {
		return nil, ErrEmptyFile
	}
	rows = rows[1:]

	stats := &ImportStats{}
	var wg sync.WaitGroup
	records := make(chan []string, batchSize)
	errors := make(chan error, 1)

	// Start worker pool
	for w := 0; w < 4; w++ {
		wg.Add(1)
		go i.worker(ctx, &wg, records, stats)
	}

	// Process rows
	go func() {
		defer close(records)
		for _, row := range rows {
			select {
			case records <- row:
				stats.TotalRecords++
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for all workers to finish
	wg.Wait()

	select {
	case err := <-errors:
		return stats, err
	default:
		return stats, nil
	}
}

func (i *Importer) worker(ctx context.Context, wg *sync.WaitGroup, records <-chan []string, stats *ImportStats) {
	defer wg.Done()

	for record := range records {
		select {
		case <-ctx.Done():
			return
		default:
			// Parse record
			parsed, err := i.parser.Parse(record)
			if err != nil {
				stats.FailedRecords++
				continue
			}

			// Normalize data
			normalized, err := i.normalizer.Normalize(parsed)
			if err != nil {
				stats.FailedRecords++
				continue
			}

			// Save to database
			if err := i.repo.Save(ctx, normalized); err != nil {
				stats.FailedRecords++
				continue
			}

			stats.ImportedRecords++
		}
	}
}
