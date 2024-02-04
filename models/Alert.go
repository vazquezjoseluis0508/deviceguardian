package models

import (
	"time"
	"gorm.io/gorm"
)

type Alert struct {
	gorm.Model
	DeviceID  uint `gorm:"not null"` // Clave foránea de Device
	UserID    uint `gorm:"not null"` // Clave foránea de User que recibe la alerta
	AlertType string `gorm:"not null"`
	Description string `gorm:"not null"`
	AlertTime  time.Time `gorm:"not null"`
}
