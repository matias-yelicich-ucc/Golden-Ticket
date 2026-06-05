package controllers

import (
	"net/http"
	"strconv"

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

// List maneja el endpoint para obtener todos los eventos (GET /events) con filtrado opcional
func (ctrl *EventController) List(c *gin.Context) {
	categoria := c.Query("categoria")
	buscar := c.Query("buscar")

	res, err := ctrl.eventService.GetAllEvents(categoria, buscar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetByID maneja el endpoint para obtener un evento por su ID (GET /events/:id)
func (ctrl *EventController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	res, err := ctrl.eventService.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Update maneja el endpoint de actualización de un evento (PUT /admin/events/:id)
func (ctrl *EventController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	var dto domain.EventCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := ctrl.eventService.UpdateEvent(uint(id), dto)
	if err != nil {
		status := http.StatusUnprocessableEntity
		if err.Error() == "evento no encontrado" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

