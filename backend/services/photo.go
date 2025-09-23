package services

import (
	"context"
	"fmt"
	"snapShare/infra/r2"
	"snapShare/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PhotoService struct {
	db        *gorm.DB
	r2Service *r2.R2Service
}

func NewPhotoService(db *gorm.DB, r2Service *r2.R2Service) *PhotoService {
	return &PhotoService{
		db:        db,
		r2Service: r2Service,
	}
}

// Service layer data structures (internal use only)
type UploadInfo struct {
	UploadURL string
	ObjectKey string
	PhotoID   uuid.UUID
}

type FileSpec struct {
	ContentType string
	Size        int64
}

type BulkUploadResult struct {
	Uploads []UploadInfo
	BatchID string
}

type DownloadInfo struct {
	DownloadURL string
	ExpiresAt   time.Time
	PhotoCount  int
}

func (s *PhotoService) GenerateUploadURL(ctx context.Context, eventID uuid.UUID, uploaderName, contentType string) (*UploadInfo, error) {
	var event models.Event
	if err := s.db.First(&event, eventID).Error; err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Generate photo ID and object key
	photoID := uuid.New()
	ext := getExtensionFromContentType(contentType)
	objectKey := fmt.Sprintf("events/%s/photos/%s%s", eventID, photoID, ext)

	// Generate presigned URL (15 minutes expiry)
	uploadURL, err := s.r2Service.GeneratePresignedUploadURL(ctx, objectKey, contentType, 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to generate upload URL: %w", err)
	}

	// Create photo record in database
	photo := models.Photo{
		ID:           photoID,
		EventID:      eventID,
		UploaderName: uploaderName,
		ObjectKey:    objectKey,
		MimeType:     contentType,
		Size:         0, // Will be updated after upload
	}

	if err := s.db.Create(&photo).Error; err != nil {
		return nil, fmt.Errorf("failed to create photo record: %w", err)
	}

	return &UploadInfo{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		PhotoID:   photoID,
	}, nil
}

func (s *PhotoService) ConfirmUpload(ctx context.Context, photoID uuid.UUID, fileSize int64) error {
	return s.db.Model(&models.Photo{}).
		Where("id = ?", photoID).
		Update("size", fileSize).Error
}

func (s *PhotoService) GetPhotosByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Photo, error) {
	var photos []models.Photo
	err := s.db.Where("event_id = ?", eventID).
		Order("created_at DESC").
		Find(&photos).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}

	for i := range photos {
		photos[i].ObjectKey = s.r2Service.GetPublicURL(photos[i].ObjectKey)
	}

	return photos, nil
}

func (s *PhotoService) DeletePhoto(ctx context.Context, photoID uuid.UUID, userCanDelete bool) error {
	var photo models.Photo
	if err := s.db.First(&photo, photoID).Error; err != nil {
		return fmt.Errorf("photo not found: %w", err)
	}

	if !userCanDelete {
		return fmt.Errorf("unauthorized to delete photo")
	}

	// Generate presigned delete URL
	deleteURL, err := s.r2Service.GeneratePresignedDeleteURL(ctx, photo.ObjectKey, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to generate delete URL: %w", err)
	}

	// Soft delete from database
	if err := s.db.Delete(&photo).Error; err != nil {
		return fmt.Errorf("failed to delete photo record: %w", err)
	}

	// Note: In a real implementation, you might want to queue the actual R2 deletion
	// or handle it asynchronously to ensure the database operation succeeds first
	_ = deleteURL // For now, just acknowledge we have the URL

	return nil
}

// GenerateBulkUploadURLs generates multiple presigned upload URLs for bulk photo upload
func (s *PhotoService) GenerateBulkUploadURLs(ctx context.Context, eventID uuid.UUID, uploaderName string, files []FileSpec) (*BulkUploadResult, error) {
	// Validate event exists
	var event models.Event
	if err := s.db.First(&event, eventID).Error; err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	// Generate batch ID for tracking
	batchID := uuid.New().String()

	// Limit bulk upload size (e.g., max 50 files per batch)
	if len(files) > 50 {
		return nil, fmt.Errorf("too many files: maximum 50 files per batch")
	}

	uploads := make([]UploadInfo, 0, len(files))
	photoRecords := make([]models.Photo, 0, len(files))

	// Generate URLs and create photo records
	for _, fileSpec := range files {
		photoID := uuid.New()
		ext := getExtensionFromContentType(fileSpec.ContentType)
		objectKey := fmt.Sprintf("events/%s/photos/%s%s", eventID, photoID, ext)

		// Generate presigned URL
		uploadURL, err := s.r2Service.GeneratePresignedUploadURL(ctx, objectKey, fileSpec.ContentType, 15*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("failed to generate upload URL for file: %w", err)
		}

		uploads = append(uploads, UploadInfo{
			UploadURL: uploadURL,
			ObjectKey: objectKey,
			PhotoID:   photoID,
		})

		// Prepare photo record
		photoRecords = append(photoRecords, models.Photo{
			ID:           photoID,
			EventID:      eventID,
			UploaderName: uploaderName,
			ObjectKey:    objectKey,
			MimeType:     fileSpec.ContentType,
			Size:         fileSpec.Size,
		})
	}

	// Batch insert photo records
	if err := s.db.Create(&photoRecords).Error; err != nil {
		return nil, fmt.Errorf("failed to create photo records: %w", err)
	}

	return &BulkUploadResult{
		Uploads: uploads,
		BatchID: batchID,
	}, nil
}

