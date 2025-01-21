package database

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-users-svc/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID         string `gorm:"primaryKey"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Preference string
	WeightUnit string
	HeightUnit string
	Weight     float64
	Height     float64
	Name       string
	ImageUri   string
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	return
}

// Setup database
var db *gorm.DB

func ConnectDatabase(env config.Environment) *gorm.DB {
	log.Println("Connect to database ....")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s search_path=%s",
		env.DATABASE_HOST,
		env.DATABASE_USER,
		env.DATABASE_PASSWORD,
		env.DATABASE_NAME,
		env.DATABASE_PORT,
		env.DATABASE_SCHEMA,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Test query
	var result string
	db.Raw("SELECT 1;").Scan(&result)

	log.Printf("Connection successfull. Result from test SQL: %s\n", result)

	// Migrations
	db.AutoMigrate(&User{})
	return db
}

func GetDB() *gorm.DB {
	return db
}
