package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Eduardo-Imanol/excel-to-mysql/db"
	"github.com/Eduardo-Imanol/excel-to-mysql/models"
	"github.com/gorilla/mux"
)

// maxBodySize define el tamaño máximo permitido para el body de una petición (1 MB).
const maxBodySize = 1 << 20

// recordRequest estructura para recibir datos de un registro vía JSON.
type recordRequest struct {
	Nombre string          `json:"nombre"` // Nombre del registro
	Data   json.RawMessage `json:"data"`   // Datos adicionales (JSON flexible)
}

// GetNamesHandler maneja GET /names.
// Retorna todos los registros almacenados en la base de datos.
func GetNamesHandler(w http.ResponseWriter, r *http.Request) {
	var records []models.Record
	db.DB.Find(&records)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&records)
}

// GetNameHandler maneja GET /names/{id}.
// Retorna un registro específico por su ID.
// Retorna 404 si no se encuentra el registro.
func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	var record models.Record
	id := mux.Vars(r)

	db.DB.Where("id = ?", id["id"]).First(&record)

	if record.ID == 0 {
		http.Error(w, "id no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&record)
}

// PostNameHandler maneja POST /names.
// Crea un nuevo registro con los datos enviados en el body JSON.
// Valida que el Content-Type sea JSON y que el campo 'nombre' no esté vacío.
func PostNameHandler(w http.ResponseWriter, r *http.Request) {
	// Validar Content-Type
	if !isJSONContent(r) {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decodificar body con límite de tamaño
	var req recordRequest
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar campo requerido
	if strings.TrimSpace(req.Nombre) == "" {
		http.Error(w, "el campo 'nombre' es requerido", http.StatusBadRequest)
		return
	}

	// Crear registro en la BD
	record := models.Record{
		Nombre: req.Nombre,
		Data:   string(req.Data),
	}

	if err := db.DB.Create(&record).Error; err != nil {
		http.Error(w, "Error al crear el registro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&record)
}

// PostNamesHandler maneja POST /names/all.
// Crea múltiples registros a partir de un array JSON.
// Valida que el array no esté vacío y que cada elemento tenga 'nombre'.
func PostNamesHandler(w http.ResponseWriter, r *http.Request) {
	// Validar Content-Type
	if !isJSONContent(r) {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	// Decodificar array de registros
	var reqs []recordRequest
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar que el array no esté vacío
	if len(reqs) == 0 {
		http.Error(w, "El array no puede estar vacío", http.StatusBadRequest)
		return
	}

	// Crear cada registro individualmente
	var records []models.Record
	for i, req := range reqs {
		if strings.TrimSpace(req.Nombre) == "" {
			http.Error(w, fmt.Sprintf("Error en elemento %d: el campo 'nombre' es requerido", i), http.StatusBadRequest)
			return
		}

		record := models.Record{
			Nombre: req.Nombre,
			Data:   string(req.Data),
		}

		if err := db.DB.Create(&record).Error; err != nil {
			http.Error(w, "Error al crear los registros", http.StatusInternalServerError)
			return
		}
		records = append(records, record)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&records)
}

// DeleteNamesHandler maneja DELETE /names/{id}.
// Elimina permanentemente un registro por su ID (hard delete con Unscoped).
func DeleteNamesHandler(w http.ResponseWriter, r *http.Request) {
	var record models.Record
	id := mux.Vars(r)

	db.DB.Where("id = ?", id["id"]).First(&record)

	if record.ID == 0 {
		http.Error(w, "id no encontrado", http.StatusNotFound)
		return
	}

	db.DB.Unscoped().Delete(&record)
	w.WriteHeader(http.StatusOK)
}

// DeleteAllRecordsHandler maneja DELETE /names.
// Elimina permanentemente todos los registros de la tabla (hard delete).
func DeleteAllRecordsHandler(w http.ResponseWriter, r *http.Request) {
	var records []models.Record

	db.DB.Unscoped().Find(&records)

	if len(records) == 0 {
		http.Error(w, "No se encontraron registros para eliminar", http.StatusNotFound)
		return
	}

	result := db.DB.Unscoped().Delete(&records)
	if result.Error != nil {
		http.Error(w, "Error al borrar los registros", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registros eliminados exitosamente"))
}

// isJSONContent verifica si el Content-Type de la petición es application/json.
func isJSONContent(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	return strings.HasPrefix(ct, "application/json")
}
