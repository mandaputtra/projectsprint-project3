package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ValidateGetAllActivitiesQuery() gin.HandlerFunc {
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

		// Parse doneAtFrom
		var doneAtFrom time.Time
		if from := c.Query("doneAtFrom"); from != "" {
			if parsedTime, err := time.Parse(time.RFC3339Nano, from); err == nil {
				doneAtFrom = parsedTime
			}
		} else {
			doneAtFrom = time.Now().Add(-30 * 24 * time.Hour).UTC()
			c.Set("doneAtFrom", doneAtFrom)
		}

		// Parse doneAtTo
		var doneAtTo time.Time
		if to := c.Query("doneAtTo"); to != "" {
			if parsedTime, err := time.Parse(time.RFC3339Nano, to); err == nil {
				doneAtTo = parsedTime
			}
		} else {
			// Parameter tidak dikirim, gunakan nilai default khusus
			doneAtTo = time.Now().UTC()
			c.Set("doneAtTo", doneAtTo)
		}

		// Parse caloriesBurnedMin
		caloriesBurnedMin, _ := strconv.Atoi(c.DefaultQuery("caloriesBurnedMin", "0"))

		// Parse caloriesBurnedMax
		caloriesBurnedMax, _ := strconv.Atoi(c.DefaultQuery("caloriesBurnedMax", "99999999"))

		// Add validated query parameters to context
		c.Set("validatedQuery", map[string]interface{}{
			"limit":             limit,
			"offset":            offset,
			"activityType":      c.Query("activityType"),
			"search":            c.Query("search"),
			"doneAtFrom":        doneAtFrom,
			"doneAtTo":          doneAtTo,
			"caloriesBurnedMin": caloriesBurnedMin,
			"caloriesBurnedMax": caloriesBurnedMax,
		})

		c.Next()
	}
}
