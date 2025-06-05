// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack
package models

import (
	"time"

	"gorm.io/gorm"
)

type BillingRecord struct {
	gorm.Model // Adicionar campos padrão do GORM

	// Definir chaves estrangeiras e relacionamentos (alguns podem ser nulos)
	PartnerID      *string `gorm:"type:uuid"`
	Partner        Partner
	CustomerID     *string `gorm:"type:uuid"`
	Customer       Customer
	ProductID      *string `gorm:"type:uuid"`
	Product        Product
	SkuID          *string `gorm:"type:uuid"`
	Sku            Sku
	PublisherID    *string `gorm:"type:uuid"`
	Publisher      Publisher
	SubscriptionID *string `gorm:"type:uuid"`
	Subscription   Subscription
	MeterID        *string `gorm:"type:uuid"`
	Meter          Meter

	// Definir campos específicos do registro de Billing (alguns IDs são opcionais)
	MpnId                         *string   `gorm:"type:uuid"`
	Tier2MpnId                    *string   `gorm:"type:uuid"`
	InvoiceNumber                 string    `gorm:"not null"`
	AvailabilityId                *string   `gorm:"type:uuid"`
	ChargeStartDate               time.Time `gorm:"type:timestamp;not null"`
	ChargeEndDate                 time.Time `gorm:"type:timestamp;not null"`
	UsageDate                     time.Time `gorm:"type:timestamp;not null"`
	Unit                          string
	ResourceLocation              string
	ConsumedService               string
	ResourceGroup                 string
	ResourceURI                   string
	ChargeType                    string
	UnitPrice                     float64 `gorm:"type:decimal(20,6);not null"`
	Quantity                      float64 `gorm:"type:decimal(20,6);not null"`
	UnitType                      string
	BillingPreTaxTotal            float64 `gorm:"type:decimal(20,6);not null"`
	BillingCurrency               string  `gorm:"type:varchar(3);not null"`
	PricingPreTaxTotal            float64 `gorm:"type:decimal(20,6);not null"`
	PricingCurrency               string  `gorm:"type:varchar(3);not null"`
	ServiceInfo1                  string
	ServiceInfo2                  string
	Tags                          string
	AdditionalInfo                string
	EffectiveUnitPrice            float64   `gorm:"type:decimal(20,6);not null"`
	PCToBCExchangeRate            float64   `gorm:"type:decimal(20,6);not null"`
	PCToBCExchangeRateDate        time.Time `gorm:"type:timestamp"`
	EntitlementId                 *string   `gorm:"type:uuid"`
	EntitlementDescription        string
	PartnerEarnedCreditPercentage float64 `gorm:"type:decimal(5,2);not null"`
	CreditPercentage              float64 `gorm:"type:decimal(5,2);not null"`
	CreditType                    string
	BenefitOrderId                *string `gorm:"type:uuid"`
	BenefitId                     *string `gorm:"type:uuid"`
	BenefitType                   string
}

// Nota: Adicionar outras entidades posteriormente se necessário
