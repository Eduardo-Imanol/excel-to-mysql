package models

import "gorm.io/gorm"

// Upload representa un archivo Excel subido al sistema.
// Contiene el nombre del archivo y sus hojas asociadas.
// Relación: Upload 1 ──> N Sheet (cascade delete)
type Upload struct {
	gorm.Model                          // ID, CreatedAt, UpdatedAt, DeletedAt
	FileName string  `json:"file"`      // Nombre original del archivo Excel
	Sheets   []Sheet `gorm:"constraint:OnDelete:CASCADE"` // Hojas del archivo
}

// Sheet representa una hoja dentro de un archivo Excel.
// Almacena el nombre, si tiene encabezados, los headers en JSON
// y las filas de datos asociadas.
// Relación: Sheet 1 ──> N Row (cascade delete)
type Sheet struct {
	gorm.Model
	UploadID   uint   `json:"-"`                                    // FK hacia Upload
	SheetName  string `json:"sheet"`                                // Nombre de la hoja
	HasHeaders bool   `json:"hasHeaders"`                           // Si la hoja tiene fila de encabezados
	Headers    string `json:"headers" gorm:"type:text"`             // Encabezados en JSON (array de strings)
	Rows       []Row  `gorm:"constraint:OnDelete:CASCADE"`         // Filas de datos
}

// Row representa una fila de datos dentro de una hoja Excel.
// Almacena los datos de la fila como un objeto JSON (clave-valor).
// Relación: Row N ──> 1 Sheet
type Row struct {
	gorm.Model
	SheetID uint   `json:"-"`                       // FK hacia Sheet
	Data    string `json:"row" gorm:"type:longtext"` // Datos de la fila en JSON
}
