package dto

type BookingStoreInput struct {
	RoomID    uint   `json:"room_id" binding:"required"`
	CheckIn   string `json:"check_in" binding:"required"`
	CheckOut  string `json:"check_out" binding:"required"`
	PromoCode string `json:"promo_code"`
}

type BookingUpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}
