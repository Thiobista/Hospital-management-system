package main

import (
	"log"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"
	"clinic-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	// Auto migrate DB tables
	config.DB.AutoMigrate(&models.User{}, &models.Patient{}, &models.Appointment{})

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("🚀 Server running on http://localhost:8081")
	r.Run(":8081")
}
