// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package main

import (
	"log"

	"backend/internal/db"
	"backend/internal/models"
)

func main() {
	// Inicializar conexão com o banco de dados
	database, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	// Migrar o banco de dados
	if err := database.AutoMigrate(&models.User{}, &models.Client{}, &models.Category{}, &models.Resource{}, &models.Billing{}); err != nil {
		log.Fatal("Erro ao migrar banco de dados:", err)
	}

	// Adicionar índices
	if err := models.AddIndexes(database); err != nil {
		log.Fatal("Erro ao adicionar índices:", err)
	}
}
