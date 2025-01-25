package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateGetAllProductsQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultLimit := 10
		defaultOffset := 0

		// Parse limit
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit <= 0 {
			limit = defaultLimit
		}

		// Parse offset
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil || offset < 0 {
			offset = defaultOffset
		}
		// Add validated query parameters to context
		c.Set("validatedQuery", map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		})

		c.Next()
	}
}
