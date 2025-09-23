package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"snapShare/models"
	"snapShare/services"
)

// Request DTOs
type CreateEventRequest struct {
	Name        string     `json:"name" validate:"required,min=1,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	EventDate   *time.Time `json:"event_date,omitempty"`
	OwnerEmail  string     `json:"owner_email" validate:"required,email"`
}

type UpdateEventRequest struct {
	Name        *string             `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=1000"`
	EventDate   *time.Time          `json:"event_date,omitempty"`
	Status      *models.EventStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive closed"`
}

// Response DTOs
type EventResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Description *string             `json:"description,omitempty"`
	EventDate   *time.Time          `json:"event_date,omitempty"`
	Status      models.EventStatus  `json:"status"`
	OwnerEmail  string              `json:"owner_email"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
	}
}

// CreateEvent creates a new event
func (h *EventHandler) CreateEvent(c echo.Context) error {
	var req CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Convert to service layer request
	serviceReq := &services.CreateEventRequest{
		Name:        req.Name,
		Description: req.Description,
		EventDate:   req.EventDate,
		OwnerEmail:  req.OwnerEmail,
	}

	event, err := h.eventService.CreateEvent(c.Request().Context(), serviceReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Code:        event.Code,
		Description: event.Description,
		EventDate:   event.EventDate,
		Status:      event.Status,
		OwnerEmail:  event.OwnerEmail,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

// GetEventByID retrieves an event by ID
func (h *EventHandler) GetEventByID(c echo.Context) error {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	event, err := h.eventService.GetEventByID(c.Request().Context(), eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	response := EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Code:        event.Code,
		Description: event.Description,
		EventDate:   event.EventDate,
		Status:      event.Status,
		OwnerEmail:  event.OwnerEmail,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// GetEventByCode retrieves an event by its unique code (for QR access)
func (h *EventHandler) GetEventByCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "event code is required")
	}

	event, err := h.eventService.GetEventByCode(c.Request().Context(), code)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	response := EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Code:        event.Code,
		Description: event.Description,
		EventDate:   event.EventDate,
		Status:      event.Status,
		OwnerEmail:  event.OwnerEmail,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// GetEventsByOwner retrieves all events owned by a user
func (h *EventHandler) GetEventsByOwner(c echo.Context) error {
	ownerEmail := c.QueryParam("owner_email")
	if ownerEmail == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "owner_email query parameter is required")
	}

	events, err := h.eventService.GetEventsByOwner(c.Request().Context(), ownerEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Convert to response DTOs
	responses := make([]EventResponse, len(events))
	for i, event := range events {
		responses[i] = EventResponse{
			ID:          event.ID.String(),
			Name:        event.Name,
			Code:        event.Code,
			Description: event.Description,
			EventDate:   event.EventDate,
			Status:      event.Status,
			OwnerEmail:  event.OwnerEmail,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, map[string]any{"events": responses})
}

// UpdateEvent updates an existing event
func (h *EventHandler) UpdateEvent(c echo.Context) error {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	var req UpdateEventRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Convert to service layer request
	serviceReq := &services.UpdateEventRequest{
		Name:        req.Name,
		Description: req.Description,
		EventDate:   req.EventDate,
		Status:      req.Status,
	}

	event, err := h.eventService.UpdateEvent(c.Request().Context(), eventID, serviceReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := EventResponse{
		ID:          event.ID.String(),
		Name:        event.Name,
		Code:        event.Code,
		Description: event.Description,
		EventDate:   event.EventDate,
		Status:      event.Status,
		OwnerEmail:  event.OwnerEmail,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteEvent deletes an event
func (h *EventHandler) DeleteEvent(c echo.Context) error {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	// TODO: Add authorization check to ensure only the owner can delete

	if err := h.eventService.DeleteEvent(c.Request().Context(), eventID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "event deleted"})
}

// CloseEvent closes an event (sets status to closed)
func (h *EventHandler) CloseEvent(c echo.Context) error {
	eventIDStr := c.Param("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	// TODO: Add authorization check to ensure only the owner can close

	if err := h.eventService.CloseEvent(c.Request().Context(), eventID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "event closed"})
}