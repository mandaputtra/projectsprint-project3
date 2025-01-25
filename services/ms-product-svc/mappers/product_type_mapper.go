package mappers

import (
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
)

func MapProductTypeModelToResponse(productTypeModel *models.ProductType) *dtos.ProductTypeResponseDTO {
	return &dtos.ProductTypeResponseDTO{
		ID:   productTypeModel.ID,
		Type: productTypeModel.Type,
	}
}
