package routes

import (
	"lumina-hotel-api/app/controllers"
	"lumina-hotel-api/app/middleware"
	"lumina-hotel-api/app/repositories"
	"lumina-hotel-api/app/services"
	"lumina-hotel-api/config"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	roomRepo := repositories.NewRoomRepository(config.DB)
	bookingRepo := repositories.NewBookingRepository(config.DB)
	weekendPricingRepo := repositories.NewWeekendPricingRepository(config.DB)

	// Services
	authService := services.NewAuthService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	weekendPricingService := services.NewWeekendPricingService(weekendPricingRepo, roomRepo)
	bookingService := services.NewBookingService(bookingRepo, roomRepo, weekendPricingService)

	// Controllers
	authController := controllers.NewAuthController(authService)
	roomController := controllers.NewRoomController(roomService)
	bookingController := controllers.NewBookingController(bookingService)
	weekendPricingController := controllers.NewWeekendPricingController(weekendPricingService)

	// ==========================================
	// 1. Separate Guard: User Auth & Protected
	// ==========================================
	userGroup := api.Group("/user")
	{
		userGroup.POST("/register", authController.Register)
		userGroup.POST("/login", authController.UserLogin)
		userGroup.POST("/logout", authController.UserLogout)

		// Protected User routes
		userProtected := userGroup.Group("/")
		userProtected.Use(middleware.UserAuthMiddleware())
		{
			// userProtected.GET("/profile", authController.MeUser)
			userProtected.GET("/bookings", bookingController.UserBookings)
			userProtected.POST("/bookings", bookingController.Store)
		}
	}

	// ==========================================
	// 2. Separate Guard: Admin Auth & Protected
	// ==========================================
	adminGroup := api.Group("/admin")
	{
		adminGroup.POST("/login", authController.AdminLogin)
		adminGroup.POST("/logout", authController.AdminLogout)

		// Protected Admin routes
		adminProtected := adminGroup.Group("/")
		adminProtected.Use(middleware.AdminAuthMiddleware())
		{
			// adminProtected.GET("/profile", authController.MeAdmin)

			// Rooms CRUD
			adminProtected.POST("/rooms", roomController.Store)
			adminProtected.PUT("/rooms/:id", roomController.Update)
			adminProtected.DELETE("/rooms/:id", roomController.Destroy)

			// Bookings Management
			adminProtected.GET("/bookings", bookingController.Index)
			adminProtected.PUT("/bookings/:id/status", bookingController.UpdateStatus)

			// Weekend pricing management
			adminProtected.GET("/weekend-pricing", weekendPricingController.Index)
			adminProtected.GET("/weekend-pricing/:id", weekendPricingController.Show)
			adminProtected.POST("/weekend-pricing", weekendPricingController.Store)
			adminProtected.PUT("/weekend-pricing/:id", weekendPricingController.Update)
			adminProtected.DELETE("/weekend-pricing/:id", weekendPricingController.Destroy)
			adminProtected.GET("/rooms/:id/weekend-pricing", weekendPricingController.GetByRoom)
		}
	}

	// Public Room Routes
	api.GET("/rooms", roomController.Index)
	api.GET("/rooms/:id", roomController.Show)
	api.GET("/rooms/:id/price", weekendPricingController.GetPrice)

	// Legacy Protected Routes
	// protected := api.Group("/")
	// protected.Use(middleware.UserAuthMiddleware())
	// {
	// 	protected.GET("/bookings", bookingController.UserBookings)
	// 	protected.POST("/bookings", bookingController.Store)
	// }
}
