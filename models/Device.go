package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Type        string `gorm:"not null"` // Nuevo: UPS, sensor, etc.
	Description string `gorm:"not null"`
	UserID      uint   `gorm:"not null"` // Clave foránea de User
	Location    string `gorm:"not null"`
	InstallationDate string `gorm:"not null"`
	Status      string `gorm:"not null"` // Nuevo: activo, inactivo
	DeviceData  []DeviceData `gorm:"foreignKey:DeviceID"` // Datos del dispositivo
	Alerts      []Alert `gorm:"foreignKey:DeviceID"` // Alertas relacionadas
}
