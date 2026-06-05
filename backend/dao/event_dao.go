package dao

import (
	"golden-ticket/backend/domain"
)

// EventDAO define las operaciones de acceso a datos para los eventos
type EventDAO interface {
	Create(event *domain.Event) error
	GetAll(categoria string, buscar string) ([]*domain.Event, error)
	GetByID(id uint) (*domain.Event, error)
}

type eventDAOImpl struct{}

// NewEventDAO crea una nueva instancia de EventDAO
func NewEventDAO() EventDAO {
	return &eventDAOImpl{}
}

// GetByID obtiene un evento por su ID precalificando la relación de Tickets
func (d *eventDAOImpl) GetByID(id uint) (*domain.Event, error) {
	var event domain.Event
	err := DB.Preload("Tickets").First(&event, id).Error
	return &event, err
}

// GetAll obtiene los eventos de la base de datos aplicando filtros de categoría y búsqueda
func (d *eventDAOImpl) GetAll(categoria string, buscar string) ([]*domain.Event, error) {
	var events []*domain.Event
	query := DB.Preload("Tickets")
	if categoria != "" {
		query = query.Where("categoria = ?", categoria)
	}
	if buscar != "" {
		query = query.Where("titulo LIKE ? OR descripcion LIKE ?", "%"+buscar+"%", "%"+buscar+"%")
	}
	err := query.Find(&events).Error
	return events, err
}

// Create inserta un nuevo evento en la base de datos
func (d *eventDAOImpl) Create(event *domain.Event) error {
	return DB.Create(event).Error
}


