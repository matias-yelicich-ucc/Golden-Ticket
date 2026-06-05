package dao

import (
	"errors"
	"fmt"
	"time"

	"golden-ticket/backend/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketDAO interface {
	BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error)
}

type ticketDAOImpl struct{}

func NewTicketDAO() TicketDAO {
	return &ticketDAOImpl{}
}

func (d *ticketDAOImpl) BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error) {
	var tickets []domain.Ticket

	err := DB.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the event row to prevent concurrent race conditions
		var event domain.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&event, eventID).Error; err != nil {
			return err
		}

		// 2. Count active tickets
		var activeTicketsCount int64
		if err := tx.Model(&domain.Ticket{}).Where("event_id = ? AND estado = ?", eventID, "activo").Count(&activeTicketsCount).Error; err != nil {
			return err
		}

		// 3. Validate capacity
		entradasDisponibles := event.Capacidad - int(activeTicketsCount)
		if entradasDisponibles < cantidad {
			return fmt.Errorf("capacidad insuficiente para realizar la compra. Cupo disponible: %d", entradasDisponibles)
		}

		// 4. Validate if event is in the future
		eventDateTimeStr := fmt.Sprintf("%sT%s:00", event.Fecha, event.HoraInicio)
		eventTime, err := time.ParseInLocation("2006-01-02T15:04:05", eventDateTimeStr, time.Local)
		if err == nil && eventTime.Before(time.Now()) {
			return errors.New("no se pueden comprar entradas para un evento que ya ocurrió o está en curso")
		}

		// 5. Create tickets
		for i := 0; i < cantidad; i++ {
			tickets = append(tickets, domain.Ticket{
				UserID:      userID,
				EventID:     eventID,
				Estado:      "activo",
				FechaCompra: time.Now(),
			})
		}

		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		return nil
	})

	return tickets, err
}
