package oauth2_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-oauth2/oauth2/v4"
)

func TestMemoryAssertionReplayStore(t *testing.T) {
	t.Run("first store succeeds", func(t *testing.T) {
		s := oauth2.NewMemoryAssertionReplayStore()
		if err := s.StoreAssertionID(context.Background(), "jti-1", time.Now().Add(5*time.Minute)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("duplicate before exp returns error", func(t *testing.T) {
		s := oauth2.NewMemoryAssertionReplayStore()
		exp := time.Now().Add(5 * time.Minute)
		if err := s.StoreAssertionID(context.Background(), "jti-dup", exp); err != nil {
			t.Fatalf("first store: %v", err)
		}
		if err := s.StoreAssertionID(context.Background(), "jti-dup", exp); err == nil {
			t.Error("expected error for duplicate jti, got nil")
		}
	})

	t.Run("same jti accepted after exp", func(t *testing.T) {
		s := oauth2.NewMemoryAssertionReplayStore()
		// Store with an already-expired exp so the sweep removes it on next call.
		pastExp := time.Now().Add(-1 * time.Millisecond)
		if err := s.StoreAssertionID(context.Background(), "jti-expired", pastExp); err != nil {
			t.Fatalf("initial store: %v", err)
		}
		// A tiny sleep ensures the exp timestamp is in the past when StoreAssertionID runs.
		time.Sleep(5 * time.Millisecond)
		if err := s.StoreAssertionID(context.Background(), "jti-expired", time.Now().Add(5*time.Minute)); err != nil {
			t.Errorf("after expiry: %v", err)
		}
	})

	t.Run("cancelled context returns error", func(t *testing.T) {
		s := oauth2.NewMemoryAssertionReplayStore()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := s.StoreAssertionID(ctx, "jti-cancel", time.Now().Add(5*time.Minute)); err == nil {
			t.Error("expected error from cancelled context, got nil")
		}
	})
}
