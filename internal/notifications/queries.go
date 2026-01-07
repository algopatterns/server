package notifications

const (
	queryCreate = `
		INSERT INTO notifications (user_id, type, title, body, data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, type, title, body, read, created_at
	`

	queryListForUser = `
		SELECT id, user_id, type, title, body, data, read, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	queryListUnreadForUser = `
		SELECT id, user_id, type, title, body, data, read, created_at
		FROM notifications
		WHERE user_id = $1 AND read = false
		ORDER BY created_at DESC
		LIMIT $2
	`

	queryMarkRead = `
		UPDATE notifications
		SET read = true
		WHERE id = $1 AND user_id = $2
	`

	queryMarkAllRead = `
		UPDATE notifications
		SET read = true
		WHERE user_id = $1 AND read = false
	`

	queryUnreadCount = `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND read = false
	`
)
