package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateAppointment(c *gin.Context) {
	var a models.Appointment
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify patient exists
	var patient models.Patient
	if err := config.DB.First(&patient, a.PatientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found"})
		return
	}

	if err := config.DB.Create(&a).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusCreated, a)
}

func GetAppointments(c *gin.Context) {
	var appointments []models.Appointment
	query := config.DB.Preload("Patient").Preload("Doctor").Order("date DESC")

	// Filter by patient if provided
	if patientID := c.Query("patientId"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Filter by doctor if provided
	if doctorID := c.Query("doctorId"); doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// If user is a doctor, only show their appointments
	userRole, exists := c.Get("userRole")
	if exists && userRole == "doctor" {
		userID, _ := c.Get("userID")
		// Note: This assumes doctor user ID matches doctor ID
		// You may need to adjust based on your user-doctor relationship
		query = query.Where("doctor_id = ?", userID)
	}

	if err := query.Find(&appointments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

func GetAppointmentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var appointment models.Appointment
	if err := config.DB.Preload("Patient").Preload("Doctor").First(&appointment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, appointment)
}

func UpdateAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var appointment models.Appointment
	if err := config.DB.First(&appointment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	if err := c.ShouldBindJSON(&appointment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	config.DB.Preload("Patient").Preload("Doctor").First(&appointment, appointment.ID)
	c.JSON(http.StatusOK, appointment)
}

func DeleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	if err := config.DB.Delete(&models.Appointment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Appointment deleted successfully"})
}
