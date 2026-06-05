package services

import (
	"golden-ticket/backend/dao"
	"golden-ticket/backend/domain"
)

type TicketService interface {
	BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error)
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
