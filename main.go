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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {

	db.DBconnection()

	db.DB.AutoMigrate(models.Cal{}, models.Upload{}, models.Sheet{}, models.Row{})

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	r.HandleFunc("/names", routes.GetNamesHandler).Methods("GET")
	r.HandleFunc("/names/{id}", routes.GetNameHandler).Methods("GET")
	r.HandleFunc("/names", routes.PostNameHandler).Methods("POST")
	r.HandleFunc("/names/all", routes.PostNamesHandler).Methods("POST")
	r.HandleFunc("/names/{id}", routes.DeleteNamesHandler).Methods("DELETE")
	r.HandleFunc("/names", routes.DeleteAllRecordsHandler).Methods("DELETE")

	r.HandleFunc("/excel/upload", routes.UploadExcelHandler).Methods("POST")
	r.HandleFunc("/excel/uploads", routes.GetUploadsHandler).Methods("GET")

	allowedOrigin := getEnv("CORS_ORIGIN", "*")
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{allowedOrigin})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	port := getEnv("PORT", "3000")
	err := http.ListenAndServe(":"+port, handlers.CORS(originsOk, headersOk, methodsOk)(r))
	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}

}
