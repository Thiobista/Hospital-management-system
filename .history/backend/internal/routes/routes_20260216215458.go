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
		// Patient routes
		auth.POST("/patients", controllers.CreatePatient)
		auth.GET("/patients", controllers.GetPatients)
		auth.GET("/patients/:id", controllers.GetPatientByID)
		auth.PUT("/patients/:id", controllers.UpdatePatient)
		auth.DELETE("/patients/:id", controllers.DeletePatient)

		// Appointment routes
		auth.POST("/appointments", controllers.CreateAppointment)
		auth.GET("/appointments", controllers.GetAppointments)
		auth.GET("/appointments/:id", controllers.GetAppointmentByID)
		auth.PUT("/appointments/:id", controllers.UpdateAppointment)
		auth.DELETE("/appointments/:id", controllers.DeleteAppointment)
	}
}
