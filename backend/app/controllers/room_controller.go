package controllers

import (
	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	service *services.RoomService
}

func NewRoomController(s *services.RoomService) *RoomController {
	return &RoomController{service: s}
}

func (ctrl *RoomController) Index(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted"})
}
