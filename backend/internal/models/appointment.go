package models

import "time"

type Appointment struct {
	ID        uint      `gorm:"primaryKey"`
	PatientID uint
	Doctor    string
	Date      time.Time
	Notes     string
}
