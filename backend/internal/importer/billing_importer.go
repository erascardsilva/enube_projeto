// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack
//
// TODO: Otimizar o processamento paralelo para melhor performance
// TODO: Adicionar mais logs para debug
// TODO: Implementar retry mechanism para falhas de banco
// TODO: Melhorar tratamento de erros nas datas
// TODO: Considerar usar batch processing para entidades relacionadas
// TODO: Avaliar uso de prepared statements para melhor performance
// TODO: Implementar validação mais robusta dos dados de entrada
// TODO: Adicionar métricas de performance
// TODO: Considerar usar transactions para garantir consistência
// TODO: Implementar mecanismo de rollback em caso de falha

package importer

import (
	"backend/internal/models"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Definir lista de colunas esperadas no arquivo XLS de billing
var ExpectedBillingHeader = []string{
	"PartnerId",
	"PartnerName",
	"CustomerId",
	"CustomerName",
	"CustomerDomainName",
	"CustomerCountry",
	"MpnId",
	"Tier2MpnId",
	"InvoiceNumber",
	"ProductId",
	"SkuId",
	"AvailabilityId",
	"SkuName",
	"ProductName",
	"PublisherName",
	"PublisherId",
	"SubscriptionDescription",
	"SubscriptionId",
	"ChargeStartDate",
	"ChargeEndDate",
	"UsageDate",
	"MeterType",
	"MeterCategory",
	"MeterId",
	"MeterSubCategory",
	"MeterName",
	"MeterRegion",
	"Unit",
	"ResourceLocation",
	"ConsumedService",
	"ResourceGroup",
	"ResourceURI",
	"ChargeType",
	"UnitPrice",
	"Quantity",
	"UnitType",
	"BillingPreTaxTotal",
	"BillingCurrency",
	"PricingPreTaxTotal",
	"PricingCurrency",
	"ServiceInfo1",
	"ServiceInfo2",
	"Tags",
	"AdditionalInfo",
	"EffectiveUnitPrice",
	"PCToBCExchangeRate",
	"PCToBCExchangeRateDate",
	"EntitlementId",
	"EntitlementDescription",
	"PartnerEarnedCreditPercentage",
	"CreditPercentage",
	"CreditType",
	"BenefitOrderId",
	"BenefitId",
	"BenefitType",
}

// Definir structs para coletar dados brutos
type PartnerData struct{ PartnerId, PartnerName string }
type CustomerData struct{ CustomerId, CustomerName, CustomerDomainName, CustomerCountry string }
type ProductData struct{ ProductId, ProductName string }
type SkuData struct{ SkuId, SkuName string }
type PublisherData struct{ PublisherId, PublisherName string }
type SubscriptionData struct{ SubscriptionId, SubscriptionDescription string }
type MeterData struct{ MeterId, MeterType, MeterCategory, MeterSubCategory, MeterName, MeterRegion string }

// Converter XLS para structs BillingRecord
func ImportBillingXLS(db *gorm.DB, filePath string) ([]models.BillingRecord, []string, []string, error) {
	log.Println("--- Iniciando ImportBillingXLS ---")
	if db == nil {
		log.Println("Erro: DB é nil")
		return nil, nil, nil, fmt.Errorf("objeto DB é nil em ImportBillingXLS")
	}

	// Abrir o arquivo Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer f.Close()

	// Obter a primeira planilha
	sheet := f.GetSheetList()[0]
	// Obter todas as linhas da planilha
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("erro ao ler planilha: %v", err)
	}

	// Validar se o arquivo tem pelo menos cabeçalho e uma linha de dados
	if len(rows) < 2 {
		return nil, nil, nil, fmt.Errorf("O arquivo XLS não contém dados além do cabeçalho.")
	}

	// Mapear colunas do cabeçalho
	header := rows[0]
	colMap := make(map[string]int)
	for i, colName := range header {
		colMap[strings.TrimSpace(colName)] = i
	}

	// Verificar colunas faltantes
	var missingColumns []string
	for _, expectedCol := range ExpectedBillingHeader {
		if _, ok := colMap[expectedCol]; !ok {
			missingColumns = append(missingColumns, expectedCol)
		}
	}

	// Opcional: Avisar se a contagem de colunas for diferente
	if len(header) != len(ExpectedBillingHeader) {
		log.Printf("Aviso: Número de colunas no cabeçalho do arquivo (%d) difere do esperado (%d).", len(header), len(ExpectedBillingHeader))
	}

	// Processar dados em paralelo
	dataRows := rows[1:]
	numWorkers := runtime.NumCPU()
	if numWorkers > len(dataRows) {
		numWorkers = len(dataRows)
	}
	if numWorkers == 0 && len(dataRows) > 0 {
		numWorkers = 1
	}

	// Configurar canais para processamento paralelo
	rowsChan := make(chan []string, len(dataRows))
	resultChan := make(chan struct {
		record       models.BillingRecord
		partner      PartnerData
		customer     CustomerData
		product      ProductData
		sku          SkuData
		publisher    PublisherData
		subscription SubscriptionData
		meter        MeterData
		warning      string
	}, len(dataRows))

	var wg sync.WaitGroup

	// Iniciar workers para processar linhas
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowsChan {
				record, partnerData, customerData, productData, skuData, publisherData, subscriptionData, meterData, warning := processRow(row, colMap)
				resultChan <- struct {
					record       models.BillingRecord
					partner      PartnerData
					customer     CustomerData
					product      ProductData
					sku          SkuData
					publisher    PublisherData
					subscription SubscriptionData
					meter        MeterData
					warning      string
				}{record, partnerData, customerData, productData, skuData, publisherData, subscriptionData, meterData, warning}
			}
		}()
	}

	// Enviar linhas para o canal de workers
	go func() {
		for _, row := range dataRows {
			rowsChan <- row
		}
		close(rowsChan)
	}()

	// Esperar workers terminarem e fechar canal de resultados
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados dos canais
	var records []models.BillingRecord
	var allPartnerData []PartnerData
	var allCustomerData []CustomerData
	var allProductData []ProductData
	var allSkuData []SkuData
	var allPublisherData []PublisherData
	var allSubscriptionData []SubscriptionData
	var allMeterData []MeterData
	var finalWarnings []string

	for result := range resultChan {
		records = append(records, result.record)
		allPartnerData = append(allPartnerData, result.partner)
		allCustomerData = append(allCustomerData, result.customer)
		allProductData = append(allProductData, result.product)
		allSkuData = append(allSkuData, result.sku)
		allPublisherData = append(allPublisherData, result.publisher)
		allSubscriptionData = append(allSubscriptionData, result.subscription)
		allMeterData = append(allMeterData, result.meter)
		if result.warning != "" {
			finalWarnings = append(finalWarnings, result.warning)
		}
	}

	// Processar e salvar entidades únicas e obter mapas de ID
	partnerIDMap, customerIDMap, productIDMap, skuIDMap, publisherIDMap, subscriptionIDMap, meterIDMap, err := processAndSaveEntities(db, allPartnerData, allCustomerData, allProductData, allSkuData, allPublisherData, allSubscriptionData, allMeterData)
	if err != nil {
		return nil, finalWarnings, missingColumns, fmt.Errorf("erro ao processar e salvar entidades: %v", err)
	}

	// Preencher IDs de chave estrangeira nos registros
	recordsWithFKs := populateBillingRecordsWithFKs(records, allPartnerData, allCustomerData, allProductData, allSkuData, allPublisherData, allSubscriptionData, allMeterData, partnerIDMap, customerIDMap, productIDMap, skuIDMap, publisherIDMap, subscriptionIDMap, meterIDMap)

	// Preparar registros válidos para salvar
	var validRecords []models.BillingRecord
	for _, rec := range recordsWithFKs {
		validRecords = append(validRecords, rec)
	}

	log.Printf("Vou salvar %d registros no banco de dados.", len(validRecords))

	// Salvar registros de billing no banco
	if err := SaveBillingRecords(db, validRecords); err != nil {
		return validRecords, finalWarnings, missingColumns, fmt.Errorf("erro ao salvar registros de billing: %v", err)
	}

	return validRecords, finalWarnings, missingColumns, nil
}

