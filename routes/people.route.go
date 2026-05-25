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

const getter = "id, nombre, math, physical, chemistry, biologi, histori, geografi, literature, spanish, english"
const maxBodySize = 1 << 20 // 1 MB

func GetNamesHandler(w http.ResponseWriter, r *http.Request) {
	var calsNames []models.Cal

	db.DB.Select(getter).Find(&calsNames)

	json.NewEncoder(w).Encode(&calsNames)
}

func GetNameHandler(w http.ResponseWriter, r *http.Request) {
	var calsName models.Cal
	id := mux.Vars(r)

	db.DB.Select(getter).Where("id = ?", id["id"]).First(&calsName)

	if calsName.ID == 0 {
		http.Error(w, "id no encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&calsName)
}

func PostNameHandler(w http.ResponseWriter, r *http.Request) {
	if !isJSONContent(r) {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	var calsName models.Cal
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	if err := json.NewDecoder(r.Body).Decode(&calsName); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := validateCal(&calsName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB.Create(&calsName).Error; err != nil {
		http.Error(w, "Error al crear el registro", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&calsName)
}

func PostNamesHandler(w http.ResponseWriter, r *http.Request) {
	if !isJSONContent(r) {
		http.Error(w, "Content-Type debe ser application/json", http.StatusUnsupportedMediaType)
		return
	}

	var calsNames []models.Cal
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	if err := json.NewDecoder(r.Body).Decode(&calsNames); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if len(calsNames) == 0 {
		http.Error(w, "El array no puede estar vacío", http.StatusBadRequest)
		return
	}

	for i := range calsNames {
		if err := validateCal(&calsNames[i]); err != nil {
			http.Error(w, fmt.Sprintf("Error en elemento %d: %s", i, err.Error()), http.StatusBadRequest)
			return
		}
		if err := db.DB.Create(&calsNames[i]).Error; err != nil {
			http.Error(w, "Error al crear los registros", http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(&calsNames)
}

func DeleteNamesHandler(w http.ResponseWriter, r *http.Request) {
	var calsName models.Cal
	id := mux.Vars(r)

	db.DB.Select(getter).Where("id = ?", id["id"]).First(&calsName)

	if calsName.ID == 0 {
		http.Error(w, "id no encontrado", http.StatusNotFound)
		return
	}

	db.DB.Unscoped().Delete(&calsName)
	w.WriteHeader(http.StatusOK)
}

func DeleteAllRecordsHandler(w http.ResponseWriter, r *http.Request) {
	var calsNames []models.Cal

	db.DB.Unscoped().Select(getter).Find(&calsNames)

	if len(calsNames) == 0 {
		http.Error(w, "No se encontraron registros para eliminar", http.StatusNotFound)
		return
	}

	result := db.DB.Unscoped().Delete(&calsNames)
	if result.Error != nil {
		http.Error(w, "Error al borrar los registros", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registros eliminados exitosamente"))
}

func isJSONContent(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	return strings.HasPrefix(ct, "application/json")
}

func validateCal(c *models.Cal) error {
	if strings.TrimSpace(c.Nombre) == "" {
		return fmt.Errorf("el campo 'nombre' es requerido")
	}
	return nil
}

