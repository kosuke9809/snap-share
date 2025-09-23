package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"snapShare/services"
)

// Request DTOs
type CreateSessionRequest struct {
	EventCode string `json:"event_code" validate:"required,len=8"`
	GuestName string `json:"guest_name" validate:"required,min=1,max=100"`
}

type RefreshSessionRequest struct {
	SessionToken string `json:"session_token" validate:"required"`
}

type RevokeSessionRequest struct {
	SessionToken string `json:"session_token" validate:"required"`
}

// Response DTOs
type SessionResponse struct {
	ID           string    `json:"id"`
	EventID      string    `json:"event_id"`
	GuestName    string    `json:"guest_name"`
	SessionToken string    `json:"session_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	Event        *EventResponse `json:"event,omitempty"`
}

type SessionsListResponse struct {
	Sessions []SessionResponse `json:"sessions"`
	Count    int              `json:"count"`
}

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

// CreateSession creates a new guest session for an event
func (h *SessionHandler) CreateSession(c echo.Context) error {
	var req CreateSessionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// First, get the event by code to validate it exists and is active
	// This could be done in the service layer, but we'll handle it here for clarity
	// In a real implementation, you might want to move this logic to the service

	// For now, we'll assume the event code maps to an event ID
	// This would typically involve another service call
	eventID := uuid.New() // TODO: Replace with actual event lookup by code

	session, err := h.sessionService.CreateSession(c.Request().Context(), eventID, req.GuestName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := SessionResponse{
		ID:           session.ID.String(),
		EventID:      session.EventID.String(),
		GuestName:    session.GuestName,
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}

	// Include event details if available
	if session.Event.ID != uuid.Nil {
		response.Event = &EventResponse{
			ID:          session.Event.ID.String(),
			Name:        session.Event.Name,
			Code:        session.Event.Code,
			Description: session.Event.Description,
			EventDate:   session.Event.EventDate,
			Status:      session.Event.Status,
			OwnerEmail:  session.Event.OwnerEmail,
			CreatedAt:   session.Event.CreatedAt,
			UpdatedAt:   session.Event.UpdatedAt,
		}
	}

	return c.JSON(http.StatusCreated, response)
}

// ValidateSession validates a session token and returns session info
func (h *SessionHandler) ValidateSession(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "session token is required")
	}

	session, err := h.sessionService.ValidateSession(c.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	response := SessionResponse{
		ID:           session.ID.String(),
		EventID:      session.EventID.String(),
		GuestName:    session.GuestName,
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}

	// Include event details
	if session.Event.ID != uuid.Nil {
		response.Event = &EventResponse{
			ID:          session.Event.ID.String(),
			Name:        session.Event.Name,
			Code:        session.Event.Code,
			Description: session.Event.Description,
			EventDate:   session.Event.EventDate,
			Status:      session.Event.Status,
			OwnerEmail:  session.Event.OwnerEmail,
			CreatedAt:   session.Event.CreatedAt,
			UpdatedAt:   session.Event.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, response)
}

// RefreshSession extends the session expiry time
func (h *SessionHandler) RefreshSession(c echo.Context) error {
	var req RefreshSessionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	session, err := h.sessionService.RefreshSession(c.Request().Context(), req.SessionToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := SessionResponse{
		ID:           session.ID.String(),
		EventID:      session.EventID.String(),
		GuestName:    session.GuestName,
		SessionToken: session.SessionToken,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// RevokeSession invalidates a session token
func (h *SessionHandler) RevokeSession(c echo.Context) error {
	var req RevokeSessionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.sessionService.RevokeSession(c.Request().Context(), req.SessionToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "session revoked"})
}

// GetSessionsByEvent retrieves all active sessions for an event (admin only)
func (h *SessionHandler) GetSessionsByEvent(c echo.Context) error {
	eventIDStr := c.Param("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	// TODO: Add authorization check to ensure only event owner can access

	sessions, err := h.sessionService.GetSessionsByEvent(c.Request().Context(), eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Convert to response DTOs
	responses := make([]SessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = SessionResponse{
			ID:           session.ID.String(),
			EventID:      session.EventID.String(),
			GuestName:    session.GuestName,
			SessionToken: session.SessionToken,
			ExpiresAt:    session.ExpiresAt,
			CreatedAt:    session.CreatedAt,
		}
	}

	result := SessionsListResponse{
		Sessions: responses,
		Count:    len(responses),
	}

	return c.JSON(http.StatusOK, result)
}

// CleanupExpiredSessions removes expired sessions (admin/system endpoint)
func (h *SessionHandler) CleanupExpiredSessions(c echo.Context) error {
	// TODO: Add proper authorization for admin/system endpoints

	if err := h.sessionService.CleanupExpiredSessions(c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "expired sessions cleaned up"})
}

// Middleware function to validate session from Authorization header
func (h *SessionHandler) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "authorization header required")
			}

			// Extract token from "Bearer <token>" format
			const bearerPrefix = "Bearer "
			if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization format")
			}

			token := authHeader[len(bearerPrefix):]
			session, err := h.sessionService.ValidateSession(c.Request().Context(), token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
			}

			// Set session info in context for use by handlers
			c.Set("session", session)
			c.Set("uploader_name", session.GuestName)
			c.Set("event_id", session.EventID.String())

			return next(c)
		}
	}
}