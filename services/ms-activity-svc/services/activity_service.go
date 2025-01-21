package services

import (
	// "errors"
	// "fmt"

	// "github.com/gin-gonic/gin"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/mappers"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/models"
	"github.com/mandaputtra/projectsprint-projects2/services/ms-activity-svc/repositories"
	"gorm.io/gorm"
	// "gorm.io/gorm"
)

type ActivityService struct {
	repo             *repositories.ActivityRepository
	activityTypeRepo *repositories.ActivityTypeRepository
}

func NewActivityService(repo *repositories.ActivityRepository, activityTypeRepo *repositories.ActivityTypeRepository) *ActivityService {
	return &ActivityService{
		repo:             repo,
		activityTypeRepo: activityTypeRepo,
	}
}

func (s *ActivityService) Create(activityReqDTO *dtos.ActivityRequestDTO, ctx *gin.Context) (*dtos.ActivityResponseDTO, error) {
	activityType, err := s.activityTypeRepo.GetOneByName(activityReqDTO.ActivityType)

	if err != nil {
		return nil, err
	}

	userId, _ := ctx.Get("userId")
	doneAt, _ := time.Parse(time.RFC3339Nano, activityReqDTO.DoneAt)

	newActivityModel := &models.Activity{
		UserID:            userId.(string),
		ActivityTypeID:    activityType.ID,
		ActivityTypeName:  activityType.ActivityType,
		CaloriesBurned:    activityType.Calories * activityReqDTO.DurationInMinutes,
		DurationInMinutes: activityReqDTO.DurationInMinutes,
		DoneAt:            doneAt,
	}

	activity, err := s.repo.Create(newActivityModel)
	if err != nil {
		return nil, err
	}

	activityResponseDTO := mappers.MapActivityModelToResponse(activity)
	return activityResponseDTO, nil
}

func (s *ActivityService) GetAll(params map[string]interface{}) ([]*dtos.ActivityResponseDTO, error) {
	// Ambil data dari repository
	activities, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Konversi model ke DTO menggunakan mapper
	var activityDTOs []*dtos.ActivityResponseDTO
	for _, activity := range activities {
		activityDTOs = append(activityDTOs, mappers.MapActivityModelToResponse(activity))
	}

	return activityDTOs, nil
}

func (s *ActivityService) GetOne(id, userId string) (*dtos.ActivityResponseDTO, error) {
	activity, err := s.repo.GetOne(id, userId)
	if err != nil {
		return nil, err
	}

	return mappers.MapActivityModelToResponse(activity), nil
}

func (s *ActivityService) UpdateActivity(id, userId string, activityDTO *dtos.ActivityRequestDTO) (*dtos.ActivityResponseDTO, error) {
	_, err := s.GetOne(id, userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("activity not found")
		}
		return nil, err // Error lain
	}

	existingActivityType, err := s.activityTypeRepo.GetOneByName(activityDTO.ActivityType)
	if err != nil {
		return nil, err
	}

	doneAt, _ := time.Parse(time.RFC3339Nano, activityDTO.DoneAt)

	updateActivityModel := &models.Activity{
		ID:                id,
		UserID:            userId,
		ActivityTypeID:    existingActivityType.ID,
		ActivityTypeName:  existingActivityType.ActivityType,
		CaloriesBurned:    existingActivityType.Calories * activityDTO.DurationInMinutes,
		DurationInMinutes: activityDTO.DurationInMinutes,
		DoneAt:            doneAt,
	}

	updatedData, err := s.repo.UpdateActivity(updateActivityModel)
	if err != nil {
		return nil, err
	}

	// Map hasil ke DTO respons
	return mappers.MapActivityModelToResponse(updatedData), nil
}

func (s *ActivityService) DeleteById(id, userId string) error {
	_, err := s.repo.GetOne(id, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("activity not found")
		}
		return err
	}

	err = s.repo.DeleteById(id)

	if err != nil {
		return err
	}
	return nil
}
