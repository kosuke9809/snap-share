package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"snapShare/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{
		db: db,
	}
}

type CreateEventRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description,omitempty"`
	EventDate   *time.Time `json:"event_date,omitempty"`
	OwnerEmail  string     `json:"owner_email" binding:"required,email"`
}

type UpdateEventRequest struct {
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	EventDate   *time.Time          `json:"event_date,omitempty"`
	Status      *models.EventStatus `json:"status,omitempty"`
}

// CreateEvent creates a new event with a unique code
func (s *EventService) CreateEvent(ctx context.Context, req *CreateEventRequest) (*models.Event, error) {
	code, err := s.generateUniqueCode(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique code: %w", err)
	}

	event := &models.Event{
		ID:          uuid.New(),
		Name:        req.Name,
		Code:        code,
		Description: req.Description,
		EventDate:   req.EventDate,
		Status:      models.EventStatusActive,
		OwnerEmail:  req.OwnerEmail,
	}

	if err := s.db.Create(event).Error; err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

// GetEventByID retrieves an event by its ID
func (s *EventService) GetEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	var event models.Event
	if err := s.db.First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

// GetEventByCode retrieves an event by its unique code
func (s *EventService) GetEventByCode(ctx context.Context, code string) (*models.Event, error) {
	var event models.Event
	if err := s.db.Where("code = ? AND status != ?", code, models.EventStatusClosed).First(&event).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found or closed")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

// GetEventsByOwner retrieves all events owned by a specific email
func (s *EventService) GetEventsByOwner(ctx context.Context, ownerEmail string) ([]models.Event, error) {
	var events []models.Event
	if err := s.db.Where("owner_email = ?", ownerEmail).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	return events, nil
}

// UpdateEvent updates an existing event
func (s *EventService) UpdateEvent(ctx context.Context, eventID uuid.UUID, req *UpdateEventRequest) (*models.Event, error) {
	var event models.Event
	if err := s.db.First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Update fields if provided
	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.EventDate != nil {
		updates["event_date"] = *req.EventDate
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&event).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update event: %w", err)
		}
	}

	return &event, nil
}

// DeleteEvent soft deletes an event
func (s *EventService) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.db.Delete(&models.Event{}, eventID).Error; err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil
}

// CloseEvent closes an event (sets status to closed)
func (s *EventService) CloseEvent(ctx context.Context, eventID uuid.UUID) error {
	if err := s.db.Model(&models.Event{}).
		Where("id = ?", eventID).
		Update("status", models.EventStatusClosed).Error; err != nil {
		return fmt.Errorf("failed to close event: %w", err)
	}

	return nil
}

// generateUniqueCode generates a unique 8-character alphanumeric code
func (s *EventService) generateUniqueCode(_ context.Context) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8
	maxAttempts := 10

	for range maxAttempts {
		code := make([]byte, codeLength)
		for i := range code {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return "", fmt.Errorf("failed to generate random number: %w", err)
			}
			code[i] = charset[num.Int64()]
		}

		codeStr := string(code)

		// Check if code already exists
		var count int64
		if err := s.db.Model(&models.Event{}).Where("code = ?", codeStr).Count(&count).Error; err != nil {
			return "", fmt.Errorf("failed to check code uniqueness: %w", err)
		}

		if count == 0 {
			return codeStr, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after %d attempts", maxAttempts)
}
