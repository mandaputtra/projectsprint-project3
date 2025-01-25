package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID               string `gorm:"primaryKey"`
	UserID           string `gorm:"type:varchar(255);null"`
	Name             string `gorm:"type:varchar(100);null"`
	CategoryID       string `gorm:"type:varchar(255);not null"`
	CategoryName     string `gorm:"type:varchar(100);not null"`
	Qty              int
	Price            int
	Sku              string `gorm:"type:varchar(100);null"`
	FileID           string `gorm:"type:varchar(255);null"`
	FileURI          string `gorm:"type:varchar(255);null"`
	FileThumbnailURI string `gorm:"type:varchar(255);null"`
	DoneAt           time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

func (Product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if Product.ID == "" {
		Product.ID = uuid.NewString()
	}
	return
}
