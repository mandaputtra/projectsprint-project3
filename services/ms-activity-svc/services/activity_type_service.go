package services

import (
	// "errors"
	// "fmt"

	// "github.com/gin-gonic/gin"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/mappers"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/repositories"
	// "gorm.io/gorm"
)

type ActivityTypeService struct {
	repo *repositories.ActivityTypeRepository
}

func NewActivityTypeService(repo *repositories.ActivityTypeRepository) *ActivityTypeService {
	return &ActivityTypeService{
		repo: repo,
	}
}

func (s *ActivityTypeService) GetAll(limit, offset int) ([]*dtos.ActivityTypeResponseDTO, error) {
	activityTypes, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map models to response DTOs
	var activityTypeDTOs []*dtos.ActivityTypeResponseDTO
	for _, activityType := range activityTypes {
		activityTypeDTOs = append(activityTypeDTOs, mappers.MapActivityTypeModelToResponse(activityType))
	}

	return activityTypeDTOs, nil
}

func (s *ActivityTypeService) GetOne(id string) (*dtos.ActivityTypeResponseDTO, error) {
	activityType, err := s.repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return mappers.MapActivityTypeModelToResponse(activityType), nil
}
