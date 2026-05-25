package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Eduardo-Imanol/excel-to-mysql/db"
	"github.com/Eduardo-Imanol/excel-to-mysql/models"
)

const excelMaxSize = 10 << 20

type uploadRequest struct {
	File   string         `json:"file"`
	Sheets []sheetRequest `json:"sheets"`
}

type sheetRequest struct {
	Sheet      string              `json:"sheet"`
	Headers    []string            `json:"headers"`
	Rows       []map[string]string `json:"rows"`
	HasHeaders bool                `json:"hasHeaders"`
}

func UploadExcelHandler(w http.ResponseWriter, r *http.Request) {
	if !isJSONContent(r) {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, excelMaxSize)

	var req uploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.File == "" {
		http.Error(w, "El campo 'file' es requerido", http.StatusBadRequest)
		return
	}

	if len(req.Sheets) == 0 {
		http.Error(w, "Debe incluir al menos una hoja", http.StatusBadRequest)
		return
	}

	upload := models.Upload{FileName: req.File}

	totalRows := 0
	for _, s := range req.Sheets {
		if strings.TrimSpace(s.Sheet) == "" {
			http.Error(w, "Todas las hojas deben tener un nombre", http.StatusBadRequest)
			return
		}

		headersJSON, err := json.Marshal(s.Headers)
		if err != nil {
			http.Error(w, "Error al procesar headers", http.StatusInternalServerError)
			return
		}

		sheet := models.Sheet{
			SheetName:  s.Sheet,
			HasHeaders: s.HasHeaders,
			Headers:    string(headersJSON),
		}

		for _, row := range s.Rows {
			rowJSON, err := json.Marshal(row)
			if err != nil {
				http.Error(w, "Error al procesar una fila", http.StatusInternalServerError)
				return
			}
			sheet.Rows = append(sheet.Rows, models.Row{Data: string(rowJSON)})
		}

		totalRows += len(s.Rows)
		upload.Sheets = append(upload.Sheets, sheet)
	}

	if err := db.DB.Create(&upload).Error; err != nil {
		http.Error(w, "Error al guardar los datos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Datos guardados correctamente",
		"upload_id":  upload.ID,
		"sheets":     len(upload.Sheets),
		"total_rows": totalRows,
	})
}

func GetUploadsHandler(w http.ResponseWriter, r *http.Request) {
	var uploads []models.Upload

	if err := db.DB.Preload("Sheets.Rows").Order("created_at desc").Find(&uploads).Error; err != nil {
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}

	if len(uploads) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	type rowResp struct {
		ID   uint                 `json:"id"`
		Data map[string]string    `json:"data"`
	}

	type sheetResp struct {
		ID         uint       `json:"id"`
		SheetName  string     `json:"sheet"`
		HasHeaders bool       `json:"hasHeaders"`
		Headers    []string   `json:"headers"`
		Rows       []rowResp  `json:"rows"`
	}

	type uploadResp struct {
		ID        uint         `json:"id"`
		FileName  string       `json:"file"`
		CreatedAt string       `json:"created_at"`
		Sheets    []sheetResp  `json:"sheets"`
	}

	var result []uploadResp
	for _, u := range uploads {
		var sheets []sheetResp
		for _, s := range u.Sheets {
			var headers []string
			json.Unmarshal([]byte(s.Headers), &headers)

			var rows []rowResp
			for _, r := range s.Rows {
				var data map[string]string
				json.Unmarshal([]byte(r.Data), &data)
				rows = append(rows, rowResp{ID: r.ID, Data: data})
			}

			sheets = append(sheets, sheetResp{
				ID:         s.ID,
				SheetName:  s.SheetName,
				HasHeaders: s.HasHeaders,
				Headers:    headers,
				Rows:       rows,
			})
		}

		result = append(result, uploadResp{
			ID:        u.ID,
			FileName:  u.FileName,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
			Sheets:    sheets,
		})
	}

	json.NewEncoder(w).Encode(result)
}
