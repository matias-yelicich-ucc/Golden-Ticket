package services

import (
	"errors"
	"strings"

	"golden-ticket/backend/dao"
	"golden-ticket/backend/domain"
	"golden-ticket/backend/utils"
)

var (
	// ErrUserAlreadyExists is returned when trying to register an email that is already taken
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrInvalidCredentials is returned on login failure
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService defines the business logic for authentication
type AuthService interface {
	Register(dto domain.UserRegisterDTO) (*domain.UserResponseDTO, error)
	Login(dto domain.UserLoginDTO) (*domain.LoginResponseDTO, error)
}

type authServiceImpl struct {
	userDAO dao.UserDAO
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userDAO dao.UserDAO) AuthService {
	return &authServiceImpl{
		userDAO: userDAO,
	}
}

// Register registers a new user if the email is not taken
func (s *authServiceImpl) Register(dto domain.UserRegisterDTO) (*domain.UserResponseDTO, error) {
	emailNormalized := strings.ToLower(strings.TrimSpace(dto.Email))

	// Check if user already exists
	existingUser, _ := s.userDAO.GetByEmail(emailNormalized)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	role := dto.Rol
	if role == "" {
		role = "cliente"
	}

	user := domain.User{
		Nombre:   dto.Nombre,
		Apellido: dto.Apellido,
		Email:    emailNormalized,
		Password: hashedPassword,
		Rol:      role,
	}

	if err := s.userDAO.Create(&user); err != nil {
		return nil, err
	}

	response := domain.UserResponseDTO{
		ID:       user.ID,
		Nombre:   user.Nombre,
		Apellido: user.Apellido,
		Email:    user.Email,
		Rol:      user.Rol,
	}

	return &response, nil
}

// Login verifies credentials and returns user details with a JWT
func (s *authServiceImpl) Login(dto domain.UserLoginDTO) (*domain.LoginResponseDTO, error) {
	emailNormalized := strings.ToLower(strings.TrimSpace(dto.Email))

	user, err := s.userDAO.GetByEmail(emailNormalized)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(dto.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(user.ID, user.Rol)
	if err != nil {
		return nil, err
	}

	response := domain.LoginResponseDTO{
		User: domain.UserResponseDTO{
			ID:       user.ID,
			Nombre:   user.Nombre,
			Apellido: user.Apellido,
			Email:    user.Email,
			Rol:      user.Rol,
		},
		Token: token,
	}

	return &response, nil
}