// Obter valor da coluna
func getColumnValue(row []string, colMap map[string]int, colName string) string {
	if index, ok := colMap[colName]; ok {
		if index < len(row) {
			return row[index]
		}
	}
	return ""
}

// Processar linha do XLS e retornar dados
func processRow(row []string, colMap map[string]int) (models.BillingRecord, PartnerData, CustomerData, ProductData, SkuData, PublisherData, SubscriptionData, MeterData, string) {
	var warning string
	record := models.BillingRecord{}

	// Função auxiliar para obter valor da coluna de forma segura
	getValue := func(colName string) string {
		return getColumnValue(row, colMap, colName)
	}

	// Processar IDs e converter para UUIDs opcionais
	partnerIdCleaned := cleanUUID(getValue("PartnerId"))
	if partnerIdCleaned != "" {
		record.PartnerID = &partnerIdCleaned
	} else {
		record.PartnerID = nil
	}

	customerIdCleaned := cleanUUID(getValue("CustomerId"))
	if customerIdCleaned != "" {
		record.CustomerID = &customerIdCleaned
	} else {
		record.CustomerID = nil
	}

	productIdCleaned := cleanUUID(getValue("ProductId"))
	if productIdCleaned != "" {
		record.ProductID = &productIdCleaned
	} else {
		record.ProductID = nil
	}

	skuIdCleaned := cleanUUID(getValue("SkuId"))
	if skuIdCleaned != "" {
		record.SkuID = &skuIdCleaned
	} else {
		record.SkuID = nil
	}

	publisherIdCleaned := cleanUUID(getValue("PublisherId"))
	if publisherIdCleaned != "" {
		record.PublisherID = &publisherIdCleaned
	} else {
		record.PublisherID = nil
	}

	subscriptionIdCleaned := cleanUUID(getValue("SubscriptionId"))
	if subscriptionIdCleaned != "" {
		record.SubscriptionID = &subscriptionIdCleaned
	} else {
		record.SubscriptionID = nil
	}

	meterIdCleaned := cleanUUID(getValue("MeterId"))
	if meterIdCleaned != "" {
		record.MeterID = &meterIdCleaned
	} else {
		record.MeterID = nil
	}

	// Criar structs temporários para dados brutos
	partnerData := PartnerData{
		PartnerId:   cleanUUID(getValue("PartnerId")),
		PartnerName: getValue("PartnerName"),
	}

	customerData := CustomerData{
		CustomerId:         cleanUUID(getValue("CustomerId")),
		CustomerName:       getValue("CustomerName"),
		CustomerDomainName: getValue("CustomerDomainName"),
		CustomerCountry:    getValue("CustomerCountry"),
	}

	productData := ProductData{
		ProductId:   cleanUUID(getValue("ProductId")),
		ProductName: getValue("ProductName"),
	}

	skuData := SkuData{
		SkuId:   cleanUUID(getValue("SkuId")),
		SkuName: getValue("SkuName"),
	}

	publisherData := PublisherData{
		PublisherId:   cleanUUID(getValue("PublisherId")),
		PublisherName: getValue("PublisherName"),
	}

	subscriptionData := SubscriptionData{
		SubscriptionId:          cleanUUID(getValue("SubscriptionId")),
		SubscriptionDescription: getValue("SubscriptionDescription"),
	}

	meterData := MeterData{
		MeterId:          cleanUUID(getValue("MeterId")),
		MeterType:        getValue("MeterType"),
		MeterCategory:    getValue("MeterCategory"),
		MeterSubCategory: getValue("MeterSubCategory"),
		MeterName:        getValue("MeterName"),
		MeterRegion:      getValue("MeterRegion"),
	}

	// Preencher campos diretos do registro
	record.MpnId = stringPtr(cleanUUID(getValue("MpnId")))
	record.Tier2MpnId = stringPtr(cleanUUID(getValue("Tier2MpnId")))
	record.InvoiceNumber = getValue("InvoiceNumber")
	record.AvailabilityId = stringPtr(cleanUUID(getValue("AvailabilityId")))

	// Converter datas
	chargeStartDateStr := getValue("ChargeStartDate")
	if date, err := parseDate(chargeStartDateStr); err == nil {
		record.ChargeStartDate = date
	} else if chargeStartDateStr != "" {
		warning = fmt.Sprintf("%s; Erro ao converter ChargeStartDate ('%s'): %s", warning, chargeStartDateStr, err.Error())
	}

	chargeEndDateStr := getValue("ChargeEndDate")
	if date, err := parseDate(chargeEndDateStr); err == nil {
		record.ChargeEndDate = date
	} else if chargeEndDateStr != "" {
		warning = fmt.Sprintf("%s; Erro ao converter ChargeEndDate ('%s'): %s", warning, chargeEndDateStr, err.Error())
	}

	usageDateStr := getValue("UsageDate")
	if date, err := parseDate(usageDateStr); err == nil {
		record.UsageDate = date
	} else if usageDateStr != "" {
		warning = fmt.Sprintf("%s; Erro ao converter UsageDate ('%s'): %s", warning, usageDateStr, err.Error())
	}

	record.Unit = getValue("Unit")
	record.ResourceLocation = getValue("ResourceLocation")
	record.ConsumedService = getValue("ConsumedService")
	record.ResourceGroup = getValue("ResourceGroup")
	record.ResourceURI = getValue("ResourceURI")
	record.ChargeType = getValue("ChargeType")

	// Converter valores numéricos
	if val, err := strconv.ParseFloat(getValue("UnitPrice"), 64); err == nil {
		record.UnitPrice = val
	} else if getValue("UnitPrice") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter UnitPrice: %s", warning, err)
	}

	if val, err := strconv.ParseFloat(getValue("Quantity"), 64); err == nil {
		record.Quantity = val
	} else if getValue("Quantity") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter Quantity: %s", warning, err)
	}

	record.UnitType = getValue("UnitType")

	if val, err := strconv.ParseFloat(getValue("BillingPreTaxTotal"), 64); err == nil {
		record.BillingPreTaxTotal = val
	} else if getValue("BillingPreTaxTotal") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter BillingPreTaxTotal: %s", warning, err)
	}

	record.BillingCurrency = getValue("BillingCurrency")

	if val, err := strconv.ParseFloat(getValue("PricingPreTaxTotal"), 64); err == nil {
		record.PricingPreTaxTotal = val
	} else if getValue("PricingPreTaxTotal") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter PricingPreTaxTotal: %s", warning, err)
	}

	record.PricingCurrency = getValue("PricingCurrency")
	record.ServiceInfo1 = getValue("ServiceInfo1")
	record.ServiceInfo2 = getValue("ServiceInfo2")
	record.Tags = getValue("Tags")
	record.AdditionalInfo = getValue("AdditionalInfo")

	if val, err := strconv.ParseFloat(getValue("EffectiveUnitPrice"), 64); err == nil {
		record.EffectiveUnitPrice = val
	} else if getValue("EffectiveUnitPrice") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter EffectiveUnitPrice: %s", warning, err)
	}

	if val, err := strconv.ParseFloat(getValue("PCToBCExchangeRate"), 64); err == nil {
		record.PCToBCExchangeRate = val
	} else if getValue("PCToBCExchangeRate") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter PCToBCExchangeRate: %s", warning, err)
	}

	if date, err := parseDate(getValue("PCToBCExchangeRateDate")); err == nil {
		record.PCToBCExchangeRateDate = date
	} else if getValue("PCToBCExchangeRateDate") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter PCToBCExchangeRateDate: %s", warning, err.Error())
	}

	record.EntitlementId = stringPtr(cleanUUID(getValue("EntitlementId")))
	record.EntitlementDescription = getValue("EntitlementDescription")

	if val, err := strconv.ParseFloat(getValue("PartnerEarnedCreditPercentage"), 64); err == nil {
		record.PartnerEarnedCreditPercentage = val
	} else if getValue("PartnerEarnedCreditPercentage") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter PartnerEarnedCreditPercentage: %s", warning, err)
	}

	if val, err := strconv.ParseFloat(getValue("CreditPercentage"), 64); err == nil {
		record.CreditPercentage = val
	} else if getValue("CreditPercentage") != "" {
		warning = fmt.Sprintf("%s; Erro ao converter CreditPercentage: %s", warning, err)
	}

	record.CreditType = getValue("CreditType")
	record.BenefitOrderId = stringPtr(cleanUUID(getValue("BenefitOrderId")))
	record.BenefitId = stringPtr(cleanUUID(getValue("BenefitId")))
	record.BenefitType = getValue("BenefitType")

	return record, partnerData, customerData, productData, skuData, publisherData, subscriptionData, meterData, warning
}

