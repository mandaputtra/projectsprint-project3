package repositories

import (
	"strconv"
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
	var products []*models.Product
	query := r.db.Model(&models.Product{})

	// Filter by productId
	if productId, ok := params["productId"].(string); ok && productId != "" {
		query = query.Where("id = ?", productId)
	}

	// Filter by sku
	if sku, ok := params["sku"].(string); ok && sku != "" {
		query = query.Where("sku = ?", sku)
	}

	// Filter by category
	if category, ok := params["category"].(string); ok && category != "" {
		query = query.Where("category_name ILIKE ?", "%"+category+"%")
	}

	// Filter by exact search (name)
	if search, ok := params["search"].(string); ok && search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// Handle sorting (newest, cheapest, sold-x)
	if sortBy, ok := params["sortBy"].(string); ok && sortBy != "" {
		switch {
		case sortBy == "newest":
			query = query.Order("created_at DESC, updated_at DESC")
		case sortBy == "cheapest":
			query = query.Order("price ASC")
		case len(sortBy) > 5 && sortBy[:5] == "sold-":
			if seconds, err := strconv.Atoi(sortBy[5:]); err == nil {
				timeLimit := time.Now().Add(-time.Duration(seconds) * time.Second)
				query = query.Where("sold_at >= ?", timeLimit).Order("sold_at DESC")
			}
		}
	}

	// Handle limit and offset
	if limit, ok := params["limit"].(int); ok && limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(5) // Default limit
	}

	if offset, ok := params["offset"].(int); ok && offset >= 0 {
		query = query.Offset(offset)
	} else {
		query = query.Offset(0) // Default offset
	}

	// Debug query for troubleshooting
	query = query.Debug()

	// Execute query
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
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
