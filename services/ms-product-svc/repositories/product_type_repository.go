package repositories

import (
	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
	"gorm.io/gorm"
)

type ProductTypeRepository struct {
	db *gorm.DB
}

func NewProductTypeRepository(db *gorm.DB) *ProductTypeRepository {
	return &ProductTypeRepository{
		db: db,
	}
}

func (r *ProductTypeRepository) GetAll(limit, offset int) ([]*models.ProductType, error) {
	var ProductTypes []*models.ProductType
	query := r.db.Limit(limit).Offset(offset)

	err := query.Find(&ProductTypes).Error
	return ProductTypes, err
}

func (r *ProductTypeRepository) GetOne(id string) (*models.ProductType, error) {
	var ProductType models.ProductType
	err := r.db.First(&ProductType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ProductType, nil
}

func (r *ProductTypeRepository) GetOneByName(name string) (*models.ProductType, error) {
	var ProductType models.ProductType
	err := r.db.First(&ProductType, "type = ?", name).Error
	if err != nil {
		return nil, err
	}
	return &ProductType, nil
}
