package controllers

import (
	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PromoController struct {
	serv *services.PromoService
}

func NewPromoController(s *services.PromoService) *PromoController {
	return &PromoController{serv: s}
}

// GET /admin/promos
func (ctrl *PromoController) Index(c *gin.Context) {
	promos, err := ctrl.serv.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, promos)
}

// GET /admin/promos/:id
// func (ctrl *PromoController) Show(c *gin.Context) {
// 	id := c.Param("id")

// 	promo, err := ctrl.serv.FindByID(id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"error": "promo not found",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, promo)
// }

// POST /admin/promos
func (ctrl *PromoController) Store(c *gin.Context) {
	var input dto.PromoCreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	if _, err := ctrl.serv.FindByCode(input.Code); err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "promo code already exists",
		})
		return
	}
	// VALIDATE DATES
	if err := ctrl.serv.Validate(input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	promo, err := ctrl.serv.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, promo)
}

// PUT /admin/promos/:id
func (ctrl *PromoController) Update(c *gin.Context) {
	id := c.Param("id")

	var input dto.PromoUpdateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	existingPromo, err := ctrl.serv.FindByID(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "promo_code record is not found",
		})
		return
	}

	disabled, reason := ctrl.serv.IsPromoDisabled(existingPromo)

	if disabled {

		if reason == "expired" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "promo_code is expired and cannot be updated",
			})
			return
		}

		if reason == "limit_reached" {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "promo_code usage limit has been reached and cannot be updated",
			})
			return
		}
	}

	if err := ctrl.serv.Validate(input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	promo, err := ctrl.serv.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, promo)
}

// DELETE /admin/promos/:id
func (ctrl *PromoController) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.serv.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "promo deleted successfully",
	})
}
