package users

import (
	"time"

	"github.com/algorave/server/internal/errors"
	"github.com/algorave/server/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getTodayUsage(db *pgxpool.Pool, c *gin.Context, userID string) (int, error) {
	var count int
	err := db.QueryRow(c.Request.Context(), `
		SELECT get_user_usage_today($1)
	`, userID).Scan(&count)

	if err != nil {
		logger.ErrorErr(err, "failed to get today's usage", "user_id", userID)
		errors.InternalError(c, "failed to fetch usage data", err)
		return 0, err
	}

	return count, nil
}

func getUserTierOrDefault(db *pgxpool.Pool, c *gin.Context, userID string) string {
	var tier string
	err := db.QueryRow(c.Request.Context(), `
		SELECT tier FROM users WHERE id = $1
	`, userID).Scan(&tier)

	if err != nil {
		logger.ErrorErr(err, "failed to get user tier", "user_id", userID)
		return "free"
	}

	return tier
}

func calculateDailyLimit(tier string) int {
	if tier == "pro" || tier == "byok" {
		return -1
	}
	return 100
}

func getUsageHistory(db *pgxpool.Pool, c *gin.Context, userID string) []DailyUsage {
	rows, err := db.Query(c.Request.Context(), `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM usage_logs
		WHERE user_id = $1
		AND is_byok = false
		AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY DATE(created_at)
		ORDER BY date DESC
	`, userID)

	if err != nil {
		logger.ErrorErr(err, "failed to get usage history", "user_id", userID)
		errors.InternalError(c, "failed to fetch usage history", err)
		return []DailyUsage{}
	}
	defer rows.Close()

	history := []DailyUsage{}
	for rows.Next() {
		var date time.Time
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			logger.ErrorErr(err, "failed to scan usage history row")
			continue
		}
		history = append(history, DailyUsage{
			Date:  date.Format("2006-01-02"),
			Count: count,
		})
	}

	return history
}
