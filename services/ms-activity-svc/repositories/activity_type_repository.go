package repositories

import (
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
	"gorm.io/gorm"
)

type ActivityTypeRepository struct {
	db *gorm.DB
}

func NewActivityTypeRepository(db *gorm.DB) *ActivityTypeRepository {
	return &ActivityTypeRepository{
		db: db,
	}
}

func (r *ActivityTypeRepository) GetAll(limit, offset int) ([]*models.ActivityType, error) {
	var activityTypes []*models.ActivityType
	query := r.db.Limit(limit).Offset(offset)

	err := query.Find(&activityTypes).Error
	return activityTypes, err
}

func (r *ActivityTypeRepository) GetOne(id string) (*models.ActivityType, error) {
	var activityType models.ActivityType
	err := r.db.First(&activityType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &activityType, nil
}

func (r *ActivityTypeRepository) GetOneByName(name string) (*models.ActivityType, error) {
	var activityType models.ActivityType
	err := r.db.First(&activityType, "activity_type = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &activityType, nil
}
