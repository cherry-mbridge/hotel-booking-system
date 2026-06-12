package controllers

import (
	"lumina-hotel-api/app/dto"
	"lumina-hotel-api/app/services"
	"lumina-hotel-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(s *services.AuthService) *AuthController {
	return &AuthController{service: s}
}

func getUserIDFromContext(c *gin.Context, key string) (uint, bool) {
	val, exists := c.Get(key)
	if !exists {
		return 0, false
	}
	switch v := val.(type) {
	case float64:
		return uint(v), true
	case int:
		return uint(v), true
	case uint:
		return v, true
	}
	return 0, false
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var input dto.UserRegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ctrl.service.FindByEmail(input.Email)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email Already Used."})
		return
	}

	if !utils.MinLengthRegex.MatchString(input.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}
	if !utils.LowerRegex.MatchString(input.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password must contain at least one lowercase letter"})
		return
	}
	if !utils.UpperRegex.MatchString(input.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password must contain at least one uppercase letter"})
		return
	}
	if !utils.DigitRegex.MatchString(input.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password must contain at least one number"})
		return
	}
	if !utils.SpecialRegex.MatchString(input.Password) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Password must contain at least one special character"})
		return
	}

	user, err := ctrl.service.Register(input.Name, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// func (ctrl *AuthController) Login(c *gin.Context)  {
// 	ctrl.UserLogin(c)
// }

func (ctrl *AuthController) UserLogin(c *gin.Context) {
	var input dto.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := ctrl.service.LoginUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (ctrl *AuthController) AdminLogin(c *gin.Context) {
	var input dto.UserLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := ctrl.service.LoginAdmin(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (ctrl *AuthController) UserLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully from User Guard"})
}

func (ctrl *AuthController) AdminLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully from Admin Guard"})
}

// func (ctrl *AuthController) MeUser(c *gin.Context) {
// 	id, ok := getUserIDFromContext(c, "userID")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
// 		return
// 	}

// 	user, err := ctrl.service.GetProfile(id)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

// func (ctrl *AuthController) MeAdmin(c *gin.Context) {
// 	id, ok := getUserIDFromContext(c, "adminID")
// 	if !ok {
// 		// Fallback to userID if adminID was set there
// 		id, ok = getUserIDFromContext(c, "userID")
// 	}
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized context"})
// 		return
// 	}

// 	user, err := ctrl.service.GetProfile(id)
// 	if err != nil || user.Role != "admin" {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }
