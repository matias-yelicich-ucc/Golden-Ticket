package domain

import "time"

// Event representa la entidad de un evento en el sistema
type Event struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo      string    `gorm:"type:varchar(255);not null" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	Categoria   string    `gorm:"type:varchar(100)" json:"categoria"`
	Fecha       string    `gorm:"type:varchar(50);not null" json:"fecha"`
	HoraInicio  string    `gorm:"type:varchar(50);not null" json:"hora_inicio"`
	HoraFin     string    `gorm:"type:varchar(50);not null" json:"hora_fin"`
	Ubicacion   string    `gorm:"type:varchar(255);not null" json:"ubicacion"`
	Coordenadas string    `gorm:"type:varchar(100)" json:"coordenadas"`
	UrlImagen   string    `gorm:"type:varchar(255)" json:"url_imagen"`
	Capacidad   int       `gorm:"not null" json:"capacidad"`
	Precio      float64   `gorm:"type:decimal(10,2);not null;default:0.0" json:"precio"`
	Tickets     []Ticket  `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"tickets,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventCreateDTO se utiliza para la solicitud de creación de un nuevo evento
type EventCreateDTO struct {
	Titulo      string  `json:"titulo" binding:"required"`
	Descripcion string  `json:"descripcion"`
	Categoria   string  `json:"categoria"`
	Fecha       string  `json:"fecha" binding:"required"`
	HoraInicio  string  `json:"hora_inicio" binding:"required"`
	HoraFin     string  `json:"hora_fin" binding:"required"`
	Ubicacion   string  `json:"ubicacion" binding:"required"`
	Coordenadas string  `json:"coordenadas"`
	UrlImagen   string  `json:"url_imagen"`
	Capacidad   int     `json:"capacidad" binding:"required,gt=0"`
	Precio      float64 `json:"precio" binding:"required,gte=0"`
}

// EventResponseDTO representa la respuesta segura tras crear o consultar un evento
type EventResponseDTO struct {
	ID                  uint    `json:"id"`
	Titulo              string  `json:"titulo"`
	Descripcion         string  `json:"descripcion"`
	Categoria           string  `json:"categoria"`
	Fecha               string  `json:"fecha"`
	HoraInicio          string  `json:"hora_inicio"`
	HoraFin             string  `json:"hora_fin"`
	Ubicacion           string  `json:"ubicacion"`
	Coordenadas         string  `json:"coordenadas"`
	UrlImagen           string  `json:"url_imagen"`
	Capacidad           int     `json:"capacidad"`
	EntradasDisponibles int     `json:"entradas_disponibles"`
	Precio              float64 `json:"precio"`
}
