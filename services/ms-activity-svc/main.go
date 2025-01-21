package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/config"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/controllers"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/database"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/middlewares"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/repositories"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/services"
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
	db.AutoMigrate(&models.Activity{})
	db.AutoMigrate(&models.ActivityType{})

	return db
}

func setupRouter(
	activityController *controllers.ActivityController,
	activityTypeController *controllers.ActivityTypeController,
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
		activity := v1.Group("/activity")
		{
			activity.GET("/", middlewares.ValidateGetAllActivitiesQuery(), activityController.GetAllActivities)
			activity.GET("/:id", activityController.GetOneActivity)
			activity.POST("/", activityController.Create)
			activity.PATCH("/:id", activityController.UpdateActivity)
			activity.DELETE("/:id", activityController.DeleteOneActivity)
		}

		// Routes untuk activity-type
		activityTypes := v1.Group("/activity-type")
		{
			activityTypes.GET("/", activityTypeController.GetAllActivityType)
			activityTypes.GET("/:id", activityTypeController.GetOneActivityType)
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
	database.SeedActivityTypes(db)

	activityTypeRepo := repositories.NewActivityTypeRepository(db)
	activityRepo := repositories.NewActivityRepository(db)

	activityTypeService := services.NewActivityTypeService(activityTypeRepo)
	activityService := services.NewActivityService(activityRepo, activityTypeRepo)

	activityController := controllers.NewActivityController(activityService)
	activityTypeController := controllers.NewActivityTypeController(activityTypeService)

	r := setupRouter(activityController, activityTypeController)
	r.Run(":" + cfg.PORT)
}
