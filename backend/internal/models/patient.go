package models

import "time"

type Patient struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	Name           string          `json:"name"`
	Age            int             `json:"age"`
	Gender         string          `json:"gender"`
	Phone          string          `json:"phone"`
	Email          string          `json:"email,omitempty"`
	Address        string          `json:"address,omitempty"`
	RoomID         *uint           `json:"roomId,omitempty"`
	Room           *Room           `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	Appointments   []Appointment   `gorm:"foreignKey:PatientID" json:"appointments,omitempty"`
	MedicalRecords []MedicalRecord `gorm:"foreignKey:PatientID" json:"medicalRecords,omitempty"`
	Prescriptions  []Prescription  `gorm:"foreignKey:PatientID" json:"prescriptions,omitempty"`
	Bills          []Bill          `gorm:"foreignKey:PatientID" json:"bills,omitempty"`
}
