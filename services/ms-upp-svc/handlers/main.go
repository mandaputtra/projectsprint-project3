package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"project3/services/ms-upp-svc/config"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIEnv struct {
	DB  *gorm.DB
	ENV *config.Environment
}

// Utils
func isValidFile(file *multipart.FileHeader) (bool, string) {
	allowedExtensions := map[string]bool{
		".jpeg": true,
		".jpg":  true,
		".png":  true,
	}
	const maxFileSize = 100 * 1024 // 100 KiB

	ext := filepath.Ext(file.Filename)
	if !allowedExtensions[ext] {
		return false, "Invalid file extension"
	}

	if file.Size > maxFileSize {
		return false, "File size exceeds the 100KiB limit"
	}

	return true, ""
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
		"fileThumbnailUri":  user.FileId,
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
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
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
		"fileThumbnailUri":  user.FileId,
		"phone":             user.Phone,
	})
}

func (a *APIEnv) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to get file from request",
		})
		return
	}

	if valid, msg := isValidFile(file); !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": msg,
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to open file",
		})
		return
	}
	defer src.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.ENV.AWS_S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			a.ENV.AWS_ACCESS_KEY_ID,
			a.ENV.AWS_SECRET_ACCESS_KEY,
			"", // a token will be created when the session it's used.
		),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create AWS session",
		})
		return
	}

	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.ENV.AWS_S3_BUCKET),
		Key:    aws.String(file.Filename),
		Body:   src,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to upload file, %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uri": result.Location,
	})
}
