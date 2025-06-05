package main

import (
	"log"
	// Remover net/http

	"backend/internal/db"
	"backend/internal/handlers" // Importar handlers

	"github.com/gin-gonic/gin"
)

// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

func main() {
	// Conectar e migrar banco de dados
	database := db.ConnectAndMigrate()

	// Criar instância do Gin
	r := gin.Default()

	// Registrar handler de importação
	r.POST("/importar-xls", handlers.ImportXLSHandler(database))

	// Iniciar servidor
	log.Println("Servidor rodando na porta 8080...")
	r.Run(":8080")
}
