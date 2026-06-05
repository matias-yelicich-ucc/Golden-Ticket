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
}
