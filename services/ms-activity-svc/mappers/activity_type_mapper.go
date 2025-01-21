package mappers

import (
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
)

func MapActivityTypeModelToResponse(activityModel *models.ActivityType) *dtos.ActivityTypeResponseDTO {
	return &dtos.ActivityTypeResponseDTO{
		ID:           activityModel.ID,
		ActivityType: activityModel.ActivityType,
		Calories:     activityModel.Calories,
	}
}
