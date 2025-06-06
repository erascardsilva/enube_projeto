// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package models

import (
	"gorm.io/gorm"
)

// AddIndexes adiciona índices para otimizar consultas
func AddIndexes(db *gorm.DB) error {
	// Índices para a tabela billing
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_billing_client_id ON billing(client_id);
		CREATE INDEX IF NOT EXISTS idx_billing_category_id ON billing(category_id);
		CREATE INDEX IF NOT EXISTS idx_billing_resource_id ON billing(resource_id);
		CREATE INDEX IF NOT EXISTS idx_billing_date ON billing(billing_date);
		CREATE INDEX IF NOT EXISTS idx_billing_year_month ON billing(year, month);
	`).Error; err != nil {
		return err
	}

	// Índices para a tabela clients
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_clients_client_id ON clients(client_id);
		CREATE INDEX IF NOT EXISTS idx_clients_name ON clients(name);
	`).Error; err != nil {
		return err
	}

	// Índices para a tabela categories
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
	`).Error; err != nil {
		return err
	}

	// Índices para a tabela resources
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_resources_name ON resources(name);
		CREATE INDEX IF NOT EXISTS idx_resources_location ON resources(location);
	`).Error; err != nil {
		return err
	}

	return nil
}
