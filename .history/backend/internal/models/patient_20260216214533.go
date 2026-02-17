package models

import "time"

type Patient struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Name       string       `json:"name"`
	Age        int          `json:"age"`
	Gender     string       `json:"gender"`
	Phone      string       `json:"phone"`
	Email      string       `json:"email,omitempty"`
	Address    string       `json:"address,omitempty"`
	CreatedAt  time.Time    `json:"createdAt"`
	Appointments []Appointment `gorm:"foreignKey:PatientID" json:"appointments,omitempty"`
}
