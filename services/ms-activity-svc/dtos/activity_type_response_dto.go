package dtos

type ActivityTypeResponseDTO struct {
	ID           string `json:"id"`
	ActivityType string `json:"activityType"`
	Calories     int    `json:"calories"`
}
