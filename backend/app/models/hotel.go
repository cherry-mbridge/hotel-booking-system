package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `json:"role" gorm:"default:user"` // user, admin
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	Rooms       []Room         `json:"rooms,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Room struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CategoryID    uint           `gorm:"not null" json:"category_id"`
	Category      *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name          string         `gorm:"not null" json:"name"`
	Description   string         `gorm:"type:varchar(255)" json:"description"`
	PricePerNight float64        `gorm:"not null" json:"price_per_night"`
	Capacity      int            `gorm:"not null" json:"capacity"`
	ImageURL      string         `json:"image_url"`
	Status        string         `gorm:"default:available;not null" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Booking struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RoomID     uint      `gorm:"not null" json:"room_id"`
	Room       Room      `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	CheckIn    time.Time `gorm:"not null" json:"check_in"`
	CheckOut   time.Time `gorm:"not null" json:"check_out"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
	Status     string    `gorm:"default:pending;not null" json:"status"` // pending, confirmed, checked_in, cancelled, completed, rejected
	CreatedAt  time.Time `json:"created_at"`
}

type WeekendPricing struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	RoomTypeID      uint       `gorm:"not null" json:"room_type_id"`
	AdjustmentType  string     `gorm:"not null" json:"adjustment_type"` // increase, decrease
	ValueType       string     `gorm:"not null" json:"value_type"`      // percentage, fixed
	AdjustmentValue float64    `gorm:"not null" json:"adjustment_value"`
	DaysOfWeek      string     `gorm:"not null" json:"days_of_week"` // e.g. "Friday,Saturday,Sunday"
	StartDate       *time.Time `json:"start_date"`                   // Nullable start date
	EndDate         *time.Time `json:"end_date"`                     // Nullable end date
	IsActive        bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (WeekendPricing) TableName() string {
	return "weekend_pricings"
}
