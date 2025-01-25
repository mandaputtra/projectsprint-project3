package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/mappers"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/repositories"
	"gorm.io/gorm"
)

type ProductService struct {
	repo            *repositories.ProductRepository
	ProductTypeRepo *repositories.ProductTypeRepository
}

func NewProductService(repo *repositories.ProductRepository, ProductTypeRepo *repositories.ProductTypeRepository) *ProductService {
	return &ProductService{
		repo:            repo,
		ProductTypeRepo: ProductTypeRepo,
	}
}

func (s *ProductService) Create(productRequestDTO *dtos.ProductRequestDTO, ctx *gin.Context) (*dtos.ProductResponseDTO, error) {
	productType, err := s.ProductTypeRepo.GetOneByName(productRequestDTO.Category)

	if err != nil {
		return nil, err
	}

	userId, _ := ctx.Get("userId")

	newProductModel := &models.Product{
		UserID:       userId.(string),
		Name:         productRequestDTO.Name,
		CategoryID:   productType.ID,
		CategoryName: productType.Type,
		Qty:          productRequestDTO.Qty,
		Price:        productRequestDTO.Price,
		Sku:          productRequestDTO.Sku,
		FileID:       productRequestDTO.FileId,
		// FileURI: productRequestDTO.,
	}

	Product, err := s.repo.Create(newProductModel)
	if err != nil {
		return nil, err
	}

	ProductResponseDTO := mappers.MapProductModelToResponse(Product)
	return ProductResponseDTO, nil
}

func (s *ProductService) GetAll(params map[string]interface{}) ([]*dtos.ProductResponseDTO, error) {
	// Ambil data dari repository
	activities, err := s.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	// Konversi model ke DTO menggunakan mapper
	var ProductDTOs []*dtos.ProductResponseDTO
	for _, Product := range activities {
		ProductDTOs = append(ProductDTOs, mappers.MapProductModelToResponse(Product))
	}

	return ProductDTOs, nil
}

func (s *ProductService) GetOne(id, userId string) (*dtos.ProductResponseDTO, error) {
	Product, err := s.repo.GetOne(id, userId)
	if err != nil {
		return nil, err
	}

	return mappers.MapProductModelToResponse(Product), nil
}

func (s *ProductService) UpdateProduct(id, userId string, productRequestDTO *dtos.ProductRequestDTO) (*dtos.ProductResponseDTO, error) {
	_, err := s.GetOne(id, userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Product not found")
		}
		return nil, err // Error lain
	}

	existingProductType, err := s.ProductTypeRepo.GetOneByName(productRequestDTO.Category)
	if err != nil {
		return nil, err
	}

	updateProductModel := &models.Product{
		ID:           id,
		UserID:       userId,
		Name:         productRequestDTO.Name,
		CategoryID:   existingProductType.ID,
		CategoryName: existingProductType.Type,
		Qty:          productRequestDTO.Qty,
		Price:        productRequestDTO.Price,
		Sku:          productRequestDTO.Sku,
		FileID:       productRequestDTO.FileId,
	}

	updatedData, err := s.repo.UpdateProduct(updateProductModel)
	if err != nil {
		return nil, err
	}

	// Map hasil ke DTO respons
	return mappers.MapProductModelToResponse(updatedData), nil
}

func (s *ProductService) DeleteById(id, userId string) error {
	_, err := s.repo.GetOne(id, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Product not found")
		}
		return err
	}

	err = s.repo.DeleteById(id)

	if err != nil {
		return err
	}
	return nil
}
