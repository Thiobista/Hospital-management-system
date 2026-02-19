package models

import "time"

type Appointment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PatientID uint      `json:"patientId"`
	Patient   Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID  uint      `json:"doctorId"`
	Doctor    Doctor    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	Date      time.Time `json:"date"`
	Time      string    `json:"time"`   // e.g., "10:00 AM"
	Status    string    `json:"status"` // Scheduled, Completed, Cancelled
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
