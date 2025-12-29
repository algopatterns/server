package websocket

import (
	"testing"
	"time"
)

// Test agent request rate limiting for free tier (10/minute)
func TestAgentRequestRateLimitFreeTier(t *testing.T) {
	client := &Client{
		Tier:                   "free",
		agentRequestTimestamps: make([]time.Time, 0, 10),
	}

	// First 10 requests should pass
	for i := 0; i < 10; i++ {
		if !client.checkAgentRequestRateLimit() {
			t.Errorf("Request %d should have been allowed, but was rate limited", i+1)
		}
	}

	// 11th request should be rate limited
	if client.checkAgentRequestRateLimit() {
		t.Error("11th request should have been rate limited, but was allowed")
	}

	if len(client.agentRequestTimestamps) != 10 {
		t.Errorf("Expected 10 timestamps, got %d", len(client.agentRequestTimestamps))
	}
}

// Test agent request rate limiting for pro tier (20/minute)
func TestAgentRequestRateLimitProTier(t *testing.T) {
	client := &Client{
		Tier:                   "pro",
		agentRequestTimestamps: make([]time.Time, 0, 20),
	}

	// First 20 requests should pass
	for i := 0; i < 20; i++ {
		if !client.checkAgentRequestRateLimit() {
			t.Errorf("Request %d should have been allowed, but was rate limited", i+1)
		}
	}

	// 21st request should be rate limited
	if client.checkAgentRequestRateLimit() {
		t.Error("21st request should have been rate limited, but was allowed")
	}

	if len(client.agentRequestTimestamps) != 20 {
		t.Errorf("Expected 20 timestamps, got %d", len(client.agentRequestTimestamps))
	}
}

// Test agent request rate limiting for BYOK tier (30/minute)
func TestAgentRequestRateLimitBYOKTier(t *testing.T) {
	client := &Client{
		Tier:                   "byok",
		agentRequestTimestamps: make([]time.Time, 0, 30),
	}

	// First 30 requests should pass
	for i := 0; i < 30; i++ {
		if !client.checkAgentRequestRateLimit() {
			t.Errorf("Request %d should have been allowed, but was rate limited", i+1)
		}
	}

	// 31st request should be rate limited
	if client.checkAgentRequestRateLimit() {
		t.Error("31st request should have been rate limited, but was allowed")
	}

	if len(client.agentRequestTimestamps) != 30 {
		t.Errorf("Expected 30 timestamps, got %d", len(client.agentRequestTimestamps))
	}
}

// Test agent request rate limit window expiration
func TestAgentRequestRateLimitWindowExpiration(t *testing.T) {
	client := &Client{
		Tier:                   "free",
		agentRequestTimestamps: make([]time.Time, 0, 10),
	}

	// Simulate 10 requests from 2 minutes ago (should be expired)
	twoMinutesAgo := time.Now().Add(-2 * time.Minute)
	for i := 0; i < 10; i++ {
		client.agentRequestTimestamps = append(client.agentRequestTimestamps, twoMinutesAgo)
	}

	// Next request should pass because old timestamps are expired
	if !client.checkAgentRequestRateLimit() {
		t.Error("Request should have been allowed after old timestamps expired")
	}

	// Old timestamps should be cleaned up, only 1 new one should remain
	if len(client.agentRequestTimestamps) != 1 {
		t.Errorf("Expected 1 timestamp after cleanup, got %d", len(client.agentRequestTimestamps))
	}
}

// Test code update rate limiting (10/second)
func TestCodeUpdateRateLimit(t *testing.T) {
	client := &Client{
		codeUpdateTimestamps: make([]time.Time, 0, maxCodeUpdatesPerSecond),
	}

	// First 10 updates should pass
	for i := 0; i < maxCodeUpdatesPerSecond; i++ {
		if !client.checkCodeUpdateRateLimit() {
			t.Errorf("Code update %d should have been allowed, but was rate limited", i+1)
		}
	}

	// 11th update should be rate limited
	if client.checkCodeUpdateRateLimit() {
		t.Error("11th code update should have been rate limited, but was allowed")
	}

	if len(client.codeUpdateTimestamps) != maxCodeUpdatesPerSecond {
		t.Errorf("Expected %d timestamps, got %d", maxCodeUpdatesPerSecond, len(client.codeUpdateTimestamps))
	}
}

