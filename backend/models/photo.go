package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Photo struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID      uuid.UUID      `json:"event_id" gorm:"type:uuid;not null;index"`
	UploaderName string         `json:"uploader_name" gorm:"not null;size:100;index"`
	ObjectKey    string         `json:"object_key" gorm:"not null;size:255;index"`
	Size         int64          `json:"file_size" gorm:"not null"`
	MimeType     string         `json:"mime_type" gorm:"not null;size:50"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty"`

	Event Event `json:"event,omitempty" gorm:"foreignKey:EventID;references:ID;constraint:OnDelete:CASCADE"`
}
