package services

import (
	"errors"
	"time"

	"golden-ticket/backend/dao"
	"golden-ticket/backend/domain"
)

// EventService define la lógica de negocio para la gestión de eventos
type EventService interface {
	CreateEvent(dto domain.EventCreateDTO) (*domain.EventResponseDTO, error)
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
	if dto.FechaHora.Before(time.Now()) {
		return nil, errors.New("la fecha del evento debe ser en el futuro")
	}

	event := domain.Event{
		Titulo:              dto.Titulo,
		Descripcion:         dto.Descripcion,
		Categoria:           dto.Categoria,
		FechaHora:           dto.FechaHora,
		Duracion:            dto.Duracion,
		Capacidad:           dto.Capacidad,
		EntradasDisponibles: dto.Capacidad, // El cupo inicial es igual a la capacidad total
	}

	if err := s.eventDAO.Create(&event); err != nil {
		return nil, err
	}

	response := domain.EventResponseDTO{
		ID:                  event.ID,
		Titulo:              event.Titulo,
		Descripcion:         event.Descripcion,
		Categoria:           event.Categoria,
		FechaHora:           event.FechaHora,
		Duracion:            event.Duracion,
		Capacidad:           event.Capacidad,
		EntradasDisponibles: event.EntradasDisponibles,
	}

	return &response, nil
}
