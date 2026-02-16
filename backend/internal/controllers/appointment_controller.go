package controllers

import (
	"net/http"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateAppointment(c *gin.Context) {
	var a models.Appointment
	c.BindJSON(&a)
	config.DB.Create(&a)
	c.JSON(http.StatusOK, a)
}

func GetAppointments(c *gin.Context) {
	var appointments []models.Appointment
	config.DB.Find(&appointments)
	c.JSON(http.StatusOK, appointments)
}
