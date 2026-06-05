package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"golden-ticket/backend/domain"
	"golden-ticket/backend/services"

	"github.com/gin-gonic/gin"
)

type mockUserDAO struct {
	users map[string]*domain.User
}

func (m *mockUserDAO) Create(user *domain.User) error {
	if _, ok := m.users[user.Email]; ok {
		return errors.New("duplicate key error")
	}
	user.ID = uint(len(m.users) + 1)
	m.users[user.Email] = user
	return nil
}

func (m *mockUserDAO) GetByEmail(email string) (*domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *mockUserDAO) GetByDNI(dni string) (*domain.User, error) {
	for _, u := range m.users {
		if u.DNI == dni {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func TestAuthController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret_for_controller")
	defer os.Unsetenv("JWT_SECRET")

	mockDAO := &mockUserDAO{users: make(map[string]*domain.User)}
	authService := services.NewAuthService(mockDAO)
	ctrl := NewAuthController(authService)

	router := gin.Default()
	router.POST("/register", ctrl.Register)
	router.POST("/login", ctrl.Login)

	// 1. Test Register Success
	regDTO := domain.UserRegisterDTO{
		Nombre:   "Juan",
		Apellido: "Perez",
		Email:    "juan.perez@example.com",
		Password: "password123",
		Rol:      "cliente",
		DNI:      "12345678",
	}
	body, _ := json.Marshal(regDTO)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected StatusCreated (201), got %d. Body: %s", w.Code, w.Body.String())
	}

	var respUser domain.UserResponseDTO
	if err := json.Unmarshal(w.Body.Bytes(), &respUser); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if respUser.Email != "juan.perez@example.com" {
		t.Errorf("Expected email juan.perez@example.com, got %s", respUser.Email)
	}

	// 2. Test Register Duplicate Conflict (409)
	req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected StatusConflict (409), got %d. Body: %s", w2.Code, w2.Body.String())
	}

	// 3. Test Login Success
	loginDTO := domain.UserLoginDTO{
		Email:    "juan.perez@example.com",
		Password: "password123",
	}
	loginBody, _ := json.Marshal(loginDTO)
	req3, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("Expected StatusOK (200), got %d. Body: %s", w3.Code, w3.Body.String())
	}

	var loginResp domain.LoginResponseDTO
	if err := json.Unmarshal(w3.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}
	if loginResp.Token == "" {
		t.Errorf("Expected token to be returned, got empty string")
	}
	if loginResp.User.Email != "juan.perez@example.com" {
		t.Errorf("Expected email juan.perez@example.com in login response, got %s", loginResp.User.Email)
	}

	// 4. Test Login Unauthorized (401)
	loginDTOInvalid := domain.UserLoginDTO{
		Email:    "juan.perez@example.com",
		Password: "wrongpassword",
	}
	loginBodyInvalid, _ := json.Marshal(loginDTOInvalid)
	req4, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBodyInvalid))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	if w4.Code != http.StatusUnauthorized {
		t.Errorf("Expected StatusUnauthorized (401), got %d. Body: %s", w4.Code, w4.Body.String())
	}
}
