package handlers

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"project3/services/ms-upp-svc/config"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"
	"strings"

	"github.com/gin-gonic/gin"
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

// Service
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
		"fileId":            user.File.FileID,
		"fileUri":           user.File.FileUri,
		"fileThumbnailUri":  user.File.FileThumbnailUri,
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
	if err := db.Where("id = ?", id).Preload("File").First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// To update
	user.FileID = userRequest.FileID
	user.BankAccountHolder = userRequest.BankAccountHolder
	user.BankAccountName = userRequest.BankAccountName
	user.BankAccountNumber = userRequest.BankAccountNumber

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
		"fileId":            user.File.FileID,
		"fileUri":           user.File.FileUri,
		"fileThumbnailUri":  user.File.FileThumbnailUri,
		"phone":             user.Phone,
	})
}

func (a *APIEnv) LinkPhone(c *gin.Context) {
	db := a.DB
	userID := c.GetString("userId")

	var userRequest dto.UserLinkPhone
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user database.User
	if err := db.Where("id = ?", userID).Preload("File").Where("phone = ?", "").First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Model(&user).Where("id = ?", userID).Update("phone", userRequest.Phone).Error; err != nil {
		print(err)
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
		"fileId":            user.File.FileID,
		"fileUri":           user.File.FileUri,
		"fileThumbnailUri":  user.File.FileThumbnailUri,
		"phone":             user.Phone,
	})
}

func (a *APIEnv) LinkEmail(c *gin.Context) {
	db := a.DB
	userID := c.GetString("userId")

	var userRequest dto.UserLinkEmail
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user database.User
	if err := db.Where("id = ?", userID).Preload("File").Where("phone is ?", nil).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := db.Model(&user).Where("id = ?", userID).Update("email", userRequest.Email).Error; err != nil {
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
		"fileId":            user.File.FileID,
		"fileUri":           user.File.FileUri,
		"fileThumbnailUri":  user.File.FileThumbnailUri,
		"phone":             user.Phone,
	})
}
