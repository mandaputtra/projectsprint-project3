package database

import (
	"log"

	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
	"gorm.io/gorm"
)

func SeedProductTypes(db *gorm.DB) {
	productTypes := []models.ProductType{
		{Type: "Food"},
		{Type: "Beverage"},
		{Type: "Clothes"},
		{Type: "Furniture"},
		{Type: "Tools"},
	}

	for _, ProductType := range productTypes {
		var existing models.ProductType
		if err := db.Where("type = ?", ProductType.Type).First(&existing).Error; err != nil {
			if err := db.Create(&ProductType).Error; err != nil {
				log.Printf("Failed to seed product type: %v", ProductType.Type)
			} else {
				log.Printf("Seeded product type: %v", ProductType.Type)
			}
		}
	}
}
