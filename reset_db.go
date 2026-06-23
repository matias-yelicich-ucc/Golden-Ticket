package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	if host == "" { host = "localhost" }
	if port == "" { port = "3306" }
	if user == "" { user = "root" }
	if name == "" { name = "golden_ticket" }

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	// Update user 1 to DNI "11111111" and user 2 to "22222222"
	res1 := db.Table("users").Where("id = ?", 1).Update("dni", "11111111")
	if res1.Error != nil {
		log.Printf("Error updating user 1: %v", res1.Error)
	} else {
		fmt.Printf("Updated user 1. Rows affected: %d\n", res1.RowsAffected)
	}

	res2 := db.Table("users").Where("id = ?", 2).Update("dni", "22222222")
	if res2.Error != nil {
		log.Printf("Error updating user 2: %v", res2.Error)
	} else {
		fmt.Printf("Updated user 2. Rows affected: %d\n", res2.RowsAffected)
	}
}
