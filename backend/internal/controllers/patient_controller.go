package controllers

import (
	"net/http"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func CreatePatient(c *gin.Context) {
	var p models.Patient
	c.BindJSON(&p)
	config.DB.Create(&p)
	c.JSON(http.StatusOK, p) // use http.StatusOK instead of 200
}

func GetPatients(c *gin.Context) {
	var patients []models.Patient
	config.DB.Find(&patients)
	c.JSON(http.StatusOK, patients) // use http.StatusOK
}
