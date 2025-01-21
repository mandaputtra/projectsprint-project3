package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID                string `gorm:"primaryKey"`
	UserID            string `gorm:"type:varchar(255);not null"`
	ActivityTypeID    string `gorm:"type:varchar(255);not null"`
	ActivityTypeName  string `gorm:"type:varchar(100);null"`
	CaloriesBurned    int
	DurationInMinutes int `gorm:"default:1"`
	DoneAt            time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

func (Activity) TableName() string {
	return "activities"
}

func (activity *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	if activity.ID == "" {
		activity.ID = uuid.NewString()
	}
	return
}
