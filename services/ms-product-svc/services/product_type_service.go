package services

import (
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/mappers"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/repositories"
)

type ProductTypeService struct {
	repo *repositories.ProductTypeRepository
}

func NewProductTypeService(repo *repositories.ProductTypeRepository) *ProductTypeService {
	return &ProductTypeService{
		repo: repo,
	}
}

func (s *ProductTypeService) GetAll(limit, offset int) ([]*dtos.ProductTypeResponseDTO, error) {
	productTypes, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Map models to response DTOs
	var productTypeDTOs []*dtos.ProductTypeResponseDTO
	for _, productType := range productTypes {
		productTypeDTOs = append(productTypeDTOs, mappers.MapProductTypeModelToResponse(productType))
	}

	return productTypeDTOs, nil
}

func (s *ProductTypeService) GetOne(id string) (*dtos.ProductTypeResponseDTO, error) {
	productType, err := s.repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return mappers.MapProductTypeModelToResponse(productType), nil
}
