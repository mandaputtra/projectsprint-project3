package handlers

import (
	"net/http"
	"project3/services/ms-upp-svc/config"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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
		"phone":  user.Phone,
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
		"phone": user.Phone,
		"token": tokenString,
	})
}

func (a *APIEnv) RegisterWithPhone(c *gin.Context) {
	db := a.DB
	cfg := config.EnvironmentConfig()
	var userRequest dto.UserCreateOrLoginWithPhoneRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := database.User{Phone: userRequest.Phone, Password: userRequest.Password}
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
		"phone":  user.Phone,
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
		"phone": user.Phone,
		"token": tokenString,
	})
}

func (a *APIEnv) LoginWithPhone(c *gin.Context) {
	db := a.DB
	cfg := config.EnvironmentConfig()
	var user database.User
	var userRequest dto.UserCreateOrLoginWithPhoneRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Where("phone =?", userRequest.Phone).First(&user)
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
		"phone":  user.Phone,
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
}
