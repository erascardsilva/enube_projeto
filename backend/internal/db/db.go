// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Obter variáveis de ambiente para conexão com o banco
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Montar DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	// Abrir a conexão com o GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
}

func ConnectAndMigrate() *gorm.DB {
	// Montar DSN usando variáveis de ambiente
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	// Abrir a conexão
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	// Verificar ambiente de desenvolvimento para limpar o banco
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env == "" || env == "development" || env == "dev" {
		log.Println("Ambiente de desenvolvimento detectado - Limpando banco de dados...")

		// Dropar tabelas existentes antes de migrar
		if err := db.Exec("DROP TABLE IF EXISTS billing_records CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela billing_records: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS partners CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela partners: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS customers CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela customers: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS products CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela products: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS skus CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela skus: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS publishers CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela publishers: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS subscriptions CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela subscriptions: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS meters CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela meters: %v", err)
		}
		if err := db.Exec("DROP TABLE IF EXISTS usuarios CASCADE").Error; err != nil {
			log.Printf("Aviso: Erro ao dropar tabela usuarios: %v", err)
		}
	}

	// Migrar modelos para criar tabelas
	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.BillingRecord{},
		&models.Partner{},
		&models.Customer{},
		&models.Product{},
		&models.Sku{},
		&models.Publisher{},
		&models.Subscription{},
		&models.Meter{},
	); err != nil {
		log.Fatalf("Erro ao migrar tabelas: %v", err)
	}

	// Criar usuário admin padrão se não existir
	var count int64
	db.Model(&models.Usuario{}).Where("email = ?", "admin@teste.com").Count(&count)
	if count == 0 {
		db.Create(&models.Usuario{
			Nome:  "Admin",
			Email: "admin@teste.com",
			Senha: "3727", // Em produção, usar hash
			Ativo: true,
		})
	}

	log.Println("Banco conectado, tabelas migradas, usuário padrão criado")
	return db
}
