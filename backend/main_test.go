package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golden-ticket/backend/controllers"
	"golden-ticket/backend/domain"
	"golden-ticket/backend/services"

	"github.com/gin-gonic/gin"
)

type testMainUserDAO struct{}

func (testMainUserDAO) Create(user *domain.User) error { return nil }
func (testMainUserDAO) GetByEmail(email string) (*domain.User, error) {
	return nil, services.ErrInvalidCredentials
}
func (testMainUserDAO) GetByDNI(dni string) (*domain.User, error) {
	return nil, services.ErrInvalidCredentials
}

type testMainEventDAO struct{}

func (testMainEventDAO) Create(event *domain.Event) error { return nil }
func (testMainEventDAO) GetAll(categoria string, buscar string) ([]*domain.Event, error) {
	return []*domain.Event{}, nil
}
func (testMainEventDAO) GetByID(id uint) (*domain.Event, error) {
	return &domain.Event{ID: id, Titulo: "Evento", Capacidad: 10}, nil
}
func (testMainEventDAO) GetAdminDashboardStats() (*domain.AdminDashboardStatsDTO, error) {
	return &domain.AdminDashboardStatsDTO{}, nil
}
func (testMainEventDAO) Update(event *domain.Event) error { return nil }
func (testMainEventDAO) Delete(id uint) error             { return nil }

type testMainTicketDAO struct{}

func (testMainTicketDAO) BuyTickets(userID uint, eventID uint, cantidad int) ([]domain.Ticket, error) {
	return []domain.Ticket{}, nil
}
func (testMainTicketDAO) GetByUserID(userID uint) ([]domain.Ticket, error) {
	return []domain.Ticket{}, nil
}
func (testMainTicketDAO) TransferTicket(userID uint, ticketID uint, destinationDNI string) error {
	return nil
}
func (testMainTicketDAO) CancelTicket(userID uint, ticketID uint) error { return nil }

func buildTestAppRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return buildRouter(
		controllers.NewAuthController(services.NewAuthService(testMainUserDAO{})),
		controllers.NewEventController(services.NewEventService(testMainEventDAO{})),
		controllers.NewTicketController(services.NewTicketService(testMainTicketDAO{})),
	)
}

func TestGetServerPort(t *testing.T) {
	t.Setenv("SERVER_PORT", "")
	if got := getServerPort(); got != "8080" {
		t.Fatalf("expected default port 8080, got %q", got)
	}

	t.Setenv("SERVER_PORT", "9090")
	if got := getServerPort(); got != "9090" {
		t.Fatalf("expected configured port 9090, got %q", got)
	}
}

func TestLoadEnvAndBuildApplication(t *testing.T) {
	loadEnv()
	app := buildApplication()
	if app == nil {
		t.Fatalf("expected buildApplication to return a router")
	}
}

func TestBuildRouterRegistersPublicAndProtectedRoutes(t *testing.T) {
	t.Setenv("JWT_SECRET", "main-test-secret")
	router := buildTestAppRouter()

	req, _ := http.NewRequest(http.MethodOptions, "/events", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected CORS preflight to return 204, got %d", rec.Code)
	}

	req, _ = http.NewRequest(http.MethodGet, "/events", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected /events to be registered, got %d", rec.Code)
	}

	req, _ = http.NewRequest(http.MethodGet, "/profile", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected protected route to require auth, got %d", rec.Code)
	}
}
