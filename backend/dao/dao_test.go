package dao

import (
	"fmt"
	"testing"
	"time"

	"golden-ticket/backend/domain"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared&_foreign_keys=1", t.Name())
	testDB, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if err := testDB.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		t.Fatalf("failed to migrate sqlite db: %v", err)
	}

	DB = testDB
}

func futureEventModel() domain.Event {
	return domain.Event{
		Titulo:      "Evento futuro",
		Descripcion: "Descripcion",
		Categoria:   "Musica",
		Fecha:       time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		HoraInicio:  "20:00",
		HoraFin:     "22:00",
		Ubicacion:   "Cordoba",
		Coordenadas: "-31,-64",
		UrlImagen:   "https://example.com/evento.jpg",
		Capacidad:   10,
		Precio:      1500,
	}
}
