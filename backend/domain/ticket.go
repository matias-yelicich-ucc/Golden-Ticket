package domain

import "time"

// Ticket representa la entidad de una entrada/ticket emitida
type Ticket struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	EventID     *uint     `json:"event_id"`
	Event       *Event    `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"event,omitempty"`
	Estado      string    `gorm:"type:varchar(50);not null;default:'activo'" json:"estado"` // 'activo' o 'cancelado'
	FechaCompra time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"fecha_compra"`
}

// TicketPurchaseDTO se utiliza para la solicitud de compra de entradas
type TicketPurchaseDTO struct {
	Cantidad int `json:"cantidad" binding:"required,gt=0"`
}
