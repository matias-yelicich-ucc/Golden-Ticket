package domain

import "time"

// User represents the user entity in the database
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre    string    `gorm:"type:varchar(100);not null" json:"nombre"`
	Apellido  string    `gorm:"type:varchar(100);not null" json:"apellido"`
	Email     string    `gorm:"type:varchar(191);uniqueIndex;not null" json:"email"` // 191 is MySQL safe index size
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Rol       string    `gorm:"type:varchar(50);not null;default:'cliente'" json:"rol"`
	Tickets   []Ticket  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"tickets,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRegisterDTO is used for incoming registration requests
type UserRegisterDTO struct {
	Nombre   string `json:"nombre" binding:"required"`
	Apellido string `json:"apellido" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Rol      string `json:"rol" binding:"omitempty,oneof=cliente administrador"`
}

// UserLoginDTO is used for incoming login requests
type UserLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponseDTO is used for returning safe user data (without password)
type UserResponseDTO struct {
	ID       uint   `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
	Rol      string `json:"rol"`
}

// LoginResponseDTO is returned upon successful authentication
type LoginResponseDTO struct {
	User  UserResponseDTO `json:"user"`
	Token string          `json:"token"`
}
