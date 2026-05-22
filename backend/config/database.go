package config

import (
	"lumina-hotel-api/app/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("hotel.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = database.AutoMigrate(&models.User{}, &models.Room{}, &models.Category{}, &models.Booking{}, &models.WeekendPricing{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	DB = database

	// Seed default admin user if none exists
	var adminCount int64
	database.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Name:     "Lumina Admin",
			Email:    "admin@lumina.com",
			Password: string(hashed),
			Role:     "admin",
		}
		if err := database.Create(&admin).Error; err != nil {
			log.Println("Failed to seed admin user:", err)
		} else {
			log.Println("Seeded default admin user: email: admin@lumina.com / password: admin123")
		}
	}
}
