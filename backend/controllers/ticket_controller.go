package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"golden-ticket/backend/domain"
	"golden-ticket/backend/services"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService services.TicketService
}

func NewTicketController(ticketService services.TicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

func (ctrl *TicketController) Buy(c *gin.Context) {
	// 1. Get user ID from Auth context
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID de usuario inválido en el contexto"})
		return
	}

	// 2. Get event ID from route
	eventIDStr := c.Param("id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	// 3. Bind request body
	var dto domain.TicketPurchaseDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Buy tickets
	tickets, err := ctrl.ticketService.BuyTickets(userID, uint(eventID), dto.Cantidad)
	if err != nil {
		status := http.StatusBadRequest
		errStr := err.Error()

		if strings.Contains(errStr, "not found") {
			status = http.StatusNotFound
		} else if strings.Contains(errStr, "capacidad insuficiente") {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{"error": errStr})
		return
	}

	c.JSON(http.StatusCreated, tickets)
}
