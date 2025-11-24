package repositories

import (
	"app/domain/models"
	"fmt"
	"strconv"
	"time"
)

func (r *appRepository) SaveMessage(message *models.Message) (*models.Message, error) {
	err := r.Conn.Create(message).Error
	return message, err
}

// GetMessageHistory returns non-deleted messages for customers (soft delete filter applied)
func (r *appRepository) GetMessageHistory(conversationID uint64, limit int, cursor string) ([]models.Message, string, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	q := r.Conn.Where("conversation_id = ? AND deleted_at IS NULL", conversationID)

	if cursor != "" {
		if curID, err := strconv.ParseUint(cursor, 10, 64); err == nil {
			q = q.Where("id < ?", curID)
		}
	}

	var msgs []models.Message
	err := q.Order("id desc").Limit(limit + 1).Find(&msgs).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(msgs) > limit {
		nextCursor = fmt.Sprintf("%d", msgs[len(msgs)-1].ID)
		msgs = msgs[:len(msgs)-1]
	}

	// Reverse to chronological order
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nextCursor, nil
}

// GetMessageHistoryForAdmin returns ALL messages including soft-deleted ones for admins
func (r *appRepository) GetMessageHistoryForAdmin(conversationID uint64, limit int, cursor string) ([]models.Message, string, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// Admin sees all messages including soft-deleted
	q := r.Conn.Where("conversation_id = ?", conversationID)

	if cursor != "" {
		if curID, err := strconv.ParseUint(cursor, 10, 64); err == nil {
			q = q.Where("id < ?", curID)
		}
	}

	var msgs []models.Message
	err := q.Order("id desc").Limit(limit + 1).Find(&msgs).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	if len(msgs) > limit {
		nextCursor = fmt.Sprintf("%d", msgs[len(msgs)-1].ID)
		msgs = msgs[:len(msgs)-1]
	}

	// Reverse to chronological order
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nextCursor, nil
}

// SoftDeleteConversationMessages marks all messages in a conversation as soft-deleted
// and sets purge timestamp for permanent deletion after X days
func (r *appRepository) SoftDeleteConversationMessages(conversationID uint64, purgeAfterDays int) error {
	now := time.Now()
	purgeAt := now.AddDate(0, 0, purgeAfterDays)

	return r.Conn.Model(&models.Message{}).
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Updates(map[string]interface{}{
			"deleted_at": now,
			"purge_at":   purgeAt,
		}).Error
}

// ResetPurgeTimestamp removes purge countdown when conversation is reopened
// but keeps soft delete for filtering (only removes purge_at, not deleted_at)
func (r *appRepository) ResetPurgeTimestamp(conversationID uint64) error {
	return r.Conn.Model(&models.Message{}).
		Where("conversation_id = ? AND purge_at IS NOT NULL", conversationID).
		Update("purge_at", nil).Error
}

// PermanentlyDeleteExpiredMessages permanently deletes messages that have passed their purge date
func (r *appRepository) PermanentlyDeleteExpiredMessages() error {
	now := time.Now()
	return r.Conn.Unscoped(). // Unscoped for permanent deletion
					Where("purge_at IS NOT NULL AND purge_at <= ?", now).
					Delete(&models.Message{}).Error
}
