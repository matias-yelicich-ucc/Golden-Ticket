package services

import (
	"golden-ticket/backend/dao"
	"golden-ticket/backend/domain"
)

type TicketService interface {
	BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error)
	GetTicketsByUserID(userID uint) ([]domain.Ticket, error)
	TransferTicket(userID uint, ticketID uint, destinationDNI string) error
}

type ticketServiceImpl struct {
	ticketDAO dao.TicketDAO
}

func NewTicketService(ticketDAO dao.TicketDAO) TicketService {
	return &ticketServiceImpl{
		ticketDAO: ticketDAO,
	}
}

func (s *ticketServiceImpl) BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error) {
	return s.ticketDAO.BuyTickets(userID, eventID, cantidad)
}

func (s *ticketServiceImpl) GetTicketsByUserID(userID uint) ([]domain.Ticket, error) {
	return s.ticketDAO.GetByUserID(userID)
}

func (s *ticketServiceImpl) TransferTicket(userID uint, ticketID uint, destinationDNI string) error {
	return s.ticketDAO.TransferTicket(userID, ticketID, destinationDNI)
}
