package validators

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryParamValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit := ctx.DefaultQuery("limit", "10")
		offset := ctx.DefaultQuery("offset", "0")
		name := ctx.Query("name")

		// Validate limit
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			ctx.Abort()
			return
		}

		// Validate offset
		offsetInt, err := strconv.Atoi(offset)
		if err != nil || offsetInt < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
			ctx.Abort()
			return
		}

		if name != "" {
			// Validasi panjang name
			if len(name) < 0 {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "'name' must be at least 3 characters long"})
				ctx.Abort()
				return
			}
		}

		// Set validated values to context for reuse
		ctx.Set("limit", limitInt)
		ctx.Set("offset", offsetInt)
		ctx.Set("name", name)

		ctx.Next()
	}
}
