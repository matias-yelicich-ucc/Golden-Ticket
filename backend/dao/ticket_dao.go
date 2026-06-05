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
	GetByUserID(userID uint) ([]domain.Ticket, error)
	TransferTicket(userID uint, ticketID uint, destinationDNI string) error
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

func (d *ticketDAOImpl) GetByUserID(userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	err := DB.Preload("Event").Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

func (d *ticketDAOImpl) TransferTicket(userID uint, ticketID uint, destinationDNI string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 1. Fetch ticket
		var ticket domain.Ticket
		if err := tx.First(&ticket, ticketID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("entrada no encontrada")
			}
			return err
		}

		// 2. Verify owner
		if ticket.UserID != userID {
			return errors.New("no eres el propietario de esta entrada")
		}

		// 3. Verify status
		if ticket.Estado != "activo" {
			return errors.New("no se puede transferir una entrada cancelada")
		}

		// 4. Fetch destination user by DNI
		var destUser domain.User
		if err := tx.Where("dni = ?", destinationDNI).First(&destUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("no existe ningún usuario registrado con el DNI ingresado")
			}
			return err
		}

		// 5. Verify it's not transferring to oneself
		if destUser.ID == userID {
			return errors.New("no podés transferirte una entrada a vos mismo")
		}

		// 6. Update ticket owner
		ticket.UserID = destUser.ID
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}

		return nil
	})
}
