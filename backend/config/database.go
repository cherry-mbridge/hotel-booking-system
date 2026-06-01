package config

import (
	"log"
	"lumina-hotel-api/app/models"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("hotel.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = database.AutoMigrate(&models.User{}, &models.Category{}, &models.Room{}, &models.Booking{}, &models.WeekendPricing{})
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

	// Seed default categories if none exist
	var catCount int64
	database.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		categories := []models.Category{
			{Name: "Deluxe Room", Description: "A spacious deluxe room featuring modern amenities and elegant decor."},
			{Name: "Executive Suite", Description: "Perfect for business travellers, containing a distinct study and lounge space."},
			{Name: "Presidential Penthouse", Description: "Our signature luxury offering with panoramic city-skyline views."},
			{Name: "Standard Room", Description: "A comfortable and affordable room with essential amenities for a pleasant stay."},
			{Name: "Family Suite", Description: "A large suite with multiple beds and a living area, ideal for families."},
			{Name: "Honeymoon Suite", Description: "A romantic suite with premium furnishings and scenic views for couples."},
			{Name: "Garden Villa", Description: "A private villa surrounded by lush gardens with direct outdoor access."},
			{Name: "Poolside Cabana", Description: "A cozy cabana right by the pool with quick access to refreshments."},
			{Name: "Ocean View Room", Description: "Wake up to the sound of waves in this stunning ocean-facing room."},
			{Name: "Mountain Retreat", Description: "A serene room with majestic mountain views and fresh alpine air."},
			{Name: "City Skyline Room", Description: "Enjoy vibrant city lights and skyline views from your window."},
			{Name: "Boutique Studio", Description: "A stylish studio with contemporary design and smart amenities."},
			{Name: "Heritage Chamber", Description: "A room adorned with classical architecture and antique furnishings."},
			{Name: "Spa Suite", Description: "An in-room spa experience with a private jacuzzi and sauna access."},
			{Name: "Eco Lodge", Description: "An environmentally friendly room built with sustainable materials."},
		}
		for _, cat := range categories {
			database.Create(&cat)
		}
		log.Println("Seeded 15 default room categories")
	}

	// Seed default rooms if none exist
	var roomCount int64
	database.Model(&models.Room{}).Count(&roomCount)
	if roomCount == 0 {
		var cats []models.Category
		database.Find(&cats)
		if len(cats) == 0 {
			log.Println("No categories found; skipping room seeding")
			return
		}

		rooms := []models.Room{
			{
				CategoryID:    cats[0].ID,
				Name:          "Deluxe Ocean Breeze 301",
				Description:   "A breath-taking deluxe room with a majestic view of the sea and high-quality premium amenities.",
				PricePerNight: 150.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1566665797739-1674de7a421a?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[0].ID,
				Name:          "Deluxe Garden View 302",
				Description:   "An elegant deluxe room facing our lush organic gardens. Quiet, calm, and perfectly appointed.",
				PricePerNight: 155.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1591088398332-8a7791972843?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[1].ID,
				Name:          "Executive Business Suite 501",
				Description:   "Designed with the business traveler in mind, this room has deep workstations and luxurious resting lounges.",
				PricePerNight: 280.0,
				Capacity:      3,
				ImageURL:      "https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[2].ID,
				Name:          "Grand Lumina Presidential Penthouse",
				Description:   "The pinnacle of comfort. Features an enormous living zone, master bedroom suite, spa bath, and scenic balcony views.",
				PricePerNight: 650.0,
				Capacity:      5,
				ImageURL:      "https://images.unsplash.com/photo-1582719508461-905c673771fd?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[3].ID,
				Name:          "Standard Cozy 101",
				Description:   "A clean and comfortable standard room with all essentials for a relaxing stay.",
				PricePerNight: 85.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1611892440504-42a792e24d32?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[4].ID,
				Name:          "Family Grand Suite 201",
				Description:   "Spacious family suite with two queen beds, a kids' corner, and a large living area.",
				PricePerNight: 320.0,
				Capacity:      6,
				ImageURL:      "https://images.unsplash.com/photo-1596394516093-501ba68a0ba6?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[5].ID,
				Name:          "Honeymoon Paradise 401",
				Description:   "Romantically styled with rose petals, champagne service, and a private balcony sunset view.",
				PricePerNight: 410.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1578683010236-d716f9a3f461?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[6].ID,
				Name:          "Garden Villa Rose",
				Description:   "A charming standalone villa with a private garden, hammock, and outdoor seating area.",
				PricePerNight: 360.0,
				Capacity:      4,
				ImageURL:      "https://images.unsplash.com/photo-1587061949409-02df41d5e562?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[7].ID,
				Name:          "Poolside Cabana A1",
				Description:   "Steps from the pool bar, this cabana offers sun loungers and a shaded daybed for ultimate relaxation.",
				PricePerNight: 190.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1575425186775-b8de9a427e67?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[8].ID,
				Name:          "Ocean Horizon 701",
				Description:   "Floor-to-ceiling windows reveal endless blue waters in this premium ocean-facing suite.",
				PricePerNight: 340.0,
				Capacity:      3,
				ImageURL:      "https://images.unsplash.com/photo-1522771739844-6a9f6d5f14af?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[9].ID,
				Name:          "Mountain Escape 801",
				Description:   "A peaceful alpine retreat with a fireplace, warm timber finishes, and sweeping valley views.",
				PricePerNight: 220.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1493809842364-78817add7ffb?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[10].ID,
				Name:          "Skyline Deluxe 901",
				Description:   "Perched high above the city, this room offers dazzling night views and a sleek modern interior.",
				PricePerNight: 275.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1493809842364-78817add7ffb?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[11].ID,
				Name:          "Boutique Studio S1",
				Description:   "Compact yet chic, this studio features smart lighting, a kitchenette, and artistic decor.",
				PricePerNight: 130.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[12].ID,
				Name:          "Heritage Chamber H1",
				Description:   "Step back in time with ornate woodwork, vintage chandeliers, and royal tapestries.",
				PricePerNight: 240.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1560448204-e02f11c3d0e2?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[13].ID,
				Name:          "Spa Sanctuary 601",
				Description:   "Rejuvenate in your own private jacuzzi, steam room, and aromatherapy lounge.",
				PricePerNight: 480.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1600585154340-be6161a56a0c?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[14].ID,
				Name:          "Eco Lodge Green",
				Description:   "Built from reclaimed bamboo and solar-powered, this lodge connects you with nature responsibly.",
				PricePerNight: 170.0,
				Capacity:      3,
				ImageURL:      "https://images.unsplash.com/photo-1510798831971-661eb04b3739?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[0].ID,
				Name:          "Deluxe Sunset 303",
				Description:   "A deluxe corner room bathed in golden sunset light with a private seating nook.",
				PricePerNight: 165.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1618773928121-c32242e63f39?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[1].ID,
				Name:          "Executive Corner 502",
				Description:   "A corner executive suite with dual-aspect windows and a boardroom-style dining table.",
				PricePerNight: 310.0,
				Capacity:      3,
				ImageURL:      "https://images.unsplash.com/photo-1631049307264-da0ec9d70304?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[4].ID,
				Name:          "Family Bunk Haven 202",
				Description:   "Fun-filled family room with bunk beds, a gaming console, and a snack bar for the kids.",
				PricePerNight: 300.0,
				Capacity:      5,
				ImageURL:      "https://images.unsplash.com/photo-1555854877-bab0e564b8d5?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[8].ID,
				Name:          "Ocean Breeze Junior 702",
				Description:   "A junior ocean suite with a private terrace and direct beach boardwalk access.",
				PricePerNight: 290.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1522708323590-d24dbb6b0267?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
			{
				CategoryID:    cats[11].ID,
				Name:          "Boutique Artist Loft S2",
				Description:   "An open-plan loft with gallery walls, creative lighting, and an inspiring city view for artists.",
				PricePerNight: 145.0,
				Capacity:      2,
				ImageURL:      "https://images.unsplash.com/photo-1505693314120-0d443867891c?auto=format&fit=crop&w=1200&q=80",
				Status:        "available",
			},
		}

		for _, r := range rooms {
			database.Create(&r)
		}
		log.Println("Seeded 20 default rooms mapped to categories")
	}
}
