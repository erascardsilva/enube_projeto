// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package models

import (
	"gorm.io/gorm"
)

// Definir struct para usuário do sistema
type Usuario struct {
	gorm.Model // Campos padrão do GORM
	Nome       string
	Email      string `gorm:"uniqueIndex"`
	Senha      string
	Ativo      bool
}

// Definir struct para Parceiro
type Partner struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	PartnerId   string `gorm:"uniqueIndex;type:uuid"`
	PartnerName string
}

// Definir struct para Cliente
type Customer struct {
	ID                 string `gorm:"primaryKey;type:uuid"`
	CustomerId         string `gorm:"uniqueIndex;type:uuid"`
	CustomerName       string
	CustomerDomainName string
	CustomerCountry    string
}

// Definir struct para Produto
type Product struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	ProductId   string `gorm:"uniqueIndex;type:uuid"`
	ProductName string
}

// Definir struct para SKU (Stock Keeping Unit)
type Sku struct {
	ID      string `gorm:"primaryKey;type:uuid"`
	SkuId   string `gorm:"uniqueIndex;type:uuid"`
	SkuName string
}

// Definir struct para Editora/Fabricante
type Publisher struct {
	ID            string `gorm:"primaryKey;type:uuid"`
	PublisherId   string `gorm:"uniqueIndex;type:uuid"`
	PublisherName string
}

// Definir struct para Assinatura
type Subscription struct {
	ID                      string `gorm:"primaryKey;type:uuid"`
	SubscriptionId          string `gorm:"uniqueIndex;type:uuid"`
	SubscriptionDescription string
}

// Definir struct para Medidor de Uso
type Meter struct {
	ID               string `gorm:"primaryKey;type:uuid"`
	MeterId          string `gorm:"uniqueIndex;type:uuid"`
	MeterType        string
	MeterCategory    string
	MeterSubCategory string
	MeterName        string
	MeterRegion      string
}


