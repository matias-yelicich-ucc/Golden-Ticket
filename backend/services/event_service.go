package services

import (
	"errors"
	"fmt"
	"time"

	"golden-ticket/backend/dao"
	"golden-ticket/backend/domain"
)

// EventService define la lógica de negocio para la gestión de eventos
type EventService interface {
	CreateEvent(dto domain.EventCreateDTO) (*domain.EventResponseDTO, error)
	GetAllEvents(categoria string, buscar string) ([]*domain.EventResponseDTO, error)
	GetEventByID(id uint) (*domain.EventResponseDTO, error)
	GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error)
	UpdateEvent(id uint, dto domain.EventCreateDTO) (*domain.EventResponseDTO, error)
	DeleteEvent(id uint) error
}

type eventServiceImpl struct {
	eventDAO dao.EventDAO
}

// NewEventService crea una nueva instancia de EventService
func NewEventService(eventDAO dao.EventDAO) EventService {
	return &eventServiceImpl{
		eventDAO: eventDAO,
	}
}

// CreateEvent registra un nuevo evento en el catálogo realizando las validaciones pertinentes
func (s *eventServiceImpl) CreateEvent(dto domain.EventCreateDTO) (*domain.EventResponseDTO, error) {
	// Validación: La fecha y hora del evento debe ser en el futuro
	eventDateTimeStr := fmt.Sprintf("%sT%s:00", dto.Fecha, dto.HoraInicio)
	eventTime, err := time.ParseInLocation("2006-01-02T15:04:05", eventDateTimeStr, time.Local)
	if err != nil {
		return nil, errors.New("formato de fecha u hora de inicio inválido")
	}

	if eventTime.Before(time.Now()) {
		return nil, errors.New("la fecha del evento debe ser en el futuro")
	}

	event := domain.Event{
		Titulo:      dto.Titulo,
		Descripcion: dto.Descripcion,
		Categoria:   dto.Categoria,
		Fecha:       dto.Fecha,
		HoraInicio:  dto.HoraInicio,
		HoraFin:     dto.HoraFin,
		Ubicacion:   dto.Ubicacion,
		Coordenadas: dto.Coordenadas,
		UrlImagen:   dto.UrlImagen,
		Capacidad:   dto.Capacidad,
		Precio:      dto.Precio,
	}

	if err := s.eventDAO.Create(&event); err != nil {
		return nil, err
	}

	response := domain.EventResponseDTO{
		ID:                  event.ID,
		Titulo:              event.Titulo,
		Descripcion:         event.Descripcion,
		Categoria:           event.Categoria,
		Fecha:               event.Fecha,
		HoraInicio:          event.HoraInicio,
		HoraFin:             event.HoraFin,
		Ubicacion:           event.Ubicacion,
		Coordenadas:         event.Coordenadas,
		UrlImagen:           event.UrlImagen,
		Capacidad:           event.Capacidad,
		EntradasDisponibles: event.Capacidad, // Al crearse, no hay entradas vendidas aún
		Precio:              event.Precio,
	}

	return &response, nil
}

// GetAllEvents obtiene todos los eventos filtrados y los mapea a EventResponseDTO
func (s *eventServiceImpl) GetAllEvents(categoria string, buscar string) ([]*domain.EventResponseDTO, error) {
	events, err := s.eventDAO.GetAll(categoria, buscar)
	if err != nil {
		return nil, err
	}

	response := make([]*domain.EventResponseDTO, 0)
	for _, event := range events {
		// Calcular entradas vendidas dinámicamente
		entradasVendidas := 0
		for _, ticket := range event.Tickets {
			if ticket.Estado == "activo" {
				entradasVendidas++
			}
		}
		entradasDisponibles := event.Capacidad - entradasVendidas

		response = append(response, &domain.EventResponseDTO{
			ID:                  event.ID,
			Titulo:              event.Titulo,
			Descripcion:         event.Descripcion,
			Categoria:           event.Categoria,
			Fecha:               event.Fecha,
			HoraInicio:          event.HoraInicio,
			HoraFin:             event.HoraFin,
			Ubicacion:           event.Ubicacion,
			Coordenadas:         event.Coordenadas,
			UrlImagen:           event.UrlImagen,
			Capacidad:           event.Capacidad,
			EntradasDisponibles: entradasDisponibles,
			Precio:              event.Precio,
		})
	}

	return response, nil
}

