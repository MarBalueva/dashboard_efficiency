package auth

import (
	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	return &repository{db: db.DB}
}

func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
