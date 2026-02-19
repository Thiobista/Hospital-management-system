package controllers

import (
	"net/http"
	"strconv"
	"time"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreatePrescription(c *gin.Context) {
	var prescription models.Prescription
	if err := c.ShouldBindJSON(&prescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify patient exists
	var patient models.Patient
	if err := config.DB.First(&patient, prescription.PatientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found"})
		return
	}

	// Verify doctor exists
	var doctor models.Doctor
	if err := config.DB.First(&doctor, prescription.DoctorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Doctor not found"})
		return
	}

	// Set date if not provided
	if prescription.Date.IsZero() {
		prescription.Date = time.Now()
	}

	if err := config.DB.Create(&prescription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create prescription"})
		return
	}

	// Load relations
	config.DB.Preload("Patient").Preload("Doctor").First(&prescription, prescription.ID)

	c.JSON(http.StatusCreated, prescription)
}

func GetPrescriptions(c *gin.Context) {
	var prescriptions []models.Prescription
	query := config.DB.Preload("Patient").Preload("Doctor").Order("date DESC")

	// Filter by patient if provided
	if patientID := c.Query("patientId"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Filter by doctor if provided
	if doctorID := c.Query("doctorId"); doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if err := query.Find(&prescriptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescriptions"})
		return
	}

	c.JSON(http.StatusOK, prescriptions)
}

func GetPrescriptionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	var prescription models.Prescription
	if err := config.DB.Preload("Patient").Preload("Doctor").First(&prescription, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found"})
		return
	}

	c.JSON(http.StatusOK, prescription)
}

func UpdatePrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	var prescription models.Prescription
	if err := config.DB.First(&prescription, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prescription not found"})
		return
	}

	if err := c.ShouldBindJSON(&prescription); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&prescription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prescription"})
		return
	}

	config.DB.Preload("Patient").Preload("Doctor").First(&prescription, prescription.ID)
	c.JSON(http.StatusOK, prescription)
}

func DeletePrescription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
		return
	}

	if err := config.DB.Delete(&models.Prescription{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete prescription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prescription deleted successfully"})
}
