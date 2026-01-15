package botdefense

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	keyTrappedIP  = "botdefense:trapped:%s"
	keyRateLimit  = "botdefense:rate:%s"
	keyTrapReason = "botdefense:reason:%s"
)

// indicates why an IP was trapped
type TrapReason string

const (
	ReasonHoneypot   TrapReason = "honeypot"
	ReasonBotPattern TrapReason = "bot_pattern"
)

// manages trapped IPs and rate limits in Redis
type Store struct {
	client *redis.Client
	config *Config
}

// creates a new bot defense store
func NewStore(client *redis.Client, config *Config) *Store {
	return &Store{
		client: client,
		config: config,
	}
}

// traps an IP with a reason
func (s *Store) TrapIP(ctx context.Context, ip string, reason TrapReason) error {
	trappedKey := fmt.Sprintf(keyTrappedIP, ip)
	reasonKey := fmt.Sprintf(keyTrapReason, ip)

	pipe := s.client.Pipeline()
	pipe.Set(ctx, trappedKey, "1", s.config.TrapTTL)
	pipe.Set(ctx, reasonKey, string(reason), s.config.TrapTTL)

	_, err := pipe.Exec(ctx)
	return err
}

// checks if an IP is currently trapped
func (s *Store) IsTrapped(ctx context.Context, ip string) (bool, TrapReason, error) {
	trappedKey := fmt.Sprintf(keyTrappedIP, ip)
	reasonKey := fmt.Sprintf(keyTrapReason, ip)

	pipe := s.client.Pipeline()
	trappedCmd := pipe.Exists(ctx, trappedKey)
	reasonCmd := pipe.Get(ctx, reasonKey)

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		if trappedCmd.Err() != nil && !errors.Is(trappedCmd.Err(), redis.Nil) {
			return false, "", trappedCmd.Err()
		}
	}

	exists := trappedCmd.Val() > 0
	if !exists {
		return false, "", nil
	}

	reason, err := reasonCmd.Result()
	if errors.Is(err, redis.Nil) {
		reason = string(ReasonBotPattern)
	} else if err != nil {
		return true, "", err
	}

	return true, TrapReason(reason), nil
}

// removes an IP from the trap (for manual intervention)
func (s *Store) ReleaseIP(ctx context.Context, ip string) error {
	trappedKey := fmt.Sprintf(keyTrappedIP, ip)
	reasonKey := fmt.Sprintf(keyTrapReason, ip)

	return s.client.Del(ctx, trappedKey, reasonKey).Err()
}

// increments the request count for an IP and returns the new count
func (s *Store) IncrementRate(ctx context.Context, ip string) (int64, error) {
	key := fmt.Sprintf(keyRateLimit, ip)

	pipe := s.client.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, s.config.RateLimitWindow)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incrCmd.Val(), nil
}
