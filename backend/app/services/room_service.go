package services

import (
	"errors"

	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
)

var ErrRoomHasActiveBookings = errors.New("cannot delete room with active bookings")

type RoomService struct {
	repo        *repositories.RoomRepository
	bookingRepo *repositories.BookingRepository
}

func NewRoomService(r *repositories.RoomRepository, br *repositories.BookingRepository) *RoomService {
	return &RoomService{repo: r, bookingRepo: br}
}

func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	return s.repo.FindAll()
}

func (s *RoomService) GetRoomsPaginated(page, perPage int) (dto.PaginationResponse[models.Room], error) {
	rooms, total, err := s.repo.FindAllPaginated(page, perPage)
	if err != nil {
		return dto.PaginationResponse[models.Room]{}, err
	}
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	return dto.PaginationResponse[models.Room]{
		Data:       rooms,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *RoomService) GetRoomByID(id string) (models.Room, error) {
	return s.repo.FindByID(id)
}

func (s *RoomService) CreateRoom(input dto.RoomInput) (models.Room, error) {
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

func (s *RoomService) UpdateRoom(id string, input dto.RoomInput) (models.Room, error) {
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
	hasActive, err := s.bookingRepo.HasActiveBookingsForRoom(id)
	if err != nil {
		return err
	}
	if hasActive {
		return ErrRoomHasActiveBookings
	}
	return s.repo.Delete(id)
}
