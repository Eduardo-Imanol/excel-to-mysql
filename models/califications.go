package models

import (
	"gorm.io/gorm"
)

// Record representa un registro genérico (calificación, nombre, etc.).
// Almacena un nombre y datos adicionales en formato JSON flexible.
// Se usa la tabla "records" en MySQL.
type Record struct {
	gorm.Model                          // ID, CreatedAt, UpdatedAt, DeletedAt (soft delete)
	Nombre string          `json:"nombre"` // Nombre del registro
	Data   string          `json:"data" gorm:"type:json"` // Datos adicionales en JSON
}
