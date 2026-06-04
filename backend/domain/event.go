package domain

import "time"

// Event representa la entidad de un evento en el sistema
type Event struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo              string    `gorm:"type:varchar(255);not null" json:"titulo"`
	Descripcion         string    `gorm:"type:text" json:"descripcion"`
	Categoria           string    `gorm:"type:varchar(100)" json:"categoria"`
	FechaHora           time.Time `gorm:"not null" json:"fecha_hora"`
	Duracion            int       `gorm:"not null" json:"duracion"` // en minutos
	Capacidad           int       `gorm:"not null" json:"capacidad"`
	EntradasDisponibles int       `gorm:"not null" json:"entradas_disponibles"`
	Tickets             []Ticket  `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"tickets,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// EventCreateDTO se utiliza para la solicitud de creación de un nuevo evento
type EventCreateDTO struct {
	Titulo      string    `json:"titulo" binding:"required"`
	Descripcion string    `json:"descripcion"`
	Categoria   string    `json:"categoria"`
	FechaHora   time.Time `json:"fecha_hora" binding:"required"`
	Duracion    int       `json:"duracion" binding:"required,gt=0"` // en minutos
	Capacidad   int       `json:"capacidad" binding:"required,gt=0"`
}

// EventResponseDTO representa la respuesta segura tras crear o consultar un evento
type EventResponseDTO struct {
	ID                  uint      `json:"id"`
	Titulo              string    `json:"titulo"`
	Descripcion         string    `json:"descripcion"`
	Categoria           string    `json:"categoria"`
	FechaHora           time.Time `json:"fecha_hora"`
	Duracion            int       `json:"duracion"`
	Capacidad           int       `json:"capacidad"`
	EntradasDisponibles int       `json:"entradas_disponibles"`
}
