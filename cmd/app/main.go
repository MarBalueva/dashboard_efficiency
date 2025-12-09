package main

import (
	"log"
	"time"

	_ "github.com/MarBalueva/dashboard_efficiency/docs"
	"github.com/MarBalueva/dashboard_efficiency/internal/config"
	"github.com/MarBalueva/dashboard_efficiency/internal/db"
	"github.com/MarBalueva/dashboard_efficiency/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Employee Dashboard API
// @version 1.0
// @description API для dashbord'а эффективности сотрудников.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatal("db connect error:", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatal("migration error:", err)
	}

	r := gin.Default()
	corsCfg := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // фронт dev
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsCfg))
	routes.Register(r)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("listening on port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
