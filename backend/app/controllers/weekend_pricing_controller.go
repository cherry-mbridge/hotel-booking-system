package controllers

import (
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/services"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type WeekendPricingController struct {
	service *services.WeekendPricingService
}

func NewWeekendPricingController(s *services.WeekendPricingService) *WeekendPricingController {
	return &WeekendPricingController{service: s}
}

// GetPrice handles GET /api/rooms/:id/price
func (ctrl *WeekendPricingController) GetPrice(c *gin.Context) {
	id := c.Param("id")
	checkIn := c.Query("check_in")
	checkOut := c.Query("check_out")
	promoCode := c.Query("promo_code")

	if checkIn == "" || checkOut == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check_in and check_out are required"})
		return
	}

	breakdown, err := ctrl.service.CalculatePriceBreakdown(id, checkIn, checkOut, promoCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

// Index handles GET /api/admin/weekend-pricing
func (ctrl *WeekendPricingController) Index(c *gin.Context) {
	pricings, err := ctrl.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pricings)
}

// Show handles GET /api/admin/weekend-pricing/:id
func (ctrl *WeekendPricingController) Show(c *gin.Context) {
	id := c.Param("id")
	pricing, err := ctrl.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weekend pricing rule not found"})
		return
	}
	c.JSON(http.StatusOK, pricing)
}

// Store handles POST /api/admin/weekend-pricing
func (ctrl *WeekendPricingController) Store(c *gin.Context) {
	var input models.WeekendPricing
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pricing, err := ctrl.service.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pricing)
}

// Update handles PUT /api/admin/weekend-pricing/:id
func (ctrl *WeekendPricingController) Update(c *gin.Context) {
	id := c.Param("id")
	var input models.WeekendPricing
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pricing, err := ctrl.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pricing)
}

// Destroy handles DELETE /api/admin/weekend-pricing/:id
func (ctrl *WeekendPricingController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Weekend pricing rule deleted"})
}

// GetByRoom handles GET /api/admin/rooms/:id/weekend-pricing
func (ctrl *WeekendPricingController) GetByRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid room ID"})
		return
	}

	pricings, err := ctrl.service.GetByRoomType(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pricings)
}
