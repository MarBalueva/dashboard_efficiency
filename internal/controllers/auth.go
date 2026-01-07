package controllers

import (
	"net/http"
	"strings"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/MarBalueva/dashboard_efficiency/internal/services"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Login    string `json:"login"`
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

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_input", "message": "Заполните все поля"})
		return
	}

	var user models.User
	if err := db.DB.Preload("AccessGroups.AccessGroup").
		Where("login = ?", req.Login).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_not_found", "message": "Пользователь с таким логином не найден"})
		return
	}

	if !services.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong_password", "message": "Неверный пароль"})
		return
	}

	groups := make([]string, 0, len(user.AccessGroups))
	for _, ug := range user.AccessGroups {
		groups = append(groups, ug.AccessGroup.Name)
	}

	token, err := services.GenerateJWT(user.ID, groups)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error", "message": "Не удалось создать токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
