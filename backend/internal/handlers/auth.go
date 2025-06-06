package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.Register(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{Token: token})
}

// GetAllUsersHandler retorna todos os usuários do sistema
func GetAllUsersHandler(db *gorm.DB) gin.HandlerFunc {
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

		var users []map[string]interface{}
		var total int64

		// Contar total de registros
		if err := db.Table("users").Count(&total).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar usuários"})
			return
		}

		// Buscar usuários (excluindo a senha)
		if err := db.Table("users").
			Select("id, username, email, active, created_at, updated_at").
			Scopes(Paginate(pageNum, limitNum)).
			Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar usuários"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": users,
			"pagination": gin.H{
				"total":       total,
				"page":        pageNum,
				"limit":       limitNum,
				"total_pages": int(math.Ceil(float64(total) / float64(limitNum))),
			},
		})
	}
}