// Processar e salvar entidades no banco
func processAndSaveEntities(db *gorm.DB, partnersData []PartnerData, customersData []CustomerData, productsData []ProductData, skusData []SkuData, publishersData []PublisherData, subscriptionsData []SubscriptionData, metersData []MeterData) (map[string]string, map[string]string, map[string]string, map[string]string, map[string]string, map[string]string, map[string]string, error) {
	log.Println("--- Iniciando processAndSaveEntities ---")
	if db == nil {
		log.Println("Erro: DB é nil")
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("objeto DB é nil em processAndSaveEntities")
	}

	// Iniciar transação
	tx := db.Begin()
	if tx == nil || tx.Error != nil {
		if tx != nil {
			log.Printf("Erro ao iniciar transação: %v", tx.Error)
		}
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao iniciar transação para entidades")
	}

	defer func() {
		if tx != nil {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			} else if tx.Error != nil {
				tx.Rollback()
			}
		}
	}()

	// Processar Partners
	log.Println("Processando Partners...")
	uniquePartners := make(map[string]PartnerData)
	for _, pData := range partnersData {
		if pData.PartnerId != "" {
			uniquePartners[pData.PartnerId] = pData
		}
	}

	partnersToSave := []models.Partner{}
	originalPartnerIDs := []string{}
	for _, pData := range uniquePartners {
		partnersToSave = append(partnersToSave, models.Partner{ID: uuid.New().String(), PartnerId: pData.PartnerId, PartnerName: pData.PartnerName})
		originalPartnerIDs = append(originalPartnerIDs, pData.PartnerId)
	}

	// Salvar partners em lote (ON CONFLICT DO NOTHING)
	if len(partnersToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(partnersToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar partners: %v", err)
		}
	}

	// Consultar partners salvos para obter IDs primários
	var savedPartners []models.Partner
	if len(originalPartnerIDs) > 0 {
		if err := tx.Where("partner_id IN ?", originalPartnerIDs).Find(&savedPartners).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar partners: %v", err)
		}
	}

	// Criar mapa: partner_original_id -> partner_primary_key_id
	partnerIDMap := make(map[string]string)
	for _, p := range savedPartners {
		partnerIDMap[p.PartnerId] = p.ID
	}

	// Processar Customers
	log.Println("Processando Customers...")
	uniqueCustomers := make(map[string]CustomerData)
	for _, cData := range customersData {
		if cData.CustomerId != "" {
			uniqueCustomers[cData.CustomerId] = cData
		}
	}
	customersToSave := []models.Customer{}
	originalCustomerIDs := []string{}
	for _, cData := range uniqueCustomers {
		customersToSave = append(customersToSave, models.Customer{ID: uuid.New().String(), CustomerId: cData.CustomerId, CustomerName: cData.CustomerName, CustomerDomainName: cData.CustomerDomainName, CustomerCountry: cData.CustomerCountry})
		originalCustomerIDs = append(originalCustomerIDs, cData.CustomerId)
	}

	// Salvar customers em lote (ON CONFLICT DO NOTHING)
	if len(customersToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(customersToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar customers: %v", err)
		}
	}

	// Consultar customers salvos para obter IDs primários
	var savedCustomers []models.Customer
	if len(originalCustomerIDs) > 0 {
		if err := tx.Where("customer_id IN ?", originalCustomerIDs).Find(&savedCustomers).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar customers: %v", err)
		}
	}

	// Criar mapa: customer_original_id -> customer_primary_key_id
	customerIDMap := make(map[string]string)
	for _, c := range savedCustomers {
		customerIDMap[c.CustomerId] = c.ID
	}

	// Processar Products
	log.Println("Processando Products...")
	uniqueProducts := make(map[string]ProductData)
	for _, pData := range productsData {
		if pData.ProductId != "" {
			uniqueProducts[pData.ProductId] = pData
		}
	}
	productsToSave := []models.Product{}
	originalProductIDs := []string{}
	for _, pData := range uniqueProducts {
		productsToSave = append(productsToSave, models.Product{ID: uuid.New().String(), ProductId: pData.ProductId, ProductName: pData.ProductName})
		originalProductIDs = append(originalProductIDs, pData.ProductId)
	}

	// Salvar products em lote (ON CONFLICT DO NOTHING)
	if len(productsToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(productsToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar products: %v", err)
		}
	}
	// Consultar products salvos para obter IDs primários
	var savedProducts []models.Product
	if len(originalProductIDs) > 0 {
		if err := tx.Where("product_id IN ?", originalProductIDs).Find(&savedProducts).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar products: %v", err)
		}
	}
	// Criar mapa: product_original_id -> product_primary_key_id
	productIDMap := make(map[string]string)
	for _, p := range savedProducts {
		productIDMap[p.ProductId] = p.ID
	}

	// Processar Skus
	log.Println("Processando Skus...")
	uniqueSkus := make(map[string]SkuData)
	for _, sData := range skusData {
		if sData.SkuId != "" {
			uniqueSkus[sData.SkuId] = sData
		}
	}
	skusToSave := []models.Sku{}
	originalSkuIDs := []string{}
	for _, sData := range uniqueSkus {
		skusToSave = append(skusToSave, models.Sku{ID: uuid.New().String(), SkuId: sData.SkuId, SkuName: sData.SkuName})
		originalSkuIDs = append(originalSkuIDs, sData.SkuId)
	}

	// Salvar skus em lote (ON CONFLICT DO NOTHING)
	if len(skusToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(skusToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar skus: %v", err)
		}
	}
	// Consultar skus salvos para obter IDs primários
	var savedSkus []models.Sku
	if len(originalSkuIDs) > 0 {
		if err := tx.Where("sku_id IN ?", originalSkuIDs).Find(&savedSkus).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar skus: %v", err)
		}
	}
	// Criar mapa: sku_original_id -> sku_primary_key_id
	skuIDMap := make(map[string]string)
	for _, s := range savedSkus {
		skuIDMap[s.SkuId] = s.ID
	}

	// Processar Publishers
	log.Println("Processando Publishers...")
	uniquePublishers := make(map[string]PublisherData)
	for _, pData := range publishersData {
		if pData.PublisherId != "" {
			uniquePublishers[pData.PublisherId] = pData
		}
	}
	publishersToSave := []models.Publisher{}
	originalPublisherIDs := []string{}
	for _, pData := range uniquePublishers {
		publishersToSave = append(publishersToSave, models.Publisher{ID: uuid.New().String(), PublisherId: pData.PublisherId, PublisherName: pData.PublisherName})
		originalPublisherIDs = append(originalPublisherIDs, pData.PublisherId)
	}

	// Salvar publishers em lote (ON CONFLICT DO NOTHING)
	if len(publishersToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(publishersToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar publishers: %v", err)
		}
	}
	// Consultar publishers salvos para obter IDs primários
	var savedPublishers []models.Publisher
	if len(originalPublisherIDs) > 0 {
		if err := tx.Where("publisher_id IN ?", originalPublisherIDs).Find(&savedPublishers).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar publishers: %v", err)
		}
	}
	// Criar mapa: publisher_original_id -> publisher_primary_key_id
	publisherIDMap := make(map[string]string)
	for _, p := range savedPublishers {
		publisherIDMap[p.PublisherId] = p.ID
	}

	// Processar Subscriptions
	log.Println("Processando Subscriptions...")
	uniqueSubscriptions := make(map[string]SubscriptionData)
	for _, sData := range subscriptionsData {
		if sData.SubscriptionId != "" {
			uniqueSubscriptions[sData.SubscriptionId] = sData
		}
	}
	subscriptionsToSave := []models.Subscription{}
	originalSubscriptionIDs := []string{}
	for _, sData := range uniqueSubscriptions {
		subscriptionsToSave = append(subscriptionsToSave, models.Subscription{ID: uuid.New().String(), SubscriptionId: sData.SubscriptionId, SubscriptionDescription: sData.SubscriptionDescription})
		originalSubscriptionIDs = append(originalSubscriptionIDs, sData.SubscriptionId)
	}

	// Salvar subscriptions em lote (ON CONFLICT DO NOTHING)
	if len(subscriptionsToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(subscriptionsToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar subscriptions: %v", err)
		}
	}
	// Consultar subscriptions salvos para obter IDs primários
	var savedSubscriptions []models.Subscription
	if len(originalSubscriptionIDs) > 0 {
		if err := tx.Where("subscription_id IN ?", originalSubscriptionIDs).Find(&savedSubscriptions).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar subscriptions: %v", err)
		}
	}
	// Criar mapa: subscription_original_id -> subscription_primary_key_id
	subscriptionIDMap := make(map[string]string)
	for _, s := range savedSubscriptions {
		subscriptionIDMap[s.SubscriptionId] = s.ID
	}

	// Processar Meters
	log.Println("Processando Meters...")
	uniqueMeters := make(map[string]MeterData)
	for _, mData := range metersData {
		if mData.MeterId != "" {
			uniqueMeters[mData.MeterId] = mData
		}
	}
	metersToSave := []models.Meter{}
	originalMeterIDs := []string{}
	for _, mData := range uniqueMeters {
		metersToSave = append(metersToSave, models.Meter{ID: uuid.New().String(), MeterId: mData.MeterId, MeterType: mData.MeterType, MeterCategory: mData.MeterCategory, MeterSubCategory: mData.MeterSubCategory, MeterName: mData.MeterName, MeterRegion: mData.MeterRegion})
		originalMeterIDs = append(originalMeterIDs, mData.MeterId)
	}

	// Salvar meters em lote (ON CONFLICT DO NOTHING)
	if len(metersToSave) > 0 {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(metersToSave, 100).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao salvar meters: %v", err)
		}
	}
	// Consultar meters salvos para obter IDs primários
	var savedMeters []models.Meter
	if len(originalMeterIDs) > 0 {
		if err := tx.Where("meter_id IN ?", originalMeterIDs).Find(&savedMeters).Error; err != nil {
			return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao consultar meters: %v", err)
		}
	}
	// Criar mapa: meter_original_id -> meter_primary_key_id
	meterIDMap := make(map[string]string)
	for _, m := range savedMeters {
		meterIDMap[m.MeterId] = m.ID
	}

	// Commitar transação
	if err := tx.Commit().Error; err != nil {
		return nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("erro ao commitar transação: %v", err)
	}

	return partnerIDMap, customerIDMap, productIDMap, skuIDMap, publisherIDMap, subscriptionIDMap, meterIDMap, nil
}

