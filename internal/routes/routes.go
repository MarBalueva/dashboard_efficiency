package routes

import (
	"github.com/MarBalueva/dashboard_efficiency/internal/controllers"
	"github.com/MarBalueva/dashboard_efficiency/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware())
	{
		apiGroup.GET("/employees", controllers.ListCurrentUserEmployees)
		apiGroup.POST("/upload", controllers.UploadEmployees)
		apiGroup.POST("/upload/confirm", controllers.ConfirmUpload)
		apiGroup.GET("/profile", controllers.GetProfile)
		apiGroup.PUT("/profile", controllers.UpdateProfile)
		dashboard := apiGroup.Group("/dashboard")
		dashboard.GET("/summary", controllers.DashboardSummary)
	}
}
