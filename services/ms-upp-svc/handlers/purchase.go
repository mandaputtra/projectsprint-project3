package handlers

import (
	"net/http"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"

	"github.com/gin-gonic/gin"
)

func (a *APIEnv) Purchase(c *gin.Context) {
	var userRequest dto.Order
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Purchased items must have minimal one items
	if len(userRequest.PurchasedItems) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimal purchased items should be one"})
		return
	}

	// Purchased quantity minimal should be 2

	// 1. Validate user request
	// 2. Get data product
	// 3. Validate product id
	// 4. Get product owner
	// 5. Validate product quantity already 0 or below return it as invalid product quantity bought
	// 6. Get product total

	c.JSON(http.StatusOK, gin.H{
		"purchaseId":     "",
		"purchasedItems": "[]",
		"totalPrice":     1,
		"paymentDetails": "[]",
	})
}

func (a *APIEnv) PurchaseVerify(c *gin.Context) {
	db := a.DB

	var userRequest dto.OrderVerify
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(userRequest.FileIds) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimum one file uploaded"})
		return
	}

	var file []database.File
	if result := db.Find(&file, userRequest.FileIds); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	// 1. Validate if really the file id used is correct
	// 2. Validate order paymentDetails if the length matches then its valid order
	if len(file) != len(userRequest.FileIds) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some of file Ids invalid"})
		return
	}

	// 3. Decrease product quantity

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success upload order, please comeback again to order more!",
	})
}
