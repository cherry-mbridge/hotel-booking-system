package services

import (
	"fmt"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BookingService struct {
	repo        *repositories.BookingRepository
	roomRepo    *repositories.RoomRepository
	pricingServ *WeekendPricingService
	promoRepo   *repositories.PromoRepository
}

func NewBookingService(r *repositories.BookingRepository, rr *repositories.RoomRepository, ps *WeekendPricingService, pr *repositories.PromoRepository) *BookingService {
	return &BookingService{repo: r, roomRepo: rr, pricingServ: ps, promoRepo: pr}
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
	// SAVE PROMO RELATION
	if promoCode != "" {

		// FIX: normalize input (THIS WAS YOUR BUG)
		code := strings.TrimSpace(strings.ToUpper(promoCode))

		var promo models.PromotionCode

		promo, err := s.promoRepo.FindActiveByPromoCode(code)

		if err != nil {
			return models.Booking{}, fmt.Errorf("invalid promo code")
		}

		bookingCount, err := s.repo.FindByPromoCodeID(promo.ID)

		if err == nil && promo.IsActive {

			now := time.Now()
			// date validation
			if promo.StartDate != nil && now.Before(*promo.StartDate) {
				return models.Booking{}, fmt.Errorf("this code cannot use yet")
			}

			if promo.EndDate != nil && now.After(*promo.EndDate) {
				return models.Booking{}, fmt.Errorf("this code expired")
			}

			if promo.MaxUses > 0 && bookingCount >= int64(promo.MaxUses) {
				return models.Booking{}, fmt.Errorf("promo fully used")
			}

			// usage limit
			if promo.MaxUses > 0 && promo.UsedCount >= promo.MaxUses {
				return models.Booking{}, fmt.Errorf("invalid promo code")
			}
			booking.PromoCodeID = &promo.ID
		}
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

	return s.repo.DB().Transaction(func(tx *gorm.DB) error {
		booking, err := s.repo.FindByIDTx(tx, id)

		if err != nil {
			return err
		}

		// normalize status (VERY IMPORTANT FIX)
		status = strings.ToLower(status)

		alreadyConfirmed := booking.Status == "confirmed"

		booking.Status = status

		if err := s.repo.UpdateStatusTx(tx, id, status); err != nil {
			return err
		}

		// =========================
		// PROMO INCREMENT (FIXED)
		// =========================
		if status == "confirmed" &&
			!alreadyConfirmed &&
			booking.PromoCodeID != nil {

			promoCodeID := fmt.Sprintf("%d", *booking.PromoCodeID)

			promo, err := s.promoRepo.FindActiveByPromoCodeIDTx(tx, promoCodeID)
			if err != nil {
				return err
			}

			now := time.Now()
			// date validation
			if promo.StartDate != nil && now.Before(*promo.StartDate) {
				return fmt.Errorf("this code cannot use yet")
			}

			if promo.EndDate != nil && now.After(*promo.EndDate) {
				return fmt.Errorf("this code expired")
			}

			// usage limit
			if promo.MaxUses > 0 && promo.UsedCount >= promo.MaxUses {
				return fmt.Errorf("promo fully used")
			}

			promo.UsedCount++

			if promo.MaxUses > 0 && promo.UsedCount >= promo.MaxUses {
				promo.IsActive = false
			}

			_, err = s.promoRepo.UpdateTx(tx, promoCodeID, promo)
			if err != nil {
				return err
			}
		}

		return nil

	})

}

func (s *BookingService) Delete(id string) error {
	booking, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if booking.Status == "pending" {
		return fmt.Errorf("cannot delete pending booking")
	}

	return s.repo.Delete(id)
}
