package services

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"golden-ticket/backend/domain"
)

type ticketServiceMockDAO struct {
	buyResponse     []domain.Ticket
	buyErr          error
	userTickets     []domain.Ticket
	getByUserErr    error
	transferErr     error
	cancelErr       error
	lastBuyUserID   uint
	lastBuyEventID  uint
	lastBuyCantidad int
	lastGetUserID   uint
	lastTransferID  uint
	lastTransferDNI string
	lastCancelID    uint
	lastActionUser  uint
}

func (m *ticketServiceMockDAO) BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error) {
	m.lastBuyUserID = userID
	m.lastBuyEventID = eventID
	m.lastBuyCantidad = cantidad
	return m.buyResponse, m.buyErr
}

func (m *ticketServiceMockDAO) GetByUserID(userID uint) ([]domain.Ticket, error) {
	m.lastGetUserID = userID
	return m.userTickets, m.getByUserErr
}

func (m *ticketServiceMockDAO) TransferTicket(userID uint, ticketID uint, destinationDNI string) error {
	m.lastActionUser = userID
	m.lastTransferID = ticketID
	m.lastTransferDNI = destinationDNI
	return m.transferErr
}

func (m *ticketServiceMockDAO) CancelTicket(userID uint, ticketID uint) error {
	m.lastActionUser = userID
	m.lastCancelID = ticketID
	return m.cancelErr
}

func TestTicketServiceDelegatesToDAO(t *testing.T) {
	now := time.Now()
	expectedTickets := []domain.Ticket{{ID: 1, FechaCompra: now}}
	dao := &ticketServiceMockDAO{
		buyResponse: expectedTickets,
		userTickets: expectedTickets,
	}
	service := NewTicketService(dao)

	bought, err := service.BuyTickets(7, 9, 2)
	if err != nil {
		t.Fatalf("expected buy success, got %v", err)
	}
	if dao.lastBuyUserID != 7 || dao.lastBuyEventID != 9 || dao.lastBuyCantidad != 2 {
		t.Fatalf("expected buy arguments to be forwarded")
	}
	if !reflect.DeepEqual(bought, expectedTickets) {
		t.Fatalf("expected buy response to be forwarded")
	}

	list, err := service.GetTicketsByUserID(7)
	if err != nil {
		t.Fatalf("expected list success, got %v", err)
	}
	if dao.lastGetUserID != 7 || !reflect.DeepEqual(list, expectedTickets) {
		t.Fatalf("expected get tickets to be forwarded")
	}

	if err := service.TransferTicket(7, 3, "44555666"); err != nil {
		t.Fatalf("expected transfer success, got %v", err)
	}
	if dao.lastActionUser != 7 || dao.lastTransferID != 3 || dao.lastTransferDNI != "44555666" {
		t.Fatalf("expected transfer args to be forwarded")
	}

	if err := service.CancelTicket(7, 3); err != nil {
		t.Fatalf("expected cancel success, got %v", err)
	}
	if dao.lastActionUser != 7 || dao.lastCancelID != 3 {
		t.Fatalf("expected cancel args to be forwarded")
	}
}

func TestTicketServiceReturnsDAOErrors(t *testing.T) {
	dao := &ticketServiceMockDAO{
		buyErr:       errors.New("buy failed"),
		getByUserErr: errors.New("list failed"),
		transferErr:  errors.New("transfer failed"),
		cancelErr:    errors.New("cancel failed"),
	}
	service := NewTicketService(dao)

	if _, err := service.BuyTickets(1, 2, 1); err == nil || err.Error() != "buy failed" {
		t.Fatalf("expected buy error, got %v", err)
	}
	if _, err := service.GetTicketsByUserID(1); err == nil || err.Error() != "list failed" {
		t.Fatalf("expected list error, got %v", err)
	}
	if err := service.TransferTicket(1, 2, "123"); err == nil || err.Error() != "transfer failed" {
		t.Fatalf("expected transfer error, got %v", err)
	}
	if err := service.CancelTicket(1, 2); err == nil || err.Error() != "cancel failed" {
		t.Fatalf("expected cancel error, got %v", err)
	}
}
