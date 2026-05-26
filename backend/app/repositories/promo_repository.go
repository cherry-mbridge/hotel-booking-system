package repositories

import (
	"lumina-hotel-api/app/models"

	"gorm.io/gorm"
)

type PromoRepository struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) *PromoRepository {
	return &PromoRepository{db: db}
}

func (r *PromoRepository) FindByCode(code string) (models.PromotionCode, error) {
	var promo models.PromotionCode

	err := r.db.Where("code = ? AND is_active = ?", code, true).First(&promo).Error

	return promo, err
}

func (r *PromoRepository) FindActiveByPromoCode(code string) (models.PromotionCode, error) {
	var promo models.PromotionCode

	err := r.db.Where("UPPER(code) = UPPER(?) AND is_active = ?", code, true).First(&promo).Error

	return promo, err
}

func (r *PromoRepository) FindActiveByPromoCodeID(id string) (models.PromotionCode, error) {
	var promo models.PromotionCode
	err := r.db.Where("id = ? AND is_active = ?", id, true).First(&promo).Error
	return promo, err
}

func (r *PromoRepository) FindActiveByPromoCodeIDTx(tx *gorm.DB, id string) (models.PromotionCode, error) {
	var promo models.PromotionCode
	err := tx.Where("id = ? AND is_active = ?", id, true).First(&promo).Error
	return promo, err
}

// =========================
// ADMIN SIDE
// =========================

func (r *PromoRepository) GetAll() ([]models.PromotionCode, error) {
	var promos []models.PromotionCode

	err := r.db.
		Order("created_at DESC").
		Find(&promos).Error

	return promos, err
}

func (r *PromoRepository) FindByID(id string) (models.PromotionCode, error) {
	var promo models.PromotionCode

	err := r.db.
		Where("id = ?", id).
		First(&promo).Error

	return promo, err
}

func (r *PromoRepository) Create(promo models.PromotionCode) (models.PromotionCode, error) {
	err := r.db.Create(&promo).Error

	return promo, err
}

func (r *PromoRepository) Update(id string, promo models.PromotionCode) (models.PromotionCode, error) {

	err := r.db.Model(&models.PromotionCode{}).Where("id = ?", id).Updates(promo).Error

	if err != nil {
		return models.PromotionCode{}, err
	}

	err = r.db.First(&promo, id).Error

	return promo, err
}

func (r *PromoRepository) UpdateTx(tx *gorm.DB, id string, input models.PromotionCode) (models.PromotionCode, error) {
	var promo models.PromotionCode
	err := tx.Model(&promo).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"used_count": input.UsedCount,
			"is_active":  input.IsActive,
		}).Error
	return promo, err
}

func (r *PromoRepository) Delete(id string) error {
	return r.db.Delete(&models.PromotionCode{}, id).Error
}

func (r *PromoRepository) CountUsage(id string) (int64, error) {
    var count int64

    err := r.db.Model(&models.Booking{}).
        Where("promo_code_id = ?", id).Where("status = ?", "pending").Count(&count).Error
    return count, err
}
