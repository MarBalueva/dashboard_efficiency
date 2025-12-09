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
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user id"})
		return
	}
	userID := uint(userIDFloat)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "файл обязателен"})
		return
	}

	tempPath := fmt.Sprintf("./tmp/%d_%s", time.Now().UnixNano(), file.Filename)
	os.MkdirAll("./tmp", os.ModePerm)
	c.SaveUploadedFile(file, tempPath)
	defer os.Remove(tempPath)

	var rows []models.Employee
	var parseErr error

	if strings.HasSuffix(file.Filename, ".csv") {
		rows, parseErr = processCSV(tempPath, userID)
	} else if strings.HasSuffix(file.Filename, ".xlsx") || strings.HasSuffix(file.Filename, ".xls") {
		rows, parseErr = processExcel(tempPath, userID)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "только CSV или Excel"})
		return
	}

	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": parseErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preview": rows,
	})
}

// @Summary Подтверждение загрузки данных сотрудников
// @Description Получает массив объектов Employee и сохраняет их в БД
// @Tags upload
// @Accept json
// @Produce json
// @Param data body []models.Employee true "Данные сотрудников для записи"
// @Success 200 {object} map[string]interface{} "Сообщение о добавленных/обновленных записях"
// @Failure 400 {object} map[string]string "Ошибка при обработке данных"
// @Failure 401 {object} map[string]string "Пользователь не авторизован"
// @Router /api/upload/confirm [post]
// @Security BearerAuth
func ConfirmUpload(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	userIDFloat, ok := userIDRaw.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid user id"})
		return
	}
	userID := uint(userIDFloat)

	var rows []models.Employee
	if err := c.ShouldBindJSON(&rows); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	added := 0
	errors := []string{}

	for _, emp := range rows {
		emp.UserID = userID
		emp.UploadedAt = time.Now()

		// Если есть EmployeeID — обновляем, иначе создаем
		if err := db.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "employee_id"}, {Name: "user_id"}},
			UpdateAll: true,
		}).Create(&emp).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Ошибка EmployeeID=%s: %v", emp.EmployeeID, err))
			continue
		}
		added++
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Добавлено/обновлено %d записей", added),
		"errors":  errors,
	})
}

func processCSV(path string, userID uint) ([]models.Employee, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть CSV")
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	_, err = reader.Read() // skip header
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать заголовок CSV")
	}

	var rows []models.Employee
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		emp, err := parseEmployeeRow(row)
		if err != nil {
			continue
		}

		emp.UserID = userID
		rows = append(rows, emp)

	}

	return rows, nil
}

// --- Excel обработка с обновлением существующих ---
func processExcel(path string, userID uint) ([]models.Employee, error) {
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть Excel")
	}

	sheet := xlFile.Sheets[0]

	var rows []models.Employee

	for i, row := range sheet.Rows {
		if i == 0 {
			continue
		}

		cells := make([]string, len(row.Cells))
		for j, cell := range row.Cells {
			cells[j] = cell.String()
		}

		emp, err := parseEmployeeRow(cells)
		if err != nil {
			continue
		}

		emp.UserID = userID
		rows = append(rows, emp)
	}

	return rows, nil
}

// --- Парсинг строки в Employee ---
func parseEmployeeRow(row []string) (models.Employee, error) {
	if len(row) < 15 {
		return models.Employee{}, fmt.Errorf("не хватает столбцов")
	}

	age, _ := strconv.Atoi(row[1])
	yearsAt, _ := strconv.Atoi(row[4])
	monthlyHours, _ := strconv.Atoi(row[5])
	meetings, _ := strconv.Atoi(row[7])
	tasks, _ := strconv.Atoi(row[8])
	overtime, _ := strconv.ParseFloat(row[9], 64)
	jobSatisfaction, _ := strconv.Atoi(row[11])
	productivity, _ := strconv.ParseFloat(row[12], 64)
	annualSalary, _ := strconv.ParseFloat(row[13], 64)
	absences, _ := strconv.Atoi(row[14])

	remote := false
	if strings.ToLower(row[6]) == "yes" ||
		strings.ToLower(row[6]) == "true" ||
		strings.ToLower(row[6]) == "hybrid" {
		remote = true
	}

	return models.Employee{
		EmployeeID:           row[0],
		Age:                  age,
		Department:           row[2],
		JobLevel:             row[3],
		YearsAtCompany:       yearsAt,
		MonthlyHoursWorked:   monthlyHours,
		RemoteWork:           remote,
		MeetingsPerWeek:      meetings,
		TasksCompletedPerDay: tasks,
		OvertimeHoursPerWeek: overtime,
		WorkLifeBalance:      row[10],
		JobSatisfaction:      jobSatisfaction,
		ProductivityScore:    productivity,
		AnnualSalary:         annualSalary,
		AbsencesPerYear:      absences,
	}, nil
}
