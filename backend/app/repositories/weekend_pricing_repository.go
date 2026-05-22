package repositories

import (
	"lumina-hotel-api/app/models"
	"gorm.io/gorm"
)

type WeekendPricingRepository struct {
	db *gorm.DB
}

func NewWeekendPricingRepository(db *gorm.DB) *WeekendPricingRepository {
	return &WeekendPricingRepository{db: db}
}

func (r *WeekendPricingRepository) Create(pricing models.WeekendPricing) (models.WeekendPricing, error) {
	err := r.db.Create(&pricing).Error
	return pricing, err
}

func (r *WeekendPricingRepository) FindAll() ([]models.WeekendPricing, error) {
	var rules []models.WeekendPricing
	err := r.db.Find(&rules).Error
	return rules, err
}

func (r *WeekendPricingRepository) FindByID(id string) (models.WeekendPricing, error) {
	var rule models.WeekendPricing
	err := r.db.First(&rule, id).Error
	return rule, err
}

func (r *WeekendPricingRepository) FindByRoomType(roomTypeID uint) ([]models.WeekendPricing, error) {
	var rules []models.WeekendPricing
	err := r.db.Where("room_type_id = ?", roomTypeID).Find(&rules).Error
	return rules, err
}

func (r *WeekendPricingRepository) FindActiveByRoomType(roomTypeID uint) ([]models.WeekendPricing, error) {
	var rules []models.WeekendPricing
	err := r.db.Where("room_type_id = ? AND is_active = ?", roomTypeID, true).Find(&rules).Error
	return rules, err
}

func (r *WeekendPricingRepository) Update(id string, pricing models.WeekendPricing) (models.WeekendPricing, error) {
	err := r.db.Model(&models.WeekendPricing{}).Where("id = ?", id).Updates(&pricing).Error
	if err != nil {
		return models.WeekendPricing{}, err
	}
	return r.FindByID(id)
}

func (r *WeekendPricingRepository) Delete(id string) error {
	return r.db.Delete(&models.WeekendPricing{}, id).Error
}
