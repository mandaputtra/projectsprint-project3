package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/config"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Type
type UserCreateOrLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserUpdateRequest struct {
	Email      string  `json:"email"`
	Preference string  `json:"preference" binding:"required,oneof=CARDIO WEIGHT"`
	WeightUnit string  `json:"weightUnit" binding:"required,oneof=KG LBS"`
	HeightUnit string  `json:"heightUnit" binding:"required,oneof=CM INCH"`
	Height     float64 `json:"height" binding:"required,min=3,max=250"`
	Weight     float64 `json:"weight" binding:"required,min=10,max=1000"`
	Name       string  `json:"name" binding:"omitempty,min=2,max=60"`
	ImageUri   string  `json:"imageUri" binding:"omitempty,uri"`
}

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
func (a *APIEnv) Login(c *gin.Context) {
	db := a.DB
	cfg := config.EnvironmentConfig()
	var user database.User
	var userRequest UserCreateOrLoginRequest

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
		"token": tokenString,
	})
	return
}

func (a *APIEnv) Register(c *gin.Context) {
	db := a.DB
	cfg := config.EnvironmentConfig()
	var userRequest UserCreateOrLoginRequest
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
		"preference": user.Preference,
		"weightUnit": user.WeightUnit,
		"heightUnit": user.HeightUnit,
		"weight":     user.Weight,
		"height":     user.Height,
		"name":       user.Name,
		"email":      user.Email,
		"imageUri":   user.ImageUri,
	})
}

func (a *APIEnv) UpdateUser(c *gin.Context) {
	db := a.DB

	var userRequest UserUpdateRequest
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

	if userRequest.Name != "" {
		user.Name = userRequest.Name
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot have empty name field"})
		return
	}

	if userRequest.ImageUri != "" {
		if !validateURIWithTLD(userRequest.ImageUri) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad image uri url"})
			return
		}
		user.ImageUri = userRequest.ImageUri
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot have empty imageUri field"})
		return
	}

	user.Preference = userRequest.Preference
	user.WeightUnit = userRequest.WeightUnit
	user.HeightUnit = userRequest.HeightUnit
	user.Weight = userRequest.Weight
	user.Height = userRequest.Height

	if err := db.Save(&user).Error; err != nil {
		if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preference": user.Preference,
		"weightUnit": user.WeightUnit,
		"heightUnit": user.HeightUnit,
		"weight":     user.Weight,
		"height":     user.Height,
		"name":       user.Name,
		"email":      user.Email,
		"imageUri":   user.ImageUri,
	})
}
