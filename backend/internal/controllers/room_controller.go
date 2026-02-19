package controllers

import (
	"net/http"
	"strconv"

	"clinic-backend/internal/config"
	"clinic-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if room number already exists
	var existingRoom models.Room
	if err := config.DB.Where("room_number = ?", room.RoomNumber).First(&existingRoom).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room number already exists"})
		return
	}

	if err := config.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	// Load relations
	if room.PatientID != nil {
		config.DB.Preload("Patient").First(&room, room.ID)
	}

	c.JSON(http.StatusCreated, room)
}

func GetRooms(c *gin.Context) {
	var rooms []models.Room
	query := config.DB.Preload("Patient").Order("room_number ASC")

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by type if provided
	if roomType := c.Query("type"); roomType != "" {
		query = query.Where("type = ?", roomType)
	}

	if err := query.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func GetRoomByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := config.DB.Preload("Patient").First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

func UpdateRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := config.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		return
	}

	config.DB.Preload("Patient").First(&room, room.ID)
	c.JSON(http.StatusOK, room)
}

func DeleteRoom(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	if err := config.DB.Delete(&models.Room{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

func AssignRoomToPatient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var room models.Room
	if err := config.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	var body struct {
		PatientID *uint `json:"patientId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If assigning to a patient, verify patient exists
	if body.PatientID != nil {
		var patient models.Patient
		if err := config.DB.First(&patient, *body.PatientID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found"})
			return
		}

		// Check if room is available
		if room.Status == "Occupied" && room.PatientID != nil && *room.PatientID != *body.PatientID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room is already occupied"})
			return
		}

		room.PatientID = body.PatientID
		room.Status = "Occupied"
	} else {
		// Unassigning room
		room.PatientID = nil
		if room.Status == "Occupied" {
			room.Status = "Available"
		}
	}

	if err := config.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign room"})
		return
	}

	config.DB.Preload("Patient").First(&room, room.ID)
	c.JSON(http.StatusOK, room)
}
