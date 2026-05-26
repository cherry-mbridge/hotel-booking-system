package services

import (
	"fmt"
	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
	"time"
)

type PromoService struct {
	repo *repositories.PromoRepository
}

func NewPromoService(r *repositories.PromoRepository) *PromoService {
	return &PromoService{repo: r}
}

func (s *PromoService) GetAll() ([]models.PromotionCode, error) {
	return s.repo.GetAll()
}

func (s *PromoService) FindByID(id string) (models.PromotionCode, error) {
	return s.repo.FindByID(id)
}

func (s *PromoService) FindActiveByPromoCodeID(id string) (models.PromotionCode, error) {
	return s.repo.FindActiveByPromoCodeID(id)
}

func (s *PromoService) Create(input dto.PromoCreateInput) (models.PromotionCode, error) {

	start, _ := parseDate(input.StartDate)
	end, _ := parseDate(input.EndDate)

	promo := models.PromotionCode{
		Code:         input.Code,
		DiscountType: input.DiscountType,
		Value:        input.Value,
		StartDate:    &start,
		EndDate:      &end,
		IsActive:     input.IsActive,
		MaxUses:      input.MaxUses,
	}

	return s.repo.Create(promo)
}

func (s *PromoService) Update(id string, input dto.PromoUpdateInput) (models.PromotionCode, error) {
	start, _ := parseDate(input.StartDate)
	end, _ := parseDate(input.EndDate)

	promo := models.PromotionCode{
		Code:         input.Code,
		DiscountType: input.DiscountType,
		Value:        input.Value,
		StartDate:    &start,
		EndDate:      &end,
		IsActive:     input.IsActive,
		MaxUses:      input.MaxUses,
	}

	return s.repo.Update(id, promo)
}

func (s *PromoService) Delete(id string) error {
    promo, err := s.repo.FindByID(id)
    if err != nil {
        return err
    }

    pdcount, err := s.repo.CountUsage(id)
    if err != nil {
        return err
    }

	if pdcount > 0 {
        return fmt.Errorf("you cannot delete : promocode is pending booking.")
    }
	// promo.UsedCount > 0 need for disable deleted
    if promo.IsActive && promo.UsedCount > 0 {
        return fmt.Errorf("you cannot delete : this promo code is currently in use.")
    }

    return s.repo.Delete(id)
}

func (s *PromoService) FindByCode(code string) (models.PromotionCode, error) {
	return s.repo.FindByCode(code)
}

func parseDate(dateStr string) (time.Time, error) {

	// Postman format
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}

	// Browser format
	if t, err := time.Parse(time.RFC3339Nano, dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}

func (s *PromoService) Validate(input dto.PromoInput) error {

	start, err := parseDate(input.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date")
	}

	end, err := parseDate(input.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date")
	}

	now := time.Now()
	loc := now.Location()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, loc)

	endDay := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, loc)

	if startDay.Before(today) {
		return fmt.Errorf("start_date must be today or a future date")
	}

	if endDay.Before(today) {
		return fmt.Errorf("end_date must be today or a future date")
	}

	if startDay.After(endDay) {
		return fmt.Errorf("start_date cannot be after end_date")
	}

	switch input.DiscountType {

	case "percentage":
		if input.Value <= 0 || input.Value > 100 {
			return fmt.Errorf("percentage discount must be between 1 and 100")
		}

	case "fixed":
		if input.Value <= 0 {
			return fmt.Errorf("fixed discount must be greater than 0")
		}

	default:
		return fmt.Errorf("discount_type must be 'percentage' or 'fixed'")
	}

	return nil
}

func (s *PromoService) IsPromoDisabled(p models.PromotionCode) (bool, string) {

	now := time.Now()

	// expire at end of day
	if p.EndDate != nil {

		ed := *p.EndDate

		endOfDay := time.Date(ed.Year(), ed.Month(), ed.Day(), 23, 59, 59, 0, now.Location())

		if now.After(endOfDay) {
			return true, "expired"
		}
	}

	if p.MaxUses > 0 && p.UsedCount >= p.MaxUses {
		return true, "limit_reached"
	}

	return false, ""
}
