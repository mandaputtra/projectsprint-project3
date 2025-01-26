package database

import (
	"fmt"
	"log"
	"project3/services/ms-upp-svc/config"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type File struct {
	FileID           string `gorm:"primaryKey"`
	FileUri          string `gorm:"required"`
	FileThumbnailUri string `gorm:"required"`
}

type User struct {
	ID                string `gorm:"primaryKey"`
	Email             string `gorm:"uniqueIndex:uq_email_and_phone_idx;"`
	Phone             string `gorm:"uniqueIndex:uq_email_and_phone_idx;"`
	Password          string `gorm:"not null"`
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
	FileID            string
	File              File `gorm:"references:FileID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	return
}

type Purchases struct {
	ID                  string `gorm:"primaryKey"`
	SenderName          string `gorm:"required"`
	SenderContactType   string `gorm:"required"` // "email"/"phone"
	SenderContactDetail string `gorm:"required"`
	PurchaseProofs      string
	OrderItems          []PurchaseItems `gorm:"constraint:OnDelete:CASCADE;foreignKey:PurchaseID;"`
}

type PurchaseItems struct {
	PurchaseID string `gorm:"required"`
	ProductID  string `gorm:"required"`
	Qty        int64  `gorm:"required"`
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
	db.AutoMigrate(&File{}, &User{}, &Purchases{}, &PurchaseItems{})
	return db
}

func GetDB() *gorm.DB {
	return db
}
