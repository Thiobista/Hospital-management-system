package models

import "time"

type Bill struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	PatientID   uint       `json:"patientId"`
	Patient     Patient    `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	Amount      float64    `json:"amount"`
	Status      string     `json:"status"` // Paid, Unpaid
	PaymentDate *time.Time `json:"paymentDate,omitempty"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
