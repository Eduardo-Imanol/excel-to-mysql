package models

import (
	"testing"
)

func TestRecordModelFields(t *testing.T) {
	record := Record{
		Nombre: "Juan Pérez",
		Data:   `{"math":"90","city":"Madrid"}`,
	}

	if record.Nombre != "Juan Pérez" {
		t.Errorf("expected Nombre 'Juan Pérez', got %q", record.Nombre)
	}
	if record.Data != `{"math":"90","city":"Madrid"}` {
		t.Errorf("expected Data '{\"math\":\"90\",\"city\":\"Madrid\"}', got %q", record.Data)
	}
}

func TestRecordModelEmptyData(t *testing.T) {
	record := Record{
		Nombre: "Sin datos",
		Data:   "{}",
	}

	if record.Nombre != "Sin datos" {
		t.Errorf("expected Nombre 'Sin datos', got %q", record.Nombre)
	}
	if record.Data != "{}" {
		t.Errorf("expected Data '{}', got %q", record.Data)
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
