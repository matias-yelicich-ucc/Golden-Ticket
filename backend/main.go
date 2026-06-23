package main

import (
	"log"
	"net/http"
	"os"

	"golden-ticket/backend/controllers"
	"golden-ticket/backend/dao"
	"golden-ticket/backend/middleware"
	"golden-ticket/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
}

func buildRouter(
	authController *controllers.AuthController,
	eventController *controllers.EventController,
	ticketController *controllers.TicketController,
) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.POST("/register", authController.Register)
	r.POST("/login", authController.Login)
	r.GET("/events", eventController.List)
	r.GET("/events/:id", eventController.GetByID)
	r.GET("/dashboard-stats", eventController.GetAdminDashboardStats)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			rol, _ := c.Get("rol")
			c.JSON(http.StatusOK, gin.H{
				"message": "Access granted to protected endpoint",
				"user_id": userID,
				"rol":     rol,
			})
		})

		protected.POST("/events/:id/tickets", ticketController.Buy)
		protected.GET("/my-tickets", ticketController.GetMyTickets)
		protected.POST("/my-tickets/:id/transfer", ticketController.Transfer)
		protected.POST("/my-tickets/:id/cancel", ticketController.Cancel)
		protected.DELETE("/my-tickets/:id", ticketController.Cancel)

		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.AuthorizeRole("administrador", "admin"))
		{
			adminOnly.GET("/dashboard", eventController.GetAdminDashboardStats)
			adminOnly.POST("/events", eventController.Create)
			adminOnly.PUT("/events/:id", eventController.Update)
			adminOnly.DELETE("/events/:id", eventController.Delete)
		}
	}

	return r
}

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return "8080"
	}
	return port
}

func buildApplication() *gin.Engine {
	userDAO := dao.NewUserDAO()
	authService := services.NewAuthService(userDAO)
	authController := controllers.NewAuthController(authService)

	eventDAO := dao.NewEventDAO()
	eventService := services.NewEventService(eventDAO)
	eventController := controllers.NewEventController(eventService)

	ticketDAO := dao.NewTicketDAO()
	ticketService := services.NewTicketService(ticketDAO)
	ticketController := controllers.NewTicketController(ticketService)

	return buildRouter(authController, eventController, ticketController)
}

func main() {
	loadEnv()
	dao.InitDB()

	r := buildApplication()
	port := getServerPort()

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
