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
)

// UploadEmployees обрабатывает загрузку CSV/Excel с данными сотрудников
func UploadEmployees(c *gin.Context) {
	// Получаем user_id из контекста (middleware авторизации)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	// Получаем файл
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "файл обязателен"})
		return
	}

	// Сохраняем временно
	tempPath := fmt.Sprintf("./tmp/%d_%s", time.Now().UnixNano(), file.Filename)
	if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "не удалось создать временную папку"})
		return
	}

	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "не удалось сохранить файл"})
		return
	}
	defer os.Remove(tempPath) // удалить файл после обработки

	// Определяем расширение
	if strings.HasSuffix(file.Filename, ".csv") {
		processCSV(tempPath, userID, c)
	} else if strings.HasSuffix(file.Filename, ".xlsx") || strings.HasSuffix(file.Filename, ".xls") {
		processExcel(tempPath, userID, c)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "только CSV или Excel"})
	}
}

// --- CSV обработка ---
func processCSV(path string, userID uint, c *gin.Context) {
	f, err := os.Open(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "не удалось открыть CSV"})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.TrimLeadingSpace = true

	// Пропускаем заголовок
	_, err = reader.Read() // пропускаем заголовок
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "не удалось прочитать заголовок CSV"})
		return
	}

	var added, skipped int
	var errors []string

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, fmt.Sprintf("строка %d: %v", added+skipped+2, err))
			continue
		}

		emp, err := parseEmployeeRow(row)
		if err != nil {
			errors = append(errors, fmt.Sprintf("строка %d: %v", added+skipped+2, err))
			skipped++
			continue
		}

		emp.UserID = userID
		emp.UploadedAt = time.Now()

		if err := db.DB.Create(&emp).Error; err != nil {
			errors = append(errors, fmt.Sprintf("строка %d: %v", added+skipped+2, err))
			skipped++
			continue
		}
		added++
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "обработка завершена",
		"added":   added,
		"skipped": skipped,
		"errors":  errors,
	})
}

// --- Excel обработка ---
func processExcel(path string, userID uint, c *gin.Context) {
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "не удалось открыть Excel"})
		return
	}

	sheet := xlFile.Sheets[0]

	var added, skipped int
	var errors []string

	for i, row := range sheet.Rows {
		if i == 0 {
			continue // заголовок
		}

		cells := make([]string, len(row.Cells))
		for j, cell := range row.Cells {
			cells[j] = cell.String()
		}

		emp, err := parseEmployeeRow(cells)
		if err != nil {
			errors = append(errors, fmt.Sprintf("строка %d: %v", i+1, err))
			skipped++
			continue
		}

		emp.UserID = userID
		emp.UploadedAt = time.Now()

		if err := db.DB.Create(&emp).Error; err != nil {
			errors = append(errors, fmt.Sprintf("строка %d: %v", i+1, err))
			skipped++
			continue
		}
		added++
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "обработка завершена",
		"added":   added,
		"skipped": skipped,
		"errors":  errors,
	})
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
	if strings.ToLower(row[6]) == "yes" || strings.ToLower(row[6]) == "true" {
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
