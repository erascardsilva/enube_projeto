// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package models

import (
	"time"

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
	ID          string `gorm:"primaryKey;type:varchar(255)"`
	PartnerId   string `gorm:"uniqueIndex;type:varchar(255)"`
	PartnerName string
}

// Definir struct para Cliente
type Customer struct {
	ID                 string `gorm:"primaryKey;type:varchar(255)"`
	CustomerId         string `gorm:"uniqueIndex;type:varchar(255)"`
	CustomerName       string
	CustomerDomainName string
	CustomerCountry    string
}

// Definir struct para Produto
type Product struct {
	ID          string `gorm:"primaryKey;type:varchar(255)"`
	ProductId   string `gorm:"uniqueIndex;type:varchar(255)"`
	ProductName string
}

// Definir struct para SKU (Stock Keeping Unit)
type Sku struct {
	ID      string `gorm:"primaryKey;type:varchar(255)"`
	SkuId   string `gorm:"uniqueIndex;type:varchar(255)"`
	SkuName string
}

// Definir struct para Editora/Fabricante
type Publisher struct {
	ID            string `gorm:"primaryKey;type:varchar(255)"`
	PublisherId   string `gorm:"uniqueIndex;type:varchar(255)"`
	PublisherName string
}

// Definir struct para Assinatura
type Subscription struct {
	ID                      string `gorm:"primaryKey;type:varchar(255)"`
	SubscriptionId          string `gorm:"uniqueIndex;type:varchar(255)"`
	SubscriptionDescription string
}

// Definir struct para Medidor de Uso
type Meter struct {
	ID               string `gorm:"primaryKey;type:varchar(255)"`
	MeterId          string `gorm:"uniqueIndex;type:varchar(255)"`
	MeterType        string
	MeterCategory    string
	MeterSubCategory string
	MeterName        string
	MeterRegion      string
}

// Client representa um registro na tabela 'clients'
type Client struct {
	gorm.Model
	ClientID   string `gorm:"type:varchar(255);unique;not null;index"` // ID do cliente do Excel
	Name       string `gorm:"not null"`
	DomainName string
	Country    string
}

// Category representa um registro na tabela 'categories'
type Category struct {
	gorm.Model
	Name        string `gorm:"unique;not null;index"`
	Type        string // MeterType
	SubCategory string // MeterSubCategory
}

// Resource representa um registro na tabela 'resources'
type Resource struct {
	gorm.Model
	Name            string `gorm:"column:name;unique;not null;index"` // ResourceURI
	Location        string `gorm:"column:location"`                   // ResourceLocation
	Group           string `gorm:"column:group"`                      // ResourceGroup
	ConsumedService string `gorm:"column:consumed_service"`           // ConsumedService
}

// Billing representa um registro na tabela 'billing'
type Billing struct {
	gorm.Model
	ClientID           uint      `gorm:"not null;index"`              // Chave estrangeira para Client
	CategoryID         uint      `gorm:"not null;index"`              // Chave estrangeira para Category
	ResourceID         uint      `gorm:"not null;index"`              // Chave estrangeira para Resource
	Amount             float64   `gorm:"type:decimal(10,2);not null"` // BillingPreTaxTotal
	BillingDate        time.Time `gorm:"not null"`                    // ChargeStartDate
	Description        string    // MeterSubCategory
	Year               int       `gorm:"not null;index"`
	Month              int       `gorm:"not null;index"`
	Day                int       `gorm:"not null;index"`
	Unit               string    `gorm:"column:unit"`        // Unit
	UnitPrice          float64   `gorm:"type:decimal(20,6)"` // UnitPrice
	Quantity           float64   `gorm:"type:decimal(20,6)"` // Quantity
	UnitType           string    `gorm:"column:unit_type"`   // UnitType
	BillingCurrency    string    `gorm:"type:varchar(3)"`    // BillingCurrency
	PricingPreTaxTotal float64   `gorm:"type:decimal(20,6)"` // PricingPreTaxTotal
	PricingCurrency    string    `gorm:"type:varchar(3)"`    // PricingCurrency

	// Relacionamentos
	Client   Client   `gorm:"foreignKey:ClientID"`
	Category Category `gorm:"foreignKey:CategoryID"`
	Resource Resource `gorm:"foreignKey:ResourceID"`
}

// TableName define o nome da tabela para o modelo Billing
func (Billing) TableName() string {
	return "billing"
}
