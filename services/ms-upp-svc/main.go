package main

import (
	"log"
	"net/http"
	"project3/libs/utils"
	"project3/services/ms-upp-svc/config"
	"project3/services/ms-upp-svc/database"
	"project3/services/ms-upp-svc/dto"
	"project3/services/ms-upp-svc/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

// Controller
func setupRouter(db *gorm.DB, cfg *config.Environment) *gin.Engine {
	r := gin.Default()

	// Custom validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phoneNumber", dto.ValidatePhoneNumber)
	}

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
		v1.POST("/login/phone", api.LoginWithPhone)
		v1.POST("/register/phone", api.RegisterWithPhone)

		v1.GET("/user", utils.Authorization, api.GetUser)
		v1.PATCH("/user", utils.Authorization, api.UpdateUser)
		v1.POST("/user/link/phone", utils.Authorization, api.LinkPhone)
		v1.POST("/user/link/email", utils.Authorization, api.LinkEmail)

		v1.POST("/file", api.UploadFile)

		v1.POST("/purchase", api.UploadFile)
		v1.POST("/purchase/:purchaseId", api.UploadFile)
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
