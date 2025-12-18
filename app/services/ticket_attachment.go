package services

import (
	"app/domain/models"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"
)

func (s *appService) CreateTicketAttachment(ticketID int, file *multipart.FileHeader) (*models.TicketAttachment, error) {
	// Check if S3 is available
	if s.s3Repo == nil {
		return nil, fmt.Errorf("file storage service is not configured")
	}

	// Upload to MinIO
	ctx := context.Background()

	objectName, err := s.s3Repo.UploadFile(ctx, file)
	if err != nil {
		log.Printf("[ticket-attachment-service] UploadFile error: %v", err)
		return nil, err
	}
	log.Printf("[ticket-attachment-service] uploaded object: %s", objectName)

	attachment := &models.TicketAttachment{
		TicketID: ticketID,
		FilePath: objectName,
	}

	if err := s.repo.CreateTicketAttachment(attachment); err != nil {
		// Rollback: delete uploaded file
		log.Printf("[ticket-attachment-service] CreateTicketAttachment DB error: %v — rolling back delete %s", err, objectName)
		s.s3Repo.DeleteFile(ctx, objectName)
		return nil, err
	}

	log.Printf("[ticket-attachment-service] attachment record created: id=%d file=%s", attachment.ID, attachment.FilePath)
	return attachment, nil
}

func (s *appService) GetTicketAttachments() ([]models.TicketAttachment, error) {
	return s.repo.GetTicketAttachments()
}

func (s *appService) GetTicketAttachmentByID(id int) (*models.TicketAttachment, string, error) {
	attachment, err := s.repo.GetTicketAttachmentByID(id)
	if err != nil {
		return nil, "", err
	}

	// Generate presigned URL for download (valid for 1 hour)
	ctx := context.Background()

	downloadURL, err := s.s3Repo.GetFileURL(ctx, attachment.FilePath, 1*time.Hour)
	if err != nil {
		log.Printf("[ticket-attachment-service] GetFileURL error: %v", err)
		return attachment, "", err
	}

	return attachment, downloadURL, nil
}

func (s *appService) GetTicketAttachmentsByTicketID(ticketID int) ([]models.TicketAttachment, error) {
	return s.repo.GetTicketAttachmentsByTicketID(ticketID)
}

func (s *appService) UpdateTicketAttachment(id int, ticketID *int, file *multipart.FileHeader) (*models.TicketAttachment, error) {
	attachment, err := s.repo.GetTicketAttachmentByID(id)
	if err != nil {
		return nil, err
	}

	// Update ticket ID if provided
	if ticketID != nil {
		attachment.TicketID = *ticketID
	}

	// Update file if provided
	if file != nil {
		ctx := context.Background()

		// Upload new file
		objectName, err := s.s3Repo.UploadFile(ctx, file)
		if err != nil {
			log.Printf("[ticket-attachment-service] UploadFile error: %v", err)
			return nil, err
		}

		// Delete old file
		oldFilePath := attachment.FilePath
		attachment.FilePath = objectName

		if err := s.repo.UpdateTicketAttachment(attachment); err != nil {
			// Rollback: delete new file
			s.s3Repo.DeleteFile(ctx, objectName)
			return nil, err
		}

		// Delete old file after successful update
		if oldFilePath != "" {
			s.s3Repo.DeleteFile(ctx, oldFilePath)
		}
	} else {
		if err := s.repo.UpdateTicketAttachment(attachment); err != nil {
			return nil, err
		}
	}

	return attachment, nil
}

func (s *appService) DeleteTicketAttachment(id int) error {
	// Get the attachment first
	attachment, err := s.repo.GetTicketAttachmentByID(id)
	if err != nil {
		return err
	}

	// Delete from DB
	err = s.repo.DeleteTicketAttachment(id)
	if err != nil {
		return err
	}

	// Delete file from MinIO
	if attachment != nil && attachment.FilePath != "" {
		ctx := context.Background()
		_ = s.s3Repo.DeleteFile(ctx, attachment.FilePath)
	}

	return nil
}
