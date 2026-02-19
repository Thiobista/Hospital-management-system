package controllers

import (
	"net/http"
	"time"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAdminDashboard(c *gin.Context) {
	var stats struct {
		TotalPatients     int64 `json:"totalPatients"`
		TotalDoctors      int64 `json:"totalDoctors"`
		TotalAppointments int64 `json:"totalAppointments"`
		TodayAppointments int64 `json:"todayAppointments"`
		PendingBills      int64 `json:"pendingBills"`
		AvailableRooms    int64 `json:"availableRooms"`
	}

	// Get counts
	config.DB.Model(&models.Patient{}).Count(&stats.TotalPatients)
	config.DB.Model(&models.Doctor{}).Count(&stats.TotalDoctors)
	config.DB.Model(&models.Appointment{}).Count(&stats.TotalAppointments)

	// Today's appointments
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	config.DB.Model(&models.Appointment{}).
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Count(&stats.TodayAppointments)

	// Pending bills
	config.DB.Model(&models.Bill{}).Where("status = ?", "Unpaid").Count(&stats.PendingBills)

	// Available rooms
	config.DB.Model(&models.Room{}).Where("status = ?", "Available").Count(&stats.AvailableRooms)

	c.JSON(http.StatusOK, stats)
}

func GetDoctorDashboard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	doctorID := userID.(uint)

	var stats struct {
		TodayAppointments    []models.Appointment `json:"todayAppointments"`
		UpcomingAppointments []models.Appointment `json:"upcomingAppointments"`
		TotalPatients        int64                `json:"totalPatients"`
	}

	// Get today's appointments for this doctor
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	config.DB.Preload("Patient").
		Where("doctor_id = ? AND date >= ? AND date < ?", doctorID, startOfDay, endOfDay).
		Order("date ASC").
		Find(&stats.TodayAppointments)

	// Get upcoming appointments (next 7 days)
	weekFromNow := today.Add(7 * 24 * time.Hour)
	config.DB.Preload("Patient").
		Where("doctor_id = ? AND date > ? AND date <= ?", doctorID, endOfDay, weekFromNow).
		Order("date ASC").
		Limit(10).
		Find(&stats.UpcomingAppointments)

	// Count unique patients
	config.DB.Model(&models.Appointment{}).
		Where("doctor_id = ?", doctorID).
		Distinct("patient_id").
		Count(&stats.TotalPatients)

	c.JSON(http.StatusOK, stats)
}

func GetReceptionistDashboard(c *gin.Context) {
	var stats struct {
		TodayAppointments   []models.Appointment `json:"todayAppointments"`
		PendingAppointments int64                `json:"pendingAppointments"`
		TotalPatients       int64                `json:"totalPatients"`
		AvailableRooms      int64                `json:"availableRooms"`
	}

	// Get today's appointments
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	config.DB.Preload("Patient").Preload("Doctor").
		Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Order("date ASC").
		Find(&stats.TodayAppointments)

	// Count scheduled appointments
	config.DB.Model(&models.Appointment{}).
		Where("status = ?", "Scheduled").
		Count(&stats.PendingAppointments)

	// Total patients
	config.DB.Model(&models.Patient{}).Count(&stats.TotalPatients)

	// Available rooms
	config.DB.Model(&models.Room{}).Where("status = ?", "Available").Count(&stats.AvailableRooms)

	c.JSON(http.StatusOK, stats)
}
