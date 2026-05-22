package services

import (
	"fmt"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
	"time"
)

type BookingService struct {
	repo        *repositories.BookingRepository
	roomRepo    *repositories.RoomRepository
	pricingServ *WeekendPricingService
}

func NewBookingService(r *repositories.BookingRepository, rr *repositories.RoomRepository, ps *WeekendPricingService) *BookingService {
	return &BookingService{repo: r, roomRepo: rr, pricingServ: ps}
}

func (s *BookingService) CreateBooking(userID uint, roomID uint, checkIn, checkOut string, promoCode string) (models.Booking, error) {
	_, err := s.roomRepo.FindByID(fmt.Sprintf("%d", roomID))
	if err != nil {
		return models.Booking{}, err
	}

	in, _ := time.Parse("2006-01-02", checkIn)
	out, _ := time.Parse("2006-01-02", checkOut)
	
	breakdown, err := s.pricingServ.CalculatePriceBreakdown(fmt.Sprintf("%d", roomID), checkIn, checkOut, promoCode)
	if err != nil {
		return models.Booking{}, err
	}

	booking := models.Booking{
		UserID:     userID,
		RoomID:     roomID,
		CheckIn:    in,
		CheckOut:   out,
		TotalPrice: breakdown.FinalPrice,
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