// Test code update rate limit window expiration (1 second window)
func TestCodeUpdateRateLimitWindowExpiration(t *testing.T) {
	client := &Client{
		codeUpdateTimestamps: make([]time.Time, 0, maxCodeUpdatesPerSecond),
	}

	// Simulate 10 updates from 2 seconds ago (should be expired)
	twoSecondsAgo := time.Now().Add(-2 * time.Second)
	for i := 0; i < maxCodeUpdatesPerSecond; i++ {
		client.codeUpdateTimestamps = append(client.codeUpdateTimestamps, twoSecondsAgo)
	}

	// Next update should pass because old timestamps are expired
	if !client.checkCodeUpdateRateLimit() {
		t.Error("Code update should have been allowed after old timestamps expired")
	}

	// Old timestamps should be cleaned up
	if len(client.codeUpdateTimestamps) != 1 {
		t.Errorf("Expected 1 timestamp after cleanup, got %d", len(client.codeUpdateTimestamps))
	}
}

// Test chat message rate limiting (20/minute)
func TestChatRateLimit(t *testing.T) {
	client := &Client{
		chatMessageTimestamps: make([]time.Time, 0, maxChatMessagesPerMinute),
	}

	// First 20 messages should pass
	for i := 0; i < maxChatMessagesPerMinute; i++ {
		if !client.checkChatRateLimit() {
			t.Errorf("Chat message %d should have been allowed, but was rate limited", i+1)
		}
	}

	// 21st message should be rate limited
	if client.checkChatRateLimit() {
		t.Error("21st chat message should have been rate limited, but was allowed")
	}

	if len(client.chatMessageTimestamps) != maxChatMessagesPerMinute {
		t.Errorf("Expected %d timestamps, got %d", maxChatMessagesPerMinute, len(client.chatMessageTimestamps))
	}
}

// Test chat rate limit window expiration
func TestChatRateLimitWindowExpiration(t *testing.T) {
	client := &Client{
		chatMessageTimestamps: make([]time.Time, 0, maxChatMessagesPerMinute),
	}

	// Simulate 20 messages from 2 minutes ago (should be expired)
	twoMinutesAgo := time.Now().Add(-2 * time.Minute)
	for i := 0; i < maxChatMessagesPerMinute; i++ {
		client.chatMessageTimestamps = append(client.chatMessageTimestamps, twoMinutesAgo)
	}

	// Next message should pass because old timestamps are expired
	if !client.checkChatRateLimit() {
		t.Error("Chat message should have been allowed after old timestamps expired")
	}

	// Old timestamps should be cleaned up
	if len(client.chatMessageTimestamps) != 1 {
		t.Errorf("Expected 1 timestamp after cleanup, got %d", len(client.chatMessageTimestamps))
	}
}

// Test that default tier falls back to free tier limits
func TestAgentRequestRateLimitDefaultTier(t *testing.T) {
	client := &Client{
		Tier:                   "", // empty tier should default to free
		agentRequestTimestamps: make([]time.Time, 0, 10),
	}

	// Should use default (free) limit of 10
	for i := 0; i < 10; i++ {
		if !client.checkAgentRequestRateLimit() {
			t.Errorf("Request %d should have been allowed with default tier", i+1)
		}
	}

	// 11th should be limited
	if client.checkAgentRequestRateLimit() {
		t.Error("11th request should have been rate limited with default tier")
	}
}

// Test getAgentRequestLimit returns correct values
func TestGetAgentRequestLimit(t *testing.T) {
	tests := []struct {
		tier     string
		expected int
	}{
		{"free", 10},
		{"pro", 20},
		{"byok", 30},
		{"", 10},        // default
		{"unknown", 10}, // unknown tier defaults to free
	}

	for _, tt := range tests {
		client := &Client{Tier: tt.tier}
		got := client.getAgentRequestLimit()
		if got != tt.expected {
			t.Errorf("getAgentRequestLimit() for tier %q = %d, want %d", tt.tier, got, tt.expected)
		}
	}
}

// Test CanWrite permission check
func TestCanWrite(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"host", true},
		{"co-author", true},
		{"viewer", false},
		{"", false},
	}

	for _, tt := range tests {
		client := &Client{Role: tt.role}
		got := client.CanWrite()
		if got != tt.expected {
			t.Errorf("CanWrite() for role %q = %v, want %v", tt.role, got, tt.expected)
		}
	}
}
