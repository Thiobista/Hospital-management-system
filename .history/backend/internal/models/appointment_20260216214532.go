package models

import "time"

type Appointment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PatientID uint      `json:"patientId"`
	Patient   Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Doctor    string    `json:"doctor"`
	Date      time.Time `json:"date"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}
