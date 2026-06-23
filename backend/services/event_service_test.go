package services

import (
	"errors"
	"testing"
	"time"

	"golden-ticket/backend/domain"
)

type eventServiceMockDAO struct {
	events      map[uint]*domain.Event
	createErr   error
	getAllErr   error
	statsErr    error
	updateErr   error
	deleteErr   error
	stats       *domain.AdminDashboardStatsDTO
	created     *domain.Event
	updated     *domain.Event
	deletedID   uint
	nextEventID uint
}

func (m *eventServiceMockDAO) Create(event *domain.Event) error {
	if m.createErr != nil {
		return m.createErr
	}
	if m.nextEventID == 0 {
		m.nextEventID = 1
	}
	event.ID = m.nextEventID
	m.nextEventID++
	m.created = event
	if m.events == nil {
		m.events = map[uint]*domain.Event{}
	}
	m.events[event.ID] = event
	return nil
}

func (m *eventServiceMockDAO) GetAll(categoria string, buscar string) ([]*domain.Event, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	response := make([]*domain.Event, 0, len(m.events))
	for _, event := range m.events {
		response = append(response, event)
	}
	return response, nil
}

func (m *eventServiceMockDAO) GetByID(id uint) (*domain.Event, error) {
	if event, ok := m.events[id]; ok {
		return event, nil
	}
	return nil, errors.New("not found")
}

func (m *eventServiceMockDAO) GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error) {
	if m.statsErr != nil {
		return nil, m.statsErr
	}
	return m.stats, nil
}

func (m *eventServiceMockDAO) Update(event *domain.Event) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updated = event
	m.events[event.ID] = event
	return nil
}

func (m *eventServiceMockDAO) Delete(id uint) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deletedID = id
	delete(m.events, id)
	return nil
}

func futureEventDTO() domain.EventCreateDTO {
	return domain.EventCreateDTO{
		Titulo:      "Evento",
		Descripcion: "Descripcion",
		Categoria:   "Musica",
		Fecha:       time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		HoraInicio:  "21:00",
		HoraFin:     "23:00",
		Ubicacion:   "Cordoba",
		Coordenadas: "-31,-64",
		UrlImagen:   "https://example.com/evento.jpg",
		Capacidad:   100,
		Precio:      2500,
	}
}

func TestEventServiceCreateEvent(t *testing.T) {
	dao := &eventServiceMockDAO{events: map[uint]*domain.Event{}}
	service := NewEventService(dao)

	response, err := service.CreateEvent(futureEventDTO())
	if err != nil {
		t.Fatalf("expected create success, got %v", err)
	}
	if response.EntradasDisponibles != response.Capacidad {
		t.Fatalf("expected entradas disponibles to start equal to capacidad")
	}

	invalidDateDTO := futureEventDTO()
	invalidDateDTO.Fecha = "bad-date"
	if _, err := service.CreateEvent(invalidDateDTO); err == nil || err.Error() != "formato de fecha u hora de inicio inválido" {
		t.Fatalf("expected invalid date error, got %v", err)
	}

	pastDTO := futureEventDTO()
	pastDTO.Fecha = time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	if _, err := service.CreateEvent(pastDTO); err == nil || err.Error() != "la fecha del evento debe ser en el futuro" {
		t.Fatalf("expected past date error, got %v", err)
	}

	dao.createErr = errors.New("db error")
	if _, err := service.CreateEvent(futureEventDTO()); err == nil || err.Error() != "db error" {
		t.Fatalf("expected DAO create error, got %v", err)
	}
}

func TestEventServiceGetAllAndGetByID(t *testing.T) {
	dao := &eventServiceMockDAO{
		events: map[uint]*domain.Event{
			1: {
				ID:         1,
				Titulo:     "Rock",
				Capacidad:  10,
				Categoria:  "Musica",
				Fecha:      "2026-01-01",
				HoraInicio: "20:00",
				HoraFin:    "22:00",
				Ubicacion:  "Cordoba",
				Precio:     1500,
				Tickets: []domain.Ticket{
					{Estado: "activo"},
					{Estado: "activo"},
					{Estado: "cancelado"},
				},
			},
		},
	}
	service := NewEventService(dao)

	all, err := service.GetAllEvents("", "")
	if err != nil {
		t.Fatalf("expected get all success, got %v", err)
	}
	if len(all) != 1 || all[0].EntradasDisponibles != 8 {
		t.Fatalf("expected 1 event with 8 available tickets, got %+v", all)
	}

	detail, err := service.GetEventByID(1)
	if err != nil {
		t.Fatalf("expected detail success, got %v", err)
	}
	if detail.EntradasDisponibles != 8 {
		t.Fatalf("expected 8 available tickets, got %d", detail.EntradasDisponibles)
	}

	dao.getAllErr = errors.New("dao list error")
	if _, err := service.GetAllEvents("", ""); err == nil || err.Error() != "dao list error" {
		t.Fatalf("expected list error, got %v", err)
	}

	if _, err := service.GetEventByID(999); err == nil {
		t.Fatalf("expected get by id error for missing event")
	}
}

