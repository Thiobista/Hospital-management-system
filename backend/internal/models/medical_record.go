package models

import "time"

type MedicalRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PatientID    uint      `json:"patientId"`
	Patient      Patient   `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID     uint      `json:"doctorId"`
	Doctor       Doctor    `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	Diagnosis    string    `json:"diagnosis"`
	Prescription string    `json:"prescription"`
	Notes        string    `json:"notes"`
	Date         time.Time `json:"date"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
