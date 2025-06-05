package handlers

import (
	"backend/internal/importer"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

// Handler para importar o arquivo XLS
func ImportXLSHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar se o objeto DB é válido
		log.Println("--- Handler de Importação Acionado ---")
		if db == nil {
			log.Println("Erro no Handler: Objeto DB capturado é nil.")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor: objeto DB não disponível"})
			return
		}
		log.Println("Handler: Objeto DB capturado não é nil.")

		// Obter o arquivo do formulário
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
			return
		}

		// Salvar arquivo temporariamente
		filePath := "/tmp/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar arquivo"})
			return
		}

		// Processar o arquivo de Billing
		recordsSalvos, avisosImportacao, missingColumns, err := importer.ImportBillingXLS(db, filePath)
		if err != nil {
			// Em caso de erro no processamento, retornar status 400
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Erro ao processar arquivo de billing: %v", err)})
			return
		}

		// Retornar sucesso com resumo da importação
		c.JSON(http.StatusOK, gin.H{
			"status":            "sucesso",
			"mensagem":          fmt.Sprintf("Importação concluída com sucesso. %d registros de billing processados e salvos.", len(recordsSalvos)),
			"avisos":            avisosImportacao,
			"colunas_faltantes": missingColumns,
		})
	}
}
