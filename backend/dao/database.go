package dao

import (
	"database/sql"
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

type dbConfig struct {
	host string
	port string
	user string
	pass string
	name string
}

func loadDBConfigFromEnv() dbConfig {
	cfg := dbConfig{
		host: os.Getenv("DB_HOST"),
		port: os.Getenv("DB_PORT"),
		user: os.Getenv("DB_USER"),
		pass: os.Getenv("DB_PASSWORD"),
		name: os.Getenv("DB_NAME"),
	}

	if cfg.host == "" {
		cfg.host = "localhost"
	}
	if cfg.port == "" {
		cfg.port = "3306"
	}
	if cfg.user == "" {
		cfg.user = "root"
	}
	if cfg.name == "" {
		cfg.name = "golden_ticket"
	}

	return cfg
}

func buildDSN(cfg dbConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.user, cfg.pass, cfg.host, cfg.port, cfg.name)
}

func openDatabase(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func setupConnectionPool(sqlDB *sql.DB) {
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func migrateAndPrepareDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		return err
	}

	db.Exec("ALTER TABLE tickets MODIFY COLUMN event_id BIGINT UNSIGNED NULL")
	if !db.Migrator().HasConstraint(&domain.Ticket{}, "fk_events_tickets") &&
		!db.Migrator().HasConstraint(&domain.Ticket{}, "fk_tickets_event") {
		db.Exec(`ALTER TABLE tickets ADD CONSTRAINT fk_events_tickets 
			FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE SET NULL ON UPDATE CASCADE`)
	}

	return nil
}

func initializeDB(dsn string, openFn func(string) (*gorm.DB, error)) error {
	var err error
	DB, err = openFn(dsn)
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	setupConnectionPool(sqlDB)
	if err := migrateAndPrepareDatabase(DB); err != nil {
		return err
	}
	return nil
}

// InitDB inicializa la conexión con MySQL utilizando GORM y realiza la automigración
func InitDB() {
	cfg := loadDBConfigFromEnv()
	dsn := buildDSN(cfg)

	if err := initializeDB(dsn, openDatabase); err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}

	log.Println("¡Conexión a la base de datos inicializada en la capa DAO!")
	log.Println("¡Esquemas de la base de datos automigrados con éxito!")
}
