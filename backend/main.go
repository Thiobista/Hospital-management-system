package main

import (
	"log"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"
	"clinic-backend/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func main() {
	config.ConnectDB()

	// Auto migrate DB tables
	config.DB.AutoMigrate(&models.User{}, &models.Patient{}, &models.Appointment{})

	r := gin.Default()

r.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:3000"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge: 12 * time.Hour,
}))

	routes.SetupRoutes(r)

	log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
