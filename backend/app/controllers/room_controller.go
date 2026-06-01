package controllers

import (
	"errors"
	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/helpers"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/services"
	"lumina-hotel-api/config"
	"lumina-hotel-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RoomController struct {
	service *services.RoomService
}

func NewRoomController(s *services.RoomService) *RoomController {
	return &RoomController{service: s}
}

func (ctrl *RoomController) Index(c *gin.Context) {
	var q dto.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	q.Normalize(utils.DefaultPaginateNumber)

	response, err := ctrl.service.GetRoomsPaginated(q.Page, q.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (ctrl *RoomController) All(c *gin.Context) {
	rooms, err := ctrl.service.GetAllRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func (ctrl *RoomController) Show(c *gin.Context) {
	id := c.Param("id")
	room, err := ctrl.service.GetRoomByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (ctrl *RoomController) Store(c *gin.Context) {
	var input dto.RoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {

			c.JSON(http.StatusBadRequest, gin.H{
				"errors": helpers.ValidationErrors(
					dto.RoomInput{},
					errs,
				),
			})

			return
		}
	}

	// DB validation (foreign key check)
	var category models.Category

	if err := config.DB.First(&category, input.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": gin.H{
				"category_id": "Selected category does not exist",
			},
		})
		return
	}

	room, err := ctrl.service.CreateRoom(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, room)
}

func (ctrl *RoomController) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.RoomInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {

			c.JSON(http.StatusBadRequest, gin.H{
				"errors": helpers.ValidationErrors(
					dto.RoomInput{},
					errs,
				),
			})

			return
		}
	}

	room, err := ctrl.service.UpdateRoom(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, room)
}

func (ctrl *RoomController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteRoom(id); err != nil {
		if errors.Is(err, services.ErrRoomHasActiveBookings) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "ROOM_HAS_ACTIVE_BOOKINGS",
				"message": "This room still has active bookings.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred. Please try again later.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}
