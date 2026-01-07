package controllers

import (
	"net/http"
	"time"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/gin-gonic/gin"
)

type DashboardSummaryResponse struct {
	AvgLoad         float64 `json:"avgLoad"`
	AvgProductivity float64 `json:"avgProductivity"`
	Overtime        float64 `json:"overtime"`
	Satisfaction    float64 `json:"satisfaction"`

	DeptEfficiency []DeptEfficiencyItem `json:"deptEfficiency"`
	TopOvertime    []TopOvertimeItem    `json:"topOvertime"`

	MonthlyStats  []MonthlyStat       `json:"monthlyStats"`
	TopEfficiency []TopEfficiencyItem `json:"topEfficiency"`
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

type MonthlyStat struct {
	Month    string  `json:"month"`
	Load     float64 `json:"load"`
	Overtime float64 `json:"overtime"`
}

type TopEfficiencyItem struct {
	Name         string  `json:"name"`
	Productivity float64 `json:"productivity"`
	Tasks        int     `json:"tasks"`
	Overtime     float64 `json:"overtime"`
}

// DashboardSummary godoc
// @Summary Получение сводных показателей эффективности
// @Description Возвращает показатели сотрудников (admin/manager — все, employee — свои)
// @Tags dashboard
// @Accept json
// @Produce json
// @Param date_from query string false "Начало периода"
// @Param date_to query string false "Конец периода"
// @Success 200 {object} DashboardSummaryResponse
// @Failure 401 {object} map[string]string
// @Router /api/dashboard/summary [get]
// @Security BearerAuth
func DashboardSummary(c *gin.Context) {
	var summary struct {
		AvgLoad         float64
		AvgProductivity float64
		Overtime        float64
		Satisfaction    float64
	}

	// Считаем средние показатели через новую модель (WorkDay → WorkProcess, SatisfactionMetric)
	db.DB.Table("work_days wd").
		Select(`
		AVG(EXTRACT(EPOCH FROM (wd.end_work_day - wd.start_work_day))/3600) AS avg_load,
		AVG(sm.productivity) AS avg_productivity,
		AVG(sm.satisfaction) AS satisfaction,
		AVG(GREATEST(EXTRACT(EPOCH FROM (wd.end_work_day - wd.start_work_day))/3600 - 8, 0)) AS overtime
	`).
		Joins("LEFT JOIN satisfaction_metrics sm ON sm.work_day_id = wd.id AND sm.deleted_at IS NULL").
		Scan(&summary)

	// Эффективность по департаменту и позиции
	var deptData []DeptEfficiencyItem
	db.DB.Table("employee_hrs ehr").
		Select(`
		d.name AS department,
		p.name AS job_level,
		AVG(sm.productivity) AS avg_prod
	`).
		Joins("JOIN departments d ON d.id = ehr.department_id").
		Joins("JOIN positions p ON p.id = ehr.position_id").
		Joins("JOIN employees e ON e.id = ehr.employee_id").
		Joins("JOIN work_days wd ON wd.employee_id = e.id AND wd.deleted_at IS NULL").
		Joins("LEFT JOIN satisfaction_metrics sm ON sm.work_day_id = wd.id AND sm.deleted_at IS NULL").
		Where("ehr.deleted_at IS NULL").
		Group("d.name, p.name").
		Scan(&deptData)

	// Топ 3 по сверхурочным
	var top []TopOvertimeItem
	db.DB.Table("employees e").
		Select(`
		e.last_name || ' ' || e.first_name || ' ' || e.middle_name AS name,
		sm.satisfaction AS job_satisfaction,
		wp.completed_tasks AS tasks_completed_per_day,
		GREATEST(EXTRACT(EPOCH FROM (wd.end_work_day - wd.start_work_day))/3600 - 8, 0) AS overtime_hours
	`).
		Joins("JOIN work_days wd ON wd.employee_id = e.id AND wd.deleted_at IS NULL").
		Joins("LEFT JOIN work_processes wp ON wp.work_day_id = wd.id AND wp.deleted_at IS NULL").
		Joins("LEFT JOIN satisfaction_metrics sm ON sm.work_day_id = wd.id AND sm.deleted_at IS NULL").
		Order("overtime_hours DESC").
		Limit(3).
		Scan(&top)

	// Статистика по месяцам (6 месяцев)
	var monthly []MonthlyStat
	now := time.Now()
	for i := 5; i >= 0; i-- {
		m := now.AddDate(0, -i, 0)
		monthStr := m.Format("2006-01")

		var stats struct {
			Load     float64
			Overtime float64
		}

		db.DB.Table("work_days wd").
			Select(`
			AVG(EXTRACT(EPOCH FROM (wd.end_work_day - wd.start_work_day))/3600) AS load,
			AVG(GREATEST(EXTRACT(EPOCH FROM (wd.end_work_day - wd.start_work_day))/3600 - 8,0)) AS overtime
		`).
			Where("TO_CHAR(wd.start_work_day,'YYYY-MM') = ?", monthStr).
			Scan(&stats)

		monthly = append(monthly, MonthlyStat{
			Month:    monthStr,
			Load:     stats.Load,
			Overtime: stats.Overtime,
		})
	}

	// Топ 3 по продуктивности
	var topEff []TopEfficiencyItem
	db.DB.Table("employees e").
		Select(`
		e.last_name || ' ' || e.first_name || ' ' || e.middle_name AS name,
		sm.productivity AS productivity,
		wp.completed_tasks AS tasks,
		GREATEST(EXTRACT(EPOCH FROM (wd.end_work_day - wd.start_work_day))/3600 - 8,0) AS overtime
	`).
		Joins("JOIN work_days wd ON wd.employee_id = e.id AND wd.deleted_at IS NULL").
		Joins("LEFT JOIN work_processes wp ON wp.work_day_id = wd.id AND wp.deleted_at IS NULL").
		Joins("LEFT JOIN satisfaction_metrics sm ON sm.work_day_id = wd.id AND sm.deleted_at IS NULL").
		Order("productivity DESC").
		Limit(3).
		Scan(&topEff)

	c.JSON(http.StatusOK, DashboardSummaryResponse{
		AvgLoad:         summary.AvgLoad,
		AvgProductivity: summary.AvgProductivity,
		Overtime:        summary.Overtime,
		Satisfaction:    summary.Satisfaction,
		DeptEfficiency:  deptData,
		TopOvertime:     top,
		MonthlyStats:    monthly,
		TopEfficiency:   topEff,
	})
}
