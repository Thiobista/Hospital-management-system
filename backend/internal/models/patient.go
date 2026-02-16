package models

import "time"

type Patient struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string
	Age       int
	Gender    string
	Phone     string
	CreatedAt time.Time
}
