package handlers

import (
	"net/http"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"
	"strconv"

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
	query := api.DB

	// Pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Filtering
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if sku := c.Query("sku"); sku != "" {
		query = query.Where("sku = ?", sku)
	}
	if id := c.Query("id"); id != "" {
		query = query.Where("id = ?", id)
	}

	// Sorting
	sortBy := c.DefaultQuery("sortBy", "newest")
	if sortBy == "cheapest" {
		query = query.Order("price asc")
	} else {
		query = query.Order("created_at desc")
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
