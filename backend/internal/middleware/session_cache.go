package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const maxCachedSessionCookies = 1024

type cachedSessionClaims struct {
	claims    SessionClaims
	expiresAt time.Time
}

type cachedSessionVerifier struct {
	delegate FirebaseSessionVerifier
	ttl      time.Duration
	now      func() time.Time

	mu      sync.Mutex
	entries map[string]cachedSessionClaims
	group   singleflight.Group
}

// NewCachedSessionVerifier wraps Firebase session verification with a short
// in-process cache. It keeps revocation checks reasonably fresh while avoiding
// repeated Firebase Admin SDK calls during fast RSC route transitions.
func NewCachedSessionVerifier(delegate FirebaseSessionVerifier, ttl time.Duration) FirebaseSessionVerifier {
	if delegate == nil || ttl <= 0 {
		return delegate
	}
	return &cachedSessionVerifier{
		delegate: delegate,
		ttl:      ttl,
		now:      time.Now,
		entries:  make(map[string]cachedSessionClaims),
	}
}

func (v *cachedSessionVerifier) VerifySessionCookie(ctx context.Context, sessionCookie string) (*SessionClaims, error) {
	key := sessionCacheKey(sessionCookie)
	if claims, ok := v.get(key); ok {
		return claims, nil
	}

	value, err, _ := v.group.Do(key, func() (any, error) {
		if claims, ok := v.get(key); ok {
			return claims, nil
		}
		claims, err := v.delegate.VerifySessionCookie(ctx, sessionCookie)
		if err != nil {
			return nil, err
		}
		v.set(key, claims)
		return cloneSessionClaims(claims), nil
	})
	if err != nil {
		return nil, err
	}
	return value.(*SessionClaims), nil
}

func (v *cachedSessionVerifier) get(key string) (*SessionClaims, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	entry, ok := v.entries[key]
	if !ok {
		return nil, false
	}
	if !v.now().Before(entry.expiresAt) {
		delete(v.entries, key)
		return nil, false
	}
	return cloneSessionClaims(&entry.claims), true
}

func (v *cachedSessionVerifier) set(key string, claims *SessionClaims) {
	if claims == nil {
		return
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	v.pruneLocked()
	v.entries[key] = cachedSessionClaims{
		claims:    *claims,
		expiresAt: v.now().Add(v.ttl),
	}
}

func (v *cachedSessionVerifier) pruneLocked() {
	if len(v.entries) < maxCachedSessionCookies {
		return
	}
	now := v.now()
	for key, entry := range v.entries {
		if !now.Before(entry.expiresAt) {
			delete(v.entries, key)
		}
	}
	for key := range v.entries {
		if len(v.entries) < maxCachedSessionCookies {
			return
		}
		delete(v.entries, key)
	}
}

func cloneSessionClaims(claims *SessionClaims) *SessionClaims {
	if claims == nil {
		return nil
	}
	cloned := *claims
	return &cloned
}

func sessionCacheKey(sessionCookie string) string {
	sum := sha256.Sum256([]byte(sessionCookie))
	return hex.EncodeToString(sum[:])
}
