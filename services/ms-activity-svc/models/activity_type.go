package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityType struct {
	ID           string `gorm:"primaryKey"`
	ActivityType string `gorm:"type:varchar(100);not null"`
	Calories     int    `gorm:"calories;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (ActivityType) TableName() string {
	return "activity_types"
}

func (activity_type *ActivityType) BeforeCreate(tx *gorm.DB) (err error) {
	if activity_type.ID == "" {
		activity_type.ID = uuid.NewString()
	}
	return
}
