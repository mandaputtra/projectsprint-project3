package middlewares

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateGetAllProductsQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Default values
		defaultLimit := 5
		defaultOffset := 0

		// Parse and validate limit
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit <= 0 {
			limit = defaultLimit
		}

		// Parse and validate offset
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil || offset < 0 {
			offset = defaultOffset
		}

		// Parse productId
		productId := c.Query("productId")
		if productId == "" {
			productId = "" // Ignore invalid productId
		}

		// Parse and validate sku
		sku := c.Query("sku")
		if sku == "" {
			sku = "" // Ignore invalid sku
		}

		// Parse and validate category
		category := c.Query("category")
		validCategories := []string{"Food", "Beverage", "Clothes", "Furniture", "Tools"} // Example categories
		isValidCategory := false
		for _, validCategory := range validCategories {
			if category == validCategory {
				isValidCategory = true
				break
			}
		}
		if !isValidCategory {
			category = "" // Ignore invalid category
		}

		// Parse and validate search (exact match)
		search := c.Query("search")
		if search == "" {
			search = "" // Ignore invalid search
		}

		// Parse and validate sortBy
		sortBy := c.Query("sortBy")
		validSortOptions := map[string]bool{
			"newest":   true,
			"cheapest": true,
		}
		if len(sortBy) > 5 && sortBy[:5] == "sold-" {
			_, err := strconv.Atoi(sortBy[5:])
			if err == nil {
				validSortOptions[sortBy] = true
			}
		}
		if !validSortOptions[sortBy] {
			sortBy = "" // Ignore invalid sortBy
		}

		// Add validated query parameters to context
		c.Set("validatedQuery", map[string]interface{}{
			"limit":     limit,
			"offset":    offset,
			"productId": productId,
			"sku":       sku,
			"category":  category,
			"search":    search,
			"sortBy":    sortBy,
		})

		c.Next()
	}
}
