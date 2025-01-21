package mappers

import (
	"time"

	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
)

func MapActivityModelToResponse(activityModel *models.Activity) *dtos.ActivityResponseDTO {
	return &dtos.ActivityResponseDTO{
		ActivityId:        activityModel.ID,
		ActivityType:      activityModel.ActivityTypeName,
		DoneAt:            activityModel.DoneAt.UTC().Format(time.RFC3339Nano),
		DurationInMinutes: activityModel.DurationInMinutes,
		CaloriesBurned:    activityModel.CaloriesBurned,
		CreatedAt:         activityModel.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:         activityModel.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
