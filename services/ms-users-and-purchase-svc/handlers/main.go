package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/config"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/database"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/dto"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIEnv struct {
	DB *gorm.DB
}

// Utils
func validateURIWithTLD(uri string) bool {
	parsedURI, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	return strings.Contains(parsedURI.Host, ".")
}

// Services
func (a *APIEnv) LoginWithEmail(c *gin.Context) {
	db := a.DB
	cfg := config.EnvironmentConfig()
	var user database.User
	var userRequest dto.UserCreateOrLoginWithEmailRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Where("email =?", userRequest.Email).First(&user)
	if user.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"userId": user.ID,
	})
	tokenString, err := token.SignedString([]byte(
		cfg.JWT_SECRET,
	))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"phone": user.Phone,
		"token": tokenString,
	})
	return
}

func (a *APIEnv) RegisterWithEmail(c *gin.Context) {
	db := a.DB
	cfg := config.EnvironmentConfig()
	var userRequest dto.UserCreateOrLoginWithEmailRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := database.User{Email: userRequest.Email, Password: userRequest.Password}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)
	result := db.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": result.Error,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"userId": user.ID,
	})

	tokenString, err := token.SignedString([]byte(cfg.JWT_SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"email": user.Email,
		"token": tokenString,
	})
	return
}

func (a *APIEnv) GetUser(c *gin.Context) {
	db := a.DB
	id := c.GetString("userId")

	var user database.User
	db.Where("id = ?", id).First(&user)

	c.JSON(http.StatusOK, gin.H{
		"bankAccountHolder": user.BankAccountHolder,
		"bankAccountName":   user.BankAccountName,
		"bankAccountNumber": user.BankAccountNumber,
		"email":             user.Email,
		"fileId":            user.FileId,
		"fileThumbnailUri":  user.FileThumbnailUri,
		"phone":             user.Phone,
	})
}

func (a *APIEnv) UpdateUser(c *gin.Context) {
	db := a.DB

	var userRequest dto.UserUpdateRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.GetString("userId")

	var user database.User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// To update
	// user.Email = ...

	if err := db.Save(&user).Error; err != nil {
		if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"bankAccountHolder": user.BankAccountHolder,
		"bankAccountName":   user.BankAccountName,
		"bankAccountNumber": user.BankAccountNumber,
		"email":             user.Email,
		"fileId":            user.FileId,
		"fileThumbnailUri":  user.FileThumbnailUri,
		"phone":             user.Phone,
	})
}
