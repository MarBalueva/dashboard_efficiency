package controllers

import (
	"net/http"
	"strings"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/MarBalueva/dashboard_efficiency/internal/services"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Регистрация создает новый аккаунт пользователя
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   payload body RegisterRequest true "User payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "message": "Неверный формат запроса"})
		return
	}

	// Простая валидация полей
	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" || req.Password == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "message": "Заполните все поля"})
		return
	}
	if len(req.Password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "message": "Пароль должен быть не менее 5 символов"})
		return
	}

	// Проверим, не занят ли email заранее
	var exists models.User
	if err := db.DB.Where("email = ?", req.Email).First(&exists).Error; err == nil {
		// пользователь уже есть
		c.JSON(http.StatusConflict, gin.H{"error": "email_exists", "message": "Пользователь с таким email уже существует"})
		return
	}

	// Хешируем пароль
	hash, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Ошибка при обработке пароля"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: hash,
		Name:     req.Name,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		// На случай гонки или прочих ошибок БД — возвращаем общий конфликт/ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя и выдача JWT-токена
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   payload body LoginRequest true "Login payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "message": "Неверный формат запроса"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "message": "Заполните все поля"})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Пользователь не найден — даём понятный код и сообщение
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_found", "message": "Пользователь с таким email не найден"})
		return
	}

	// Проверяем пароль
	if !services.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong_password", "message": "Неверный пароль"})
		return
	}

	// Генерируем JWT
	token, err := services.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Не удалось создать токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
