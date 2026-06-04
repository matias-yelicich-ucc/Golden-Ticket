package dao

import (
	"golden-ticket/backend/domain"
)

// EventDAO define las operaciones de acceso a datos para los eventos
type EventDAO interface {
	Create(event *domain.Event) error
}

type eventDAOImpl struct{}

// NewEventDAO crea una nueva instancia de EventDAO
func NewEventDAO() EventDAO {
	return &eventDAOImpl{}
}

// Create inserta un nuevo evento en la base de datos
func (d *eventDAOImpl) Create(event *domain.Event) error {
	return DB.Create(event).Error
}
