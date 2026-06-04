package controllers

import (
	"net/http"

	"golden-ticket/backend/domain"
	"golden-ticket/backend/services"

	"github.com/gin-gonic/gin"
)

// AuthController handles the HTTP requests for authentication
type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new instance of AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration (POST /register)
func (ctrl *AuthController) Register(c *gin.Context) {
	var dto domain.UserRegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctrl.authService.Register(dto)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "El email ingresado ya se encuentra registrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Login handles user login (POST /login)
func (ctrl *AuthController) Login(c *gin.Context) {
	var dto domain.UserLoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctrl.authService.Login(dto)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
		return
	}

	c.JSON(http.StatusOK, res)
}
