package main

import (
	"log"
	"net/http"
	"project3/libs/utils"
	"project3/services/ms-upp-svc/config"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

// Controller
func setupRouter(db *gorm.DB, cfg *config.Environment) *gin.Engine {
	r := gin.Default()

	api := &handlers.APIEnv{
		DB:  db,
		ENV: cfg,
	}

	r.Use(utils.CheckContentType)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/login/email", api.LoginWithEmail)
		v1.POST("/register/email", api.RegisterWithEmail)
		v1.POST("/login/phone", api.LoginWithEmail)
		v1.POST("/register/phone", api.RegisterWithEmail)
		v1.GET("/user", utils.Authorization, api.GetUser)
		v1.PATCH("/user", utils.Authorization, api.UpdateUser)
	}

	return r
}

func main() {
	// Load .env
	cfg := config.EnvironmentConfig()
	db := database.ConnectDatabase(cfg)

	r := setupRouter(db, &cfg)
	if err := r.Run(":" + cfg.PORT); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
