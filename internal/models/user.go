package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint       `gorm:"primaryKey;column:id"`
	Login      string     `gorm:"column:login;unique;not null"`
	Password   string     `gorm:"column:password;not null"`
	EmployeeID uint       `gorm:"column:employee_id;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`

	AccessGroups []UserAccessGroup `gorm:"foreignKey:UserID"`
}

type UserAccessGroup struct {
	UserID uint
	User   User `gorm:"foreignKey:UserID"`

	AccessGroupID uint
	AccessGroup   AccessGroup `gorm:"foreignKey:AccessGroupID"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
