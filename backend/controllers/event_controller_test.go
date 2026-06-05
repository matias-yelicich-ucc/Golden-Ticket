package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func (m *mockEventDAO) GetByID(id uint) (*domain.Event, error) {
	for _, e := range m.events {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, errors.New("event not found")
}

func (m *mockEventDAO) Update(event *domain.Event) error {
	for i, e := range m.events {
		if e.ID == event.ID {
			m.events[i] = event
			return nil
		}
	}
	return errors.New("event not found")
}

func (m *mockEventDAO) GetAll(categoria string, buscar string) ([]*domain.Event, error) {
	if categoria == "" && buscar == "" {
		return m.events, nil
	}
	var res []*domain.Event
	for _, e := range m.events {
		match := true
		if categoria != "" && e.Categoria != categoria {
			match = false
		}
		if buscar != "" {
			titleMatch := false
			descMatch := false
			if e.Titulo != "" {
				titleMatch = strings.Contains(strings.ToLower(e.Titulo), strings.ToLower(buscar))
			}
			if e.Descripcion != "" {
				descMatch = strings.Contains(strings.ToLower(e.Descripcion), strings.ToLower(buscar))
			}
			if !titleMatch && !descMatch {
				match = false
			}
		}
		if match {
			res = append(res, e)
		}
	}
	return res, nil
}

func TestEventController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret_for_event_controller")
	defer os.Unsetenv("JWT_SECRET")

	mockDAO := &mockEventDAO{events: make([]*domain.Event, 0)}
	eventService := services.NewEventService(mockDAO)
	ctrl := NewEventController(eventService)

	router := gin.Default()
	router.GET("/events", ctrl.List)
	router.GET("/events/:id", ctrl.GetByID)

	// Setup JWT auth middleware on protected group
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.AuthorizeRole("administrador", "admin"))
		{
			adminOnly.POST("/events", ctrl.Create)
			adminOnly.PUT("/events/:id", ctrl.Update)
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
		Precio:      1500.00,
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
	if respEvent.Precio != 1500.00 {
		t.Errorf("Expected Precio 1500.00, got %f", respEvent.Precio)
	}

	// 2. Error: Required fields missing (empty title)
	invalidEventMissingFields := domain.EventCreateDTO{
		Titulo:     "",
		Fecha:      futureDate,
		HoraInicio: "21:00",
		HoraFin:    "23:30",
		Ubicacion:  "Estadio Mario Alberto Kempes",
		Capacidad:  500,
		Precio:     1500.00,
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
		Precio:     1500.00,
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

	// 6. Seed another event to mock DAO directly for testing GET /events
	theaterEvent := &domain.Event{
		Titulo:      "Obra de Teatro",
		Descripcion: "Comedia dramática en tres actos",
		Categoria:   "Teatro",
		Fecha:       futureDate,
		HoraInicio:  "20:00",
		HoraFin:     "22:00",
		Ubicacion:   "Teatro Libertador",
		Capacidad:   200,
		Precio:      800.00,
	}
	_ = mockDAO.Create(theaterEvent)

	// 7. GET /events (No filters)
	req6, _ := http.NewRequest("GET", "/events", nil)
	w6 := httptest.NewRecorder()
	router.ServeHTTP(w6, req6)

	if w6.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", w6.Code, w6.Body.String())
	}

	var eventsAll []domain.EventResponseDTO
	_ = json.Unmarshal(w6.Body.Bytes(), &eventsAll)
	if len(eventsAll) != 2 {
		t.Errorf("Expected 2 events, got %d", len(eventsAll))
	}

	// 8. GET /events?categoria=Teatro (Filter by category)
	req7, _ := http.NewRequest("GET", "/events?categoria=Teatro", nil)
	w7 := httptest.NewRecorder()
	router.ServeHTTP(w7, req7)

	if w7.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w7.Code)
	}

	var eventsTeatro []domain.EventResponseDTO
	_ = json.Unmarshal(w7.Body.Bytes(), &eventsTeatro)
	if len(eventsTeatro) != 1 || eventsTeatro[0].Titulo != "Obra de Teatro" {
		t.Errorf("Expected 1 event ('Obra de Teatro'), got %d", len(eventsTeatro))
	}

	// 9. GET /events?buscar=Rock (Filter by search keyword)
	req8, _ := http.NewRequest("GET", "/events?buscar=Rock", nil)
	w8 := httptest.NewRecorder()
	router.ServeHTTP(w8, req8)

	if w8.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w8.Code)
	}

	var eventsRock []domain.EventResponseDTO
	_ = json.Unmarshal(w8.Body.Bytes(), &eventsRock)
	if len(eventsRock) != 1 || eventsRock[0].Titulo != "Concierto de Rock" {
		t.Errorf("Expected 1 event ('Concierto de Rock'), got %d", len(eventsRock))
	}

	// 10. GET /events?categoria=Música&buscar=Obra (Combined filtering, no match)
	req9, _ := http.NewRequest("GET", "/events?categoria=M\u00fasica&buscar=Obra", nil)
	w9 := httptest.NewRecorder()
	router.ServeHTTP(w9, req9)

	if w9.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w9.Code)
	}

	var eventsCombined []domain.EventResponseDTO
	_ = json.Unmarshal(w9.Body.Bytes(), &eventsCombined)
	if len(eventsCombined) != 0 {
		t.Errorf("Expected 0 events, got %d", len(eventsCombined))
	}

	// 11. GET /events/1 (Success: single event by ID)
	req10, _ := http.NewRequest("GET", "/events/1", nil)
	w10 := httptest.NewRecorder()
	router.ServeHTTP(w10, req10)

	if w10.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", w10.Code, w10.Body.String())
	}

	var eventDetail domain.EventResponseDTO
	_ = json.Unmarshal(w10.Body.Bytes(), &eventDetail)
	if eventDetail.ID != 1 || eventDetail.Titulo != "Concierto de Rock" {
		t.Errorf("Expected event with ID 1 and title 'Concierto de Rock', got ID %d and title '%s'", eventDetail.ID, eventDetail.Titulo)
	}

	// 12. GET /events/999 (Error: event not found)
	req11, _ := http.NewRequest("GET", "/events/999", nil)
	w11 := httptest.NewRecorder()
	router.ServeHTTP(w11, req11)

	if w11.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", w11.Code)
	}

	// 13. GET /events/abc (Error: invalid ID format)
	req12, _ := http.NewRequest("GET", "/events/abc", nil)
	w12 := httptest.NewRecorder()
	router.ServeHTTP(w12, req12)

	if w12.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", w12.Code)
	}

	// 14. Success: update event 1 by admin
	futureDateStr := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	updateDTO := domain.EventCreateDTO{
		Titulo:      "Concierto de Rock Recargado",
		Descripcion: "Nueva descripción para el concierto",
		Categoria:   "Música",
		Fecha:       futureDateStr,
		HoraInicio:  "20:00",
		HoraFin:     "23:00",
		Ubicacion:   "Estadio River Plate",
		Coordenadas: "-34.5453,-58.4497",
		UrlImagen:   "http://example.com/river.jpg",
		Capacidad:   200,
		Precio:      3500.00,
	}
	bodyUpdate, _ := json.Marshal(updateDTO)
	req14, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyUpdate))
	req14.Header.Set("Authorization", "Bearer "+adminToken)
	req14.Header.Set("Content-Type", "application/json")

	w14 := httptest.NewRecorder()
	router.ServeHTTP(w14, req14)

	if w14.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", w14.Code, w14.Body.String())
	}

	// Verify update in DAO
	if mockDAO.events[0].Titulo != "Concierto de Rock Recargado" || mockDAO.events[0].Capacidad != 200 {
		t.Errorf("Expected title 'Concierto de Rock Recargado' and capacity 200, got '%s' and %d", mockDAO.events[0].Titulo, mockDAO.events[0].Capacidad)
	}

	// 15. Error: update event without admin token (e.g. client token)
	req15, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyUpdate))
	req15.Header.Set("Authorization", "Bearer "+clientToken)
	req15.Header.Set("Content-Type", "application/json")

	w15 := httptest.NewRecorder()
	router.ServeHTTP(w15, req15)

	if w15.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", w15.Code)
	}

	// 16. Error: update non-existent event
	req16, _ := http.NewRequest("PUT", "/admin/events/999", bytes.NewBuffer(bodyUpdate))
	req16.Header.Set("Authorization", "Bearer "+adminToken)
	req16.Header.Set("Content-Type", "application/json")

	w16 := httptest.NewRecorder()
	router.ServeHTTP(w16, req16)

	if w16.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", w16.Code)
	}

	// 17. Error: update event with past date
	pastDateStr := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	updatePastDTO := updateDTO
	updatePastDTO.Fecha = pastDateStr
	bodyPastUpdate, _ := json.Marshal(updatePastDTO)
	req17, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyPastUpdate))
	req17.Header.Set("Authorization", "Bearer "+adminToken)
	req17.Header.Set("Content-Type", "application/json")

	w17 := httptest.NewRecorder()
	router.ServeHTTP(w17, req17)

	if w17.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 Unprocessable Entity, got %d. Body: %s", w17.Code, w17.Body.String())
	}

	// 18. Error: update event capacity less than active sold tickets
	// Inject 2 active tickets to event 1 in mockDAO
	mockDAO.events[0].Tickets = []domain.Ticket{
		{Estado: "activo"},
		{Estado: "activo"},
	}
	updateLowCapDTO := updateDTO
	updateLowCapDTO.Capacidad = 1 // Less than 2
	bodyLowCapUpdate, _ := json.Marshal(updateLowCapDTO)
	req18, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyLowCapUpdate))
	req18.Header.Set("Authorization", "Bearer "+adminToken)
	req18.Header.Set("Content-Type", "application/json")

	w18 := httptest.NewRecorder()
	router.ServeHTTP(w18, req18)

	if w18.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 Unprocessable Entity, got %d. Body: %s", w18.Code, w18.Body.String())
	}
}