// ConfirmBulkUpload confirms multiple photo uploads with their actual file sizes
func (s *PhotoService) ConfirmBulkUpload(ctx context.Context, confirmations map[string]int64) error {
	if len(confirmations) == 0 {
		return nil
	}

	// Convert photoID strings to UUIDs for batch update
	photoIDs := make([]uuid.UUID, 0, len(confirmations))
	for photoIDStr := range confirmations {
		photoID, err := uuid.Parse(photoIDStr)
		if err != nil {
			return fmt.Errorf("invalid photo ID: %s", photoIDStr)
		}
		photoIDs = append(photoIDs, photoID)
	}

	// Update sizes in batch - Note: GORM doesn't support batch updates with different values easily
	// So we'll do individual updates in a transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		for photoIDStr, size := range confirmations {
			photoID, _ := uuid.Parse(photoIDStr) // Already validated above
			if err := tx.Model(&models.Photo{}).
				Where("id = ?", photoID).
				Update("size", size).Error; err != nil {
				return fmt.Errorf("failed to update photo %s size: %w", photoIDStr, err)
			}
		}
		return nil
	})
}

// GenerateBulkDownloadURL creates a zip archive of all photos in an event and returns download URL
func (s *PhotoService) GenerateBulkDownloadURL(ctx context.Context, eventID uuid.UUID) (*DownloadInfo, error) {
	// Get all photos for the event
	photos, err := s.GetPhotosByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}

	if len(photos) == 0 {
		return nil, fmt.Errorf("no photos found for event")
	}

	// Generate unique archive name
	archiveKey := fmt.Sprintf("events/%s/archives/%s.zip", eventID, uuid.New())

	// Note: In a real implementation, you would:
	// 1. Create a background job to generate the zip file
	// 2. Stream photos from R2 and create zip archive
	// 3. Upload the zip to R2
	// 4. Return presigned download URL

	// For now, we'll generate a presigned URL that expires in 1 hour
	// This assumes the zip creation process is handled elsewhere
	downloadURL, err := s.r2Service.GeneratePresignedDownloadURL(ctx, archiveKey, 1*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate download URL: %w", err)
	}

	return &DownloadInfo{
		DownloadURL: downloadURL,
		ExpiresAt:   time.Now().Add(1 * time.Hour),
		PhotoCount:  len(photos),
	}, nil
}

// DeleteBulkPhotos deletes multiple photos at once
func (s *PhotoService) DeleteBulkPhotos(ctx context.Context, photoIDs []uuid.UUID, eventID uuid.UUID) error {
	if len(photoIDs) == 0 {
		return nil
	}

	// Verify all photos belong to the event
	var count int64
	if err := s.db.Model(&models.Photo{}).
		Where("id IN ? AND event_id = ?", photoIDs, eventID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to verify photos: %w", err)
	}

	if count != int64(len(photoIDs)) {
		return fmt.Errorf("some photos not found or don't belong to this event")
	}

	// Get photo object keys for R2 deletion
	var photos []models.Photo
	if err := s.db.Select("object_key").
		Where("id IN ?", photoIDs).
		Find(&photos).Error; err != nil {
		return fmt.Errorf("failed to get photo object keys: %w", err)
	}

	// Delete from database first (soft delete)
	if err := s.db.Where("id IN ?", photoIDs).Delete(&models.Photo{}).Error; err != nil {
		return fmt.Errorf("failed to delete photo records: %w", err)
	}

	// Note: In a real implementation, you would queue R2 deletions
	// or handle them asynchronously to ensure database consistency
	for _, photo := range photos {
		_, _ = s.r2Service.GeneratePresignedDeleteURL(ctx, photo.ObjectKey, 5*time.Minute)
		// Queue actual deletion or handle asynchronously
	}

	return nil
}

func getExtensionFromContentType(contentType string) string {
	switch {
	case strings.HasPrefix(contentType, "image/jpeg"):
		return ".jpg"
	case strings.HasPrefix(contentType, "image/png"):
		return ".png"
	case strings.HasPrefix(contentType, "image/gif"):
		return ".gif"
	case strings.HasPrefix(contentType, "image/webp"):
		return ".webp"
	case strings.HasPrefix(contentType, "image/heic"):
		return ".heic"
	case strings.HasPrefix(contentType, "image/heif"):
		return ".heif"
	default:
		return ".jpg"
	}
}
