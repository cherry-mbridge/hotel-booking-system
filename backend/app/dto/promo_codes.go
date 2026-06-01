package dto

type PromoInput struct {
	Code         string  `json:"code"`
	DiscountType string  `json:"discount_type"` // percentage | fixed
	Value        float64 `json:"value"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	MaxUses      int     `json:"max_uses"`
	IsActive     bool    `json:"is_active"`
}

type PromoCreateInput = PromoInput
type PromoUpdateInput = PromoInput
