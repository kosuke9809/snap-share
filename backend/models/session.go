package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID      uuid.UUID      `json:"event_id" gorm:"type:uuid;not null;index"`
	GuestName    string         `json:"guest_name" gorm:"not null;size:255;index"`
	SessionToken string         `json:"session_token" gorm:"not null;unique;size:128"`
	ExpiresAt    time.Time      `json:"expires_at" gorm:"not null;index"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty"`

	Event Event `json:"event,omitempty" gorm:"foreignKey:EventID;references:ID;constraint:OnDelete:CASCADE"`
}
