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

func (r *BookingRepository) Create(booking models.Booking) (models.Booking, error) {
	err := r.db.Create(&booking).Error
	return booking, err
}

func (r *BookingRepository) FindByUser(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("User").Preload("Room").Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindAll() ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Preload("Room").Preload("User").Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", status).Error
}
