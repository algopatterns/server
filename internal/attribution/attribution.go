package attribution

import (
	"context"

	"github.com/algrv/server/internal/logger"
	"github.com/algrv/server/internal/retriever"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

// records that examples were used as RAG context
// runs asynchronously to not block the agent response
func (s *Service) RecordAttributions(
	ctx context.Context,
	examples []retriever.ExampleResult,
	requestingUserID string,
	targetStrudelID *string,
) {
	go func() {
		for _, ex := range examples {
			if ex.UserID == "" || ex.ID == "" {
				continue
			}

			// don't record self-attribution
			if ex.UserID == requestingUserID {
				continue
			}

			_, err := s.db.Exec(
				context.Background(),
				queryRecordAttribution,
				ex.ID,
				targetStrudelID,
				requestingUserID,
				ex.Similarity,
			)

			if err != nil {
				logger.Warn("failed to record attribution", "error", err, "source_strudel_id", ex.ID)
			}
		}
	}()
}

// gets attribution stats for a user's strudels
func (s *Service) GetUserAttributionStats(ctx context.Context, userID string) (*AttributionStats, error) {
	var stats AttributionStats

	err := s.db.QueryRow(ctx, queryGetUserAttributionStats, userID).Scan(
		&stats.TotalUses,
		&stats.UniqueStrudels,
		&stats.LastUsedAt,
	)

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// gets recent attributions for a user (their strudels being used)
func (s *Service) GetRecentAttributions(ctx context.Context, userID string, limit int) ([]Attribution, error) {
	rows, err := s.db.Query(ctx, queryGetRecentAttributions, userID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var attributions []Attribution

	for rows.Next() {
		var a Attribution
		err := rows.Scan(
			&a.ID,
			&a.SourceStrudelID,
			&a.SourceStrudelTitle,
			&a.RequestingUserID,
			&a.RequestingDisplayName,
			&a.SimilarityScore,
			&a.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		attributions = append(attributions, a)
	}

	return attributions, rows.Err()
}
