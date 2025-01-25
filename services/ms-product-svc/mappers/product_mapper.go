package mappers

import (
	"time"

	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/dtos"
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
)

func MapProductModelToResponse(productModel *models.Product) *dtos.ProductResponseDTO {
	return &dtos.ProductResponseDTO{
		ProductId: productModel.ID,
		Name:      productModel.Name,
		Category:  productModel.CategoryName,
		Qty:       productModel.Qty,
		Price:     productModel.Price,
		Sku:       productModel.Sku,
		FileId:    productModel.FileID,
		//harus hit ke ms-upload-svc untuk mendapatkan fileUri dan fileThumbnailUri
		CreatedAt: productModel.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt: productModel.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
}
