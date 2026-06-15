package middleware

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestCachedSessionVerifierCachesWithinTTL(t *testing.T) {
	calls := 0
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	verifier := &cachedSessionVerifier{
		delegate: FirebaseSessionVerifierFunc(func(context.Context, string) (*SessionClaims, error) {
			calls++
			return &SessionClaims{UID: "firebase-uid"}, nil
		}),
		ttl:     time.Minute,
		now:     func() time.Time { return now },
		entries: make(map[string]cachedSessionClaims),
	}

	first, err := verifier.VerifySessionCookie(context.Background(), "session-cookie")
	if err != nil {
		t.Fatalf("first VerifySessionCookie: %v", err)
	}
	second, err := verifier.VerifySessionCookie(context.Background(), "session-cookie")
	if err != nil {
		t.Fatalf("second VerifySessionCookie: %v", err)
	}

	if first.UID != "firebase-uid" || second.UID != "firebase-uid" {
		t.Fatalf("UIDs = %q, %q; want firebase-uid", first.UID, second.UID)
	}
	if calls != 1 {
		t.Fatalf("delegate calls = %d, want 1", calls)
	}
}

func TestCachedSessionVerifierExpires(t *testing.T) {
	calls := 0
	now := time.Date(2026, 6, 15, 12, 0, 0, 0, time.UTC)
	verifier := &cachedSessionVerifier{
		delegate: FirebaseSessionVerifierFunc(func(context.Context, string) (*SessionClaims, error) {
			calls++
			return &SessionClaims{UID: "firebase-uid"}, nil
		}),
		ttl:     time.Second,
		now:     func() time.Time { return now },
		entries: make(map[string]cachedSessionClaims),
	}

	if _, err := verifier.VerifySessionCookie(context.Background(), "session-cookie"); err != nil {
		t.Fatalf("first VerifySessionCookie: %v", err)
	}
	now = now.Add(2 * time.Second)
	if _, err := verifier.VerifySessionCookie(context.Background(), "session-cookie"); err != nil {
		t.Fatalf("second VerifySessionCookie: %v", err)
	}

	if calls != 2 {
		t.Fatalf("delegate calls = %d, want 2", calls)
	}
}

func TestCachedSessionVerifierSingleflightsConcurrentMisses(t *testing.T) {
	var calls int
	started := make(chan struct{})
	release := make(chan struct{})
	verifier := NewCachedSessionVerifier(FirebaseSessionVerifierFunc(func(context.Context, string) (*SessionClaims, error) {
		calls++
		close(started)
		<-release
		return &SessionClaims{UID: "firebase-uid"}, nil
	}), time.Minute)

	var wg sync.WaitGroup
	wg.Add(2)
	for range 2 {
		go func() {
			defer wg.Done()
			claims, err := verifier.VerifySessionCookie(context.Background(), "session-cookie")
			if err != nil {
				t.Errorf("VerifySessionCookie: %v", err)
				return
			}
			if claims.UID != "firebase-uid" {
				t.Errorf("UID = %q, want firebase-uid", claims.UID)
			}
		}()
	}

	<-started
	close(release)
	wg.Wait()

	if calls != 1 {
		t.Fatalf("delegate calls = %d, want 1", calls)
	}
}

type FirebaseSessionVerifierFunc func(context.Context, string) (*SessionClaims, error)

func (f FirebaseSessionVerifierFunc) VerifySessionCookie(ctx context.Context, sessionCookie string) (*SessionClaims, error) {
	return f(ctx, sessionCookie)
}
