package repositories

import (
	"app/domain/models"
	"fmt"
	"strconv"
)

func (r *appRepository) SaveMessage(message *models.Message) (*models.Message, error) {
	err := r.Conn.Create(message).Error
	return message, err
}

// GetMessageHistory returns messages for a conversation using cursor pagination.
// - conversationID: conversation to fetch messages for
// - limit: maximum number of messages to return (server caps applied)
// - cursor: optional string cursor representing last seen message ID; if provided, returns messages with id < cursor
// Returns messages (chronological ascending), nextCursor (string, empty if no more), error
func (r *appRepository) GetMessageHistory(conversationID uint64, limit int, cursor string) ([]models.Message, string, error) {
	// cap and defaults
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	q := r.Conn.Where("conversation_id = ?", conversationID)

	// apply cursor if provided
	if cursor != "" {
		if curID, err := strconv.ParseUint(cursor, 10, 64); err == nil {
			q = q.Where("id < ?", curID)
		}
	}

	// fetch limit+1 to detect if there's a next page
	var msgs []models.Message
	err := q.Order("id desc").Limit(limit + 1).Find(&msgs).Error
	if err != nil {
		return nil, "", err
	}

	var nextCursor string
	// if we fetched more than limit, there is a next cursor
	if len(msgs) > limit {
		nextCursor = fmt.Sprintf("%d", msgs[len(msgs)-1].ID)
		msgs = msgs[:len(msgs)-1]
	}

	// currently msgs are in descending order (newest first) — reverse to chronological ascending
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	return msgs, nextCursor, nil
}
