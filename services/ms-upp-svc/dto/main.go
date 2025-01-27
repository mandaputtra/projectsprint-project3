package dto

import "github.com/go-playground/validator/v10"

// Type
type UserCreateOrLoginWithEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserCreateOrLoginWithPhoneRequest struct {
	Phone    string `json:"phone" binding:"required,phoneNumber"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type UserUpdateRequest struct {
	FileID            string `json:"fileId"`
	BankAccountName   string `json:"bankAccountName" binding:"required,min=4,max=32"`
	BankAccountHolder string `json:"bankAccountHolder" binding:"required,min=4,max=32"`
	BankAccountNumber string `json:"bankAccountNumber" binding:"required,min=4,max=32"`
}

type UserLinkPhone struct {
	Phone string `json:"phone" binding:"required,phoneNumber"`
}

type UserLinkEmail struct {
	Email string `json:"email" binding:"required,email"`
}

// Custom validation function
func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Check if the string starts with '+'
	if len(value) == 0 || value[0] != '+' {
		return false
	}

	// Check the length of the string
	length := len(value)
	if length < 10 || length > 13 {
		return false
	}

	return true
}
