package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"seesharpsi/htmx_quickstart/logger"
	"seesharpsi/htmx_quickstart/services"
	"seesharpsi/htmx_quickstart/templ"
)

type Handler struct {
	Service services.Service
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Info("handling index request", "request_id", requestID)

	// Execute business logic
	pageData, err := h.Service.RenderIndexPage(r.Context())
	if err != nil {
		slog.Error("failed to execute index page business logic", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get or create session
	sess, cookie := h.Service.GetOrCreateSession(r)
	http.SetCookie(w, &cookie)

	// Update page data with session info
	if sess != nil {
		pageData.UserID = sess.ID
		pageData.IsLoggedIn = true
	}

	// Render template
	if err := templ.Index().Render(r.Context(), w); err != nil {
		slog.Error("failed to render index template", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Info("handling test request", "request_id", requestID)

	// Execute business logic
	pageData, err := h.Service.RenderTestPage(r.Context())
	if err != nil {
		slog.Error("failed to execute test page business logic", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get or create session
	sess, cookie := h.Service.GetOrCreateSession(r)
	http.SetCookie(w, &cookie)

	// Update page data with session info
	if sess != nil {
		pageData.UserID = sess.ID
		pageData.IsLoggedIn = true
	}

	// Render template
	if err := templ.Test().Render(r.Context(), w); err != nil {
		slog.Error("failed to render test template", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Info("handling 404 request", "request_id", requestID, "path", r.URL.Path)

	// Execute business logic
	pageData, err := h.Service.RenderNotFoundPage(r.Context())
	if err != nil {
		slog.Error("failed to execute not found page business logic", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get or create session
	sess, cookie := h.Service.GetOrCreateSession(r)
	http.SetCookie(w, &cookie)

	// Update page data with session info
	if sess != nil {
		pageData.UserID = sess.ID
		pageData.IsLoggedIn = true
	}

	w.WriteHeader(http.StatusNotFound)
	if err := templ.Error404().Render(r.Context(), w); err != nil {
		slog.Error("failed to render 404 template", "error", err, "request_id", requestID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	requestID := logger.RequestIDFromContext(r.Context())
	slog.Debug("handling health check", "request_id", requestID)

	// Return simple JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}
