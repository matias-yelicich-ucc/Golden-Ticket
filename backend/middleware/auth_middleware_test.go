package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golden-ticket/backend/utils"

	"github.com/gin-gonic/gin"
)

func buildMiddlewareRouter(mw gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userID": c.GetUint("userID"),
			"rol":    c.GetString("rol"),
		})
	})
	return router
}

func TestAuthMiddlewareScenarios(t *testing.T) {
	os.Setenv("JWT_SECRET", "middleware-secret")
	defer os.Unsetenv("JWT_SECRET")
	os.Setenv("JWT_EXPIRATION_MINUTES", "1")
	defer os.Unsetenv("JWT_EXPIRATION_MINUTES")

	router := buildMiddlewareRouter(AuthMiddleware())

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d", rec.Code)
	}

	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token abc")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid bearer format, got %d", rec.Code)
	}

	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", rec.Code)
	}

	validToken, err := utils.GenerateToken(7, "admin")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	req, _ = http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid token, got %d", rec.Code)
	}
}

func TestAuthMiddlewareRejectsExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "middleware-secret")
	defer os.Unsetenv("JWT_SECRET")
	os.Setenv("JWT_EXPIRATION_MINUTES", "-1")
	defer os.Unsetenv("JWT_EXPIRATION_MINUTES")

	expiredToken, err := utils.GenerateToken(8, "cliente")
	if err != nil {
		t.Fatalf("failed to generate expired token: %v", err)
	}

	// tiny pause so IssuedAt/ExpiresAt aren't identical edge cases
	time.Sleep(10 * time.Millisecond)

	router := buildMiddlewareRouter(AuthMiddleware())
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", rec.Code)
	}
}

func TestAuthorizeRoleScenarios(t *testing.T) {
	gin.SetMode(gin.TestMode)

	makeRouter := func(seedContext func(*gin.Context)) *gin.Engine {
		router := gin.New()
		router.GET("/admin", func(c *gin.Context) {
			if seedContext != nil {
				seedContext(c)
			}
			c.Next()
		}, AuthorizeRole("admin", "administrador"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		return router
	}

	tests := []struct {
		name       string
		seed       func(*gin.Context)
		statusCode int
	}{
		{
			name:       "missing role in context",
			seed:       nil,
			statusCode: http.StatusForbidden,
		},
		{
			name: "invalid role type",
			seed: func(c *gin.Context) {
				c.Set("rol", 123)
			},
			statusCode: http.StatusForbidden,
		},
		{
			name: "forbidden role",
			seed: func(c *gin.Context) {
				c.Set("rol", "cliente")
			},
			statusCode: http.StatusForbidden,
		},
		{
			name: "allowed role",
			seed: func(c *gin.Context) {
				c.Set("rol", "admin")
			},
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := makeRouter(tt.seed)
			req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tt.statusCode {
				t.Fatalf("expected status %d, got %d", tt.statusCode, rec.Code)
			}
		})
	}
}
