package database

import (
	"log"

	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
	"gorm.io/gorm"
)

func SeedActivityTypes(db *gorm.DB) {
	activityTypes := []models.ActivityType{
		{ActivityType: "Walking", Calories: 4},
		{ActivityType: "Yoga", Calories: 4},
		{ActivityType: "Stretching", Calories: 4},
		{ActivityType: "Cycling", Calories: 8},
		{ActivityType: "Swimming", Calories: 8},
		{ActivityType: "Dancing", Calories: 8},
		{ActivityType: "Hiking", Calories: 10},
		{ActivityType: "Running", Calories: 10},
		{ActivityType: "HIIT", Calories: 10},
		{ActivityType: "JumpRope", Calories: 10},
	}

	for _, activityType := range activityTypes {
		var existing models.ActivityType
		if err := db.Where("activity_type = ?", activityType.ActivityType).First(&existing).Error; err != nil {
			if err := db.Create(&activityType).Error; err != nil {
				log.Printf("Failed to seed activity: %v", activityType.ActivityType)
			} else {
				log.Printf("Seeded activity: %v", activityType.ActivityType)
			}
		}
	}
}
