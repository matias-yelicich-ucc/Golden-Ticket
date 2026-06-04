package controllers

import (
	"bytes"
	"encoding/json"
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

type mockEventDAO struct {
	events []*domain.Event
}

func (m *mockEventDAO) Create(event *domain.Event) error {
	event.ID = uint(len(m.events) + 1)
	m.events = append(m.events, event)
	return nil
}

func (m *mockEventDAO) GetAll() ([]*domain.Event, error) {
	return m.events, nil
}

func TestEventController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret_for_event_controller")
	defer os.Unsetenv("JWT_SECRET")

	mockDAO := &mockEventDAO{events: make([]*domain.Event, 0)}
	eventService := services.NewEventService(mockDAO)
	ctrl := NewEventController(eventService)

	router := gin.Default()

	// Setup JWT auth middleware on protected group
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.AuthorizeRole("administrador"))
		{
			adminOnly.POST("/events", ctrl.Create)
		}
	}

	// Helper to generate a valid admin token
	adminToken, err := utils.GenerateToken(1, "administrador")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Helper to generate a client token
	clientToken, err := utils.GenerateToken(2, "cliente")
	if err != nil {
		t.Fatalf("Failed to generate client token: %v", err)
	}

	// 1. Success: Admin creates valid event
	futureDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	validEvent := domain.EventCreateDTO{
		Titulo:      "Concierto de Rock",
		Descripcion: "Un show espectacular",
		Categoria:   "Música",
		Fecha:       futureDate,
		HoraInicio:  "21:00",
		HoraFin:     "23:30",
		Ubicacion:   "Estadio Mario Alberto Kempes",
		Coordenadas: "-31.4201, -64.1888",
		UrlImagen:   "https://example.com/imagen.jpg",
		Capacidad:   500,
	}
	body, _ := json.Marshal(validEvent)

	req, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var respEvent domain.EventResponseDTO
	_ = json.Unmarshal(w.Body.Bytes(), &respEvent)

	if respEvent.Titulo != "Concierto de Rock" {
		t.Errorf("Expected title 'Concierto de Rock', got '%s'", respEvent.Titulo)
	}
	if respEvent.EntradasDisponibles != 500 {
		t.Errorf("Expected EntradasDisponibles 500, got %d", respEvent.EntradasDisponibles)
	}
	if respEvent.Ubicacion != "Estadio Mario Alberto Kempes" {
		t.Errorf("Expected Ubicacion 'Estadio Mario Alberto Kempes', got '%s'", respEvent.Ubicacion)
	}

	// 2. Error: Required fields missing (empty title)
	invalidEventMissingFields := domain.EventCreateDTO{
		Titulo:     "",
		Fecha:      futureDate,
		HoraInicio: "21:00",
		HoraFin:    "23:30",
		Ubicacion:  "Estadio Mario Alberto Kempes",
		Capacidad:  500,
	}
	body2, _ := json.Marshal(invalidEventMissingFields)
	req2, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+adminToken)

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	// 3. Error: Past date
	pastDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	invalidEventPastDate := domain.EventCreateDTO{
		Titulo:     "Evento del Pasado",
		Fecha:      pastDate,
		HoraInicio: "12:00",
		HoraFin:    "14:00",
		Ubicacion:  "Estadio Mario Alberto Kempes",
		Capacidad:  100,
	}
	body3, _ := json.Marshal(invalidEventPastDate)
	req3, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body3))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "Bearer "+adminToken)

	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 Unprocessable Entity, got %d. Body: %s", w3.Code, w3.Body.String())
	}

	// 4. Error: Unauthorized (no token)
	req4, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body))
	req4.Header.Set("Content-Type", "application/json")

	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d. Body: %s", w4.Code, w4.Body.String())
	}

	// 5. Error: Forbidden (client role, not admin)
	req5, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("Authorization", "Bearer "+clientToken)

	w5 := httptest.NewRecorder()
	router.ServeHTTP(w5, req5)

	if w5.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d. Body: %s", w5.Code, w5.Body.String())
	}
}
