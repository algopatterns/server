package attribution

import "time"

type Attribution struct {
	ID                    string
	SourceStrudelID       string
	SourceStrudelTitle    string
	RequestingUserID      *string
	RequestingDisplayName *string
	SimilarityScore       *float32
	CreatedAt             time.Time
}

type AttributionStats struct {
	TotalUses      int
	UniqueStrudels int
	LastUsedAt     *time.Time
}
