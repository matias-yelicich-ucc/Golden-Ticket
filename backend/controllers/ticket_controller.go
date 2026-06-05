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

// GetMyTickets obtiene todos los tickets del usuario logueado (GET /my-tickets)
func (ctrl *TicketController) GetMyTickets(c *gin.Context) {
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

	tickets, err := ctrl.ticketService.GetTicketsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// TicketTransferDTO represents the payload for transferring a ticket
type TicketTransferDTO struct {
	DNI string `json:"dni" binding:"required"`
}

// Transfer transfers a ticket to another user by DNI (POST /my-tickets/:id/transfer)
func (ctrl *TicketController) Transfer(c *gin.Context) {
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

	// 2. Get ticket ID from route
	ticketIDStr := c.Param("id")
	ticketID, err := strconv.ParseUint(ticketIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de entrada inválido"})
		return
	}

	// 3. Bind JSON body
	var dto TicketTransferDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DNI es requerido"})
		return
	}

	// 4. Perform transfer
	err = ctrl.ticketService.TransferTicket(userID, uint(ticketID), dto.DNI)
	if err != nil {
		status := http.StatusBadRequest
		errStr := err.Error()

		if errStr == "entrada no encontrada" || errStr == "no existe ningún usuario registrado con el DNI ingresado" {
			status = http.StatusNotFound
		} else if errStr == "no eres el propietario de esta entrada" {
			status = http.StatusForbidden
		}

		c.JSON(status, gin.H{"error": errStr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entrada transferida con éxito"})
}

