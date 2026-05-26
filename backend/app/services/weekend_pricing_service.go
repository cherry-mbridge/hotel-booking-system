package services

import (
	"fmt"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
	"strings"
	"time"
)

type WeekendPricingService struct {
	repo     *repositories.WeekendPricingRepository
	roomRepo *repositories.RoomRepository
}

func NewWeekendPricingService(r *repositories.WeekendPricingRepository, rr *repositories.RoomRepository) *WeekendPricingService {
	return &WeekendPricingService{repo: r, roomRepo: rr}
}

func (s *WeekendPricingService) GetAll() ([]models.WeekendPricing, error) {
	return s.repo.FindAll()
}

func (s *WeekendPricingService) GetByID(id string) (models.WeekendPricing, error) {
	return s.repo.FindByID(id)
}

func (s *WeekendPricingService) Create(pricing models.WeekendPricing) (models.WeekendPricing, error) {
	// Validate overlaps
	if err := s.ValidateOverlap(pricing, 0); err != nil {
		return models.WeekendPricing{}, err
	}
	// Verify that the adjustment would not cause negative final price
	// (For safety, we clamp final price to 0 anyway, but let's run simple logical validation)
	if pricing.AdjustmentValue < 0 {
		return models.WeekendPricing{}, fmt.Errorf("adjustment value must be non-negative")
	}
	return s.repo.Create(pricing)
}

func (s *WeekendPricingService) Update(id string, pricing models.WeekendPricing) (models.WeekendPricing, error) {
	var current uint
	if existing, err := s.repo.FindByID(id); err == nil {
		current = existing.ID
		pricing.RoomTypeID = existing.RoomTypeID // keep room type ID unchanged or updated
	} else {
		return models.WeekendPricing{}, err
	}

	// Validate overlaps excluding current rule
	if err := s.ValidateOverlap(pricing, current); err != nil {
		return models.WeekendPricing{}, err
	}
	if pricing.AdjustmentValue < 0 {
		return models.WeekendPricing{}, fmt.Errorf("adjustment value must be non-negative")
	}

	return s.repo.Update(id, pricing)
}

func (s *WeekendPricingService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *WeekendPricingService) GetByRoomType(roomTypeID uint) ([]models.WeekendPricing, error) {
	return s.repo.FindByRoomType(roomTypeID)
}

// ValidateOverlap checks for active overlapping rules for the same RoomTypeID
func (s *WeekendPricingService) ValidateOverlap(pricing models.WeekendPricing, excludeID uint) error {
	if !pricing.IsActive {
		return nil
	}

	rules, err := s.repo.FindActiveByRoomType(pricing.RoomTypeID)
	if err != nil {
		return err
	}

	for _, existing := range rules {
		if excludeID > 0 && existing.ID == excludeID {
			continue
		}

		// Days of week check
		daysOverlap := false
		daysA := parseDays(pricing.DaysOfWeek)
		daysB := parseDays(existing.DaysOfWeek)
		for _, dA := range daysA {
			for _, dB := range daysB {
				if dA == dB {
					daysOverlap = true
					break
				}
			}
			if daysOverlap {
				break
			}
		}

		if !daysOverlap {
			continue
		}

		// Dates overlap check
		startA := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		if pricing.StartDate != nil {
			startA = *pricing.StartDate
		}
		endA := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
		if pricing.EndDate != nil {
			endA = *pricing.EndDate
		}

		startB := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		if existing.StartDate != nil {
			startB = *existing.StartDate
		}
		endB := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
		if existing.EndDate != nil {
			endB = *existing.EndDate
		}

		// Ranges intersect if A starts before or at B's end AND B starts before or at A's end
		if (startA.Before(endB) || startA.Equal(endB)) && (startB.Before(endA) || startB.Equal(endA)) {
			return fmt.Errorf("overlapping weekend pricing rule exists with ID %d", existing.ID)
		}
	}
	return nil
}

type PriceBreakdown struct {
	BasePrice         float64 `json:"base_price"`
	WeekendAdjustment float64 `json:"weekend_adjustment"`
	SeasonalPricing   float64 `json:"seasonal_pricing"`
	PromotionDiscount float64 `json:"promotion_discount"`
	FinalPrice        float64 `json:"final_price"`
}

