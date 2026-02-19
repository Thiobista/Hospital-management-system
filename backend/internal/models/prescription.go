package models

import "time"

type Prescription struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PatientID    uint      `json:"patientId"`
	Patient      Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID     uint      `json:"doctorId"`
	Doctor       Doctor    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	MedicineName string    `json:"medicineName"`
	Dosage       string    `json:"dosage"`
	Instructions string    `json:"instructions"`
	Date         time.Time `json:"date"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
