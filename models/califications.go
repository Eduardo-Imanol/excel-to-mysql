package models

import (
	// "time"

	"gorm.io/gorm"
)

type Cal struct {
	gorm.Model

	Nombre     string `json:"nombre"`
	Math       string `json:"math"`
	Physical   string `json:"physical"`
	Chemistry  string `json:"chemistry"`
	Biologi    string `json:"biologi"`
	Histori    string `json:"histori"`
	Geografi   string `json:"geografi"`
	Literature string `json:"literature"`
	Spanish    string `json:"spanish"`
	English    string `json:"english"`
}
