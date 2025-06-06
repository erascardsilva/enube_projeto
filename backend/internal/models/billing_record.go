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
	PartnerID      *string `gorm:"type:varchar(255)"`
	Partner        Partner
	CustomerID     *string `gorm:"type:varchar(255)"`
	Customer       Customer
	ProductID      *string `gorm:"type:varchar(255)"`
	Product        Product
	SkuID          *string `gorm:"type:varchar(255)"`
	Sku            Sku
	PublisherID    *string `gorm:"type:varchar(255)"`
	Publisher      Publisher
	SubscriptionID *string `gorm:"type:varchar(255)"`
	Subscription   Subscription
	MeterID        *string `gorm:"type:varchar(255)"`
	Meter          Meter

	// Definir campos específicos do registro de Billing (alguns IDs são opcionais)
	MpnId                         *string   `gorm:"type:varchar(255)"`
	Tier2MpnId                    *string   `gorm:"type:varchar(255)"`
	InvoiceNumber                 string    `gorm:"not null"`
	AvailabilityId                *string   `gorm:"type:varchar(255)"`
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
	EntitlementId                 *string   `gorm:"type:varchar(255)"`
	EntitlementDescription        string
	PartnerEarnedCreditPercentage float64 `gorm:"type:decimal(5,2);not null"`
	CreditPercentage              float64 `gorm:"type:decimal(5,2);not null"`
	CreditType                    string
	BenefitOrderId                *string `gorm:"type:varchar(255)"`
	BenefitId                     *string `gorm:"type:varchar(255)"`
	BenefitType                   string
}
