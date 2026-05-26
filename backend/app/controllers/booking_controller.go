package controllers

import (
	"fmt"
	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	service *services.BookingService
}

func NewBookingController(s *services.BookingService) *BookingController {
	return &BookingController{service: s}
}

func (ctrl *BookingController) Store(c *gin.Context) {
	var input dto.BookingStoreInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(float64) // JWT returns claim as float64
	fmt.Println(userID)
	booking, err := ctrl.service.CreateBooking(uint(userID), input.RoomID, input.CheckIn, input.CheckOut, input.PromoCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

func (ctrl *BookingController) UserBookings(c *gin.Context) {
	userID := c.MustGet("userID")
	bookings, err := ctrl.service.GetUserBookings(fmt.Sprintf("%v", userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (ctrl *BookingController) Index(c *gin.Context) {
	bookings, err := ctrl.service.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func (ctrl *BookingController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var input dto.BookingUpdateStatusInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.service.UpdateBookingStatus(id, input.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated"})
}

func (ctrl *BookingController) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking deleted successfully",
	})
}

