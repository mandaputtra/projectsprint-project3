package repositories

import (
	"time"

	"github.com/mandaputtra/projectsprint-projects3/services/ms-product-svc/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(Product *models.Product) (*models.Product, error) {
	if err := r.db.Create(Product).Error; err != nil {
		return nil, err
	}
	return Product, nil
}

func (r *ProductRepository) GetAll(params map[string]interface{}) ([]*models.Product, error) {
	var activities []*models.Product
	query := r.db.Model(&models.Product{})

	if ProductType, ok := params["ProductType"].(string); ok && ProductType != "" {
		query = query.Where("Product_type_name = ?", ProductType)
	}

	if doneAtFrom, ok := params["doneAtFrom"].(time.Time); ok {
		doneAtFrom = doneAtFrom.UTC()
		query = query.Where("done_at >= ?", doneAtFrom)
	}

	if doneAtTo, ok := params["doneAtTo"].(time.Time); ok {
		doneAtTo = doneAtTo.UTC()
		query = query.Where("done_at <= ?", doneAtTo)
	}

	if min, ok := params["caloriesBurnedMin"].(int); ok && min > 0 {
		query = query.Where("calories_burned >= ?", min)
	}

	if max, ok := params["caloriesBurnedMax"].(int); ok && max > 0 {
		query = query.Where("calories_burned <= ?", max)
	}

	// Handle limit and offset
	limit := params["limit"].(int)
	offset := params["offset"].(int)
	query = query.Limit(limit).Offset(offset)
	query = query.Debug()

	// Execute query
	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *ProductRepository) GetOne(id, userId string) (*models.Product, error) {
	var Product models.Product
	err := r.db.First(&Product, "id = ? AND user_id = ?", id, userId).Error
	if err != nil {
		return nil, err
	}
	return &Product, nil
}

func (r *ProductRepository) UpdateProduct(data *models.Product) (*models.Product, error) {

	updateErr := r.db.Model(&models.Product{}).
		Where("id = ?", data.ID).
		Updates(data).
		First(data).Error

	if updateErr != nil {
		return nil, updateErr
	}
	return data, nil
}

func (r *ProductRepository) DeleteById(id string) error {
	deleteErr := r.db.Where("id = ? ", id).Delete(&models.Product{}).Error
	return deleteErr
}
