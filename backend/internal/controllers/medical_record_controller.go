package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateMedicalRecord(c *gin.Context) {
	var record models.MedicalRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify patient exists
	var patient models.Patient
	if err := config.DB.First(&patient, record.PatientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found"})
		return
	}

	// Verify doctor exists
	var doctor models.Doctor
	if err := config.DB.First(&doctor, record.DoctorID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Doctor not found"})
		return
	}

	if err := config.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
		return
	}

	// Load relations
	config.DB.Preload("Patient").Preload("Doctor").First(&record, record.ID)

	c.JSON(http.StatusCreated, record)
}

func GetMedicalRecords(c *gin.Context) {
	var records []models.MedicalRecord
	query := config.DB.Preload("Patient").Preload("Doctor").Order("date DESC")

	// Filter by patient if provided
	if patientID := c.Query("patientId"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Filter by doctor if provided
	if doctorID := c.Query("doctorId"); doctorID != "" {
		query = query.Where("doctor_id = ?", doctorID)
	}

	if err := query.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medical records"})
		return
	}

	c.JSON(http.StatusOK, records)
}

func GetMedicalRecordByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medical record ID"})
		return
	}

	var record models.MedicalRecord
	if err := config.DB.Preload("Patient").Preload("Doctor").First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func UpdateMedicalRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medical record ID"})
		return
	}

	var record models.MedicalRecord
	if err := config.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical record"})
		return
	}

	config.DB.Preload("Patient").Preload("Doctor").First(&record, record.ID)
	c.JSON(http.StatusOK, record)
}

func DeleteMedicalRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medical record ID"})
		return
	}

	if err := config.DB.Delete(&models.MedicalRecord{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete medical record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medical record deleted successfully"})
}
