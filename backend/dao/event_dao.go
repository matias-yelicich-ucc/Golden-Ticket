package dao

import (
	"golden-ticket/backend/domain"

	"gorm.io/gorm"
)

// EventDAO define las operaciones de acceso a datos para los eventos
type EventDAO interface {
	Create(event *domain.Event) error
	GetAll(categoria string, buscar string) ([]*domain.Event, error)
	GetByID(id uint) (*domain.Event, error)
	GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error)
	Update(event *domain.Event) error
	Delete(id uint) error
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

// GetAdminDashboardStats obtiene las metricas agregadas del panel de administracion
func (d *eventDAOImpl) GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error) {
	var totalEventos int64
	if err := DB.Model(&domain.Event{}).Count(&totalEventos).Error; err != nil {
		return nil, err
	}

	var entradasVendidas int64
	if err := DB.Model(&domain.Ticket{}).
		Where("estado = ? AND event_id IS NOT NULL", "activo").
		Count(&entradasVendidas).Error; err != nil {
		return nil, err
	}

	var capacidadTotal int64
	if err := DB.Model(&domain.Event{}).
		Select("COALESCE(SUM(capacidad), 0)").
		Scan(&capacidadTotal).Error; err != nil {
		return nil, err
	}

	var recaudacionTotal float64
	if err := DB.Model(&domain.Ticket{}).
		Joins("JOIN events ON events.id = tickets.event_id").
		Where("tickets.estado = ? AND tickets.event_id IS NOT NULL", "activo").
		Select("COALESCE(SUM(events.precio), 0)").
		Scan(&recaudacionTotal).Error; err != nil {
		return nil, err
	}

	ocupacionMedia := 0.0
	if capacidadTotal > 0 {
		ocupacionMedia = (float64(entradasVendidas) / float64(capacidadTotal)) * 100
	}

	return &domain.AdminDashboardStatsDTO{
		TotalEventos:     int(totalEventos),
		EntradasVendidas: int(entradasVendidas),
		OcupacionMedia:   ocupacionMedia,
		RecaudacionTotal: recaudacionTotal,
	}, nil
}

func (d *eventDAOImpl) Update(event *domain.Event) error {
	return DB.Save(event).Error
}

// Delete cancels all active tickets for the event and then deletes the event
func (d *eventDAOImpl) Delete(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Verify the event exists
		var event domain.Event
		if err := tx.First(&event, id).Error; err != nil {
			return err
		}

		// Cancel all active tickets for this event
		if err := tx.Model(&domain.Ticket{}).
			Where("event_id = ? AND estado = ?", id, "activo").
			Update("estado", "cancelado").Error; err != nil {
			return err
		}

		// Delete the event — FK SET NULL will nullify event_id on tickets
		return tx.Delete(&domain.Event{}, id).Error
	})
}
