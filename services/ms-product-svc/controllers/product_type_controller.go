package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/services"
)

type ProductTypeController struct {
	service *services.ProductTypeService
}

func NewProductTypeController(service *services.ProductTypeService) *ProductTypeController {
	return &ProductTypeController{
		service: service,
	}
}

func (c *ProductTypeController) GetAllProductType(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		// Jika terjadi error, beri nilai default
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		// Jika terjadi error, beri nilai default
		offset = 0
	}

	Product_types, err := c.service.GetAll(limit, offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product_types"})
		return
	}

	if len(Product_types) <= 0 {
		ctx.JSON(http.StatusNoContent, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, Product_types)
}

func (c *ProductTypeController) GetOneProductType(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.service.GetOne(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ID is not found"})
		return
	}
	ctx.JSON(http.StatusOK, product)
}
