package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventStatus string

const (
	EventStatusActive   EventStatus = "active"
	EventStatusInactive EventStatus = "inactive"
	EventStatusClosed   EventStatus = "closed"
)

type Event struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"size:255;not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:8;not null"`
	Description *string        `json:"description,omitempty" gorm:"type:text"`
	EventDate   *time.Time     `json:"event_date,omitempty" gorm:"type:date"`
	Status      EventStatus    `json:"status" gorm:"not null;default:'active'"`
	OwnerEmail  string         `json:"owner_email" gorm:"not null;size:255"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty"`

	Photos []Photo `json:"photos,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}
