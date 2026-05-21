package services

import (
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
)

type RoomService struct {
	repo *repositories.RoomRepository
}

type RoomInput struct {
	Name          string  `json:"name"`
	CategoryID    uint    `json:"category_id"`
	Description   string  `json:"description"`
	PricePerNight float64 `json:"price_per_night"`
	Capacity      int     `json:"capacity"`
	ImageURL      string  `json:"image_url"`
}

func NewRoomService(r *repositories.RoomRepository) *RoomService {
	return &RoomService{repo: r}
}

func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	return s.repo.FindAll()
}

func (s *RoomService) GetRoomByID(id string) (models.Room, error) {
	return s.repo.FindByID(id)
}

func (s *RoomService) CreateRoom(input RoomInput) (models.Room, error) {
	room := models.Room{
		Name:          input.Name,
		CategoryID:    input.CategoryID,
		Description:   input.Description,
		PricePerNight: input.PricePerNight,
		Capacity:      input.Capacity,
		ImageURL:      input.ImageURL,
	}
	return s.repo.Create(room)
}

func (s *RoomService) UpdateRoom(id string, input RoomInput) (models.Room, error) {
	room := models.Room{
		Name:          input.Name,
		CategoryID:    input.CategoryID,
		Description:   input.Description,
		PricePerNight: input.PricePerNight,
		Capacity:      input.Capacity,
		ImageURL:      input.ImageURL,
	}
	return s.repo.Update(id, room)
}

func (s *RoomService) DeleteRoom(id string) error {
	return s.repo.Delete(id)
}
