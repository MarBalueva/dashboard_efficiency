package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	EmployeeID           string  `gorm:"uniqueIndex;not null" json:"employee_id"`
	Age                  int     `json:"age"`
	Department           string  `json:"department"`
	JobLevel             string  `json:"job_level"`
	YearsAtCompany       int     `json:"years_at_company"`
	MonthlyHoursWorked   int     `json:"monthly_hours_worked"`
	RemoteWork           bool    `json:"remote_work"`
	MeetingsPerWeek      int     `json:"meetings_per_week"`
	TasksCompletedPerDay int     `json:"tasks_completed_per_day"`
	OvertimeHoursPerWeek float64 `json:"overtime_hours_per_week"`
	WorkLifeBalance      string  `json:"work_life_balance"`
	JobSatisfaction      int     `json:"job_satisfaction"`
	ProductivityScore    float64 `json:"productivity_score"`
	AnnualSalary         float64 `json:"annual_salary"`
	AbsencesPerYear      int     `json:"absences_per_year"`

	UserID     uint      `json:"user_id"`     // кто загрузил
	UploadedAt time.Time `json:"uploaded_at"` // дата загрузки
}
