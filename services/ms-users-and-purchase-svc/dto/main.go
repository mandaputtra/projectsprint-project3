package dto

// Type
type UserCreateOrLoginWithEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserCreateOrLoginWithPhoneRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserUpdateRequest struct {
	Email      string  `json:"email"`
	Preference string  `json:"preference" binding:"required,oneof=CARDIO WEIGHT"`
	WeightUnit string  `json:"weightUnit" binding:"required,oneof=KG LBS"`
	HeightUnit string  `json:"heightUnit" binding:"required,oneof=CM INCH"`
	Height     float64 `json:"height" binding:"required,min=3,max=250"`
	Weight     float64 `json:"weight" binding:"required,min=10,max=1000"`
	Name       string  `json:"name" binding:"omitempty,min=2,max=60"`
	ImageUri   string  `json:"imageUri" binding:"omitempty,uri"`
}
