package models

import "time"

type Doctor struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `json:"name"`
	Email          string    `gorm:"unique" json:"email"`
	Phone          string    `json:"phone"`
	Specialization string    `json:"specialization"`
	Availability   string    `json:"availability"` // e.g., "Monday-Friday, 9AM-5PM"
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Relations
	Appointments   []Appointment   `gorm:"foreignKey:DoctorID" json:"appointments,omitempty"`
	MedicalRecords []MedicalRecord `gorm:"foreignKey:DoctorID" json:"medicalRecords,omitempty"`
	Prescriptions  []Prescription  `gorm:"foreignKey:DoctorID" json:"prescriptions,omitempty"`
}