// Preencher IDs de chave estrangeira nos registros
func populateBillingRecordsWithFKs(records []models.BillingRecord, partnersData []PartnerData, customersData []CustomerData, productsData []ProductData, skusData []SkuData, publishersData []PublisherData, subscriptionsData []SubscriptionData, metersData []MeterData, partnerIDMap map[string]string, customerIDMap map[string]string, productIDMap map[string]string, skuIDMap map[string]string, publisherIDMap map[string]string, subscriptionIDMap map[string]string, meterIDMap map[string]string) []models.BillingRecord {
	log.Println("--- Preenchendo IDs de FK nos registros ---")
	// Iterar sobre registros e preencher IDs de FK
	for i := range records {
		// Preencher PartnerID
		if i < len(partnersData) {
			if partnerID, ok := partnerIDMap[partnersData[i].PartnerId]; ok && partnerID != "" {
				records[i].PartnerID = &partnerID
			} else {
				records[i].PartnerID = nil
			}
		} else {
			records[i].PartnerID = nil
		}

		// Preencher CustomerID
		if i < len(customersData) {
			if customerID, ok := customerIDMap[customersData[i].CustomerId]; ok && customerID != "" {
				records[i].CustomerID = &customerID
			} else {
				records[i].CustomerID = nil
			}
		}

		// Preencher ProductID
		if i < len(productsData) {
			if productID, ok := productIDMap[productsData[i].ProductId]; ok && productID != "" {
				records[i].ProductID = &productID
			} else {
				records[i].ProductID = nil
			}
		}

		// Preencher SkuID
		if i < len(skusData) {
			if skuID, ok := skuIDMap[skusData[i].SkuId]; ok && skuID != "" {
				records[i].SkuID = &skuID
			} else {
				records[i].SkuID = nil
			}
		}

		// Preencher PublisherID
		if i < len(publishersData) {
			if publisherID, ok := publisherIDMap[publishersData[i].PublisherId]; ok && publisherID != "" {
				records[i].PublisherID = &publisherID
			} else {
				records[i].PublisherID = nil
			}
		}

		// Preencher SubscriptionID
		if i < len(subscriptionsData) {
			if subscriptionID, ok := subscriptionIDMap[subscriptionsData[i].SubscriptionId]; ok && subscriptionID != "" {
				records[i].SubscriptionID = &subscriptionID
			} else {
				records[i].SubscriptionID = nil
			}
		}

		// Preencher MeterID
		if i < len(metersData) {
			if meterID, ok := meterIDMap[metersData[i].MeterId]; ok && meterID != "" {
				records[i].MeterID = &meterID
			} else {
				records[i].MeterID = nil
			}
		}
	}

	return records
}

