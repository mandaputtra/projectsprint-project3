package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"github.com/mandaputtra/projectsprint-projects2/libs/utils"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/config"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/database"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/handlers"
)

// Controller
func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	api := &handlers.APIEnv{
		DB: db,
	}

	r.Use(utils.CheckContentType)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/login", api.Login)
		v1.POST("/register", api.Register)
		v1.GET("/user", utils.Authorization, api.GetUser)
		v1.PATCH("/user", utils.Authorization, api.UpdateUser)
	}

	return r
}

func main() {
	// Load .env
	cfg := config.EnvironmentConfig()
	db := database.ConnectDatabase(cfg)

	r := setupRouter(db)
	if err := r.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
