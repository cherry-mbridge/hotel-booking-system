package services

import (
	"errors"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(r *repositories.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(name, email, password string) (models.User, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	return s.repo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, models.User, error) {
	return s.LoginUser(email, password)
}

func (s *AuthService) FindByEmail(email string) (models.User, error) {
	user, err := s.repo.FindByEmail(email)
	return user, err
}

func (s *AuthService) LoginUser(email, password string) (string, models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", models.User{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", models.User{}, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"guard": "user",
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", models.User{}, err
	}

	return tokenString, user, nil
}

func (s *AuthService) LoginAdmin(email, password string) (string, models.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", models.User{}, errors.New("invalid credentials")
	}

	if user.Role != "admin" {
		return "", models.User{}, errors.New("unauthorized: administrative access only")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", models.User{}, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"guard": "admin",
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", models.User{}, err
	}

	return tokenString, user, nil
}

func (s *AuthService) GetProfile(id uint) (models.User, error) {
	return s.repo.FindByID(id)
}
