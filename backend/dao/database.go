package dao

import (
	"fmt"
	"log"
	"os"
	"time"

	"golden-ticket/backend/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la instancia global de conexión a la base de datos
var DB *gorm.DB

// InitDB inicializa la conexión con MySQL utilizando GORM y realiza la automigración
func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	// Valores por defecto para desarrollo local
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3306"
	}
	if user == "" {
		user = "root"
	}
	if name == "" {
		name = "golden_ticket"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error al obtener sqlDB estándar desde GORM: %v", err)
	}

	// Configuración del connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("¡Conexión a la base de datos inicializada en la capa DAO!")

	// Automigración de esquemas
	if err := DB.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		log.Fatalf("Error durante la automigración: %v", err)
	}
	log.Println("¡Esquemas de la base de datos automigrados con éxito!")

	// Mantiene event_id nullable y evita recrear la FK en cada arranque.
	DB.Exec("ALTER TABLE tickets MODIFY COLUMN event_id BIGINT UNSIGNED NULL")
	if !DB.Migrator().HasConstraint(&domain.Ticket{}, "fk_events_tickets") &&
		!DB.Migrator().HasConstraint(&domain.Ticket{}, "fk_tickets_event") {
		DB.Exec(`ALTER TABLE tickets ADD CONSTRAINT fk_events_tickets 
			FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE SET NULL ON UPDATE CASCADE`)
	}
}
