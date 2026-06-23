package dao

import (
	"testing"
	"time"

	"golden-ticket/backend/domain"
)

func seedTicketTestData(t *testing.T) (domain.User, domain.User, domain.Event, domain.Event) {
	t.Helper()

	owner := domain.User{
		Nombre:   "Mati",
		Apellido: "Owner",
		Email:    "owner@ucc.edu.ar",
		Password: "hash",
		Rol:      "cliente",
		DNI:      "11111111",
	}
	dest := domain.User{
		Nombre:   "Ana",
		Apellido: "Dest",
		Email:    "dest@ucc.edu.ar",
		Password: "hash",
		Rol:      "cliente",
		DNI:      "22222222",
	}
	if err := DB.Create(&owner).Error; err != nil {
		t.Fatalf("failed to seed owner: %v", err)
	}
	if err := DB.Create(&dest).Error; err != nil {
		t.Fatalf("failed to seed destination user: %v", err)
	}

	future := futureEventModel()
	if err := DB.Create(&future).Error; err != nil {
		t.Fatalf("failed to seed future event: %v", err)
	}

	past := futureEventModel()
	past.Titulo = "Evento pasado"
	past.Fecha = time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	if err := DB.Create(&past).Error; err != nil {
		t.Fatalf("failed to seed past event: %v", err)
	}

	return owner, dest, future, past
}

func TestTicketDAOBuyAndGetByUser(t *testing.T) {
	setupTestDB(t)
	ticketDAO := NewTicketDAO()
	owner, _, future, past := seedTicketTestData(t)

	bought, err := ticketDAO.BuyTickets(owner.ID, future.ID, 2)
	if err != nil || len(bought) != 2 {
		t.Fatalf("expected two bought tickets, got len=%d err=%v", len(bought), err)
	}

	if _, err := ticketDAO.BuyTickets(owner.ID, future.ID, 20); err == nil {
		t.Fatalf("expected insufficient capacity error")
	}

	if _, err := ticketDAO.BuyTickets(owner.ID, past.ID, 1); err == nil {
		t.Fatalf("expected past event buy error")
	}

	userTickets, err := ticketDAO.GetByUserID(owner.ID)
	if err != nil || len(userTickets) != 2 {
		t.Fatalf("expected get by user to return 2 tickets, got len=%d err=%v", len(userTickets), err)
	}
	if userTickets[0].Event == nil {
		t.Fatalf("expected event to be preloaded")
	}
}

func TestTicketDAOTransferTicket(t *testing.T) {
	setupTestDB(t)
	ticketDAO := NewTicketDAO()
	owner, dest, future, _ := seedTicketTestData(t)

	ticket := domain.Ticket{UserID: owner.ID, EventID: &future.ID, Estado: "activo", FechaCompra: time.Now()}
	if err := DB.Create(&ticket).Error; err != nil {
		t.Fatalf("failed to seed ticket: %v", err)
	}

	if err := ticketDAO.TransferTicket(owner.ID, ticket.ID, dest.DNI); err != nil {
		t.Fatalf("expected transfer success, got %v", err)
	}
	var updated domain.Ticket
	if err := DB.First(&updated, ticket.ID).Error; err != nil {
		t.Fatalf("failed to reload ticket: %v", err)
	}
	if updated.UserID != dest.ID {
		t.Fatalf("expected owner to change to destination user")
	}

	if err := ticketDAO.TransferTicket(owner.ID, 999, dest.DNI); err == nil {
		t.Fatalf("expected missing ticket error")
	}
	if err := ticketDAO.TransferTicket(owner.ID, ticket.ID, "99999999"); err == nil {
		t.Fatalf("expected missing destination user error")
	}
	if err := ticketDAO.TransferTicket(owner.ID, ticket.ID, owner.DNI); err == nil {
		t.Fatalf("expected transfer to self error")
	}

	canceledTicket := domain.Ticket{UserID: owner.ID, EventID: &future.ID, Estado: "cancelado", FechaCompra: time.Now()}
	if err := DB.Create(&canceledTicket).Error; err != nil {
		t.Fatalf("failed to seed canceled ticket: %v", err)
	}
	if err := ticketDAO.TransferTicket(owner.ID, canceledTicket.ID, dest.DNI); err == nil {
		t.Fatalf("expected canceled ticket transfer error")
	}
}

func TestTicketDAOCancelTicket(t *testing.T) {
	setupTestDB(t)
	ticketDAO := NewTicketDAO()
	owner, dest, future, past := seedTicketTestData(t)

	activeTicket := domain.Ticket{UserID: owner.ID, EventID: &future.ID, Estado: "activo", FechaCompra: time.Now()}
	if err := DB.Create(&activeTicket).Error; err != nil {
		t.Fatalf("failed to seed active ticket: %v", err)
	}

	if err := ticketDAO.CancelTicket(owner.ID, activeTicket.ID); err != nil {
		t.Fatalf("expected cancel success, got %v", err)
	}
	var updated domain.Ticket
	if err := DB.Preload("Event").First(&updated, activeTicket.ID).Error; err != nil {
		t.Fatalf("failed to reload canceled ticket: %v", err)
	}
	if updated.Estado != "cancelado" {
		t.Fatalf("expected ticket to be canceled")
	}

	if err := ticketDAO.CancelTicket(owner.ID, 999); err == nil {
		t.Fatalf("expected missing ticket error")
	}
	if err := ticketDAO.CancelTicket(dest.ID, activeTicket.ID); err == nil {
		t.Fatalf("expected ownership error")
	}
	if err := ticketDAO.CancelTicket(owner.ID, activeTicket.ID); err == nil {
		t.Fatalf("expected already canceled error")
	}

	pastTicket := domain.Ticket{UserID: owner.ID, EventID: &past.ID, Estado: "activo", FechaCompra: time.Now()}
	if err := DB.Create(&pastTicket).Error; err != nil {
		t.Fatalf("failed to seed past ticket: %v", err)
	}
	if err := ticketDAO.CancelTicket(owner.ID, pastTicket.ID); err == nil {
		t.Fatalf("expected past event cancel error")
	}
}
