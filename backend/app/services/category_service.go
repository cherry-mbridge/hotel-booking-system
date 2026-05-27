package services

import (
	"errors"

	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/models"
	"lumina-hotel-api/app/repositories"
)

var ErrCategoryHasRooms = errors.New("cannot delete category that still has rooms")

type CategoryService struct {
	repo    *repositories.CategoryRepository
	roomRepo *repositories.RoomRepository
}

func NewCategoryService(r *repositories.CategoryRepository, roomRepo *repositories.RoomRepository) *CategoryService {
	return &CategoryService{repo: r, roomRepo: roomRepo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.FindAll()
}

func (s *CategoryService) GetCategoriesPaginated(page, perPage int) (dto.PaginationResponse[models.Category], error) {
	categories, total, err := s.repo.FindAllPaginated(page, perPage)
	if err != nil {
		return dto.PaginationResponse[models.Category]{}, err
	}
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))
	return dto.PaginationResponse[models.Category]{
		Data:       categories,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *CategoryService) GetCategoryByID(id string) (models.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) CreateCategory(input dto.CategoryInput) (models.Category, error) {
	cat := models.Category{
		Name:        input.Name,
		Description: input.Description,
	}
	return s.repo.Create(cat)
}

func (s *CategoryService) UpdateCategory(id string, input dto.CategoryInput) (models.Category, error) {
	cat := models.Category{
		Name:        input.Name,
		Description: input.Description,
	}
	return s.repo.Update(id, cat)
}

func (s *CategoryService) DeleteCategory(id string) error {
	count, err := s.repo.CountRoomsByCategoryID(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryHasRooms
	}
	return s.repo.Delete(id)
}
