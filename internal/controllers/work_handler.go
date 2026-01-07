package controllers

import (
	"net/http"
	"strconv"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/MarBalueva/dashboard_efficiency/internal/services"
	"github.com/gin-gonic/gin"
)

// GetEmployeesTable godoc
//
//	@Summary		Получить таблицу сотрудников
//	@Description	employee видит только свои данные, admin и manager — все
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.EmployeeWorkSummary
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/employees/work [get]
//
// @Security BearerAuth
func GetWork(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	query := db.DB.
		Table("work_days wd").
		Select(`
		wd.employee_id,
		CONCAT(e.last_name, ' ', e.first_name, ' ', e.middle_name) AS full_name,
		wd.start_work_day,
		wd.end_work_day,
		wp.calls_count,
		wp.completed_tasks,
		sm.work_life_balance,
		sm.satisfaction,
		sm.productivity
	`).
		Joins(`
		JOIN work_processes wp 
		  ON wp.work_day_id = wd.id 
		 AND wp.deleted_at IS NULL
	`).
		Joins(`
		JOIN satisfaction_metrics sm 
		  ON sm.work_day_id = wd.id 
		 AND sm.deleted_at IS NULL
	`).
		Joins("JOIN employees e ON e.id = wd.employee_id").
		Where("wd.deleted_at IS NULL")

	if services.HasAnyGroup(user, "employee") && !services.HasAnyGroup(user, "admin") && !services.HasAnyGroup(user, "manager") {
		query = query.Where("wd.employee_id = ?", user.EmployeeID)
	}

	var result []models.EmployeeWorkSummary
	if err := query.Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteEmployee godoc
//
//	@Summary		Удалить сотрудника
//	@Description	Удаляет рабочий день, процессы и метрики сотрудника по employee_id (soft delete)
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			employee_id	path	int	true	"ID сотрудника"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/employees/work/{employee_id} [delete]
//
// @Security BearerAuth
func DeleteWork(c *gin.Context) {
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDAny.(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	paramID, err := strconv.ParseUint(c.Param("employee_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	if services.HasAnyGroup(user, "employee") && user.EmployeeID != uint(paramID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	tx := db.DB.Begin()

	tx.Where("employee_id = ?", paramID).Delete(&models.WorkDay{})
	tx.Where("employee_id = ?", paramID).Delete(&models.WorkProcess{})
	tx.Where("employee_id = ?", paramID).Delete(&models.SatisfactionMetric{})

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
