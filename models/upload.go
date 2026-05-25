package models

import "gorm.io/gorm"

type Upload struct {
	gorm.Model
	FileName string  `json:"file"`
	Sheets   []Sheet `gorm:"constraint:OnDelete:CASCADE"`
}

type Sheet struct {
	gorm.Model
	UploadID   uint   `json:"-"`
	SheetName  string `json:"sheet"`
	HasHeaders bool   `json:"hasHeaders"`
	Headers    string `json:"headers" gorm:"type:text"`
	Rows       []Row  `gorm:"constraint:OnDelete:CASCADE"`
}

type Row struct {
	gorm.Model
	SheetID uint   `json:"-"`
	Data    string `json:"row" gorm:"type:longtext"`
}
