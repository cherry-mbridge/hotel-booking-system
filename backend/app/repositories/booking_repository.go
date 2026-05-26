package repositories

import (
	"lumina-hotel-api/app/models"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) DB() *gorm.DB {
	return r.db
}

func (r *BookingRepository) Create(booking models.Booking) (models.Booking, error) {
	err := r.db.Create(&booking).Error
	return booking, err
}

func (r *BookingRepository) FindByUser(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("User").Preload("Room").Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindByID(id string) (models.Booking, error) {
	var booking models.Booking
	err := r.db.First(&booking, id).Error
	return booking, err
}

func (r *BookingRepository) FindByIDTx(tx *gorm.DB, id string) (models.Booking, error) {
	var booking models.Booking
	err := tx.First(&booking, id).Error
	return booking, err
}

func (r *BookingRepository) FindAll() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Room").Preload("User").Preload("PromoCode",
		func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) HasActiveBookingsForRoom(roomID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("room_id = ? AND status IN ?", roomID, []string{"pending", "confirmed", "checked_in"}).
		Count(&count).Error
	return count > 0, err
}

func (r *BookingRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookingRepository) UpdateStatusTx(tx *gorm.DB, id string, status string) error {
	return tx.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookingRepository) FindByPromoCodeID(promoCodeID uint) (int64, error) {
	var count int64

	err := r.db.
		Model(&models.Booking{}).
		Where("promo_code_id = ? and status != ?", promoCodeID, "rejected").
		Count(&count).Error

	return count, err
}

func (r *BookingRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Booking{}).Error
}
