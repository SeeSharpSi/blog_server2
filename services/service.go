package services

import (
	"context"
	"log/slog"
	"net/http"

	"seesharpsi/htmx_quickstart/logger"
	"seesharpsi/htmx_quickstart/session"
)

// Service defines the interface for business logic operations
type Service interface {
	// Page operations
	RenderIndexPage(ctx context.Context) (*PageData, error)
	RenderTestPage(ctx context.Context) (*PageData, error)
	RenderNotFoundPage(ctx context.Context) (*PageData, error)

	// Session operations
	GetOrCreateSession(r *http.Request) (*session.Session, http.Cookie)

	// Business logic operations
	ProcessUserAction(ctx context.Context, action string) (*ActionResult, error)
	ValidateAndProcessInput(ctx context.Context, input map[string]interface{}) (*ValidationResult, error)
}

// PageData represents data needed for page rendering
type PageData struct {
	Title      string
	UserID     string
	IsLoggedIn bool
	CustomData map[string]interface{}
	Timestamp  string
}

// ActionResult represents the result of a user action
type ActionResult struct {
	Success bool
	Message string
	Data    map[string]interface{}
}

// ValidationResult represents the result of input validation
type ValidationResult struct {
	Valid  bool
	Errors []string
	Data   map[string]interface{}
}

// service implements the Service interface
type service struct {
	sessionManager *session.Manager
	logger         *slog.Logger
}

// NewService creates a new service instance with dependencies
func NewService(sessionManager *session.Manager, logger *slog.Logger) Service {
	return &service{
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// RenderIndexPage handles the business logic for rendering the index page
func (s *service) RenderIndexPage(ctx context.Context) (*PageData, error) {
	requestID := logger.RequestIDFromContext(ctx)
	s.logger.Info("rendering index page", "request_id", requestID)

	// Business logic for index page
	pageData := &PageData{
		Title:      "Welcome",
		IsLoggedIn: false, // In a real app, check session
		CustomData: map[string]interface{}{
			"welcomeMessage": "Hello from the service layer!",
		},
		Timestamp: "2024-01-01T00:00:00Z", // In real app, use current time
	}

	// Here you could add business logic like:
	// - Fetch user data from database
	// - Check user permissions
	// - Load dynamic content
	// - Apply business rules

	return pageData, nil
}

// RenderTestPage handles the business logic for rendering the test page
func (s *service) RenderTestPage(ctx context.Context) (*PageData, error) {
	requestID := logger.RequestIDFromContext(ctx)
	s.logger.Info("rendering test page", "request_id", requestID)

	// Business logic for test page
	pageData := &PageData{
		Title:      "Test Page",
		IsLoggedIn: false,
		CustomData: map[string]interface{}{
			"testMessage": "This is a test page from the service layer!",
		},
		Timestamp: "2024-01-01T00:00:00Z",
	}

	// Here you could add business logic like:
	// - Run tests
	// - Validate data
	// - Process forms
	// - etc.

	return pageData, nil
}

// RenderNotFoundPage handles the business logic for rendering the 404 page
func (s *service) RenderNotFoundPage(ctx context.Context) (*PageData, error) {
	requestID := logger.RequestIDFromContext(ctx)
	s.logger.Info("rendering not found page", "request_id", requestID)

	// Business logic for 404 page
	pageData := &PageData{
		Title:      "Page Not Found",
		IsLoggedIn: false,
		CustomData: map[string]interface{}{
			"errorCode": 404,
		},
		Timestamp: "2024-01-01T00:00:00Z",
	}

	// Here you could add business logic like:
	// - Log 404 attempts
	// - Track missing pages
	// - Suggest alternatives
	// - etc.

	return pageData, nil
}

// ProcessUserAction handles user actions with business logic
func (s *service) ProcessUserAction(ctx context.Context, action string) (*ActionResult, error) {
	requestID := logger.RequestIDFromContext(ctx)
	s.logger.Info("processing user action", "action", action, "request_id", requestID)

	// Business logic for user actions
	result := &ActionResult{
		Success: true,
		Message: "Action processed successfully",
		Data: map[string]interface{}{
			"action":      action,
			"processedAt": "2024-01-01T00:00:00Z",
		},
	}

	return result, nil
}

// ValidateAndProcessInput validates and processes user input
func (s *service) ValidateAndProcessInput(ctx context.Context, input map[string]interface{}) (*ValidationResult, error) {
	requestID := logger.RequestIDFromContext(ctx)
	s.logger.Info("validating and processing input", "request_id", requestID)

	// Business logic for input validation and processing
	result := &ValidationResult{
		Valid:  true,
		Errors: []string{},
		Data:   input, // Echo back the input for now
	}

	return result, nil
}

// GetOrCreateSession delegates to the session manager
func (s *service) GetOrCreateSession(r *http.Request) (*session.Session, http.Cookie) {
	return s.sessionManager.GetOrCreateSession(r)
}
