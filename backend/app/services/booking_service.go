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
	mailServ    *MailService
}

func NewBookingService(r *repositories.BookingRepository, rr *repositories.RoomRepository, ps *WeekendPricingService, ms *MailService) *BookingService {
	return &BookingService{repo: r, roomRepo: rr, pricingServ: ps, mailServ: ms}
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

func (s *BookingService) GetBookingByID(id string) (models.Booking, error) {
	return s.repo.FindByID(id)
}

func (s *BookingService) UpdateBookingStatus(id string, status string) (models.Booking, error) {
	booking, err := s.repo.FindByID(id)
	if err != nil {
		return models.Booking{}, err
	}

	if err := s.repo.UpdateStatus(id, status); err != nil {
		return models.Booking{}, err
	}

	booking.Status = status

	if s.mailServ != nil && booking.User.Email != "" {
		go s.mailServ.SendBookingStatusNotification(
			booking.User.Email,
			booking.User.Name,
			booking.Room.Name,
			status,
			booking.CheckIn.Format("2006-01-02"),
			booking.CheckOut.Format("2006-01-02"),
			booking.TotalPrice,
		)
	}

	return booking, nil
}
