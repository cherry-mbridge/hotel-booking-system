package controllers

import (
	"errors"
	"net/http"

	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/services"
	"lumina-hotel-api/utils"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service *services.CategoryService
}

func NewCategoryController(s *services.CategoryService) *CategoryController {
	return &CategoryController{service: s}
}

func (ctrl *CategoryController) Index(c *gin.Context) {
	var q dto.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	q.Normalize(utils.DefaultPaginateNumber)

	response, err := ctrl.service.GetCategoriesPaginated(q.Page, q.PerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (ctrl *CategoryController) All(c *gin.Context) {
	categories, err := ctrl.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (ctrl *CategoryController) Show(c *gin.Context) {
	id := c.Param("id")
	category, err := ctrl.service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func (ctrl *CategoryController) Store(c *gin.Context) {
	var input dto.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := ctrl.service.CreateCategory(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (ctrl *CategoryController) Update(c *gin.Context) {
	id := c.Param("id")
	var input dto.CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := ctrl.service.UpdateCategory(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (ctrl *CategoryController) Destroy(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteCategory(id); err != nil {
		if errors.Is(err, services.ErrCategoryHasRooms) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "CATEGORY_HAS_ROOMS",
				"message": "This category still has rooms assigned to it.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred. Please try again later.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
