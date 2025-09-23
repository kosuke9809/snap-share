package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"snapShare/services"
)

// Request DTOs
type UploadURLRequest struct {
	EventID     string `json:"event_id" validate:"required,uuid"`
	ContentType string `json:"content_type" validate:"required"`
}

type BulkUploadRequest struct {
	EventID string     `json:"event_id" validate:"required,uuid"`
	Files   []FileInfo `json:"files" validate:"required,min=1,max=50"`
}

type FileInfo struct {
	ContentType string `json:"content_type" validate:"required"`
	Size        int64  `json:"size,omitempty"`
}

type ConfirmUploadRequest struct {
	FileSize int64 `json:"file_size" validate:"required,min=1"`
}

type BulkConfirmRequest struct {
	Confirmations map[string]int64 `json:"confirmations" validate:"required"`
}

type DeleteBulkRequest struct {
	PhotoIDs []string `json:"photo_ids" validate:"required,min=1"`
	EventID  string   `json:"event_id" validate:"required,uuid"`
}

// Response DTOs
type UploadURLResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	PhotoID   string `json:"photo_id"`
}

type BulkUploadResponse struct {
	Uploads []UploadURLResponse `json:"uploads"`
	BatchID string              `json:"batch_id"`
}

type BulkDownloadResponse struct {
	DownloadURL string    `json:"download_url"`
	ExpiresAt   time.Time `json:"expires_at"`
	PhotoCount  int       `json:"photo_count"`
}

type PhotoHandler struct {
	photoService *services.PhotoService
}

func NewPhotoHandler(photoService *services.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		photoService: photoService,
	}
}

// GenerateUploadURL generates a presigned upload URL for a single photo
func (h *PhotoHandler) GenerateUploadURL(c echo.Context) error {
	var req UploadURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	// Get uploader name from session or context
	uploaderName := c.Get("uploader_name")
	if uploaderName == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "uploader name required")
	}

	uploadInfo, err := h.photoService.GenerateUploadURL(c.Request().Context(), eventID, uploaderName.(string), req.ContentType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := UploadURLResponse{
		UploadURL: uploadInfo.UploadURL,
		ObjectKey: uploadInfo.ObjectKey,
		PhotoID:   uploadInfo.PhotoID.String(),
	}

	return c.JSON(http.StatusOK, response)
}

// GenerateBulkUploadURLs generates multiple presigned upload URLs
func (h *PhotoHandler) GenerateBulkUploadURLs(c echo.Context) error {
	var req BulkUploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	uploaderName := c.Get("uploader_name")
	if uploaderName == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "uploader name required")
	}

	// Convert DTOs to service layer types
	files := make([]services.FileSpec, len(req.Files))
	for i, file := range req.Files {
		files[i] = services.FileSpec{
			ContentType: file.ContentType,
			Size:        file.Size,
		}
	}

	result, err := h.photoService.GenerateBulkUploadURLs(c.Request().Context(), eventID, uploaderName.(string), files)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Convert service layer response to DTO
	uploads := make([]UploadURLResponse, len(result.Uploads))
	for i, upload := range result.Uploads {
		uploads[i] = UploadURLResponse{
			UploadURL: upload.UploadURL,
			ObjectKey: upload.ObjectKey,
			PhotoID:   upload.PhotoID.String(),
		}
	}

	response := BulkUploadResponse{
		Uploads: uploads,
		BatchID: result.BatchID,
	}

	return c.JSON(http.StatusOK, response)
}

// ConfirmUpload confirms a single photo upload
func (h *PhotoHandler) ConfirmUpload(c echo.Context) error {
	photoIDStr := c.Param("id")
	photoID, err := uuid.Parse(photoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid photo ID")
	}

	var req ConfirmUploadRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.photoService.ConfirmUpload(c.Request().Context(), photoID, req.FileSize); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "upload confirmed"})
}

// ConfirmBulkUpload confirms multiple photo uploads
func (h *PhotoHandler) ConfirmBulkUpload(c echo.Context) error {
	var req BulkConfirmRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.photoService.ConfirmBulkUpload(c.Request().Context(), req.Confirmations); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "bulk upload confirmed"})
}

// GetPhotosByEvent retrieves all photos for an event
func (h *PhotoHandler) GetPhotosByEvent(c echo.Context) error {
	eventIDStr := c.Param("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	photos, err := h.photoService.GetPhotosByEvent(c.Request().Context(), eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"photos": photos})
}

// GenerateBulkDownloadURL creates a download URL for all photos in an event
func (h *PhotoHandler) GenerateBulkDownloadURL(c echo.Context) error {
	eventIDStr := c.Param("event_id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	downloadInfo, err := h.photoService.GenerateBulkDownloadURL(c.Request().Context(), eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := BulkDownloadResponse{
		DownloadURL: downloadInfo.DownloadURL,
		ExpiresAt:   downloadInfo.ExpiresAt,
		PhotoCount:  downloadInfo.PhotoCount,
	}

	return c.JSON(http.StatusOK, response)
}

// DeletePhoto deletes a single photo
func (h *PhotoHandler) DeletePhoto(c echo.Context) error {
	photoIDStr := c.Param("id")
	photoID, err := uuid.Parse(photoIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid photo ID")
	}

	// Check if user can delete (implement proper authorization)
	userCanDelete := true // TODO: Implement proper authorization logic

	if err := h.photoService.DeletePhoto(c.Request().Context(), photoID, userCanDelete); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "photo deleted"})
}

// DeleteBulkPhotos deletes multiple photos
func (h *PhotoHandler) DeleteBulkPhotos(c echo.Context) error {
	var req DeleteBulkRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid event ID")
	}

	// Convert string IDs to UUIDs
	photoIDs := make([]uuid.UUID, len(req.PhotoIDs))
	for i, idStr := range req.PhotoIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid photo ID: "+idStr)
		}
		photoIDs[i] = id
	}

	if err := h.photoService.DeleteBulkPhotos(c.Request().Context(), photoIDs, eventID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "photos deleted", "count": len(photoIDs)})
}
