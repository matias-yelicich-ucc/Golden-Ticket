package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golden-ticket/backend/domain"
	"golden-ticket/backend/middleware"
	"golden-ticket/backend/services"
	"golden-ticket/backend/utils"

	"github.com/gin-gonic/gin"
)

type mockTicketDAO struct {
	events  map[uint]*domain.Event
	tickets []domain.Ticket
}

func (m *mockTicketDAO) BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error) {
	event, exists := m.events[eventID]
	if !exists {
		return nil, errors.New("record not found")
	}

	// Count active tickets in memory
	activeCount := 0
	for _, t := range m.tickets {
		if t.EventID == eventID && t.Estado == "activo" {
			activeCount++
		}
	}

	available := event.Capacidad - activeCount
	if available < cantidad {
		return nil, errors.New("capacidad insuficiente para realizar la compra")
	}

	// Validate date
	eventDateTimeStr := fmt.Sprintf("%sT%s:00", event.Fecha, event.HoraInicio)
	eventTime, err := time.ParseInLocation("2006-01-02T15:04:05", eventDateTimeStr, time.Local)
	if err == nil && eventTime.Before(time.Now()) {
		return nil, errors.New("no se pueden comprar entradas para un evento que ya ocurrió")
	}

	var newTickets []domain.Ticket
	for i := 0; i < cantidad; i++ {
		t := domain.Ticket{
			ID:          uint(len(m.tickets) + 1),
			UserID:      userID,
			EventID:     eventID,
			Estado:      "activo",
			FechaCompra: time.Now(),
		}
		newTickets = append(newTickets, t)
		m.tickets = append(m.tickets, t)
	}

	return newTickets, nil
}

func (m *mockTicketDAO) GetByUserID(userID uint) ([]domain.Ticket, error) {
	var res []domain.Ticket
	for _, t := range m.tickets {
		if t.UserID == userID {
			if ev, exists := m.events[t.EventID]; exists {
				t.Event = ev
			}
			res = append(res, t)
		}
	}
	return res, nil
}

func (m *mockTicketDAO) TransferTicket(userID uint, ticketID uint, destinationDNI string) error {
	var ticket *domain.Ticket
	idx := -1
	for i, t := range m.tickets {
		if t.ID == ticketID {
			ticket = &m.tickets[i]
			idx = i
			break
		}
	}
	if ticket == nil {
		return errors.New("entrada no encontrada")
	}

	if ticket.UserID != userID {
		return errors.New("no eres el propietario de esta entrada")
	}

	if ticket.Estado != "activo" {
		return errors.New("no se puede transferir una entrada cancelada")
	}

	var destUserID uint
	if destinationDNI == "87654321" {
		destUserID = 3
	} else if destinationDNI == "12345678" {
		destUserID = 2
	} else {
		return errors.New("no existe ningún usuario registrado con el DNI ingresado")
	}

	if destUserID == userID {
		return errors.New("no podés transferirte una entrada a vos mismo")
	}

	ticket.UserID = destUserID
	m.tickets[idx].UserID = destUserID
	return nil
}

func (m *mockTicketDAO) CancelTicket(userID uint, ticketID uint) error {
	var ticket *domain.Ticket
	for i := range m.tickets {
		if m.tickets[i].ID == ticketID {
			ticket = &m.tickets[i]
			break
		}
	}
	if ticket == nil {
		return errors.New("entrada no encontrada")
	}
	if ticket.UserID != userID {
		return errors.New("no eres el propietario de esta entrada")
	}
	if ticket.Estado != "activo" {
		return errors.New("la entrada ya se encuentra cancelada")
	}

	if ev, exists := m.events[ticket.EventID]; exists && ev != nil {
		eventDateTimeStr := fmt.Sprintf("%sT%s:00", ev.Fecha, ev.HoraInicio)
		eventTime, err := time.ParseInLocation("2006-01-02T15:04:05", eventDateTimeStr, time.Local)
		if err == nil && eventTime.Before(time.Now()) {
			return errors.New("no se pueden cancelar entradas para un evento que ya ocurrió o está en curso")
		}
	}

	ticket.Estado = "cancelado"
	return nil
}

