package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/MarBalueva/dashboard_efficiency/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListEmployees godoc
// @Summary Список сотрудников
// @Description admin и manager видят всех, employee — только себя
// @Tags employees
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.EmployeeFullResponse
// @Failure 401 {object} map[string]string
// @Router /api/employees [get]
func ListEmployees(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := db.DB.
		Preload("AccessGroups.AccessGroup").
		First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	query := db.DB.
		Model(&models.EmployeeHR{}).
		Preload("Employee").
		Preload("Department").
		Preload("Position")

	if !services.HasAnyGroup(user, "admin", "manager") {
		query = query.Where("employee_id = ?", user.EmployeeID)
	}

	var hrList []models.EmployeeHR
	if err := query.Find(&hrList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	result := make([]models.EmployeeFullResponse, 0, len(hrList))
	for _, hr := range hrList {
		result = append(result, models.EmployeeFullResponse{
			ID:         hr.Employee.ID,
			LastName:   hr.Employee.LastName,
			FirstName:  hr.Employee.FirstName,
			MiddleName: hr.Employee.MiddleName,
			Department: hr.Department.Name,
			Position:   hr.Position.Name,
			IsRemote:   hr.IsRemote,
			BirthDate:  hr.BirthDate,
			HireDate:   hr.HireDate,
			FireDate:   hr.FireDate,
			Salary:     hr.Salary,
			CreatedAt:  hr.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetEmployeeByID godoc
// @Summary Получить сотрудника по ID
// @Description
// admin и manager могут получать любого сотрудника,
// employee — только себя
// @Tags employees
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID сотрудника"
// @Success 200 {object} models.EmployeeFullResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/employees/{id} [get]
func GetEmployeeByID(c *gin.Context) {
	userID := c.GetUint("user_id")

	employeeIDParam := c.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var user models.User
	if err := db.DB.
		Preload("AccessGroups.AccessGroup").
		First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	if !services.HasAnyGroup(user, "admin", "manager") {
		if user.EmployeeID != uint(employeeID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

	var hr models.EmployeeHR
	if err := db.DB.
		Preload("Employee").
		Preload("Department").
		Preload("Position").
		Where("employee_id = ?", employeeID).
		First(&hr).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	response := models.EmployeeFullResponse{
		ID:         hr.Employee.ID,
		LastName:   hr.Employee.LastName,
		FirstName:  hr.Employee.FirstName,
		MiddleName: hr.Employee.MiddleName,
		Department: hr.Department.Name,
		Position:   hr.Position.Name,
		IsRemote:   hr.IsRemote,
		BirthDate:  hr.BirthDate,
		HireDate:   hr.HireDate,
		FireDate:   hr.FireDate,
		Salary:     hr.Salary,
		CreatedAt:  hr.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// CreateEmployee godoc
// @Summary Создать сотрудника
// @Description Создаёт сотрудника и связанные кадровые данные (Department, Position, HR-поля)
// @Tags employees
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.EmployeeCreateRequest true "Данные сотрудника и кадровые данные"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employees [post]
func CreateEmployee(c *gin.Context) {
	var input models.EmployeeCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_date format, expected YYYY-MM-DD"})
		return
	}
	hireDate, err := time.Parse("2006-01-02", input.HireDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hire_date format, expected YYYY-MM-DD"})
		return
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		employee := models.Employee{
			LastName:   input.LastName,
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
		}
		if err := tx.Create(&employee).Error; err != nil {
			return err
		}

		hr := models.EmployeeHR{
			EmployeeID:   employee.ID,
			DepartmentID: input.DepartmentID,
			PositionID:   input.PositionID,
			IsRemote:     input.IsRemote,
			BirthDate:    birthDate,
			HireDate:     hireDate,
			Salary:       input.Salary,
		}
		return tx.Create(&hr).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Сотрудник создан"})
}

// UpdateEmployee godoc
// @Summary Обновить сотрудника
// @Description Обновляет персональные и кадровые данные сотрудника
// @Tags employees
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID сотрудника"
// @Param data body models.EmployeeCreateRequest true "Новые данные сотрудника"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employees/{id} [put]
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	var input models.EmployeeCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid birth_date format, expected YYYY-MM-DD"})
		return
	}
	hireDate, err := time.Parse("2006-01-02", input.HireDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hire_date format, expected YYYY-MM-DD"})
		return
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Employee{}).
			Where("id = ?", id).
			Updates(models.Employee{
				LastName:   input.LastName,
				FirstName:  input.FirstName,
				MiddleName: input.MiddleName,
			}).Error; err != nil {
			return err
		}

		return tx.Model(&models.EmployeeHR{}).
			Where("employee_id = ?", id).
			Updates(models.EmployeeHR{
				DepartmentID: input.DepartmentID,
				PositionID:   input.PositionID,
				IsRemote:     input.IsRemote,
				BirthDate:    birthDate,
				HireDate:     hireDate,
				Salary:       input.Salary,
			}).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Данные обновлены"})
}

// DeleteEmployee godoc
// @Summary Удалить сотрудника
// @Description Помечает сотрудника и его кадровые данные как удалённые (soft delete)
// @Tags employees
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID сотрудника"
// @Success 200 {object} map[string]string{message=string}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/employees/{id} [delete]
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.EmployeeHR{}, "employee_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Employee{}, id).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Сотрудник удалён"})
}
