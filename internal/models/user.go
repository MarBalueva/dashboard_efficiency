package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
}

// ProfileUpdateRequest используется для обновления профиля
type ProfileUpdateRequest struct {
	Name     string `json:"name" example:"Иван"`
	Password string `json:"password" example:"secret123"`
}
