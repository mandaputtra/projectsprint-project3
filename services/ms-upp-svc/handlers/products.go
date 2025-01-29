package handlers

import (
	"net/http"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"

	"github.com/gin-gonic/gin"
)

func (api *APIEnv) CreateProduct(c *gin.Context) {
	var req dto.ProductCreateOrUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := database.Product{
		Name:     req.Name,
		Category: req.Category,
		Qty:      req.Qty,
		Price:    req.Price,
		SKU:      req.SKU,
		FileID:   req.FileID,
	}

	if err := api.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (api *APIEnv) UpdateProduct(c *gin.Context) {
	var req dto.ProductCreateOrUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product database.Product
	if err := api.DB.First(&product, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	product.Name = req.Name
	product.Category = req.Category
	product.Qty = req.Qty
	product.Price = req.Price
	product.SKU = req.SKU
	product.FileID = req.FileID

	if err := api.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (api *APIEnv) DeleteProduct(c *gin.Context) {
	var product database.Product
	if err := api.DB.First(&product, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := api.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (api *APIEnv) GetProducts(c *gin.Context) {
	var products []database.Product
	if err := api.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
