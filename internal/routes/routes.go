package routes

import (
	"github.com/MarBalueva/dashboard_efficiency/internal/api"
	"github.com/MarBalueva/dashboard_efficiency/internal/api/controllers"
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
		apiGroup.GET("/health", api.GetHealth)
		apiGroup.GET("/employees", api.ListEmployees)
		apiGroup.POST("/employees", api.CreateEmployee)
		apiGroup.POST("/upload", controllers.UploadEmployees)
	}
}
