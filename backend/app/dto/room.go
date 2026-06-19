package dto

type RoomInput struct {
	Name          string  `json:"name" binding:"required" label:"Room Name"`
	CategoryID    uint    `json:"category_id" binding:"required" label:"Category ID"`
	Description   string  `json:"description" binding:"required" label:"Description"`
	PricePerNight float64 `json:"price_per_night" binding:"required" label:"Price/Night"`
	Capacity      int     `json:"capacity" binding:"required" label:"Max Capacity"`
	ImageURL      string  `json:"image_url" label:"Image URL"`
}
