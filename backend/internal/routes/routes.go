package routes

import (
	"clinic-backend/internal/controllers"
	"clinic-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/patients", controllers.CreatePatient)
		auth.GET("/patients", controllers.GetPatients)

		auth.POST("/appointments", controllers.CreateAppointment)
		auth.GET("/appointments", controllers.GetAppointments)
	}
}
