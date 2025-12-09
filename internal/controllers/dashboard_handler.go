package controllers

import (
	"net/http"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
)

type DashboardSummaryResponse struct {
	AvgLoad         float64 `json:"avgLoad"`
	AvgProductivity float64 `json:"avgProductivity"`
	Overtime        float64 `json:"overtime"`
	Satisfaction    float64 `json:"satisfaction"`

	DeptEfficiency []DeptEfficiencyItem `json:"deptEfficiency"`
	TopOvertime    []TopOvertimeItem    `json:"topOvertime"`
}

type DeptEfficiencyItem struct {
	Department string  `json:"department"`
	JobLevel   string  `json:"job_level"`
	AvgProd    float64 `json:"avg_productivity"`
}

type TopOvertimeItem struct {
	Name            string  `json:"name"`
	JobSatisfaction int     `json:"job_satisfaction"`
	TasksPerDay     int     `json:"tasks_completed_per_day"`
	OvertimeHours   float64 `json:"overtime_hours"`
}

// DashboardSummary godoc
// @Summary Получение сводных показателей эффективности
// @Description Возвращает средние показатели сотрудников текущего пользователя
// @Tags dashboard
// @Accept json
// @Produce json
// @Param date_from query string false "Начало периода"
// @Param date_to query string false "Конец периода"
// @Success 200 {object} map[string]float64
// @Failure 401 {object} map[string]string
// @Router /api/dashboard/summary [get]
// @Security BearerAuth
func DashboardSummary(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	userID := uint(userIDRaw.(float64))

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	query := db.DB.Model(&models.Employee{}).Where("user_id = ?", userID)

	if dateFrom != "" {
		query = query.Where("uploaded_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("uploaded_at <= ?", dateTo)
	}

	// Общие метрики
	var summary struct {
		AvgLoad         float64
		AvgProductivity float64
		Overtime        float64
		Satisfaction    float64
	}

	query.Select(`
        AVG(monthly_hours_worked) AS avg_load,
        AVG(productivity_score) AS avg_productivity,
        AVG(overtime_hours_per_week) AS overtime,
        AVG(job_satisfaction) AS satisfaction
    `).Scan(&summary)

	// Эффективность по отделам + должностям
	var deptData []DeptEfficiencyItem
	query.Select(`
        department,
        job_level,
        AVG(productivity_score) AS avg_prod
    `).Group("department, job_level").Scan(&deptData)

	// Топ переработок (3 человека) — отдельный запрос
	var top []TopOvertimeItem
	db.DB.Model(&models.Employee{}).Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Select(`
			employee_id AS name,
			job_satisfaction,
			tasks_completed_per_day,
			overtime_hours_per_week AS overtime_hours
		`).
		Order("overtime_hours_per_week DESC").
		Limit(3).
		Scan(&top)

	c.JSON(http.StatusOK, DashboardSummaryResponse{
		AvgLoad:         summary.AvgLoad,
		AvgProductivity: summary.AvgProductivity,
		Overtime:        summary.Overtime,
		Satisfaction:    summary.Satisfaction,
		DeptEfficiency:  deptData,
		TopOvertime:     top,
	})
}
