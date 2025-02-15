package database

import (
	"fmt"
	"log"
	"project3/services/ms-upp-svc/config"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// File
type File struct {
	FileID           string `gorm:"primaryKey"`
	FileUri          string `gorm:"required"`
	FileThumbnailUri string `gorm:"required"`
}

func (file *File) BeforeCreate(tx *gorm.DB) (err error) {
	if file.FileID == "" {
		file.FileID = uuid.NewString()
	}
	return
}

// User
type User struct {
	ID                string `gorm:"primaryKey"`
	Email             string `gorm:"uniqueIndex:uq_email_and_phone_idx;"`
	Phone             string `gorm:"uniqueIndex:uq_email_and_phone_idx;"`
	Password          string `gorm:"not null"`
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
	FileID            *string `gorm:"default:null"`
	File              File    `gorm:"references:FileID"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	return
}

// Purchases
type Purchases struct {
	ID                  string `gorm:"primaryKey"`
	SenderName          string `gorm:"required"`
	SenderContactType   string `gorm:"required"` // "email"/"phone"
	SenderContactDetail string `gorm:"required"`
	PurchaseProofs      string
	OrderItems          []PurchaseItems `gorm:"constraint:OnDelete:CASCADE;foreignKey:PurchaseID;"`
}

func (purchase *Purchases) BeforeCreate(tx *gorm.DB) (err error) {
	if purchase.ID == "" {
		purchase.ID = uuid.NewString()
	}
	return
}

type PurchaseItems struct {
	PurchaseID string `gorm:"required"`
	ProductID  string `gorm:"required"`
	Qty        int64  `gorm:"required"`
}

type Product struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(32);null"`
	Category  string `gorm:"type:varchar(255);not null"`
	Qty       int    `gorm:"not null;default:1"`
	Price     int    `gorm:"not null;default:100"`
	SKU       string `gorm:"type:varchar(32);null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID string `gorm:"not null"`
	User   User   `gorm:"references:ID"`

	FileID string `gorm:"default:null"`
	File   File   `gorm:"references:FileID"`
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if product.ID == "" {
		product.ID = uuid.NewString()
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
	db.AutoMigrate(&File{}, &User{}, &Purchases{}, &PurchaseItems{}, Product{})
	return db
}

func GetDB() *gorm.DB {
	return db
}
