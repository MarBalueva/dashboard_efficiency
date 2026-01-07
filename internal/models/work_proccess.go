package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkDay struct {
	ID uint `gorm:"primaryKey"`

	EmployeeID uint
	Employee   Employee `gorm:"foreignKey:EmployeeID"`

	StartWorkDay time.Time `gorm:"not null;uniqueIndex:idx_employee_day"`
	EndWorkDay   time.Time `gorm:"not null;uniqueIndex:idx_employee_day"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type WorkProcess struct {
	ID uint `gorm:"primaryKey"`

	WorkDayID uint    `gorm:"uniqueIndex"`
	WorkDay   WorkDay `gorm:"foreignKey:WorkDayID"`

	CallsCount     int `gorm:"not null"`
	CompletedTasks int `gorm:"not null"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type SatisfactionMetric struct {
	ID uint `gorm:"primaryKey"`

	WorkDayID uint    `gorm:"uniqueIndex"`
	WorkDay   WorkDay `gorm:"foreignKey:WorkDayID"`

	WorkLifeBalance int `gorm:"not null"`
	Satisfaction    int `gorm:"not null"`
	Productivity    int `gorm:"not null"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type EmployeeWorkSummary struct {
	EmployeeID uint   `json:"employee_id"`
	FullName   string `json:"full_name"`

	StartWorkDay time.Time `json:"start_work_day"`
	EndWorkDay   time.Time `json:"end_work_day"`

	CallsCount     int `json:"calls_count"`
	CompletedTasks int `json:"completed_tasks"`

	WorkLifeBalance int `json:"work_life_balance"`
	Satisfaction    int `json:"satisfaction"`
	Productivity    int `json:"productivity"`
}
