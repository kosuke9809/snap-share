package database

import (
	"context"
	"log"
	"snapShare/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	// Check if data already exists
	var eventCount int64
	if err := db.Model(&models.Event{}).Count(&eventCount).Error; err != nil {
		return err
	}

	if eventCount > 0 {
		log.Println("Seed data already exists, skipping...")
		return nil
	}

	ctx := context.Background()

	// Helper function to create string pointer
	stringPtr := func(s string) *string { return &s }
	timePtr := func(t time.Time) *time.Time { return &t }

	// Create events
	events := []models.Event{
		{
			ID:          uuid.New(),
			Name:        "山田太郎・花子 結婚式",
			Code:        "WEDDING1",
			Description: stringPtr("2025年春の結婚式です。皆様からの写真をお待ちしています！"),
			EventDate:   timePtr(time.Date(2025, 4, 15, 14, 0, 0, 0, time.UTC)),
			Status:      models.EventStatusActive,
			OwnerEmail:  "yamada@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "田中家 家族旅行",
			Code:        "TRAVEL02",
			Description: stringPtr("沖縄旅行の思い出を共有しましょう"),
			EventDate:   timePtr(time.Date(2025, 5, 1, 9, 0, 0, 0, time.UTC)),
			Status:      models.EventStatusActive,
			OwnerEmail:  "tanaka@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "高校同窓会 2025",
			Code:        "REUNION2",
			Description: stringPtr("卒業から10年！懐かしい仲間たちとの再会"),
			EventDate:   timePtr(time.Date(2025, 8, 20, 18, 0, 0, 0, time.UTC)),
			Status:      models.EventStatusActive,
			OwnerEmail:  "alumni@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, event := range events {
		if err := db.WithContext(ctx).Create(&event).Error; err != nil {
			return err
		}
	}

	// Create sample sessions
	sessions := []models.Session{
		{
			ID:           uuid.New(),
			EventID:      events[0].ID, // Wedding
			GuestName:    "佐藤一郎",
			SessionToken: "sample_token_1_sato",
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			EventID:      events[0].ID, // Wedding
			GuestName:    "鈴木花子",
			SessionToken: "sample_token_2_suzuki",
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			EventID:      events[1].ID, // Travel
			GuestName:    "田中次郎",
			SessionToken: "sample_token_3_tanaka",
			ExpiresAt:    time.Now().Add(24 * time.Hour),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, session := range sessions {
		if err := db.WithContext(ctx).Create(&session).Error; err != nil {
			return err
		}
	}

	// Create sample photos
	photos := []models.Photo{
		{
			ID:           uuid.New(),
			EventID:      events[0].ID, // Wedding
			UploaderName: "佐藤一郎",
			ObjectKey:    "events/" + events[0].ID.String() + "/photos/sample1.jpg",
			Size:         1024000,
			MimeType:     "image/jpeg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			EventID:      events[0].ID, // Wedding
			UploaderName: "鈴木花子",
			ObjectKey:    "events/" + events[0].ID.String() + "/photos/sample2.jpg",
			Size:         2048000,
			MimeType:     "image/jpeg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			EventID:      events[1].ID, // Travel
			UploaderName: "田中次郎",
			ObjectKey:    "events/" + events[1].ID.String() + "/photos/sample3.jpg",
			Size:         1536000,
			MimeType:     "image/jpeg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			EventID:      events[2].ID, // Reunion
			UploaderName: "山田太郎",
			ObjectKey:    "events/" + events[2].ID.String() + "/photos/sample4.jpg",
			Size:         1800000,
			MimeType:     "image/jpeg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	for _, photo := range photos {
		if err := db.WithContext(ctx).Create(&photo).Error; err != nil {
			return err
		}
	}

	log.Printf("Seed data created successfully:")
	log.Printf("- %d events", len(events))
	log.Printf("- %d sessions", len(sessions))
	log.Printf("- %d photos", len(photos))
	log.Println("Sample event codes: WEDDING1, TRAVEL02, REUNION2")

	return nil
}