package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBconnection() {
	var err error

	user := getEnv("MYSQL_USER", getEnv("DB_USER", "root"))
	pass := getEnv("MYSQL_PASSWORD", getEnv("DB_PASS", ""))
	host := getEnv("MYSQL_HOST", getEnv("DB_HOST", "localhost"))
	port := getEnv("MYSQL_PORT", getEnv("DB_PORT", "3306"))
	name := getEnv("MYSQL_DATABASE", getEnv("DB_NAME", "test_web_excel1"))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error conectando a la DB: ", err)
	} else {
		log.Println("DB connected")
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