// Salvar registros no banco
func SaveBillingRecords(db *gorm.DB, records []models.BillingRecord) error {
	// Verificar se há registros para salvar
	if len(records) == 0 {
		return nil
	}

	// Iniciar transação
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", tx.Error)
	}

	// Salvar em lotes
	batchSize := 100
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}

		batch := records[i:end]
		// Criar registros no banco
		if err := tx.Create(&batch).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao salvar lote de registros: %v", err)
		}
	}

	// Commitar transação
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}

// Limpar caracteres suspeitos
func Sanitize(s string) string {
	// Remover caracteres indesejados
	toRemove := []string{"'", ";", "\\"}
	for _, r := range toRemove {
		s = strings.ReplaceAll(s, r, "")
	}
	return s
}

// Limpar e validar UUID
func cleanUUID(input string) string {
	// Remover espaços
	cleaned := strings.ReplaceAll(input, " ", "")
	// Validar formato UUID
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(cleaned) {
		return ""
	}
	return cleaned
}

// Criar ponteiro para string se não for vazia
func stringPtr(s string) *string {
	// Retornar nil se string vazia, senão retornar ponteiro
	if s == "" {
		return nil
	}
	return &s
}

// Converter datas com vários formatos possíveis
func parseDate(dateStr string) (time.Time, error) {
	// Tentar analisar a data com vários layouts
	layouts := []string{
		"2006-01-02",
		"01-02-06",
		"01/02/2006",
		"2006/01/02",
		"02-01-2006",
		"02/01/2006",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	// Retornar erro se nenhum layout funcionar
	return time.Time{}, fmt.Errorf("não foi possível analisar a data: %s", dateStr)
}
