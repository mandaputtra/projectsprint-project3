package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/config"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/controllers"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/database"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/middlewares"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/repositories"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectDatabase(env config.Environment) *gorm.DB {
	log.Println("Connect to database ....")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s search_path=%s",
		env.DATABASE_HOST,
		env.DATABASE_USER,
		env.DATABASE_PASSWORD,
		env.DATABASE_NAME,
		env.DATABASE_PORT,
		env.DATABASE_SCHEMA,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Test query
	var result string
	db.Raw("SELECT 1;").Scan(&result)

	log.Printf("Connection successfull. Result from test SQL: %s\n", result)

	// Migrations
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.ProductType{})

	return db
}

func setupRouter(
	productController *controllers.ProductController,
	productTypeController *controllers.ProductTypeController,
) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	v1 := r.Group("/v1")

	// Routes untuk activity
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		product := v1.Group("/product")
		{
			product.GET("/", middlewares.ValidateGetAllProductsQuery(), productController.GetAllProducts)
			product.GET("/:id", productController.GetOneProduct)
			product.POST("/", productController.Create)
			product.PATCH("/:id", productController.UpdateProduct)
			product.DELETE("/:id", productController.DeleteOneProduct)
		}

		// Routes untuk activity-type
		productTypes := v1.Group("/product-type")
		{
			productTypes.GET("/", productTypeController.GetAllProductType)
			productTypes.GET("/:id", productTypeController.GetOneProductType)
		}
	}

	return r
}

func main() {
	// Load .env
	cfg := config.EnvironmentConfig()

	// connect databases
	db := connectDatabase(cfg)

	// seeder
	database.SeedProductTypes(db)

	productTypeRepo := repositories.NewProductTypeRepository(db)
	productRepo := repositories.NewProductRepository(db)

	productTypeService := services.NewProductTypeService(productTypeRepo)
	productService := services.NewProductService(productRepo, productTypeRepo)

	productController := controllers.NewProductController(productService)
	productTypeController := controllers.NewProductTypeController(productTypeService)

	r := setupRouter(productController, productTypeController)
	r.Run(":" + cfg.PORT)
}
