package repositories

import (
	"lumina-hotel-api/app/models"
	"gorm.io/gorm"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) FindAll() ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.Preload("Category").Find(&rooms).Error
	return rooms, err
}

func (r *RoomRepository) FindByID(id string) (models.Room, error) {
	var room models.Room
	err := r.db.Preload("Category").First(&room, id).Error
	return room, err
}

func (r *RoomRepository) Create(room models.Room) (models.Room, error) {
	err := r.db.Create(&room).Error
	return room, err
}

func (r *RoomRepository) Update(id string, room models.Room) (models.Room, error) {
	err := r.db.Model(&models.Room{}).Where("id = ?", id).Updates(room).Error
	return room, err
}

func (r *RoomRepository) Delete(id string) error {
	return r.db.Delete(&models.Room{}, id).Error
}
