// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack

package routes

import (
	"backend/internal/auth"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter configura todas as rotas da aplicação
func SetupRouter(db *gorm.DB, jwtService *auth.JWTService) *gin.Engine {
	router := gin.Default()

	// Inicializar handlers
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

	// Rotas públicas
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	// Grupo de rotas protegidas
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtService))
	{
		// Rota de importação
		api.POST("/import", handlers.ImportXLSHandler(db))

		// Rotas de Listagem
		api.GET("/clients", handlers.GetAllClientsHandler(db))
		api.GET("/categories", handlers.GetAllCategoriesHandler(db))
		api.GET("/resources", handlers.GetAllResourcesHandler(db))
		api.GET("/users", handlers.GetAllUsersHandler(db))

		// Rotas de Faturamento
		api.GET("/billing", handlers.GetAllBillingHandler(db))

		// Rotas de agrupamento
		api.GET("/billing/summary/categories", handlers.GetBillingSummaryByCategoriesHandler(db))
		api.GET("/billing/summary/resources", handlers.GetBillingSummaryByResources(db))
		api.GET("/billing/summary/clients", handlers.GetBillingSummaryByClients(db))
		api.GET("/billing/summary/months", handlers.GetBillingSummaryByMonths(db))
	}

	return router
}
