package controllers

import (
	"net/http"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile godoc
// @Summary Получить профиль текущего пользователя
// @Description Возвращает email и имя пользователя, авторизованного через JWT
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string{email=string,name=string}
// @Failure 401 {object} map[string]string{message=string}
// @Router /api/profile [get]
func GetProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	userID := uint(userIDFloat)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"name":  user.Name,
	})
}

// UpdateProfile godoc
// @Summary Обновить профиль пользователя
// @Description Позволяет изменить имя и/или пароль текущего пользователя
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param profile body models.ProfileUpdateRequest true "Данные для обновления"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string{message=string}
// @Failure 401 {object} map[string]string{message=string}
// @Failure 500 {object} map[string]string{message=string}
// @Router /api/profile [put]
func UpdateProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user id"})
		return
	}
	userID := uint(userIDFloat)

	var input models.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверные данные"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль успешно обновлён"})
}
