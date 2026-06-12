package repositories

import (
	"fmt"
	"lumina-hotel-api/app/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where(fmt.Sprintf("email = '%s'", email)).First(&user).Error
	return user, err
}

func (r *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return user, err
}
