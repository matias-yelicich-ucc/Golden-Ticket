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
	for _, event := range m.events {
		if event.ID == id {
			return event, nil
		}
	}
	return nil, errors.New("event not found")
}

func (m *mockEventDAO) GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error) {
	totalEventos := len(m.events)
	entradasVendidas := 0
	capacidadTotal := 0
	recaudacionTotal := 0.0

	for _, event := range m.events {
		capacidadTotal += event.Capacidad
		for _, ticket := range event.Tickets {
			if ticket.Estado == "activo" {
				entradasVendidas++
				recaudacionTotal += event.Precio
			}
		}
	}

	ocupacionMedia := 0.0
	if capacidadTotal > 0 {
		ocupacionMedia = (float64(entradasVendidas) / float64(capacidadTotal)) * 100
	}

	return &domain.AdminDashboardStatsDTO{
		TotalEventos:     totalEventos,
		EntradasVendidas: entradasVendidas,
		OcupacionMedia:   ocupacionMedia,
		RecaudacionTotal: recaudacionTotal,
	}, nil
}

func (m *mockEventDAO) Update(event *domain.Event) error {
	for index, currentEvent := range m.events {
		if currentEvent.ID == event.ID {
			m.events[index] = event
			return nil
		}
	}
	return errors.New("event not found")
}

func (m *mockEventDAO) Delete(id uint) error {
	for index, event := range m.events {
		if event.ID == id {
			m.events = append(m.events[:index], m.events[index+1:]...)
			return nil
		}
	}
	return errors.New("event not found")
}

func (m *mockEventDAO) GetAll(categoria string, buscar string) ([]*domain.Event, error) {
	if categoria == "" && buscar == "" {
		return m.events, nil
	}

	var response []*domain.Event
	for _, event := range m.events {
		match := true
		if categoria != "" && event.Categoria != categoria {
			match = false
		}
		if buscar != "" {
			titleMatch := false
			descriptionMatch := false
			categoryMatch := false
			locationMatch := false
			if event.Titulo != "" {
				titleMatch = strings.Contains(strings.ToLower(event.Titulo), strings.ToLower(buscar))
			}
			if event.Descripcion != "" {
				descriptionMatch = strings.Contains(strings.ToLower(event.Descripcion), strings.ToLower(buscar))
			}
			if event.Categoria != "" {
				categoryMatch = strings.Contains(strings.ToLower(event.Categoria), strings.ToLower(buscar))
			}
			if event.Ubicacion != "" {
				locationMatch = strings.Contains(strings.ToLower(event.Ubicacion), strings.ToLower(buscar))
			}
			if !titleMatch && !descriptionMatch && !categoryMatch && !locationMatch {
				match = false
			}
		}
		if match {
			response = append(response, event)
		}
	}

	return response, nil
}

