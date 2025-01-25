package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/services"
)

type ProductController struct {
	service *services.ProductService
}

func NewProductController(service *services.ProductService) *ProductController {
	return &ProductController{
		service: service,
	}
}

func (c *ProductController) Create(ctx *gin.Context) {
	if ctx.ContentType() != "application/json" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Content-Type, expected application/json",
		})
		return
	}

	var productDTO dtos.ProductRequestDTO

	// Bind JSON request body ke struct department
	if err := ctx.ShouldBindJSON(&productDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validasi ProductRequestDTO
	if err := dtos.ValidateProductRequest(productDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdProduct, err := c.service.Create(&productDTO, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	// Kirim response dengan data department yang berhasil dibuat
	ctx.JSON(http.StatusCreated, createdProduct)
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	// Ambil query yang sudah divalidasi dari context
	validatedQuery := ctx.MustGet("validatedQuery").(map[string]interface{})

	// Panggil service dengan parameter
	Products, err := c.service.GetAll(validatedQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Products)
}

func (c *ProductController) GetOneProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	product, err := c.service.GetOne(id, userId.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ID is not found"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	if ctx.ContentType() != "application/json" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Content-Type, expected application/json",
		})
		return
	}

	// Bind input dari request body ke DTO
	var productDTO dtos.ProductRequestDTO
	if err := ctx.ShouldBindJSON(&productDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := dtos.ValidateProductRequest(productDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	// Panggil service untuk update
	updatedProduct, err := c.service.UpdateProduct(id, userId.(string), &productDTO)
	if err != nil {
		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Product with the given ID not found"})
			return
		}
		// Error lain saat update
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update Product"})
		return
	}

	// Berhasil update
	ctx.JSON(http.StatusOK, updatedProduct)
}

func (c *ProductController) DeleteOneProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, _ := ctx.Get("userId")

	err := c.service.DeleteById(id, userId.(string))
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"status": "delete is successful"})
		return
	} else if err.Error() == "product not found" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product with the given ID not found"})
		return
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
}
