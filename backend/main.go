package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"backend/internal/auth"
	"backend/internal/db"
	"backend/internal/routes"
)

// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Inicializar conexão com o banco de dados
	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Inicializar serviços
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"))

	// Configurar rotas
	router := routes.SetupRouter(db, jwtService)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
