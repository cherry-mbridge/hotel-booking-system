package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `gorm:"unique" json:"email" binding:"required,email"`
	Password  string    `json:"-"`
	Role      string    `json:"role" gorm:"default:user"` // user, admin
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Rooms       []Room    `json:"rooms,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Room struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CategoryID    uint      `json:"category_id"`
	Category      Category  `json:"category,omitempty"`
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description"`
	PricePerNight float64   `json:"price_per_night"`
	Capacity      int       `json:"capacity"`
	ImageURL      string    `json:"image_url"`
	Status        string    `json:"status" gorm:"default:available"`
	CreatedAt     time.Time `json:"created_at"`
}

type Booking struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user,omitempty"`
	RoomID     uint      `json:"room_id"`
	Room       Room      `json:"room,omitempty"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status" gorm:"default:pending"` // pending, confirmed, cancelled
	CreatedAt  time.Time `json:"created_at"`
}

type WeekendPricing struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	RoomTypeID      uint       `gorm:"column:room_type_id" json:"room_type_id" binding:"required"`
	AdjustmentType  string     `gorm:"column:adjustment_type" json:"adjustment_type" binding:"required"` // increase, decrease
	ValueType       string     `gorm:"column:value_type" json:"value_type" binding:"required"`          // percentage, fixed
	AdjustmentValue float64    `gorm:"column:adjustment_value" json:"adjustment_value" binding:"required"`
	DaysOfWeek      string     `gorm:"column:days_of_week" json:"days_of_week" binding:"required"`       // e.g. "Friday,Saturday,Sunday"
	StartDate       *time.Time `gorm:"column:start_date" json:"start_date"`                             // Nullable start date
	EndDate         *time.Time `gorm:"column:end_date" json:"end_date"`                                 // Nullable end date
	IsActive        bool       `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
}

func (WeekendPricing) TableName() string {
	return "weekend_pricings"
}
