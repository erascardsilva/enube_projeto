// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package handlers

import (
	"backend/internal/models"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Scopes para queries comuns
func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// GetAllClientsHandler retorna todos os clientes com paginação
func GetAllClientsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "50")

		pageNum, _ := strconv.Atoi(page)
		limitNum, _ := strconv.Atoi(limit)

		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 || limitNum > 100 {
			limitNum = 50
		}

		var clients []map[string]interface{}
		var total int64

		if err := db.Table("clients").Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar registros"})
			return
		}

		if err := db.Table("clients").
			Select("id, client_id, name, domain_name, country, created_at, updated_at").
			Scopes(Paginate(pageNum, limitNum)).
			Find(&clients).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar clientes"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": clients,
			"pagination": gin.H{
				"total":       total,
				"page":        pageNum,
				"limit":       limitNum,
				"total_pages": int(math.Ceil(float64(total) / float64(limitNum))),
			},
		})
	}
}

// GetAllCategoriesHandler retorna todas as categorias com paginação
func GetAllCategoriesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "50")

		pageNum, _ := strconv.Atoi(page)
		limitNum, _ := strconv.Atoi(limit)

		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 || limitNum > 100 {
			limitNum = 50
		}

		var categories []map[string]interface{}
		var total int64

		if err := db.Table("categories").Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar registros"})
			return
		}

		if err := db.Table("categories").
			Select("id, name, type, sub_category, created_at, updated_at").
			Scopes(Paginate(pageNum, limitNum)).
			Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar categorias"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": categories,
			"pagination": gin.H{
				"total":       total,
				"page":        pageNum,
				"limit":       limitNum,
				"total_pages": int(math.Ceil(float64(total) / float64(limitNum))),
			},
		})
	}
}

// GetAllResourcesHandler retorna todos os recursos com paginação
func GetAllResourcesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "50")

		pageNum, _ := strconv.Atoi(page)
		limitNum, _ := strconv.Atoi(limit)

		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 || limitNum > 100 {
			limitNum = 50
		}

		var resources []models.Resource
		var total int64

		// Primeiro, verificar se a tabela existe
		if !db.Migrator().HasTable(&models.Resource{}) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Tabela de recursos não encontrada"})
			return
		}

		// Contar total de registros
		if err := db.Model(&models.Resource{}).Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar recursos: " + err.Error()})
			return
		}

		// Buscar recursos usando o modelo
		if err := db.Model(&models.Resource{}).
			Scopes(Paginate(pageNum, limitNum)).
			Find(&resources).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar recursos: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": resources,
			"pagination": gin.H{
				"total":       total,
				"page":        pageNum,
				"limit":       limitNum,
				"total_pages": int(math.Ceil(float64(total) / float64(limitNum))),
			},
		})
	}
}

// GetAllBillingHandler retorna todos os faturamentos com paginação
func GetAllBillingHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "50")

		pageNum, _ := strconv.Atoi(page)
		limitNum, _ := strconv.Atoi(limit)

		if pageNum < 1 {
			pageNum = 1
		}
		if limitNum < 1 || limitNum > 100 {
			limitNum = 50
		}

		var billings []map[string]interface{}
		var total int64

		if err := db.Table("billing").Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar registros"})
			return
		}

		// Consulta otimizada usando Preload
		if err := db.Table("billing").
			Select("billing.*, clients.name as client_name, categories.name as category_name, resources.name as resource_name").
			Joins("LEFT JOIN clients ON billing.client_id = clients.id").
			Joins("LEFT JOIN categories ON billing.category_id = categories.id").
			Joins("LEFT JOIN resources ON billing.resource_id = resources.id").
			Scopes(Paginate(pageNum, limitNum)).
			Find(&billings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": billings,
			"pagination": gin.H{
				"total":       total,
				"page":        pageNum,
				"limit":       limitNum,
				"total_pages": int(math.Ceil(float64(total) / float64(limitNum))),
			},
		})
	}
}

// GetBillingSummaryByCategoriesHandler retorna resumo de faturamento por categoria
func GetBillingSummaryByCategoriesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var summaries []struct {
			Category string  `json:"category"`
			Total    float64 `json:"total"`
			Count    int64   `json:"count"`
		}

		if err := db.Table("billing").
			Select("categories.name as category, SUM(billing.amount) as total, COUNT(*) as count").
			Joins("LEFT JOIN categories ON billing.category_id = categories.id").
			Group("categories.name").
			Order("total DESC").
			Find(&summaries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, summaries)
	}
}

// GetBillingSummaryByResources retorna o resumo de faturamento agrupado por recursos
func GetBillingSummaryByResources(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var summaries []struct {
			ResourceName string  `json:"resource_name"`
			Total        float64 `json:"total"`
			Count        int64   `json:"count"`
		}

		if err := db.Table("billing").
			Select("resources.name as resource_name, SUM(billing.amount) as total, COUNT(*) as count").
			Joins("LEFT JOIN resources ON billing.resource_id = resources.id").
			Group("resources.name").
			Order("total DESC").
			Find(&summaries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, summaries)
	}
}

// GetBillingSummaryByClients retorna o resumo de faturamento agrupado por clientes
func GetBillingSummaryByClients(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var summaries []struct {
			ClientName string  `json:"client_name"`
			Total      float64 `json:"total"`
			Count      int64   `json:"count"`
		}

		if err := db.Table("billing").
			Select("clients.name as client_name, SUM(billing.amount) as total, COUNT(*) as count").
			Joins("LEFT JOIN clients ON billing.client_id = clients.id").
			Group("clients.name").
			Order("total DESC").
			Find(&summaries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, summaries)
	}
}

// GetBillingSummaryByMonths retorna o resumo de faturamento agrupado por meses
func GetBillingSummaryByMonths(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var summaries []struct {
			Month time.Time `json:"month"`
			Total float64   `json:"total"`
			Count int64     `json:"count"`
		}

		if err := db.Table("billing").
			Select("DATE_TRUNC('month', billing_date) as month, SUM(amount) as total, COUNT(*) as count").
			Group("DATE_TRUNC('month', billing_date)").
			Order("month DESC").
			Find(&summaries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, summaries)
	}
}
