package repositories

import (
	"time"

	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{
		db: db,
	}
}

func (r *ActivityRepository) Create(activity *models.Activity) (*models.Activity, error) {
	if err := r.db.Create(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *ActivityRepository) GetAll(params map[string]interface{}) ([]*models.Activity, error) {
	var activities []*models.Activity
	query := r.db.Model(&models.Activity{})

	if activityType, ok := params["activityType"].(string); ok && activityType != "" {
		query = query.Where("activity_type_name = ?", activityType)
	}

	if doneAtFrom, ok := params["doneAtFrom"].(time.Time); ok {
		doneAtFrom = doneAtFrom.UTC()
		query = query.Where("done_at >= ?", doneAtFrom)
	}

	if doneAtTo, ok := params["doneAtTo"].(time.Time); ok {
		doneAtTo = doneAtTo.UTC()
		query = query.Where("done_at <= ?", doneAtTo)
	}

	if min, ok := params["caloriesBurnedMin"].(int); ok && min > 0 {
		query = query.Where("calories_burned >= ?", min)
	}

	if max, ok := params["caloriesBurnedMax"].(int); ok && max > 0 {
		query = query.Where("calories_burned <= ?", max)
	}

	// Handle limit and offset
	limit := params["limit"].(int)
	offset := params["offset"].(int)
	query = query.Limit(limit).Offset(offset)
	query = query.Debug()

	// Execute query
	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *ActivityRepository) GetOne(id, userId string) (*models.Activity, error) {
	var activity models.Activity
	err := r.db.First(&activity, "id = ? AND user_id = ?", id, userId).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *ActivityRepository) UpdateActivity(data *models.Activity) (*models.Activity, error) {

	updateErr := r.db.Model(&models.Activity{}).
		Where("id = ?", data.ID).
		Updates(data).
		First(data).Error

	if updateErr != nil {
		return nil, updateErr
	}
	return data, nil
}

func (r *ActivityRepository) DeleteById(id string) error {
	deleteErr := r.db.Where("id = ? ", id).Delete(&models.Activity{}).Error
	return deleteErr
}
