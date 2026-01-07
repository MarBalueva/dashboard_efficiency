package controllers

import (
	"net/http"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetProfile godoc
// @Summary Получить профиль текущего пользователя
// @Description Возвращает данные авторизованного пользователя и связанного с ним сотрудника.
// @Description
// @Description Данные пользователя:
// @Description - login
// @Description
// @Description Данные сотрудника:
// @Description - фамилия
// @Description - имя
// @Description - отчество
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Профиль пользователя"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 404 {object} map[string]string "Пользователь или сотрудник не найден"
// @Router /api/profile [get]
func GetProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDRaw.(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var employee models.Employee
	if err := db.DB.First(&employee, user.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"login": user.Login,
		"employee": gin.H{
			"last_name":   employee.LastName,
			"first_name":  employee.FirstName,
			"middle_name": employee.MiddleName,
		},
	})
}

// UpdateProfile godoc
// @Summary Обновить профиль текущего пользователя
// @Description Позволяет изменить пароль пользователя и/или ФИО сотрудника.
// @Description
// @Description Можно передавать только те поля, которые требуется изменить.
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param profile body models.ProfileUpdateRequest true "Данные для обновления профиля"
// @Success 200 {object} map[string]string "Профиль успешно обновлён"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 404 {object} map[string]string "Пользователь или сотрудник не найден"
// @Failure 500 {object} map[string]string "Ошибка при обновлении профиля"
// @Router /api/profile [put]
func UpdateProfile(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDRaw.(uint)

	var input models.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var employee models.Employee
	if err := db.DB.First(&employee, user.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}

	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password hash error"})
			return
		}
		user.Password = string(hashed)
	}

	if input.LastName != "" {
		employee.LastName = input.LastName
	}
	if input.FirstName != "" {
		employee.FirstName = input.FirstName
	}
	if input.MiddleName != "" {
		employee.MiddleName = input.MiddleName
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		if err := tx.Save(&employee).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления профиля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль успешно обновлён"})
}
