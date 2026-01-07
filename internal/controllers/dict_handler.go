package controllers

import (
	"net/http"
	"strconv"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListDepartments(c *gin.Context) {
	var items []models.Department
	ListDictionary(c, &items)
}

func CreateDepartment(c *gin.Context) {
	CreateDictionary(c, &models.Department{})
}

func UpdateDepartment(c *gin.Context) {
	UpdateDictionary(c, &models.Department{})
}

func DeleteDepartment(c *gin.Context) {
	DeleteDictionary(c, &models.Department{})
}

func ListPositions(c *gin.Context) {
	var items []models.Position
	ListDictionary(c, &items)
}

func CreatePosition(c *gin.Context) {
	CreateDictionary(c, &models.Position{})
}

func UpdatePosition(c *gin.Context) {
	UpdateDictionary(c, &models.Position{})
}

func DeletePosition(c *gin.Context) {
	DeleteDictionary(c, &models.Position{})
}

func ListAccessGroups(c *gin.Context) {
	var items []models.AccessGroup
	ListDictionary(c, &items)
}

func CreateAccessGroup(c *gin.Context) {
	CreateDictionary(c, &models.AccessGroup{})
}

func UpdateAccessGroup(c *gin.Context) {
	UpdateDictionary(c, &models.AccessGroup{})
}

func DeleteAccessGroup(c *gin.Context) {
	DeleteDictionary(c, &models.AccessGroup{})
}

// ListDepartments godoc
// @Summary Получить список отделов
// @Description Возвращает список отделов (без удалённых)
// @Tags dictionary
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Department
// @Failure 401 {object} map[string]string
// @Router /api/dict/departments [get]
// ListPositions godoc
// @Summary Получить список должностей
// @Description Возвращает список должностей (без удалённых)
// @Tags dictionary
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Position
// @Failure 401 {object} map[string]string
// @Router /api/dict/positions [get]
// ListAccessGroups godoc
// @Summary Получить список групп доступа
// @Description Возвращает групп доступа (без удалённых)
// @Tags dictionary
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.AccessGroup
// @Failure 401 {object} map[string]string
// @Router /api/dict/access-groups [get]
func ListDictionary(c *gin.Context, model interface{}) {
	if err := db.DB.Find(model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения справочника"})
		return
	}
	c.JSON(http.StatusOK, model)
}

// CreateDepartment godoc
// @Summary Создать отдел
// @Description Создание новой записи в справочнике отделов
// @Tags dictionary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.DictionaryRequest true "Данные справочника"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/dict/departments [post]
// CreatePosition godoc
// @Summary Создать должность
// @Description Добавляет новую должность в справочник
// @Tags dictionary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.DictionaryRequest true "Данные должности"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/dict/positions [post]
// CreateAccessGroup godoc
// @Summary Создать группу доступа
// @Description Создаёт новую группу доступа
// @Tags dictionary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.DictionaryRequest true "Данные группы доступа"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/dict/access-groups [post]
func CreateDictionary(c *gin.Context, model interface{}) {
	var input models.DictionaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	switch m := model.(type) {
	case *models.Department:
		m.Name = input.Name
		m.Code = input.Code
	case *models.Position:
		m.Name = input.Name
		m.Code = input.Code
	case *models.AccessGroup:
		m.Name = input.Name
		m.Code = input.Code
	}

	if err := db.DB.Create(model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания записи"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись успешно создана"})
}

// UpdateDepartment godoc
// @Summary Обновить отдел
// @Description Обновляет существующую запись справочника
// @Tags dictionary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID записи"
// @Param data body models.DictionaryRequest true "Данные справочника"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/dict/departments/{id} [put]
// UpdatePosition godoc
// @Summary Обновить должность
// @Description Обновляет данные должности
// @Tags dictionary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID должности"
// @Param data body models.DictionaryRequest true "Новые данные должности"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dict/positions/{id} [put]
// UpdateAccessGroup godoc
// @Summary Обновить группу доступа
// @Description Обновляет данные группы доступа
// @Tags dictionary
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID группы доступа"
// @Param data body models.DictionaryRequest true "Новые данные группы доступа"
// @Success 200 {object} map[string]string{message=string}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dict/access-groups/{id} [put]
func UpdateDictionary(c *gin.Context, model interface{}) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := db.DB.First(model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Запись не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения записи"})
		return
	}

	var input models.DictionaryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	switch m := model.(type) {
	case *models.Department:
		m.Name = input.Name
		m.Code = input.Code
	case *models.Position:
		m.Name = input.Name
		m.Code = input.Code
	case *models.AccessGroup:
		m.Name = input.Name
		m.Code = input.Code
	}

	if err := db.DB.Save(model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись обновлена"})
}

// DeleteDepartment godoc
// @Summary Удалить отдел
// @Description Помечает запись как удалённую (soft delete)
// @Tags dictionary
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID записи"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/dict/departments/{id} [delete]
// DeletePosition godoc
// @Summary Удалить должность
// @Description Помечает должность как удалённую
// @Tags dictionary
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID должности"
// @Success 200 {object} map[string]string{message=string}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dict/positions/{id} [delete]
// DeleteAccessGroup godoc
// @Summary Удалить должность
// @Description Помечает группу доступа как удалённую
// @Tags dictionary
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID группы доступа"
// @Success 200 {object} map[string]string{message=string}
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/dict/access-groups/{id} [delete]
func DeleteDictionary(c *gin.Context, model interface{}) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := db.DB.Delete(model, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Запись удалена"})
}