func TestEventServiceGetAdminDashboardStats(t *testing.T) {
	expected := &domain.AdminDashboardStatsDTO{
		TotalEventos:     2,
		EntradasVendidas: 6,
		OcupacionMedia:   50,
		RecaudacionTotal: 12000,
	}
	dao := &eventServiceMockDAO{stats: expected}
	service := NewEventService(dao)

	stats, err := service.GetAdminDashboardStats()
	if err != nil {
		t.Fatalf("expected stats success, got %v", err)
	}
	if stats.TotalEventos != expected.TotalEventos {
		t.Fatalf("expected total eventos %d, got %d", expected.TotalEventos, stats.TotalEventos)
	}

	dao.statsErr = errors.New("stats error")
	if _, err := service.GetAdminDashboardStats(); err == nil || err.Error() != "stats error" {
		t.Fatalf("expected stats error, got %v", err)
	}
}

func TestEventServiceUpdateEvent(t *testing.T) {
	existing := &domain.Event{
		ID:         1,
		Titulo:     "Viejo",
		Capacidad:  10,
		Fecha:      time.Now().Add(24 * time.Hour).Format("2006-01-02"),
		HoraInicio: "19:00",
		HoraFin:    "20:00",
		Ubicacion:  "Cordoba",
		Tickets: []domain.Ticket{
			{Estado: "activo"},
			{Estado: "cancelado"},
			{Estado: "activo"},
		},
	}
	dao := &eventServiceMockDAO{events: map[uint]*domain.Event{1: existing}}
	service := NewEventService(dao)

	dto := futureEventDTO()
	dto.Capacidad = 15
	dto.Titulo = "Nuevo"
	response, err := service.UpdateEvent(1, dto)
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if response.EntradasDisponibles != 13 {
		t.Fatalf("expected 13 available tickets, got %d", response.EntradasDisponibles)
	}
	if dao.updated == nil || dao.updated.Titulo != "Nuevo" {
		t.Fatalf("expected updated event to be saved")
	}

	if _, err := service.UpdateEvent(999, dto); err == nil || err.Error() != "evento no encontrado" {
		t.Fatalf("expected not found error, got %v", err)
	}

	lowCapacityDTO := dto
	lowCapacityDTO.Capacidad = 1
	if _, err := service.UpdateEvent(1, lowCapacityDTO); err == nil {
		t.Fatalf("expected capacity validation error")
	}

	badDateDTO := dto
	badDateDTO.Fecha = "invalid"
	if _, err := service.UpdateEvent(1, badDateDTO); err == nil || err.Error() != "formato de fecha u hora de inicio inválido" {
		t.Fatalf("expected invalid date error, got %v", err)
	}

	pastDateDTO := dto
	pastDateDTO.Fecha = time.Now().Add(-48 * time.Hour).Format("2006-01-02")
	if _, err := service.UpdateEvent(1, pastDateDTO); err == nil || err.Error() != "la fecha del evento debe ser en el futuro" {
		t.Fatalf("expected past date error, got %v", err)
	}

	dao.updateErr = errors.New("update failed")
	if _, err := service.UpdateEvent(1, dto); err == nil || err.Error() != "update failed" {
		t.Fatalf("expected dao update error, got %v", err)
	}
}

func TestEventServiceDeleteEvent(t *testing.T) {
	dao := &eventServiceMockDAO{
		events: map[uint]*domain.Event{
			1: {ID: 1, Titulo: "Evento"},
		},
	}
	service := NewEventService(dao)

	if err := service.DeleteEvent(1); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
	if dao.deletedID != 1 {
		t.Fatalf("expected delete to target id 1, got %d", dao.deletedID)
	}

	if err := service.DeleteEvent(999); err == nil || err.Error() != "evento no encontrado" {
		t.Fatalf("expected not found error, got %v", err)
	}

	dao.events[2] = &domain.Event{ID: 2}
	dao.deleteErr = errors.New("delete failed")
	if err := service.DeleteEvent(2); err == nil || err.Error() != "delete failed" {
		t.Fatalf("expected dao delete error, got %v", err)
	}
}
