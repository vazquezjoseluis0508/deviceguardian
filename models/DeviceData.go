package models

import (
	"time"
	"gorm.io/gorm"
)


type DeviceData struct {
	gorm.Model
	DeviceID uint `gorm:"not null"` // Clave foránea de Device
	Timestamp time.Time `gorm:"not null"`
	Temperature float64 `gorm:"not null"`
	BatteryStatus string `gorm:"not null"`
	// Añade aquí otros sensores o datos relevantes
}
