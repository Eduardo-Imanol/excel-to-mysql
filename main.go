package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Eduardo-Imanol/excel-to-mysql/db"
	"github.com/Eduardo-Imanol/excel-to-mysql/models"
	"github.com/Eduardo-Imanol/excel-to-mysql/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// getEnv obtiene el valor de una variable de entorno.
// Si la variable no existe, retorna el valor por defecto (fallback).
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// main punto de entrada de la aplicación.
// Establece la conexión a la BD, crea las tablas vía AutoMigrate,
// configura las rutas HTTP y arranca el servidor en el puerto indicado.
func main() {

	// Conectar a la base de datos MySQL
	db.DBconnection()

	// Auto-migrar modelos: crea las tablas si no existen y actualiza el esquema
	db.DB.AutoMigrate(models.Record{}, models.Upload{}, models.Sheet{}, models.Row{})

	// Crear router con Gorilla Mux
	r := mux.NewRouter()

	// ── Rutas de Health Check ──
	r.HandleFunc("/", routes.HomeHandler)

	// ── Rutas CRUD de Registros (calificaciones/nombres) ──
	r.HandleFunc("/names", routes.GetNamesHandler).Methods("GET")         // Listar todos
	r.HandleFunc("/names/{id}", routes.GetNameHandler).Methods("GET")    // Obtener uno por ID
	r.HandleFunc("/names", routes.PostNameHandler).Methods("POST")       // Crear uno
	r.HandleFunc("/names/all", routes.PostNamesHandler).Methods("POST")  // Crear varios
	r.HandleFunc("/names/{id}", routes.DeleteNamesHandler).Methods("DELETE") // Eliminar uno
	r.HandleFunc("/names", routes.DeleteAllRecordsHandler).Methods("DELETE") // Eliminar todos

	// ── Rutas de Carga de Excel ──
	r.HandleFunc("/excel/upload", routes.UploadExcelHandler).Methods("POST")   // Subir estructura Excel
	r.HandleFunc("/excel/uploads", routes.GetUploadsHandler).Methods("GET")   // Listar subidas

	// ── Configuración de CORS ──
	allowedOrigin := getEnv("CORS_ORIGIN", "*")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{allowedOrigin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// ── Arrancar servidor HTTP ──
	port := getEnv("PORT", "3000")
	log.Printf("Servidor escuchando en http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r))
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}

}
