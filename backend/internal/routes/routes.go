package routes

import (
	"clinic-backend/internal/controllers"
	"clinic-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes - require authentication
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// Dashboard routes
		auth.GET("/dashboard/admin", middleware.AdminOnly(), controllers.GetAdminDashboard)
		auth.GET("/dashboard/doctor", middleware.DoctorOnly(), controllers.GetDoctorDashboard)
		auth.GET("/dashboard/receptionist", middleware.ReceptionistOnly(), controllers.GetReceptionistDashboard)

		// Doctor routes - Admin only for create/update/delete
		auth.POST("/doctors", middleware.AdminOnly(), controllers.CreateDoctor)
		auth.GET("/doctors", controllers.GetDoctors)
		auth.GET("/doctors/:id", controllers.GetDoctorByID)
		auth.PUT("/doctors/:id", middleware.AdminOnly(), controllers.UpdateDoctor)
		auth.DELETE("/doctors/:id", middleware.AdminOnly(), controllers.DeleteDoctor)

		// Patient routes - Admin and Receptionist can manage
		auth.POST("/patients", middleware.AdminOrReceptionist(), controllers.CreatePatient)
		auth.GET("/patients", controllers.GetPatients)
		auth.GET("/patients/:id", controllers.GetPatientByID)
		auth.PUT("/patients/:id", middleware.AdminOrReceptionist(), controllers.UpdatePatient)
		auth.DELETE("/patients/:id", middleware.AdminOnly(), controllers.DeletePatient)

		// Appointment routes
		auth.POST("/appointments", middleware.AdminOrReceptionist(), controllers.CreateAppointment)
		auth.GET("/appointments", controllers.GetAppointments)
		auth.GET("/appointments/:id", controllers.GetAppointmentByID)
		auth.PUT("/appointments/:id", controllers.UpdateAppointment)
		auth.DELETE("/appointments/:id", middleware.AdminOrReceptionist(), controllers.DeleteAppointment)

		// Medical Records routes - Doctor and Admin can manage
		auth.POST("/medical-records", middleware.AdminOrDoctor(), controllers.CreateMedicalRecord)
		auth.GET("/medical-records", controllers.GetMedicalRecords)
		auth.GET("/medical-records/:id", controllers.GetMedicalRecordByID)
		auth.PUT("/medical-records/:id", middleware.AdminOrDoctor(), controllers.UpdateMedicalRecord)
		auth.DELETE("/medical-records/:id", middleware.AdminOnly(), controllers.DeleteMedicalRecord)

		// Prescription routes - Doctor and Admin can manage
		auth.POST("/prescriptions", middleware.AdminOrDoctor(), controllers.CreatePrescription)
		auth.GET("/prescriptions", controllers.GetPrescriptions)
		auth.GET("/prescriptions/:id", controllers.GetPrescriptionByID)
		auth.PUT("/prescriptions/:id", middleware.AdminOrDoctor(), controllers.UpdatePrescription)
		auth.DELETE("/prescriptions/:id", middleware.AdminOnly(), controllers.DeletePrescription)

		// Bill routes - Admin and Receptionist can manage
		auth.POST("/bills", middleware.AdminOrReceptionist(), controllers.CreateBill)
		auth.GET("/bills", controllers.GetBills)
		auth.GET("/bills/:id", controllers.GetBillByID)
		auth.PUT("/bills/:id", middleware.AdminOrReceptionist(), controllers.UpdateBill)
		auth.DELETE("/bills/:id", middleware.AdminOnly(), controllers.DeleteBill)

		// Room routes - Admin and Receptionist can manage
		auth.POST("/rooms", middleware.AdminOnly(), controllers.CreateRoom)
		auth.GET("/rooms", controllers.GetRooms)
		auth.GET("/rooms/:id", controllers.GetRoomByID)
		auth.PUT("/rooms/:id", middleware.AdminOrReceptionist(), controllers.UpdateRoom)
		auth.POST("/rooms/:id/assign", middleware.AdminOrReceptionist(), controllers.AssignRoomToPatient)
		auth.DELETE("/rooms/:id", middleware.AdminOnly(), controllers.DeleteRoom)
	}
}
