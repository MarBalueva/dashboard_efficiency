package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"gorm.io/gorm/clause"
)

// UploadEmployees обрабатывает загрузку CSV/Excel для предпросмотра
// @Summary Предпросмотр загруженных данных
// @Description Загружает CSV или Excel файл и возвращает данные для предпросмотра, без записи в БД
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл CSV/Excel"
// @Success 200 {object} map[string]interface{} "Данные для предпросмотра"
// @Failure 400 {object} map[string]string "Ошибка при обработке файла"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Router /api/upload [post]
// @Security BearerAuth
func UploadEmployees(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "файл обязателен"})
		return
	}

	tempPath := fmt.Sprintf("./tmp/%d_%s", time.Now().UnixNano(), file.Filename)
	os.MkdirAll("./tmp", os.ModePerm)
	c.SaveUploadedFile(file, tempPath)
	defer os.Remove(tempPath)

	var preview []map[string]interface{}
	var parseErr error

	if strings.HasSuffix(file.Filename, ".csv") {
		preview, parseErr = processCSV(tempPath)
	} else if strings.HasSuffix(file.Filename, ".xlsx") || strings.HasSuffix(file.Filename, ".xls") {
		preview, parseErr = processExcel(tempPath)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "только CSV или Excel"})
		return
	}

	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": parseErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preview": preview})
}

// ConfirmUpload сохраняет данные сотрудников и кадровых метрик
// @Summary Подтверждение загрузки данных сотрудников
// @Description Получает массив объектов с данными сотрудников и сохраняет их в БД
// @Tags upload
// @Accept json
// @Produce json
// @Param data body []object true "Данные сотрудников для записи"
// @Success 200 {object} map[string]interface{} "Сообщение о добавленных/обновленных записях"
// @Failure 400 {object} map[string]string "Ошибка при обработке данных"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Router /api/upload/confirm [post]
// @Security BearerAuth
func ConfirmUpload(c *gin.Context) {
	if _, exists := c.Get("user_id"); !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	type UploadRow struct {
		EmployeeID      uint      `json:"employee_id"`
		StartWorkDay    time.Time `json:"start_work_day"`
		EndWorkDay      time.Time `json:"end_work_day"`
		CallsCount      int       `json:"calls_count"`
		CompletedTasks  int       `json:"completed_tasks"`
		WorkLifeBalance int       `json:"work_life_balance"`
		Satisfaction    int       `json:"satisfaction"`
		Productivity    int       `json:"productivity"`
	}

	var rows []UploadRow
	if err := c.ShouldBindJSON(&rows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	added := 0
	errors := []string{}
	now := time.Now()

	for _, r := range rows {

		workDay := models.WorkDay{
			EmployeeID:   r.EmployeeID,
			StartWorkDay: r.StartWorkDay,
			EndWorkDay:   r.EndWorkDay,
			CreatedAt:    now,
		}

		if err := db.DB.
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "employee_id"},
					{Name: "start_work_day"},
					{Name: "end_work_day"},
				},
				UpdateAll: true,
			}).
			Create(&workDay).Error; err != nil {

			errors = append(errors, fmt.Sprintf(
				"WorkDay EmployeeID=%d: %v", r.EmployeeID, err,
			))
			continue
		}

		workDayID := workDay.ID

		workProcess := models.WorkProcess{
			WorkDayID:      workDayID,
			CallsCount:     r.CallsCount,
			CompletedTasks: r.CompletedTasks,
			CreatedAt:      now,
		}

		if err := db.DB.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "work_day_id"}},
				UpdateAll: true,
			}).
			Create(&workProcess).Error; err != nil {

			errors = append(errors, fmt.Sprintf(
				"WorkProcess WorkDayID=%d: %v", workDayID, err,
			))
			continue
		}

		satMetric := models.SatisfactionMetric{
			WorkDayID:       workDayID,
			WorkLifeBalance: r.WorkLifeBalance,
			Satisfaction:    r.Satisfaction,
			Productivity:    r.Productivity,
			CreatedAt:       now,
		}

		if err := db.DB.
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "work_day_id"}},
				UpdateAll: true,
			}).
			Create(&satMetric).Error; err != nil {

			errors = append(errors, fmt.Sprintf(
				"SatisfactionMetric WorkDayID=%d: %v", workDayID, err,
			))
			continue
		}

		added++
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Добавлено/обновлено %d записей", added),
		"errors":  errors,
	})
}

func processCSV(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть CSV")
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	// Заголовок
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать заголовок CSV")
	}

	var preview []map[string]interface{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		data, err := parseRow(row)
		if err != nil {
			continue
		}
		preview = append(preview, data)
	}

	return preview, nil
}

func processExcel(path string) ([]map[string]interface{}, error) {
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть Excel")
	}

	sheet := xlFile.Sheets[0]
	var preview []map[string]interface{}

	for i, row := range sheet.Rows {
		if i == 0 {
			continue
		}

		cells := make([]string, len(row.Cells))
		for j, cell := range row.Cells {
			cells[j] = cell.String()
		}

		data, err := parseRow(cells)
		if err != nil {
			continue
		}
		preview = append(preview, data)
	}

	return preview, nil
}

func parseRow(row []string) (map[string]interface{}, error) {
	if len(row) < 8 {
		return nil, fmt.Errorf("не хватает столбцов")
	}

	employeeID, _ := strconv.Atoi(row[0])
	startWork, _ := time.Parse(time.RFC3339, row[1])
	endWork, _ := time.Parse(time.RFC3339, row[2])
	calls, _ := strconv.Atoi(row[3])
	tasks, _ := strconv.Atoi(row[4])
	wlb, _ := strconv.Atoi(row[5])
	satisfaction, _ := strconv.Atoi(row[6])
	productivity, _ := strconv.Atoi(row[7])

	return map[string]interface{}{
		"employee_id":       employeeID,
		"start_work_day":    startWork,
		"end_work_day":      endWork,
		"calls_count":       calls,
		"completed_tasks":   tasks,
		"work_life_balance": wlb,
		"satisfaction":      satisfaction,
		"productivity":      productivity,
	}, nil
}
