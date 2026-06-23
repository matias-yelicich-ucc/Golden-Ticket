package services

import (
	"errors"
	"os"
	"testing"

	"golden-ticket/backend/domain"
	"golden-ticket/backend/utils"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

type authServiceMockUserDAO struct {
	usersByEmail map[string]*domain.User
	usersByDNI   map[string]*domain.User
	createErr    error
	createdUser  *domain.User
}

func (m *authServiceMockUserDAO) Create(user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = 99
	m.createdUser = user
	if m.usersByEmail == nil {
		m.usersByEmail = map[string]*domain.User{}
	}
	if m.usersByDNI == nil {
		m.usersByDNI = map[string]*domain.User{}
	}
	m.usersByEmail[user.Email] = user
	m.usersByDNI[user.DNI] = user
	return nil
}

func (m *authServiceMockUserDAO) GetByEmail(email string) (*domain.User, error) {
	if user, ok := m.usersByEmail[email]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *authServiceMockUserDAO) GetByDNI(dni string) (*domain.User, error) {
	if user, ok := m.usersByDNI[dni]; ok {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func TestAuthServiceRegisterSuccessDefaultsRoleAndNormalizesFields(t *testing.T) {
	dao := &authServiceMockUserDAO{
		usersByEmail: map[string]*domain.User{},
		usersByDNI:   map[string]*domain.User{},
	}
	service := NewAuthService(dao)

	response, err := service.Register(domain.UserRegisterDTO{
		Nombre:   "Mati",
		Apellido: "Yelicich",
		Email:    "  MATI@UCC.EDU.AR ",
		Password: "123456",
		DNI:      " 44555666 ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Email != "mati@ucc.edu.ar" {
		t.Fatalf("expected normalized email, got %q", response.Email)
	}
	if response.Rol != "cliente" {
		t.Fatalf("expected default role cliente, got %q", response.Rol)
	}
	if response.DNI != "44555666" {
		t.Fatalf("expected trimmed DNI, got %q", response.DNI)
	}
	if dao.createdUser == nil || dao.createdUser.Password == "123456" {
		t.Fatalf("expected password to be hashed before create")
	}
}

func TestAuthServiceRegisterRejectsDuplicateEmail(t *testing.T) {
	dao := &authServiceMockUserDAO{
		usersByEmail: map[string]*domain.User{
			"mati@ucc.edu.ar": {Email: "mati@ucc.edu.ar"},
		},
		usersByDNI: map[string]*domain.User{},
	}
	service := NewAuthService(dao)

	_, err := service.Register(domain.UserRegisterDTO{
		Nombre:   "Mati",
		Apellido: "Yelicich",
		Email:    "mati@ucc.edu.ar",
		Password: "123456",
		DNI:      "44555666",
	})

	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthServiceRegisterRejectsDuplicateDNI(t *testing.T) {
	dao := &authServiceMockUserDAO{
		usersByEmail: map[string]*domain.User{},
		usersByDNI: map[string]*domain.User{
			"44555666": {DNI: "44555666"},
		},
	}
	service := NewAuthService(dao)

	_, err := service.Register(domain.UserRegisterDTO{
		Nombre:   "Mati",
		Apellido: "Yelicich",
		Email:    "mati@ucc.edu.ar",
		Password: "123456",
		DNI:      "44555666",
	})

	if !errors.Is(err, ErrDNIAlreadyExists) {
		t.Fatalf("expected ErrDNIAlreadyExists, got %v", err)
	}
}

func TestAuthServiceRegisterMapsMySQLDuplicateErrors(t *testing.T) {
	tests := []struct {
		name      string
		createErr error
		expected  error
	}{
		{
			name:      "duplicate dni",
			createErr: &mysqlDriver.MySQLError{Number: 1062, Message: "Duplicate entry '44555666' for key 'users.dni'"},
			expected:  ErrDNIAlreadyExists,
		},
		{
			name:      "duplicate email",
			createErr: &mysqlDriver.MySQLError{Number: 1062, Message: "Duplicate entry 'mati@ucc.edu.ar' for key 'users.email'"},
			expected:  ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := &authServiceMockUserDAO{
				usersByEmail: map[string]*domain.User{},
				usersByDNI:   map[string]*domain.User{},
				createErr:    tt.createErr,
			}
			service := NewAuthService(dao)

			_, err := service.Register(domain.UserRegisterDTO{
				Nombre:   "Mati",
				Apellido: "Yelicich",
				Email:    "mati@ucc.edu.ar",
				Password: "123456",
				DNI:      "44555666",
			})

			if !errors.Is(err, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, err)
			}
		})
	}
}

func TestAuthServiceRegisterReturnsUnknownCreateError(t *testing.T) {
	expectedErr := errors.New("db down")
	dao := &authServiceMockUserDAO{
		usersByEmail: map[string]*domain.User{},
		usersByDNI:   map[string]*domain.User{},
		createErr:    expectedErr,
	}
	service := NewAuthService(dao)

	_, err := service.Register(domain.UserRegisterDTO{
		Nombre:   "Mati",
		Apellido: "Yelicich",
		Email:    "mati@ucc.edu.ar",
		Password: "123456",
		DNI:      "44555666",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestAuthServiceLoginScenarios(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		t.Fatalf("failed to hash password for test: %v", err)
	}

	service := NewAuthService(&authServiceMockUserDAO{
		usersByEmail: map[string]*domain.User{
			"mati@ucc.edu.ar": {
				ID:       5,
				Nombre:   "Mati",
				Apellido: "Yelicich",
				Email:    "mati@ucc.edu.ar",
				Password: hashedPassword,
				Rol:      "admin",
				DNI:      "44555666",
			},
		},
		usersByDNI: map[string]*domain.User{},
	})

	loginResponse, err := service.Login(domain.UserLoginDTO{
		Email:    "  MATI@UCC.EDU.AR ",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("expected successful login, got %v", err)
	}
	if loginResponse.Token == "" {
		t.Fatalf("expected token to be generated")
	}
	if loginResponse.User.Email != "mati@ucc.edu.ar" {
		t.Fatalf("expected normalized email in response, got %q", loginResponse.User.Email)
	}

	_, err = service.Login(domain.UserLoginDTO{
		Email:    "mati@ucc.edu.ar",
		Password: "wrong",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for bad password, got %v", err)
	}

	_, err = NewAuthService(&authServiceMockUserDAO{
		usersByEmail: map[string]*domain.User{},
		usersByDNI:   map[string]*domain.User{},
	}).Login(domain.UserLoginDTO{
		Email:    "missing@ucc.edu.ar",
		Password: "123456",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for missing user, got %v", err)
	}
}
