package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB es la instancia global de conexión a la base de datos (GORM).
// Se utiliza en toda la aplicación para consultas y operaciones CRUD.
var DB *gorm.DB

// DBconnection establece la conexión con MySQL usando las variables de entorno.
// Prioriza las variables MYSQL_* (Railway) y usa DB_* como fallback.
// Termina la ejecución si falla la conexión.
func DBconnection() {
	var err error

	// Lectura de variables de entorno con fallback a valores por defecto
	user := getEnv("MYSQL_USER", getEnv("DB_USER", "root"))
	pass := getEnv("MYSQL_PASSWORD", getEnv("DB_PASS", ""))
	host := getEnv("MYSQL_HOST", getEnv("DB_HOST", "localhost"))
	port := getEnv("MYSQL_PORT", getEnv("DB_PORT", "3306"))
	name := getEnv("MYSQL_DATABASE", getEnv("DB_NAME", "test_web_excel1"))

	// Construir el DSN (Data Source Name) para la conexión MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	// Abrir conexión con GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error conectando a la DB: ", err)
	} else {
		log.Println("DB connected")
	}
}

// getEnv obtiene el valor de una variable de entorno.
// Si no existe, retorna el valor por defecto (fallback).
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
