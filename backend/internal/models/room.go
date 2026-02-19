package models

import "time"

type Room struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RoomNumber string    `gorm:"unique" json:"roomNumber"`
	Type       string    `json:"type"`   // Single, Double, ICU, Emergency, etc.
	Status     string    `json:"status"` // Available, Occupied, Maintenance
	PatientID  *uint     `json:"patientId,omitempty"`
	Patient    *Patient  `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
