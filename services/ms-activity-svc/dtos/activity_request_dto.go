package dtos

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type ActivityRequestDTO struct {
	ActivityType      string `json:"activityType" binding:"required"`
	DoneAt            string `json:"doneAt" binding:"required"`
	DurationInMinutes int    `json:"durationInMinutes" binding:"omitempty"`
}

func ValidateActivityRequest(dto ActivityRequestDTO) error {
	// Validasi activityType
    if dto.ActivityType == "" {
        return errors.New("activityType is required")
    }
    
    // Cek apakah activityType valid dengan membandingkan ke daftar ActivityTypeDTO
    isValidActivity := false
    for _, activity := range ActivityValues {
        if activity.Name == dto.ActivityType {
            isValidActivity = true
            break
        }
    }
    if !isValidActivity {
        validActivities := make([]string, len(ActivityValues))
        for i, activity := range ActivityValues {
            validActivities[i] = activity.Name
        }
        return fmt.Errorf("activityType must be one of %v", validActivities)
    }

	// Validasi doneAt
	if dto.DoneAt == "" {
		return errors.New("doneAt is required")
	}
	if !isValidISODate(dto.DoneAt) {
		return errors.New("doneAt must be a valid ISO date")
	}

	// Validasi durationInMinutes
	if dto.DurationInMinutes < 1 {
		return errors.New("durationInMinutes must be at least 1")
	}

	return nil
}



func isValidISODate(date string) bool {
	isoDatePattern := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$`
	matched, _ := regexp.MatchString(isoDatePattern, date)
	if !matched {
		return false
	}
	_, err := time.Parse(time.RFC3339, date)
	return err == nil
}
