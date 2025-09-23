package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"snapShare/models"
)

type SessionService struct {
	db *gorm.DB
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) CreateSession(ctx context.Context, eventID uuid.UUID, guestName string) (*models.Session, error) {
	// Validate event exists and is active
	var event models.Event
	if err := s.db.Where("id = ? AND status = ?", eventID, models.EventStatusActive).First(&event).Error; err != nil {
		return nil, fmt.Errorf("event not found or inactive: %w", err)
	}

	// Generate session token
	token, err := s.generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session with 24 hour expiry
	session := models.Session{
		ID:           uuid.New(),
		EventID:      eventID,
		GuestName:    guestName,
		SessionToken: token,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Load the event relation
	session.Event = event

	return &session, nil
}

func (s *SessionService) ValidateSession(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	err := s.db.Preload("Event").
		Where("session_token = ? AND expires_at > ?", token, time.Now()).
		First(&session).Error

	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	// Check if event is still active
	if session.Event.Status != models.EventStatusActive {
		return nil, fmt.Errorf("event is no longer active")
	}

	return &session, nil
}

func (s *SessionService) RefreshSession(ctx context.Context, token string) (*models.Session, error) {
	session, err := s.ValidateSession(ctx, token)
	if err != nil {
		return nil, err
	}

	// Extend expiry by 24 hours
	newExpiresAt := time.Now().Add(24 * time.Hour)
	if err := s.db.Model(&session).Update("expires_at", newExpiresAt).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	session.ExpiresAt = newExpiresAt
	return session, nil
}

func (s *SessionService) RevokeSession(ctx context.Context, token string) error {
	result := s.db.Where("session_token = ?", token).Delete(&models.Session{})
	if result.Error != nil {
		return fmt.Errorf("failed to revoke session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (s *SessionService) GetSessionsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Session, error) {
	var sessions []models.Session
	err := s.db.Where("event_id = ? AND expires_at > ?", eventID, time.Now()).
		Order("created_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	return sessions, nil
}

func (s *SessionService) CleanupExpiredSessions(ctx context.Context) error {
	result := s.db.Where("expires_at <= ?", time.Now()).Delete(&models.Session{})
	return result.Error
}

func (s *SessionService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}