func TestTicketController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret_for_ticket_controller")
	defer os.Unsetenv("JWT_SECRET")

	futureDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	pastDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	mockEvents := map[uint]*domain.Event{
		1: {
			ID:         1,
			Titulo:     "Concierto de Rock",
			Fecha:      futureDate,
			HoraInicio: "21:00",
			Capacidad:  10,
		},
		2: {
			ID:         2,
			Titulo:     "Charla del Pasado",
			Fecha:      pastDate,
			HoraInicio: "10:00",
			Capacidad:  100,
		},
	}

	mockDAO := &mockTicketDAO{
		events:  mockEvents,
		tickets: make([]domain.Ticket, 0),
	}
	ticketService := services.NewTicketService(mockDAO)
	ctrl := NewTicketController(ticketService)

	router := gin.Default()
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/events/:id/tickets", ctrl.Buy)
		protected.GET("/my-tickets", ctrl.GetMyTickets)
		protected.POST("/my-tickets/:id/transfer", ctrl.Transfer)
		protected.POST("/my-tickets/:id/cancel", ctrl.Cancel)
	}

	// Helper to generate a valid client token
	clientToken, err := utils.GenerateToken(2, "cliente")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 1. Success: buy 2 tickets for valid event
	reqBody1 := domain.TicketPurchaseDTO{Cantidad: 2}
	body1, _ := json.Marshal(reqBody1)

	req1, _ := http.NewRequest("POST", "/events/1/tickets", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+clientToken)

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d. Body: %s", w1.Code, w1.Body.String())
	}

	var resTickets []domain.Ticket
	_ = json.Unmarshal(w1.Body.Bytes(), &resTickets)
	if len(resTickets) != 2 {
		t.Errorf("Expected 2 tickets in response, got %d", len(resTickets))
	}

	// 2. Error: buy too many tickets (insufficient capacity)
	reqBody2 := domain.TicketPurchaseDTO{Cantidad: 15}
	body2, _ := json.Marshal(reqBody2)

	req2, _ := http.NewRequest("POST", "/events/1/tickets", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+clientToken)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected 409 Conflict, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	// 3. Error: buy tickets for a past event
	reqBody3 := domain.TicketPurchaseDTO{Cantidad: 1}
	body3, _ := json.Marshal(reqBody3)

	req3, _ := http.NewRequest("POST", "/events/2/tickets", bytes.NewBuffer(body3))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "Bearer "+clientToken)

	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", w3.Code, w3.Body.String())
	}

	// 4. Error: unauthorized request (missing token)
	req4, _ := http.NewRequest("POST", "/events/1/tickets", bytes.NewBuffer(body1))
	req4.Header.Set("Content-Type", "application/json")

	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w4.Code)
	}

	// 5. Error: event not found (ID 999)
	req5, _ := http.NewRequest("POST", "/events/999/tickets", bytes.NewBuffer(body1))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("Authorization", "Bearer "+clientToken)

	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)

	if w5.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d. Body: %s", w5.Code, w5.Body.String())
	}

	// 6. Success: get tickets of current user (clientToken belongs to userID 2)
	req6, _ := http.NewRequest("GET", "/my-tickets", nil)
	req6.Header.Set("Authorization", "Bearer "+clientToken)

	w6 := httptest.NewRecorder()
	router.ServeHTTP(w6, req6)

	if w6.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", w6.Code, w6.Body.String())
	}

	var myTickets []domain.Ticket
	_ = json.Unmarshal(w6.Body.Bytes(), &myTickets)
	if len(myTickets) != 2 {
		t.Errorf("Expected 2 tickets, got %d", len(myTickets))
	}
	if myTickets[0].Event == nil || myTickets[0].Event.Titulo != "Concierto de Rock" {
		t.Errorf("Expected Event details populated in ticket")
	}

	// 7. Success: get tickets for a user with no tickets (userID 99)
	otherToken, _ := utils.GenerateToken(99, "cliente")
	req7, _ := http.NewRequest("GET", "/my-tickets", nil)
	req7.Header.Set("Authorization", "Bearer "+otherToken)

	w7 := httptest.NewRecorder()
	router.ServeHTTP(w7, req7)

	if w7.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w7.Code)
	}

	var otherTickets []domain.Ticket
	_ = json.Unmarshal(w7.Body.Bytes(), &otherTickets)
	if len(otherTickets) != 0 {
		t.Errorf("Expected 0 tickets, got %d", len(otherTickets))
	}

	// 8. Error: get tickets unauthorized (missing token)
	req8, _ := http.NewRequest("GET", "/my-tickets", nil)
	w8 := httptest.NewRecorder()
	router.ServeHTTP(w8, req8)

	if w8.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w8.Code)
	}

	// 9. Success: transfer ticket 1 to Maria (DNI 87654321)
	transferPayload := TicketTransferDTO{DNI: "87654321"}
	bodyTransfer, _ := json.Marshal(transferPayload)
	req9, _ := http.NewRequest("POST", "/my-tickets/1/transfer", bytes.NewBuffer(bodyTransfer))
	req9.Header.Set("Authorization", "Bearer "+clientToken)
	req9.Header.Set("Content-Type", "application/json")

	w9 := httptest.NewRecorder()
	router.ServeHTTP(w9, req9)

	if w9.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", w9.Code, w9.Body.String())
	}

	// Verify ownership changed in database
	if mockDAO.tickets[0].UserID != 3 {
		t.Errorf("Expected owner of ticket 1 to be user 3, got user %d", mockDAO.tickets[0].UserID)
	}

	// 10. Error: transfer ticket you don't own (otherToken belongs to user 99, tries to transfer ticket 1, now owned by user 3)
	req10, _ := http.NewRequest("POST", "/my-tickets/1/transfer", bytes.NewBuffer(bodyTransfer))
	req10.Header.Set("Authorization", "Bearer "+otherToken)
	req10.Header.Set("Content-Type", "application/json")

	w10 := httptest.NewRecorder()
	router.ServeHTTP(w10, req10)

	if w10.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", w10.Code)
	}

	// 11. Error: transfer to non-existent DNI
	transferInvalidPayload := TicketTransferDTO{DNI: "99999999"}
	bodyInvalidTransfer, _ := json.Marshal(transferInvalidPayload)
	// Maria (user 3) now owns ticket 1. Let's generate a token for user 3
	mariaToken, _ := utils.GenerateToken(3, "cliente")
	req11, _ := http.NewRequest("POST", "/my-tickets/1/transfer", bytes.NewBuffer(bodyInvalidTransfer))
	req11.Header.Set("Authorization", "Bearer "+mariaToken)
	req11.Header.Set("Content-Type", "application/json")

	w11 := httptest.NewRecorder()
	router.ServeHTTP(w11, req11)

	if w11.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d. Body: %s", w11.Code, w11.Body.String())
	}

	// 12. Error: transfer to oneself
	transferSelfPayload := TicketTransferDTO{DNI: "87654321"} // Maria (user 3) DNI
	bodySelfTransfer, _ := json.Marshal(transferSelfPayload)
	req12, _ := http.NewRequest("POST", "/my-tickets/1/transfer", bytes.NewBuffer(bodySelfTransfer))
	req12.Header.Set("Authorization", "Bearer "+mariaToken)
	req12.Header.Set("Content-Type", "application/json")

	w12 := httptest.NewRecorder()
	router.ServeHTTP(w12, req12)

	if w12.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", w12.Code, w12.Body.String())
	}

	// 13. Error: transfer cancelled ticket
	// Cancel ticket 1 in memory
	mockDAO.tickets[0].Estado = "cancelado"
	req13, _ := http.NewRequest("POST", "/my-tickets/1/transfer", bytes.NewBuffer(bodyTransfer))
	req13.Header.Set("Authorization", "Bearer "+mariaToken)
	req13.Header.Set("Content-Type", "application/json")

	w13 := httptest.NewRecorder()
	router.ServeHTTP(w13, req13)

	if w13.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", w13.Code, w13.Body.String())
	}

	// 14. Error: cancel ticket not owned by user
	// Ticket 2 is owned by user 2 (clientToken). Try to cancel using mariaToken (user 3)
	req14, _ := http.NewRequest("POST", "/my-tickets/2/cancel", nil)
	req14.Header.Set("Authorization", "Bearer "+mariaToken)

	w14 := httptest.NewRecorder()
	router.ServeHTTP(w14, req14)

	if w14.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d. Body: %s", w14.Code, w14.Body.String())
	}

	// 15. Success: cancel ticket 2 owned by user 2 (clientToken)
	req15, _ := http.NewRequest("POST", "/my-tickets/2/cancel", nil)
	req15.Header.Set("Authorization", "Bearer "+clientToken)

	w15 := httptest.NewRecorder()
	router.ServeHTTP(w15, req15)

	if w15.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", w15.Code, w15.Body.String())
	}

	if mockDAO.tickets[1].Estado != "cancelado" {
		t.Errorf("Expected ticket 2 Estado to be cancelado, got %s", mockDAO.tickets[1].Estado)
	}

	// 16. Error: cancel already cancelled ticket
	req16, _ := http.NewRequest("POST", "/my-tickets/2/cancel", nil)
	req16.Header.Set("Authorization", "Bearer "+clientToken)

	w16 := httptest.NewRecorder()
	router.ServeHTTP(w16, req16)

	if w16.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", w16.Code, w16.Body.String())
	}

	// 17. Error: cancel ticket of event that has already occurred (event 2 is in the past)
	// First let's inject a ticket for event 2 (past event) owned by user 2 in mockDAO
	mockDAO.tickets = append(mockDAO.tickets, domain.Ticket{
		ID:          3,
		UserID:      2,
		EventID:     2,
		Estado:      "activo",
		FechaCompra: time.Now(),
	})
	req17, _ := http.NewRequest("POST", "/my-tickets/3/cancel", nil)
	req17.Header.Set("Authorization", "Bearer "+clientToken)

	w17 := httptest.NewRecorder()
	router.ServeHTTP(w17, req17)

	if w17.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", w17.Code, w17.Body.String())
	}

	// 18. Error: cancel non-existent ticket
	req18, _ := http.NewRequest("POST", "/my-tickets/999/cancel", nil)
	req18.Header.Set("Authorization", "Bearer "+clientToken)

	w18 := httptest.NewRecorder()
	router.ServeHTTP(w18, req18)

	if w18.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d. Body: %s", w18.Code, w18.Body.String())
	}
}
