package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"snapShare/config"
	"snapShare/handlers"
	"snapShare/infra/database"
	"snapShare/infra/r2"
	"snapShare/services"
)

// CustomValidator wraps the validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize R2 service
	r2Service, err := r2.NewR2Service(cfg.R2AccountID, cfg.R2AccessKey, cfg.R2SecretAccessKey, cfg.R2BucketName, cfg.R2PublicDomain)
	if err != nil {
		log.Fatal("Failed to initialize R2 service:", err)
	}

	// Initialize services
	sessionService := services.NewSessionService(db)
	eventService := services.NewEventService(db)
	photoService := services.NewPhotoService(db, r2Service)

	// Initialize handlers
	sessionHandler := handlers.NewSessionHandler(sessionService, eventService)
	eventHandler := handlers.NewEventHandler(eventService)
	photoHandler := handlers.NewPhotoHandler(photoService)

	// Initialize Echo
	e := echo.New()

	// Set validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	api := e.Group("/api")

	// Session routes
	api.POST("/sessions", sessionHandler.CreateSession)
	api.POST("/sessions/refresh", sessionHandler.RefreshSession)
	api.DELETE("/sessions", sessionHandler.RevokeSession)

	// Event routes
	api.POST("/events", eventHandler.CreateEvent)
	api.GET("/events/:code", eventHandler.GetEventByCode)

	// Photo routes (using actual handler methods) - require authentication
	photoAPI := api.Group("/photos", sessionHandler.AuthMiddleware())
	photoAPI.POST("/upload-url", photoHandler.GenerateUploadURL)
	photoAPI.POST("/confirm/:id", photoHandler.ConfirmUpload)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}
