package dtos

import "time"

type GetAllActivitiesParamRequest struct{
	Limit             int       `form:"limit"`
	Offset            int       `form:"offset"`
	ActivityType      string    `form:"activityType"`
	Search            string    `form:"search"`
	DoneAtFrom        time.Time `form:"doneAtFrom"`
	DoneAtTo          time.Time `form:"doneAtTo"`
	CaloriesBurnedMin int       `form:"caloriesBurnedMin"`
	CaloriesBurnedMax int       `form:"caloriesBurnedMax"`
}