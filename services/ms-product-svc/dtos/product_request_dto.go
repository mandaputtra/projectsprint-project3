package dtos

import (
	"errors"
	"fmt"
	"strings"
)

type ProductRequestDTO struct {
	Name     string `json:"name" binding:"required"`
	Category string `json:"category" binding:"required"`
	Qty      int    `json:"qty" binding:"required"`
	Price    int    `json:"price" binding:"required"`
	Sku      string `json:"sku" binding:"required"`
	FileId   string `json:"fileId" binding:"required"`
}

var ValidCategories = []string{"Food", "Beverages", "Clothes", "Furniture", "Tools"}

func ValidateProductRequest(dto ProductRequestDTO) error {
	// Validasi Name
	if strings.TrimSpace(dto.Name) == "" {
		return errors.New("name is required")
	}
	if len(dto.Name) < 4 {
		return errors.New("name must be at least 4 characters")
	}
	if len(dto.Name) > 32 {
		return errors.New("name must not exceed 32 characters")
	}

	// Validasi Category
	if strings.TrimSpace(dto.Category) == "" {
		return errors.New("category is required")
	}
	isValidCategory := false
	for _, category := range ValidCategories {
		if dto.Category == category {
			isValidCategory = true
			break
		}
	}
	if !isValidCategory {
		return fmt.Errorf("category must be one of %v", ValidCategories)
	}

	// Validasi Qty
	if dto.Qty < 1 {
		return errors.New("qty must be at least 1")
	}

	// Validasi Price
	if dto.Price < 100 {
		return errors.New("price must be at least 100")
	}

	// Validasi Sku
	if len(dto.Sku) > 32 {
		return errors.New("sku must not exceed 32 characters")
	}

	// Validasi FileId
	if strings.TrimSpace(dto.FileId) == "" {
		return errors.New("fileId is required")
	}

	return nil
}
