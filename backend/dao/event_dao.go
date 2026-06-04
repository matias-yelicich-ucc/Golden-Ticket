package dao

import (
	"golden-ticket/backend/domain"
)

// EventDAO define las operaciones de acceso a datos para los eventos
type EventDAO interface {
	Create(event *domain.Event) error
	GetAll() ([]*domain.Event, error)
}

type eventDAOImpl struct{}

// NewEventDAO crea una nueva instancia de EventDAO
func NewEventDAO() EventDAO {
	return &eventDAOImpl{}
}

// GetAll obtiene todos los eventos de la base de datos precalificando la relación de Tickets
func (d *eventDAOImpl) GetAll() ([]*domain.Event, error) {
	var events []*domain.Event
	err := DB.Preload("Tickets").Find(&events).Error
	return events, err
}

// Create inserta un nuevo evento en la base de datos
func (d *eventDAOImpl) Create(event *domain.Event) error {
	return DB.Create(event).Error
}

// GetAll obtiene todos los eventos de la base de datos
func (d *eventDAOImpl) GetAll() ([]*domain.Event, error) {
	var events []*domain.Event
	err := DB.Find(&events).Error
	return events, err
}
