package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName             string `gorm:"not null"`
	LastName              string `gorm:"not null"`
	Email                 string `gorm:"not null; unique"`
	Password              string `gorm:"not null"`
	Role                  string `gorm:"not null"` // Nuevo: cliente, técnico, administrador
	Status                string `gorm:"not null"` // Nuevo: activo, inactivo
	Devices               []Device `gorm:"foreignKey:UserID"` // Relación con dispositivos
}

