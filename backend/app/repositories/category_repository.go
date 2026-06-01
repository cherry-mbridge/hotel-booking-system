package repositories

import (
	"lumina-hotel-api/app/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Preload("Rooms").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindAllPaginated(page, perPage int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	err := r.db.Model(&models.Category{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Rooms").Offset((page - 1) * perPage).Limit(perPage).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *CategoryRepository) FindByID(id string) (models.Category, error) {
	var category models.Category
	err := r.db.Preload("Rooms").First(&category, id).Error
	return category, err
}

func (r *CategoryRepository) Create(cat models.Category) (models.Category, error) {
	err := r.db.Create(&cat).Error
	return cat, err
}

func (r *CategoryRepository) Update(id string, cat models.Category) (models.Category, error) {
	err := r.db.Model(&models.Category{}).Where("id = ?", id).Updates(cat).Error
	return cat, err
}

func (r *CategoryRepository) Delete(id string) error {
	return r.db.Delete(&models.Category{}, id).Error
}

func (r *CategoryRepository) CountRoomsByCategoryID(id string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Room{}).Where("category_id = ?", id).Count(&count).Error
	return count, err
}
