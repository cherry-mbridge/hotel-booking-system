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
	
	// Services
	authService := services.NewAuthService(userRepo)
	roomService := services.NewRoomService(roomRepo)
	bookingService := services.NewBookingService(bookingRepo, roomRepo)

	// Controllers
	authController := controllers.NewAuthController(authService)
	roomController := controllers.NewRoomController(roomService)
	bookingController := controllers.NewBookingController(bookingService)

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
			userProtected.GET("/profile", authController.MeUser)
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
			adminProtected.GET("/profile", authController.MeAdmin)
			
			// Rooms CRUD
			adminProtected.POST("/rooms", roomController.Store)
			adminProtected.PUT("/rooms/:id", roomController.Update)
			adminProtected.DELETE("/rooms/:id", roomController.Destroy)

			// Bookings Management
			adminProtected.GET("/bookings", bookingController.Index)
			adminProtected.PUT("/bookings/:id/status", bookingController.UpdateStatus)
		}
	}

	// ==========================================
	// 3. Backward Compatibility & Public Endpoints
	// ==========================================
	// Legacy basic auth mapped to User endpoints
	api.POST("/register", authController.Register)
	api.POST("/login", authController.UserLogin)

	// Public Room Routes
	api.GET("/rooms", roomController.Index)
	api.GET("/rooms/:id", roomController.Show)

	// Legacy Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.UserAuthMiddleware())
	{
		protected.GET("/bookings", bookingController.UserBookings)
		protected.POST("/bookings", bookingController.Store)
	}

	// Legacy Admin Routes (pointing to AdminAuthMiddleware for unified safety)
	admin := api.Group("/admin_legacy")
	admin.Use(middleware.AdminAuthMiddleware())
	{
		// Rooms
		admin.POST("/rooms", roomController.Store)
		admin.PUT("/rooms/:id", roomController.Update)
		admin.DELETE("/rooms/:id", roomController.Destroy)
		
		// Bookings
		admin.GET("/bookings", bookingController.Index)
		admin.PUT("/bookings/:id/status", bookingController.UpdateStatus)
	}
}
