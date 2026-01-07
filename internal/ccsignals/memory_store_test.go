package ccsignals

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryLockStore_SetLock(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	err := store.SetLock(ctx, "session1", "baseline code", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state, err := store.GetLock(ctx, "session1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !state.Locked {
		t.Error("expected session to be locked")
	}
	if state.BaselineCode != "baseline code" {
		t.Errorf("expected baseline 'baseline code', got %q", state.BaselineCode)
	}
}

func TestMemoryLockStore_GetLock_NotExists(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	state, err := store.GetLock(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.Locked {
		t.Error("expected session to be unlocked")
	}
}

func TestMemoryLockStore_GetLock_Expired(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	// set lock with very short TTL
	err := store.SetLock(ctx, "session1", "baseline", 1*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// wait for expiration
	time.Sleep(5 * time.Millisecond)

	state, err := store.GetLock(ctx, "session1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.Locked {
		t.Error("expected expired lock to return unlocked")
	}
}

func TestMemoryLockStore_RemoveLock(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	if err := store.SetLock(ctx, "session1", "baseline", time.Hour); err != nil {
		t.Fatalf("failed to set lock: %v", err)
	}

	err := store.RemoveLock(ctx, "session1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	state, err := store.GetLock(ctx, "session1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Locked {
		t.Error("expected session to be unlocked after removal")
	}
}

func TestMemoryLockStore_RemoveLock_NotExists(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	// should not error when removing non-existent lock
	err := store.RemoveLock(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMemoryLockStore_RefreshTTL(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	// set lock with short TTL
	if err := store.SetLock(ctx, "session1", "baseline", 10*time.Millisecond); err != nil {
		t.Fatalf("failed to set lock: %v", err)
	}

	// refresh with longer TTL
	err := store.RefreshTTL(ctx, "session1", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// wait past original TTL
	time.Sleep(20 * time.Millisecond)

	state, err := store.GetLock(ctx, "session1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !state.Locked {
		t.Error("expected lock to still be valid after TTL refresh")
	}
}

func TestMemoryLockStore_RefreshTTL_NotExists(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	// should not error when refreshing non-existent lock
	err := store.RefreshTTL(ctx, "nonexistent", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMemoryLockStore_Close(t *testing.T) {
	store := NewMemoryLockStore()

	err := store.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// second close should not panic
	err = store.Close()
	if err != nil {
		t.Fatalf("unexpected error on second close: %v", err)
	}
}

func TestMemoryLockStore_Concurrent(_ *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	var wg sync.WaitGroup
	sessions := 100

	// concurrent writes - stress test, errors not checked intentionally
	for i := 0; i < sessions; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sessionID := string(rune('a' + id%26))
			_ = store.SetLock(ctx, sessionID, "baseline", time.Hour) //nolint:errcheck // stress test
			_, _ = store.GetLock(ctx, sessionID)                     //nolint:errcheck // stress test
			_ = store.RefreshTTL(ctx, sessionID, time.Hour)          //nolint:errcheck // stress test
		}(i)
	}

	wg.Wait()
}

func TestMemoryLockStore_LockedAt(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryLockStore()
	defer func() { _ = store.Close() }() //nolint:errcheck // test cleanup

	before := time.Now()
	if err := store.SetLock(ctx, "session1", "baseline", time.Hour); err != nil {
		t.Fatalf("failed to set lock: %v", err)
	}
	after := time.Now()

	state, err := store.GetLock(ctx, "session1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if state.LockedAt.Before(before) || state.LockedAt.After(after) {
		t.Errorf("LockedAt %v should be between %v and %v", state.LockedAt, before, after)
	}
}
