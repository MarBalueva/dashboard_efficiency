package routes

import (
	"github.com/MarBalueva/dashboard_efficiency/internal/controllers"
	"github.com/MarBalueva/dashboard_efficiency/internal/middleware"
	"github.com/MarBalueva/dashboard_efficiency/internal/services"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware())
	{
		// Профиль
		apiGroup.GET("/profile", controllers.GetProfile)
		apiGroup.PUT("/profile", controllers.UpdateProfile)

		// Сотрудники
		employees := apiGroup.Group("/employees")
		{
			employees.GET("", services.RequireGroup("admin", "manager", "employee"), controllers.ListEmployees)
			employees.POST("", services.RequireGroup("admin", "manager"), controllers.CreateEmployee)
			employees.GET("/:id", services.RequireGroup("admin", "manager", "employee"), controllers.GetEmployeeByID)
			employees.PUT("/:id", services.RequireGroup("admin", "manager"), controllers.UpdateEmployee)
			employees.DELETE("/:id", services.RequireGroup("admin", "manager"), controllers.DeleteEmployee)
			employees.GET("/work", services.RequireGroup("admin", "manager", "employee"), controllers.GetWork)
			employees.DELETE("/work/:id", services.RequireGroup("admin", "manager"), controllers.DeleteWork)
		}

		// Загрузка данных
		upload := apiGroup.Group("/upload")
		{
			upload.POST("", services.RequireGroup("admin", "manager"), controllers.UploadEmployees)
			upload.POST("/confirm", services.RequireGroup("admin", "manager"), controllers.ConfirmUpload)
		}

		// Справочники
		dict := apiGroup.Group("/dict")
		{
			// Departments
			departments := dict.Group("/departments")
			{
				departments.GET("", controllers.ListDepartments)
				departments.POST("", services.RequireGroup("admin", "manager"), controllers.CreateDepartment)
				departments.PUT("/:id", services.RequireGroup("admin", "manager"), controllers.UpdateDepartment)
				departments.DELETE("/:id", services.RequireGroup("admin", "manager"), controllers.DeleteDepartment)
			}

			// Positions
			positions := dict.Group("/positions")
			{
				positions.GET("", controllers.ListPositions)
				positions.POST("", services.RequireGroup("admin", "manager"), controllers.CreatePosition)
				positions.PUT("/:id", services.RequireGroup("admin", "manager"), controllers.UpdatePosition)
				positions.DELETE("/:id", services.RequireGroup("admin", "manager"), controllers.DeletePosition)
			}

			// AccessGroups
			accessGroups := dict.Group("/access-groups")
			{
				accessGroups.GET("", controllers.ListAccessGroups)
				accessGroups.POST("", services.RequireGroup("admin"), controllers.CreateAccessGroup)
				accessGroups.PUT("/:id", services.RequireGroup("admin"), controllers.UpdateAccessGroup)
				accessGroups.DELETE("/:id", services.RequireGroup("admin"), controllers.DeleteAccessGroup)
			}
		}

		// Dashboard
		dashboard := apiGroup.Group("/dashboard")
		{
			dashboard.GET("/summary", controllers.DashboardSummary)
		}
	}
}
