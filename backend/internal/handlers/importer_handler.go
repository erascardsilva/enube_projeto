package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

// ImportXLSHandler lida com a importação de arquivos Excel
func ImportXLSHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obter o arquivo do request
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nenhum arquivo enviado"})
			return
		}

		// Abrir o arquivo
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao abrir arquivo"})
			return
		}
		defer src.Close()

		// Ler o arquivo Excel
		f, err := excelize.OpenReader(src)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo Excel inválido"})
			return
		}
		defer f.Close()

		// Obter a primeira planilha
		sheet := f.GetSheetList()[0]
		rows, err := f.GetRows(sheet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler planilha"})
			return
		}

		// Verificar se há dados (pelo menos 2 linhas para ter cabeçalho e dados)
		if len(rows) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Planilha vazia ou com poucas linhas"})
			return
		}

		// Iniciar transação
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar transação"})
			return
		}

		// Contadores para estatísticas
		stats := struct {
			TotalRows    int
			ImportedRows int
			SkippedRows  int
			HeaderRows   int
			Categories   map[string]int
		}{
			Categories: make(map[string]int),
		}

		// Definir os índices das colunas relevantes (baseado na ANÁLISE DO CABEÇALHO COMPLETO fornecido)
		const (
			chargeStartDateIndex  = 18 // Coluna S (Data de Uso/Faturamento)
			customerIdIndex       = 2  // Coluna C (ID do Cliente do Excel)
			customerNameIndex     = 3  // Coluna D (Nome do Cliente do Excel)
			customerDomainIndex   = 4  // Coluna E (Domain Name)
			customerCountryIndex  = 5  // Coluna F (Country)
			meterCategoryIndex    = 22 // Coluna W (Categoria)
			meterTypeIndex        = 21 // Coluna V (Meter Type)
			meterSubCatIndex      = 24 // Coluna Y (Descrição)
			billingAmountIndex    = 36 // Coluna AK (Valor)
			resourceURIIndex      = 31 // Coluna AF (URI do Recurso)
			resourceLocationIndex = 29 // Coluna AD (Resource Location)
			resourceGroupIndex    = 30 // Coluna AE (Resource Group)
			consumedServiceIndex  = 28 // Coluna AC (Consumed Service)
			unitIndex             = 25 // Coluna Z (Unit)
			unitPriceIndex        = 33 // Coluna AG (Unit Price)
			quantityIndex         = 34 // Coluna AH (Quantity)
			unitTypeIndex         = 26 // Coluna AA (Unit Type)
			billingCurrencyIndex  = 37 // Coluna AL (Billing Currency)
			pricingPreTaxIndex    = 38 // Coluna AM (Pricing Pre-Tax Total)
			pricingCurrencyIndex  = 39 // Coluna AN (Pricing Currency)
		)

		// Verificar se o arquivo tem colunas suficientes (baseado no maior índice necessário)
		requiredColumns := max(chargeStartDateIndex, max(customerIdIndex, max(customerNameIndex, max(customerDomainIndex, max(customerCountryIndex, max(meterCategoryIndex, max(meterTypeIndex, max(meterSubCatIndex, max(billingAmountIndex, max(resourceURIIndex, max(resourceLocationIndex, max(resourceGroupIndex, max(consumedServiceIndex, max(unitIndex, max(unitPriceIndex, max(quantityIndex, max(unitTypeIndex, max(billingCurrencyIndex, max(pricingPreTaxIndex, pricingCurrencyIndex))))))))))))))))))) + 1

		if len(rows[0]) < requiredColumns {
			log.Printf("Cabeçalho com dados insuficientes. Esperado pelo menos %d colunas, encontrado %d.", requiredColumns, len(rows[0]))
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Arquivo com colunas insuficientes. Esperado pelo menos %d colunas.", requiredColumns)})
			return
		}

		// --- Lógica para encontrar a primeira linha de dados ---
		firstDataRowIndex := -1
		for i := 0; i < len(rows); i++ {
			row := rows[i]
			// Verificar se a linha tem colunas suficientes e se campos chave parecem válidos (não vazios)
			if len(row) > billingAmountIndex &&
				strings.TrimSpace(row[customerIdIndex]) != "" &&
				strings.TrimSpace(row[meterCategoryIndex]) != "" &&
				strings.TrimSpace(row[resourceURIIndex]) != "" &&
				strings.TrimSpace(row[chargeStartDateIndex]) != "" &&
				strings.TrimSpace(row[billingAmountIndex]) != "" {

				// Tentativa de parsear a data para uma verificação mais robusta
				dateStrCheck := strings.TrimSpace(row[chargeStartDateIndex])
				var tempDateParseErr error
				formatsToTryCheck := []string{"02/01/2006", "01/02/2006", "2006-01-02", "02-01-2006", "01-02-2006", "2006/01/02", "01-02-06", "02-01-06"}

				for _, format := range formatsToTryCheck {
					_, tempDateParseErr = time.Parse(format, dateStrCheck)
					if tempDateParseErr == nil {
						break // Data parece válida, consideramos como linha de dados
					}
				}

				if tempDateParseErr == nil { // Se a data parece válida, esta é provavelmente a primeira linha de dados
					firstDataRowIndex = i
					stats.HeaderRows = i // Contar as linhas antes como cabeçalho
					log.Printf("Primeira linha de dados detectada na linha %d (índice %d).", i+1, i)
					break
				}
			}
		}

		// Se não encontrou nenhuma linha que pareça ser de dados
		if firstDataRowIndex == -1 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Não foi possível encontrar a primeira linha de dados válida no arquivo."})
			return
		}

		// Processar cada linha a partir da primeira linha de dados detectada
		for i := firstDataRowIndex; i < len(rows); i++ {
			row := rows[i]

			stats.TotalRows++ // Contabiliza as linhas processadas (a partir da primeira linha de dados)

			// Verificar se a linha tem dados suficientes para todas as colunas necessárias
			if len(row) < requiredColumns {
				log.Printf("Linha %d ignorada APÓS DETECÇÃO: dados insuficientes. Esperado pelo menos %d colunas, encontrado %d.", i+1, requiredColumns, len(row))
				stats.SkippedRows++
				continue
			}

			// Extrair dados das colunas
			dateStr := strings.TrimSpace(row[chargeStartDateIndex])
			customerIDStr := strings.TrimSpace(row[customerIdIndex])
			customerNameStr := strings.TrimSpace(row[customerNameIndex])
			customerDomainStr := strings.TrimSpace(row[customerDomainIndex])
			customerCountryStr := strings.TrimSpace(row[customerCountryIndex])
			categoryNameStr := strings.TrimSpace(row[meterCategoryIndex])
			meterTypeStr := strings.TrimSpace(row[meterTypeIndex])
			descriptionStr := strings.TrimSpace(row[meterSubCatIndex])
			amountStrRaw := strings.TrimSpace(row[billingAmountIndex])
			resourceURIStr := strings.TrimSpace(row[resourceURIIndex])
			resourceLocationStr := strings.TrimSpace(row[resourceLocationIndex])
			resourceGroupStr := strings.TrimSpace(row[resourceGroupIndex])
			consumedServiceStr := strings.TrimSpace(row[consumedServiceIndex])
			unitStr := strings.TrimSpace(row[unitIndex])
			unitPriceStr := strings.TrimSpace(row[unitPriceIndex])
			quantityStr := strings.TrimSpace(row[quantityIndex])
			unitTypeStr := strings.TrimSpace(row[unitTypeIndex])
			billingCurrencyStr := strings.TrimSpace(row[billingCurrencyIndex])
			pricingPreTaxStr := strings.TrimSpace(row[pricingPreTaxIndex])
			pricingCurrencyStr := strings.TrimSpace(row[pricingCurrencyIndex])

			// Validar se os campos chave essenciais não estão vazios APÓS DETECÇÃO
			// Esta validação ainda é útil para linhas em branco ou incompletas DENTRO dos dados
			if customerIDStr == "" || categoryNameStr == "" || resourceURIStr == "" || dateStr == "" || amountStrRaw == "" {
				log.Printf("Linha %d ignorada APÓS DETECÇÃO: campos chave essenciais vazios. Cliente ID: '%s', Categoria: '%s', Recurso URI: '%s', Data: '%s', Valor: '%s'.", i+1, customerIDStr, categoryNameStr, resourceURIStr, dateStr, amountStrRaw)
				stats.SkippedRows++
				continue
			}

			// Converter data (tentar formatos comuns)
			var billingDate time.Time
			var dateParseErr error

			formatsToTry := []string{
				"02/01/2006", // DD/MM/YYYY
				"01/02/2006", // MM/DD/YYYY
				"2006-01-02", // YYYY-MM-DD
				"02-01-2006", // DD-MM-YYYY
				"01-02-2006", // MM-DD-YYYY
				"2006/01/02", // YYYY/MM/DD
				"01-02-06",   // MM-DD-YY
				"02-01-06",   // DD-MM-YY
			}

			for _, format := range formatsToTry {
				billingDate, dateParseErr = time.Parse(format, dateStr)
				if dateParseErr == nil {
					break // Parar no primeiro formato que funcionar
				}
			}

			if dateParseErr != nil {
				log.Printf("Erro ao converter data '%s' na linha %d APÓS DETECÇÃO. Tentados formatos: %v. Erro: %v", dateStr, i+1, formatsToTry, dateParseErr)
				stats.SkippedRows++
				continue
			}

			// Converter valor (substituir vírgula por ponto e converter para float64). Tratar erro sem pular a linha.
			amountStr := strings.ReplaceAll(amountStrRaw, ",", ".")
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				log.Printf("Aviso: Erro ao converter valor '%s' na linha %d APÓS DETECÇÃO: %v. Definindo valor para 0.0.", amountStrRaw, i+1, err)
				amount = 0.0 // Definir valor como zero se a conversão falhar
			}

			// --- Lógica de UPSERT e obtenção de IDs ---

			// Upsert Cliente (usando CustomerID como chave única e CustomerName para nome se disponível)
			var client models.Client
			if err := tx.Where("client_id = ?", customerIDStr).First(&client).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					client = models.Client{
						ClientID:   customerIDStr,
						Name:       customerNameStr,
						DomainName: customerDomainStr,
						Country:    customerCountryStr,
					}
					if err := tx.Create(&client).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar cliente '%s' na linha %d: %v", customerIDStr, i+1, err)})
						return
					}
				} else {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar cliente '%s' na linha %d: %v", customerIDStr, i+1, err)})
					return
				}
			}

			// Upsert Categoria (usando MeterCategory como chave única)
			var category models.Category
			if err := tx.Where("name = ?", categoryNameStr).First(&category).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					category = models.Category{
						Name:        categoryNameStr,
						Type:        meterTypeStr,
						SubCategory: descriptionStr,
					}
					if err := tx.Create(&category).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar categoria '%s' na linha %d: %v", categoryNameStr, i+1, err)})
						return
					}
				} else {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar categoria '%s' na linha %d: %v", categoryNameStr, i+1, err)})
					return
				}
			}

			// Upsert Recurso (usando ResourceURI como chave única)
			var resource models.Resource
			if err := tx.Where("name = ?", resourceURIStr).First(&resource).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					resource = models.Resource{
						Name:            resourceURIStr,
						Location:        resourceLocationStr,
						Group:           resourceGroupStr,
						ConsumedService: consumedServiceStr,
					}
					if err := tx.Create(&resource).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao criar recurso '%s' na linha %d: %v", resourceURIStr, i+1, err)})
						return
					}
				} else {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao buscar recurso '%s' na linha %d: %v", resourceURIStr, i+1, err)})
					return
				}
			}

			// Converter valores numéricos
			unitPrice, _ := strconv.ParseFloat(strings.ReplaceAll(unitPriceStr, ",", "."), 64)
			quantity, _ := strconv.ParseFloat(strings.ReplaceAll(quantityStr, ",", "."), 64)
			pricingPreTax, _ := strconv.ParseFloat(strings.ReplaceAll(pricingPreTaxStr, ",", "."), 64)

			// Criar registro de faturamento
			billing := models.Billing{
				ClientID:           uint(client.ID),
				CategoryID:         uint(category.ID),
				ResourceID:         uint(resource.ID),
				Amount:             amount,
				BillingDate:        billingDate,
				Description:        descriptionStr,
				Year:               billingDate.Year(),
				Month:              int(billingDate.Month()),
				Day:                billingDate.Day(),
				Unit:               unitStr,
				UnitPrice:          unitPrice,
				Quantity:           quantity,
				UnitType:           unitTypeStr,
				BillingCurrency:    billingCurrencyStr,
				PricingPreTaxTotal: pricingPreTax,
				PricingCurrency:    pricingCurrencyStr,
			}

			// Inserir registro de billing no banco
			if err := tx.Create(&billing).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erro ao salvar registro de billing na linha %d: %v", i+1, err)})
				return
			}

			stats.ImportedRows++
			stats.Categories[categoryNameStr]++ // Contar por nome da categoria do Excel
		}

		// Commit da transação
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao finalizar importação"})
			return
		}

		// Preparar resumo da importação
		summary := map[string]interface{}{
			"message": "Importação concluída com sucesso",
			"stats": map[string]interface{}{
				"total_rows":    stats.TotalRows,
				"imported_rows": stats.ImportedRows,
				"skipped_rows":  stats.SkippedRows,
				"categories":    stats.Categories,
				"header_rows":   stats.HeaderRows,
			},
		}

		// Não buscar dados importados na resposta para evitar sobrecarga com arquivos grandes
		c.JSON(http.StatusOK, summary)
	}
}

// max retorna o maior de dois inteiros
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
