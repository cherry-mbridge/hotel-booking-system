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
