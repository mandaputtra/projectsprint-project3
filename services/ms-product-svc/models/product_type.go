package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductType struct {
	ID        string `gorm:"primaryKey"`
	Type      string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ProductType) TableName() string {
	return "product_types"
}

func (product_type *ProductType) BeforeCreate(tx *gorm.DB) (err error) {
	if product_type.ID == "" {
		product_type.ID = uuid.NewString()
	}
	return
}
