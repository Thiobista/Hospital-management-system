package controllers

import (
	"net/http"
	"strconv"
	"time"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateBill(c *gin.Context) {
	var bill models.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify patient exists
	var patient models.Patient
	if err := config.DB.First(&patient, bill.PatientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found"})
		return
	}

	// Set payment date if status is Paid
	if bill.Status == "Paid" && bill.PaymentDate == nil {
		now := time.Now()
		bill.PaymentDate = &now
	}

	if err := config.DB.Create(&bill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bill"})
		return
	}

	// Load relations
	config.DB.Preload("Patient").First(&bill, bill.ID)

	c.JSON(http.StatusCreated, bill)
}

func GetBills(c *gin.Context) {
	var bills []models.Bill
	query := config.DB.Preload("Patient").Order("created_at DESC")

	// Filter by patient if provided
	if patientID := c.Query("patientId"); patientID != "" {
		query = query.Where("patient_id = ?", patientID)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&bills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bills"})
		return
	}

	c.JSON(http.StatusOK, bills)
}

func GetBillByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	var bill models.Bill
	if err := config.DB.Preload("Patient").First(&bill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	c.JSON(http.StatusOK, bill)
}

func UpdateBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	var bill models.Bill
	if err := config.DB.First(&bill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bill not found"})
		return
	}

	if err := c.ShouldBindJSON(&bill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set payment date if status changed to Paid
	if bill.Status == "Paid" && bill.PaymentDate == nil {
		now := time.Now()
		bill.PaymentDate = &now
	}

	if err := config.DB.Save(&bill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bill"})
		return
	}

	config.DB.Preload("Patient").First(&bill, bill.ID)
	c.JSON(http.StatusOK, bill)
}

func DeleteBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bill ID"})
		return
	}

	if err := config.DB.Delete(&models.Bill{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bill deleted successfully"})
}