func TestEventController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret_for_event_controller")
	defer os.Unsetenv("JWT_SECRET")

	mockDAO := &mockEventDAO{events: make([]*domain.Event, 0)}
	eventService := services.NewEventService(mockDAO)
	controller := NewEventController(eventService)

	router := gin.Default()
	router.GET("/events", controller.List)
	router.GET("/events/:id", controller.GetByID)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.AuthorizeRole("administrador", "admin"))
		{
			adminOnly.GET("/dashboard", controller.GetAdminDashboardStats)
			adminOnly.POST("/events", controller.Create)
			adminOnly.PUT("/events/:id", controller.Update)
			adminOnly.DELETE("/events/:id", controller.Delete)
		}
	}

	adminToken, err := utils.GenerateToken(1, "administrador")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	clientToken, err := utils.GenerateToken(2, "cliente")
	if err != nil {
		t.Fatalf("Failed to generate client token: %v", err)
	}

	futureDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	validEvent := domain.EventCreateDTO{
		Titulo:      "Concierto de Rock",
		Descripcion: "Un show espectacular",
		Categoria:   "Musica",
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

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, req)

	if responseRecorder.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %d. Body: %s", responseRecorder.Code, responseRecorder.Body.String())
	}

	var responseEvent domain.EventResponseDTO
	_ = json.Unmarshal(responseRecorder.Body.Bytes(), &responseEvent)

	if responseEvent.Titulo != "Concierto de Rock" {
		t.Errorf("Expected title 'Concierto de Rock', got '%s'", responseEvent.Titulo)
	}
	if responseEvent.EntradasDisponibles != 500 {
		t.Errorf("Expected EntradasDisponibles 500, got %d", responseEvent.EntradasDisponibles)
	}
	if responseEvent.Ubicacion != "Estadio Mario Alberto Kempes" {
		t.Errorf("Expected Ubicacion 'Estadio Mario Alberto Kempes', got '%s'", responseEvent.Ubicacion)
	}
	if responseEvent.Precio != 1500.00 {
		t.Errorf("Expected Precio 1500.00, got %f", responseEvent.Precio)
	}

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

	responseRecorder2 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder2, req2)

	if responseRecorder2.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d. Body: %s", responseRecorder2.Code, responseRecorder2.Body.String())
	}

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

	responseRecorder3 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder3, req3)

	if responseRecorder3.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 Unprocessable Entity, got %d. Body: %s", responseRecorder3.Code, responseRecorder3.Body.String())
	}

	req4, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body))
	req4.Header.Set("Content-Type", "application/json")

	responseRecorder4 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder4, req4)

	if responseRecorder4.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d. Body: %s", responseRecorder4.Code, responseRecorder4.Body.String())
	}

	req5, _ := http.NewRequest("POST", "/admin/events", bytes.NewBuffer(body))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("Authorization", "Bearer "+clientToken)

	responseRecorder5 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder5, req5)

	if responseRecorder5.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d. Body: %s", responseRecorder5.Code, responseRecorder5.Body.String())
	}

	theaterEvent := &domain.Event{
		Titulo:      "Obra de Teatro",
		Descripcion: "Comedia dramatica en tres actos",
		Categoria:   "Teatro",
		Fecha:       futureDate,
		HoraInicio:  "20:00",
		HoraFin:     "22:00",
		Ubicacion:   "Teatro Libertador",
		Capacidad:   200,
		Precio:      800.00,
	}
	_ = mockDAO.Create(theaterEvent)

	req6, _ := http.NewRequest("GET", "/events", nil)
	responseRecorder6 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder6, req6)

	if responseRecorder6.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", responseRecorder6.Code, responseRecorder6.Body.String())
	}

	var eventsAll []domain.EventResponseDTO
	_ = json.Unmarshal(responseRecorder6.Body.Bytes(), &eventsAll)
	if len(eventsAll) != 2 {
		t.Errorf("Expected 2 events, got %d", len(eventsAll))
	}

	mockDAO.events[0].Tickets = []domain.Ticket{
		{Estado: "activo"},
		{Estado: "activo"},
		{Estado: "cancelado"},
	}

	reqStats, _ := http.NewRequest("GET", "/admin/dashboard", nil)
	reqStats.Header.Set("Authorization", "Bearer "+adminToken)

	responseRecorderStats := httptest.NewRecorder()
	router.ServeHTTP(responseRecorderStats, reqStats)

	if responseRecorderStats.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", responseRecorderStats.Code, responseRecorderStats.Body.String())
	}

	var statsResponse domain.AdminDashboardStatsDTO
	_ = json.Unmarshal(responseRecorderStats.Body.Bytes(), &statsResponse)
	if statsResponse.TotalEventos != 2 {
		t.Errorf("Expected TotalEventos 2, got %d", statsResponse.TotalEventos)
	}
	if statsResponse.EntradasVendidas != 2 {
		t.Errorf("Expected EntradasVendidas 2, got %d", statsResponse.EntradasVendidas)
	}
	if statsResponse.RecaudacionTotal != 3000 {
		t.Errorf("Expected RecaudacionTotal 3000, got %f", statsResponse.RecaudacionTotal)
	}
	if statsResponse.OcupacionMedia <= 0 {
		t.Errorf("Expected OcupacionMedia greater than 0, got %f", statsResponse.OcupacionMedia)
	}

	req7, _ := http.NewRequest("GET", "/events?categoria=Teatro", nil)
	responseRecorder7 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder7, req7)

	if responseRecorder7.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", responseRecorder7.Code)
	}

	var eventsTeatro []domain.EventResponseDTO
	_ = json.Unmarshal(responseRecorder7.Body.Bytes(), &eventsTeatro)
	if len(eventsTeatro) != 1 || eventsTeatro[0].Titulo != "Obra de Teatro" {
		t.Errorf("Expected 1 event ('Obra de Teatro'), got %d", len(eventsTeatro))
	}

	req8, _ := http.NewRequest("GET", "/events?buscar=Rock", nil)
	responseRecorder8 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder8, req8)

	if responseRecorder8.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", responseRecorder8.Code)
	}

	var eventsRock []domain.EventResponseDTO
	_ = json.Unmarshal(responseRecorder8.Body.Bytes(), &eventsRock)
	if len(eventsRock) != 1 || eventsRock[0].Titulo != "Concierto de Rock" {
		t.Errorf("Expected 1 event ('Concierto de Rock'), got %d", len(eventsRock))
	}

	// Test search by location (Kempes)
	reqSearchLocation, _ := http.NewRequest("GET", "/events?buscar=Kempes", nil)
	recSearchLocation := httptest.NewRecorder()
	router.ServeHTTP(recSearchLocation, reqSearchLocation)
	if recSearchLocation.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", recSearchLocation.Code)
	}
	var eventsLocation []domain.EventResponseDTO
	_ = json.Unmarshal(recSearchLocation.Body.Bytes(), &eventsLocation)
	if len(eventsLocation) != 1 || eventsLocation[0].Titulo != "Concierto de Rock" {
		t.Errorf("Expected 1 event ('Concierto de Rock') when searching by location, got %d", len(eventsLocation))
	}

	// Test search by category (Teatro)
	reqSearchCategory, _ := http.NewRequest("GET", "/events?buscar=Teatro", nil)
	recSearchCategory := httptest.NewRecorder()
	router.ServeHTTP(recSearchCategory, reqSearchCategory)
	if recSearchCategory.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", recSearchCategory.Code)
	}
	var eventsCategory []domain.EventResponseDTO
	_ = json.Unmarshal(recSearchCategory.Body.Bytes(), &eventsCategory)
	if len(eventsCategory) != 1 || eventsCategory[0].Titulo != "Obra de Teatro" {
		t.Errorf("Expected 1 event ('Obra de Teatro') when searching by category, got %d", len(eventsCategory))
	}

	req9, _ := http.NewRequest("GET", "/events?categoria=Musica&buscar=Obra", nil)
	responseRecorder9 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder9, req9)

	if responseRecorder9.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", responseRecorder9.Code)
	}

	var eventsCombined []domain.EventResponseDTO
	_ = json.Unmarshal(responseRecorder9.Body.Bytes(), &eventsCombined)
	if len(eventsCombined) != 0 {
		t.Errorf("Expected 0 events, got %d", len(eventsCombined))
	}

	req10, _ := http.NewRequest("GET", "/events/1", nil)
	responseRecorder10 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder10, req10)

	if responseRecorder10.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", responseRecorder10.Code, responseRecorder10.Body.String())
	}

	var eventDetail domain.EventResponseDTO
	_ = json.Unmarshal(responseRecorder10.Body.Bytes(), &eventDetail)
	if eventDetail.ID != 1 || eventDetail.Titulo != "Concierto de Rock" {
		t.Errorf("Expected event with ID 1 and title 'Concierto de Rock', got ID %d and title '%s'", eventDetail.ID, eventDetail.Titulo)
	}

	req11, _ := http.NewRequest("GET", "/events/999", nil)
	responseRecorder11 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder11, req11)

	if responseRecorder11.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", responseRecorder11.Code)
	}

	req12, _ := http.NewRequest("GET", "/events/abc", nil)
	responseRecorder12 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder12, req12)

	if responseRecorder12.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", responseRecorder12.Code)
	}

	futureDateStr := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	updateDTO := domain.EventCreateDTO{
		Titulo:      "Concierto de Rock Recargado",
		Descripcion: "Nueva descripcion para el concierto",
		Categoria:   "Musica",
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

	responseRecorder14 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder14, req14)

	if responseRecorder14.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", responseRecorder14.Code, responseRecorder14.Body.String())
	}

	if mockDAO.events[0].Titulo != "Concierto de Rock Recargado" || mockDAO.events[0].Capacidad != 200 {
		t.Errorf("Expected title 'Concierto de Rock Recargado' and capacity 200, got '%s' and %d", mockDAO.events[0].Titulo, mockDAO.events[0].Capacidad)
	}

	req15, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyUpdate))
	req15.Header.Set("Authorization", "Bearer "+clientToken)
	req15.Header.Set("Content-Type", "application/json")

	responseRecorder15 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder15, req15)

	if responseRecorder15.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", responseRecorder15.Code)
	}

	req16, _ := http.NewRequest("PUT", "/admin/events/999", bytes.NewBuffer(bodyUpdate))
	req16.Header.Set("Authorization", "Bearer "+adminToken)
	req16.Header.Set("Content-Type", "application/json")

	responseRecorder16 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder16, req16)

	if responseRecorder16.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", responseRecorder16.Code)
	}

	pastDateStr := time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	updatePastDTO := updateDTO
	updatePastDTO.Fecha = pastDateStr
	bodyPastUpdate, _ := json.Marshal(updatePastDTO)
	req17, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyPastUpdate))
	req17.Header.Set("Authorization", "Bearer "+adminToken)
	req17.Header.Set("Content-Type", "application/json")

	responseRecorder17 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder17, req17)

	if responseRecorder17.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 Unprocessable Entity, got %d. Body: %s", responseRecorder17.Code, responseRecorder17.Body.String())
	}

	mockDAO.events[0].Tickets = []domain.Ticket{
		{Estado: "activo"},
		{Estado: "activo"},
	}
	updateLowCapDTO := updateDTO
	updateLowCapDTO.Capacidad = 1
	bodyLowCapUpdate, _ := json.Marshal(updateLowCapDTO)
	req18, _ := http.NewRequest("PUT", "/admin/events/1", bytes.NewBuffer(bodyLowCapUpdate))
	req18.Header.Set("Authorization", "Bearer "+adminToken)
	req18.Header.Set("Content-Type", "application/json")

	responseRecorder18 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder18, req18)

	if responseRecorder18.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected 422 Unprocessable Entity, got %d. Body: %s", responseRecorder18.Code, responseRecorder18.Body.String())
	}

	mockDAO.events[0].Tickets = nil

	req19, _ := http.NewRequest("DELETE", "/admin/events/1", nil)
	req19.Header.Set("Authorization", "Bearer "+clientToken)

	responseRecorder19 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder19, req19)

	if responseRecorder19.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", responseRecorder19.Code)
	}

	req20, _ := http.NewRequest("DELETE", "/admin/events/999", nil)
	req20.Header.Set("Authorization", "Bearer "+adminToken)

	responseRecorder20 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder20, req20)

	if responseRecorder20.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d. Body: %s", responseRecorder20.Code, responseRecorder20.Body.String())
	}

	req21, _ := http.NewRequest("DELETE", "/admin/events/abc", nil)
	req21.Header.Set("Authorization", "Bearer "+adminToken)

	responseRecorder21 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder21, req21)

	if responseRecorder21.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", responseRecorder21.Code)
	}

	eventCountBefore := len(mockDAO.events)
	req22, _ := http.NewRequest("DELETE", "/admin/events/1", nil)
	req22.Header.Set("Authorization", "Bearer "+adminToken)

	responseRecorder22 := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder22, req22)

	if responseRecorder22.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d. Body: %s", responseRecorder22.Code, responseRecorder22.Body.String())
	}

	if len(mockDAO.events) != eventCountBefore-1 {
		t.Errorf("Expected %d events after delete, got %d", eventCountBefore-1, len(mockDAO.events))
	}
}
