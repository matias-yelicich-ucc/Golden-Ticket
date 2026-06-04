package controllers

import (
	"net/http"

	"golden-ticket/backend/domain"
	"golden-ticket/backend/services"

	"github.com/gin-gonic/gin"
)

// EventController maneja las peticiones HTTP relativas a los eventos
type EventController struct {
	eventService services.EventService
}

// NewEventController crea una nueva instancia de EventController
func NewEventController(eventService services.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// Create maneja el endpoint de creación de un evento (POST /admin/events)
func (ctrl *EventController) Create(c *gin.Context) {
	var dto domain.EventCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctrl.eventService.CreateEvent(dto)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
