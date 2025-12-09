package controllers

import (
	"net/http"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
)

// ListCurrentUserEmployees возвращает данные сотрудников текущего пользователя
// @Summary Список сотрудников текущего пользователя
// @Description Возвращает все загруженные данные пользователя
// @Tags employees
// @Produce json
// @Success 200 {array} models.EmployeeSwagger
// @Failure 401 {object} map[string]string
// @Router /api/employees [get]
// @Security BearerAuth
func ListCurrentUserEmployees(c *gin.Context) {
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

	var employees []models.Employee
	if err := db.DB.Where("user_id = ?", userID).Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}
