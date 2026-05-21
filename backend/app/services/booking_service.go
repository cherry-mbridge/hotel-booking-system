package services

import (
	"fmt"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
	"time"
)

type BookingService struct {
	repo     *repositories.BookingRepository
	roomRepo *repositories.RoomRepository
}

func NewBookingService(r *repositories.BookingRepository, rr *repositories.RoomRepository) *BookingService {
	return &BookingService{repo: r, roomRepo: rr}
}

func (s *BookingService) CreateBooking(userID uint, roomID uint, checkIn, checkOut string) (models.Booking, error) {
	room, err := s.roomRepo.FindByID(fmt.Sprintf("%d", roomID))
	if err != nil {
		return models.Booking{}, err
	}

	in, _ := time.Parse("2006-01-02", checkIn)
	out, _ := time.Parse("2006-01-02", checkOut)
	
	duration := out.Sub(in).Hours() / 24
	totalPrice := duration * room.PricePerNight

	booking := models.Booking{
		UserID:     userID,
		RoomID:     roomID,
		CheckIn:    in,
		CheckOut:   out,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	return s.repo.Create(booking)
}

func (s *BookingService) GetUserBookings(userID string) ([]models.Booking, error) {
	return s.repo.FindByUser(userID)
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	return s.repo.FindAll()
}

func (s *BookingService) UpdateBookingStatus(id string, status string) error {
	return s.repo.UpdateStatus(id, status)
}