func (s *WeekendPricingService) CalculatePriceBreakdown(roomID, checkInStr, checkOutStr, promoCode string) (PriceBreakdown, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return PriceBreakdown{}, err
	}

	in, err := time.Parse("2006-01-02", checkInStr)
	if err != nil {
		return PriceBreakdown{}, fmt.Errorf("invalid check_in format")
	}
	out, err := time.Parse("2006-01-02", checkOutStr)
	if err != nil {
		return PriceBreakdown{}, fmt.Errorf("invalid check_out format")
	}

	if !in.Before(out) {
		return PriceBreakdown{}, fmt.Errorf("check_in must be before check_out")
	}

	// Fetch active rules for this Room (both Room ID & CategoryID supported for highest flexibility)
	rules, err := s.repo.FindActiveByRoomType(room.ID)
	if err != nil {
		return PriceBreakdown{}, err
	}

	// Also support rules mapped to category just in case
	catRules, err := s.repo.FindActiveByRoomType(room.CategoryID)
	if err == nil {
		rules = append(rules, catRules...)
	}

	var baseTotal float64
	var weekendTotal float64
	var seasonalTotal float64

	// For calculations, we step day-by-day from check-in to check-out (excluding checkout day)
	for d := in; d.Before(out); d = d.AddDate(0, 0, 1) {
		basePrice := room.PricePerNight
		baseTotal += basePrice

		// 1. Weekend pricing
		weekdayName := d.Weekday().String() // e.g. "Friday", "Saturday", "Sunday"
		var dayAdj float64

		for _, r := range rules {
			// Check if weekday matches
			matchedDay := false
			for _, dayStr := range parseDays(r.DaysOfWeek) {
				if strings.EqualFold(dayStr, weekdayName) {
					matchedDay = true
					break
				}
			}

			if !matchedDay {
				continue
			}

			// Check date range
			start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
			if r.StartDate != nil {
				start = *r.StartDate
			}
			end := time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)
			if r.EndDate != nil {
				end = *r.EndDate
			}

			// check if day 'd' is within [start, end]
			// we compare normalized dates without time for perfect match
			dNorm := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
			startNorm := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
			endNorm := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

			if (dNorm.After(startNorm) || dNorm.Equal(startNorm)) && (dNorm.Before(endNorm) || dNorm.Equal(endNorm)) {
				var diff float64
				if r.ValueType == "percentage" {
					diff = basePrice * (r.AdjustmentValue / 100.0)
				} else {
					diff = r.AdjustmentValue
				}

				if r.AdjustmentType == "decrease" {
					dayAdj -= diff
				} else {
					dayAdj += diff
				}
				// Apply one matching rule per day
				break
			}
		}
		weekendTotal += dayAdj

		// 2. Seasonal pricing: high seasons can have a markup.
		// Let's define June-August and December as warm/festive seasons with +10% seasonal dynamic markups
		month := d.Month()
		if month == time.June || month == time.July || month == time.August || month == time.December {
			seasonalTotal += basePrice * 0.10 // 10% seasonal markup
		}
	}

	// 3. Combine with promotion discount
	var promoDiscount float64
	pCode := strings.ToUpper(strings.TrimSpace(promoCode))
	subtotalBeforePromo := baseTotal + weekendTotal + seasonalTotal

	if pCode == "WELCOME10" {
		promoDiscount = subtotalBeforePromo * 0.10
	} else if pCode == "SUMMER20" {
		promoDiscount = subtotalBeforePromo * 0.20
	} else if pCode == "LUMINA30" {
		promoDiscount = subtotalBeforePromo * 0.30
	}

	finalPrice := subtotalBeforePromo - promoDiscount
	if finalPrice < 0 {
		finalPrice = 0
	}

	// Final check to guarantee "no negative final prices" validation
	if finalPrice < 0 {
		return PriceBreakdown{}, fmt.Errorf("final calculated price cannot be negative")
	}

	return PriceBreakdown{
		BasePrice:         baseTotal,
		WeekendAdjustment: weekendTotal,
		SeasonalPricing:   seasonalTotal,
		PromotionDiscount: promoDiscount,
		FinalPrice:        finalPrice,
	}, nil
}

func parseDays(daysStr string) []string {
	parts := strings.Split(daysStr, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(p))
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
