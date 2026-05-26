package dto

type RoomInput struct {
	Name          string  `json:"name" binding:"required"`
	CategoryID    uint    `json:"category_id" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PricePerNight float64 `json:"price_per_night" binding:"required"`
	Capacity      int     `json:"capacity" binding:"required"`
	ImageURL      string  `json:"image_url"`
}
