package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Eduardo-Imanol/excel-to-mysql/models"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	HomeHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	expected := "Hello Go, This is MY first API rest 40"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestIsJSONContent(t *testing.T) {
	tests := []struct {
		name     string
		ct       string
		expected bool
	}{
		{"application/json", "application/json", true},
		{"with charset", "application/json; charset=utf-8", true},
		{"text plain", "text/plain", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Content-Type", tt.ct)
			if got := isJSONContent(req); got != tt.expected {
				t.Errorf("isJSONContent() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestValidateCal(t *testing.T) {
	tests := []struct {
		name    string
		nombre  string
		wantErr bool
		errMsg  string
	}{
		{"valid name", "Juan", false, ""},
		{"empty name", "", true, "el campo 'nombre' es requerido"},
		{"spaces only", "   ", true, "el campo 'nombre' es requerido"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cal := &models.Cal{Nombre: tt.nombre}
			err := validateCal(cal)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestPostNameHandler_InvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/names", nil)
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	PostNameHandler(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415, got %d", rec.Code)
	}
}

func TestPostNameHandler_InvalidJSON(t *testing.T) {
	body := bytes.NewReader([]byte(`invalid json`))
	req := httptest.NewRequest(http.MethodPost, "/names", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	PostNameHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostNameHandler_MissingNombre(t *testing.T) {
	cal := models.Cal{
		Math: "90",
	}
	b, _ := json.Marshal(cal)
	body := bytes.NewReader(b)
	req := httptest.NewRequest(http.MethodPost, "/names", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	PostNameHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostNamesHandler_InvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/names/all", nil)
	req.Header.Set("Content-Type", "text/plain")
	rec := httptest.NewRecorder()

	PostNamesHandler(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415, got %d", rec.Code)
	}
}

func TestPostNamesHandler_InvalidJSON(t *testing.T) {
	body := bytes.NewReader([]byte(`invalid json`))
	req := httptest.NewRequest(http.MethodPost, "/names/all", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	PostNamesHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestPostNamesHandler_EmptyArray(t *testing.T) {
	body := bytes.NewReader([]byte(`[]`))
	req := httptest.NewRequest(http.MethodPost, "/names/all", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	PostNamesHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	expected := "El array no puede estar vacío\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestUploadExcelHandler_InvalidContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/excel/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	rec := httptest.NewRecorder()

	UploadExcelHandler(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415, got %d", rec.Code)
	}
}

func TestUploadExcelHandler_InvalidJSON(t *testing.T) {
	body := bytes.NewReader([]byte(`invalid json`))
	req := httptest.NewRequest(http.MethodPost, "/excel/upload", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	UploadExcelHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestUploadExcelHandler_MissingFile(t *testing.T) {
	reqBody := uploadRequest{
		File:   "",
		Sheets: []sheetRequest{{Sheet: "Sheet1", Rows: []map[string]string{{"name": "test"}}}},
	}
	b, _ := json.Marshal(reqBody)
	body := bytes.NewReader(b)
	req := httptest.NewRequest(http.MethodPost, "/excel/upload", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	UploadExcelHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	expected := "El campo 'file' es requerido\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestUploadExcelHandler_EmptySheets(t *testing.T) {
	reqBody := uploadRequest{
		File:   "test.xlsx",
		Sheets: []sheetRequest{},
	}
	b, _ := json.Marshal(reqBody)
	body := bytes.NewReader(b)
	req := httptest.NewRequest(http.MethodPost, "/excel/upload", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	UploadExcelHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	expected := "Debe incluir al menos una hoja\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestUploadExcelHandler_SheetWithoutName(t *testing.T) {
	reqBody := uploadRequest{
		File:   "test.xlsx",
		Sheets: []sheetRequest{{Sheet: "  ", Rows: []map[string]string{{"name": "test"}}}},
	}
	b, _ := json.Marshal(reqBody)
	body := bytes.NewReader(b)
	req := httptest.NewRequest(http.MethodPost, "/excel/upload", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	UploadExcelHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	expected := "Todas las hojas deben tener un nombre\n"
	if rec.Body.String() != expected {
		t.Errorf("expected body %q, got %q", expected, rec.Body.String())
	}
}
