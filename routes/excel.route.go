package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Eduardo-Imanol/excel-to-mysql/db"
	"github.com/Eduardo-Imanol/excel-to-mysql/models"
)

// excelMaxSize define el tamaño máximo para el body de subida de Excel (10 MB).
const excelMaxSize = 10 << 20

// uploadRequest estructura para recibir la estructura de un archivo Excel.
type uploadRequest struct {
	File   string         `json:"file"`   // Nombre del archivo Excel
	Sheets []sheetRequest `json:"sheets"` // Hojas del archivo
}

// sheetRequest estructura para recibir los datos de una hoja Excel.
type sheetRequest struct {
	Sheet      string              `json:"sheet"`      // Nombre de la hoja
	Headers    []string            `json:"headers"`    // Encabezados de columnas
	Rows       []map[string]string `json:"rows"`       // Filas de datos (clave-valor)
	HasHeaders bool                `json:"hasHeaders"` // Si tiene fila de encabezados
}

// UploadExcelHandler maneja POST /excel/upload.
// Recibe la estructura de un archivo Excel (nombre, hojas, headers, filas)
// y la almacena en la base de datos con relaciones Upload -> Sheet -> Row.
func UploadExcelHandler(w http.ResponseWriter, r *http.Request) {
	// Validar Content-Type
	if !isJSONContent(r) {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Limitar tamaño del body
	r.Body = http.MaxBytesReader(w, r.Body, excelMaxSize)

	// Decodificar JSON
	var req uploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar campo requerido: nombre del archivo
	if req.File == "" {
		http.Error(w, "El campo 'file' es requerido", http.StatusBadRequest)
		return
	}

	// Validar que haya al menos una hoja
	if len(req.Sheets) == 0 {
		http.Error(w, "Debe incluir al menos una hoja", http.StatusBadRequest)
		return
	}

	// Construir el objeto Upload con sus hojas y filas
	upload := models.Upload{FileName: req.File}

	totalRows := 0
	for _, s := range req.Sheets {
		// Validar que cada hoja tenga nombre
		if strings.TrimSpace(s.Sheet) == "" {
			http.Error(w, "Todas las hojas deben tener un nombre", http.StatusBadRequest)
			return
		}

		// Convertir headers a JSON string
		headersJSON, err := json.Marshal(s.Headers)
		if err != nil {
			http.Error(w, "Error al procesar headers", http.StatusInternalServerError)
			return
		}

		// Crear modelo Sheet
		sheet := models.Sheet{
			SheetName:  s.Sheet,
			HasHeaders: s.HasHeaders,
			Headers:    string(headersJSON),
		}

		// Agregar filas a la hoja
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

	// Guardar en la base de datos (crea Upload, Sheets y Rows en cascada)
	if err := db.DB.Create(&upload).Error; err != nil {
		http.Error(w, "Error al guardar los datos", http.StatusInternalServerError)
		return
	}

	// Retornar resumen de la operación
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Datos guardados correctamente",
		"upload_id":  upload.ID,
		"sheets":     len(upload.Sheets),
		"total_rows": totalRows,
	})
}

// GetUploadsHandler maneja GET /excel/uploads.
// Retorna todas las subidas de Excel con sus hojas y filas.
// Incluye preloading de relaciones (Sheets.Rows) para evitar N+1 queries.
func GetUploadsHandler(w http.ResponseWriter, r *http.Request) {
	var uploads []models.Upload

	// Cargar uploads con hojas y filas en una sola consulta optimizada
	if err := db.DB.Preload("Sheets.Rows").Order("created_at desc").Find(&uploads).Error; err != nil {
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}

	// Si no hay subidas, retornar array vacío
	if len(uploads) == 0 {
		json.NewEncoder(w).Encode([]interface{}{})
		return
	}

	// ── Estructuras de respuesta formateada ──

	// Fila individual con su data como map
	type rowResp struct {
		ID   uint              `json:"id"`
		Data map[string]string `json:"data"`
	}

	// Hoja con headers y filas
	type sheetResp struct {
		ID         uint      `json:"id"`
		SheetName  string    `json:"sheet"`
		HasHeaders bool      `json:"hasHeaders"`
		Headers    []string  `json:"headers"`
		Rows       []rowResp `json:"rows"`
	}

	// Upload con todas sus hojas
	type uploadResp struct {
		ID        uint        `json:"id"`
		FileName  string      `json:"file"`
		CreatedAt string      `json:"created_at"`
		Sheets    []sheetResp `json:"sheets"`
	}

	// Convertir modelos GORM a estructuras de respuesta limpias
	var result []uploadResp
	for _, u := range uploads {
		var sheets []sheetResp
		for _, s := range u.Sheets {
			// Deserializar headers de JSON string a array
			var headers []string
			json.Unmarshal([]byte(s.Headers), &headers)

			// Deserializar cada fila de JSON string a map
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
