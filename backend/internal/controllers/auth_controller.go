package controllers

import (
	"net/http"
	"regexp"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService services.AuthService
	Validate    *validator.Validate
}

func NewAuthController(authService services.AuthService) *AuthController {
	v := validator.New()

	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}
		re := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
		return re.MatchString(password)
	})

	return &AuthController{authService, v}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.authService.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration successful"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var creds struct {
		Identifier string `json:"identifier" binding:"required"` // Accept both email and username.
		Password   string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := ctrl.authService.Login(creds.Identifier, creds.Password)
	if err != nil {
		if err.Error() == "email not verified" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email before logging in."})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ctrl *AuthController) VerifyEmail(c *gin.Context) {
	var input struct {
		Username         string `json:"username" binding:"required"`
		VerificationCode string `json:"verification_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.authService.VerifyEmail(input.Username, input.VerificationCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