// GetEventByID obtiene un evento por su ID y calcula la disponibilidad de entradas
func (s *eventServiceImpl) GetEventByID(id uint) (*domain.EventResponseDTO, error) {
	event, err := s.eventDAO.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Calcular entradas vendidas dinámicamente
	entradasVendidas := 0
	for _, ticket := range event.Tickets {
		if ticket.Estado == "activo" {
			entradasVendidas++
		}
	}
	entradasDisponibles := event.Capacidad - entradasVendidas

	response := domain.EventResponseDTO{
		ID:                  event.ID,
		Titulo:              event.Titulo,
		Descripcion:         event.Descripcion,
		Categoria:           event.Categoria,
		Fecha:               event.Fecha,
		HoraInicio:          event.HoraInicio,
		HoraFin:             event.HoraFin,
		Ubicacion:           event.Ubicacion,
		Coordenadas:         event.Coordenadas,
		UrlImagen:           event.UrlImagen,
		Capacidad:           event.Capacidad,
		EntradasDisponibles: entradasDisponibles,
		Precio:              event.Precio,
	}

	return &response, nil
}

// GetAdminDashboardStats obtiene las metricas agregadas del panel de administracion
func (s *eventServiceImpl) GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error) {
	return s.eventDAO.GetAdminDashboardStats()
}

// UpdateEvent actualiza un evento existente validando capacidad y fecha futura
func (s *eventServiceImpl) UpdateEvent(id uint, dto domain.EventCreateDTO) (*domain.EventResponseDTO, error) {
	// 1. Fetch event from DAO
	event, err := s.eventDAO.GetByID(id)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}

	// 2. Count active tickets
	activeTicketsCount := 0
	for _, ticket := range event.Tickets {
		if ticket.Estado == "activo" {
			activeTicketsCount++
		}
	}

	// 3. Validate capacity limit
	if dto.Capacidad < activeTicketsCount {
		return nil, fmt.Errorf("la nueva capacidad (%d) no puede ser menor a las entradas ya vendidas (%d)", dto.Capacidad, activeTicketsCount)
	}

	// 4. Validate future date
	eventDateTimeStr := fmt.Sprintf("%sT%s:00", dto.Fecha, dto.HoraInicio)
	eventTime, err := time.ParseInLocation("2006-01-02T15:04:05", eventDateTimeStr, time.Local)
	if err != nil {
		return nil, errors.New("formato de fecha u hora de inicio inválido")
	}

	if eventTime.Before(time.Now()) {
		return nil, errors.New("la fecha del evento debe ser en el futuro")
	}

	// 5. Update fields
	event.Titulo = dto.Titulo
	event.Descripcion = dto.Descripcion
	event.Categoria = dto.Categoria
	event.Fecha = dto.Fecha
	event.HoraInicio = dto.HoraInicio
	event.HoraFin = dto.HoraFin
	event.Ubicacion = dto.Ubicacion
	event.Coordenadas = dto.Coordenadas
	event.UrlImagen = dto.UrlImagen
	event.Capacidad = dto.Capacidad
	event.Precio = dto.Precio

	// 6. Save in DAO
	if err := s.eventDAO.Update(event); err != nil {
		return nil, err
	}

	// 7. Map to DTO
	entradasDisponibles := event.Capacidad - activeTicketsCount
	response := domain.EventResponseDTO{
		ID:                  event.ID,
		Titulo:              event.Titulo,
		Descripcion:         event.Descripcion,
		Categoria:           event.Categoria,
		Fecha:               event.Fecha,
		HoraInicio:          event.HoraInicio,
		HoraFin:             event.HoraFin,
		Ubicacion:           event.Ubicacion,
		Coordenadas:         event.Coordenadas,
		UrlImagen:           event.UrlImagen,
		Capacidad:           event.Capacidad,
		EntradasDisponibles: entradasDisponibles,
		Precio:              event.Precio,
	}

	return &response, nil
}

// DeleteEvent verifies the event exists and delegates its deletion to the DAO
func (s *eventServiceImpl) DeleteEvent(id uint) error {
	_, err := s.eventDAO.GetByID(id)
	if err != nil {
		return errors.New("evento no encontrado")
	}

	return s.eventDAO.Delete(id)
}
