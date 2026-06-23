package dao

import (
	"testing"

	"golden-ticket/backend/domain"
)

func TestEventDAOCRUDAndQueries(t *testing.T) {
	setupTestDB(t)
	eventDAO := NewEventDAO()

	eventOne := futureEventModel()
	eventOne.Titulo = "Festival Rock"
	eventOne.Categoria = "Musica"
	eventOne.Ubicacion = "Cordoba"
	if err := eventDAO.Create(&eventOne); err != nil {
		t.Fatalf("expected create success, got %v", err)
	}

	eventTwo := futureEventModel()
	eventTwo.Titulo = "Obra Teatro"
	eventTwo.Categoria = "Teatro"
	eventTwo.Ubicacion = "Carlos Paz"
	eventTwo.Precio = 2000
	if err := eventDAO.Create(&eventTwo); err != nil {
		t.Fatalf("expected create success, got %v", err)
	}

	user := domain.User{
		Nombre:   "Mati",
		Apellido: "Yelicich",
		Email:    "mati@ucc.edu.ar",
		Password: "hash",
		Rol:      "cliente",
		DNI:      "44555666",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	tickets := []domain.Ticket{
		{UserID: user.ID, EventID: &eventOne.ID, Estado: "activo"},
		{UserID: user.ID, EventID: &eventOne.ID, Estado: "activo"},
		{UserID: user.ID, EventID: &eventOne.ID, Estado: "cancelado"},
	}
	if err := DB.Create(&tickets).Error; err != nil {
		t.Fatalf("failed to seed tickets: %v", err)
	}

	all, err := eventDAO.GetAll("", "")
	if err != nil || len(all) != 2 {
		t.Fatalf("expected 2 events, got len=%d err=%v", len(all), err)
	}

	filtered, err := eventDAO.GetAll("Teatro", "")
	if err != nil || len(filtered) != 1 || filtered[0].Titulo != "Obra Teatro" {
		t.Fatalf("expected teatro filter to return one event, got %+v err=%v", filtered, err)
	}

	search, err := eventDAO.GetAll("", "Cordoba")
	if err != nil || len(search) != 1 || search[0].Titulo != "Festival Rock" {
		t.Fatalf("expected search to return Festival Rock, got %+v err=%v", search, err)
	}

	byID, err := eventDAO.GetByID(eventOne.ID)
	if err != nil || len(byID.Tickets) != 3 {
		t.Fatalf("expected preload tickets, got tickets=%d err=%v", len(byID.Tickets), err)
	}

	stats, err := eventDAO.GetAdminDashboardStats()
	if err != nil {
		t.Fatalf("expected stats success, got %v", err)
	}
	if stats.TotalEventos != 2 || stats.EntradasVendidas != 2 || stats.RecaudacionTotal != 3000 {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	eventOne.Titulo = "Festival Rock Deluxe"
	if err := eventDAO.Update(&eventOne); err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	updated, _ := eventDAO.GetByID(eventOne.ID)
	if updated.Titulo != "Festival Rock Deluxe" {
		t.Fatalf("expected updated title, got %q", updated.Titulo)
	}

	if err := eventDAO.Delete(eventOne.ID); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
	if _, err := eventDAO.GetByID(eventOne.ID); err == nil {
		t.Fatalf("expected deleted event to be missing")
	}

	var canceledTickets int64
	if err := DB.Model(&domain.Ticket{}).Where("estado = ?", "cancelado").Count(&canceledTickets).Error; err != nil {
		t.Fatalf("failed to count canceled tickets: %v", err)
	}
	if canceledTickets < 3 {
		t.Fatalf("expected active tickets to be canceled on delete, got %d canceled", canceledTickets)
	}
}
