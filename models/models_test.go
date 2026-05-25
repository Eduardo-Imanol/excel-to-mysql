package models

import (
	"testing"
)

func TestCalModelFields(t *testing.T) {
	cal := Cal{
		Nombre:     "Juan Pérez",
		Math:       "90",
		Physical:   "85",
		Chemistry:  "88",
		Biologi:    "92",
		Histori:    "76",
		Geografi:   "81",
		Literature: "95",
		Spanish:    "87",
		English:    "93",
	}

	if cal.Nombre != "Juan Pérez" {
		t.Errorf("expected Nombre 'Juan Pérez', got %q", cal.Nombre)
	}
	if cal.Math != "90" {
		t.Errorf("expected Math '90', got %q", cal.Math)
	}
	if cal.Physical != "85" {
		t.Errorf("expected Physical '85', got %q", cal.Physical)
	}
	if cal.Chemistry != "88" {
		t.Errorf("expected Chemistry '88', got %q", cal.Chemistry)
	}
	if cal.Biologi != "92" {
		t.Errorf("expected Biologi '92', got %q", cal.Biologi)
	}
	if cal.Histori != "76" {
		t.Errorf("expected Histori '76', got %q", cal.Histori)
	}
	if cal.Geografi != "81" {
		t.Errorf("expected Geografi '81', got %q", cal.Geografi)
	}
	if cal.Literature != "95" {
		t.Errorf("expected Literature '95', got %q", cal.Literature)
	}
	if cal.Spanish != "87" {
		t.Errorf("expected Spanish '87', got %q", cal.Spanish)
	}
	if cal.English != "93" {
		t.Errorf("expected English '93', got %q", cal.English)
	}
}

func TestUploadModelRelations(t *testing.T) {
	upload := Upload{
		FileName: "calificaciones.xlsx",
		Sheets: []Sheet{
			{
				SheetName:  "Notas",
				HasHeaders: true,
				Headers:    `["nombre","math","physical"]`,
				Rows: []Row{
					{Data: `{"nombre":"Juan","math":"90","physical":"85"}`},
					{Data: `{"nombre":"María","math":"95","physical":"88"}`},
				},
			},
			{
				SheetName:  "Resumen",
				HasHeaders: false,
				Headers:    `[]`,
				Rows: []Row{
					{Data: `{"col1":"promedio","col2":"91.5"}`},
				},
			},
		},
	}

	if upload.FileName != "calificaciones.xlsx" {
		t.Errorf("expected filename 'calificaciones.xlsx', got %q", upload.FileName)
	}

	if len(upload.Sheets) != 2 {
		t.Errorf("expected 2 sheets, got %d", len(upload.Sheets))
	}

	sheet1 := upload.Sheets[0]
	if sheet1.SheetName != "Notas" {
		t.Errorf("expected SheetName 'Notas', got %q", sheet1.SheetName)
	}
	if !sheet1.HasHeaders {
		t.Errorf("expected HasHeaders to be true")
	}
	if len(sheet1.Rows) != 2 {
		t.Errorf("expected 2 rows in sheet1, got %d", len(sheet1.Rows))
	}

	sheet2 := upload.Sheets[1]
	if sheet2.SheetName != "Resumen" {
		t.Errorf("expected SheetName 'Resumen', got %q", sheet2.SheetName)
	}
	if sheet2.HasHeaders {
		t.Errorf("expected HasHeaders to be false")
	}
	if len(sheet2.Rows) != 1 {
		t.Errorf("expected 1 row in sheet2, got %d", len(sheet2.Rows))
	}
}

func TestSheetModelConstraints(t *testing.T) {
	sheet := Sheet{
		UploadID:   1,
		SheetName:  "Test",
		HasHeaders: true,
		Headers:    `["a","b"]`,
	}

	if sheet.UploadID != 1 {
		t.Errorf("expected UploadID 1, got %d", sheet.UploadID)
	}

	if sheet.Headers != `["a","b"]` {
		t.Errorf("expected Headers '[\"a\",\"b\"]', got %q", sheet.Headers)
	}
}
