package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/services"
)

type ActivityTypeController struct {
	service *services.ActivityTypeService
}

func NewActivityTypeController(service *services.ActivityTypeService) *ActivityTypeController {
	return &ActivityTypeController{
		service: service,
	}
}

func (c *ActivityTypeController) GetAllActivityType(ctx *gin.Context) {
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

	activity_types, err := c.service.GetAll(limit, offset)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity_types"})
		return
	}

	if len(activity_types) <= 0 {
		ctx.JSON(http.StatusNoContent, gin.H{})
		return
	}
	ctx.JSON(http.StatusOK, activity_types)
}

func (c *ActivityTypeController) GetOneActivityType(ctx *gin.Context) {
	id := ctx.Param("id")

	department, err := c.service.GetOne(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ID is not found"})
		return
	}
	ctx.JSON(http.StatusOK, department)
}