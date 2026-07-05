package supabaseauth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

const testIssuer = "https://project-ref.supabase.co/auth/v1"

func TestVerifierVerifyBearerTokenResolvesSupabaseIdentity(t *testing.T) {
	key := newRSAKey(t)
	jwksServer := newJWKSServer(t, key, "test-key")
	defer jwksServer.Close()

	ctx := context.Background()
	userID := entity.NewUserID()
	identityRepo := inmemory.NewExternalIdentityRepository()
	if err := identityRepo.Save(ctx, entity.NewExternalIdentity(userID, value.AuthProviderSupabase(), "supabase-user-id")); err != nil {
		t.Fatalf("Save identity failed: %v", err)
	}

	verifier := newTestVerifier(t, jwksServer.URL, identityRepo)
	rawToken := signSupabaseJWT(t, key, "test-key", testIssuer, defaultAudience, "supabase-user-id")

	got, err := verifier.VerifyBearerToken(ctx, rawToken)
	if err != nil {
		t.Fatalf("VerifyBearerToken failed: %v", err)
	}
	if got != userID {
		t.Fatalf("UserID = %v, want %v", got, userID)
	}
}

func TestVerifierVerifyBearerTokenSyncsUnknownSupabaseUser(t *testing.T) {
	key := newRSAKey(t)
	jwksServer := newJWKSServer(t, key, "test-key")
	defer jwksServer.Close()

	ctx := context.Background()
	identityRepo := inmemory.NewExternalIdentityRepository()
	syncedUserID := entity.NewUserID()
	var syncedInfo UserInfo
	verifier := newTestVerifierWithSync(t, jwksServer.URL, identityRepo, func(ctx context.Context, info UserInfo) (entity.UserID, error) {
		syncedInfo = info
		if err := identityRepo.Save(ctx, entity.NewExternalIdentity(syncedUserID, value.AuthProviderSupabase(), info.Subject)); err != nil {
			return entity.UserID{}, err
		}
		return syncedUserID, nil
	})
	rawToken := signSupabaseJWTWithUser(t, key, "test-key", testIssuer, defaultAudience, "supabase-user-id", "student@example.com", "学生 太郎")

	got, err := verifier.VerifyBearerToken(ctx, rawToken)
	if err != nil {
		t.Fatalf("VerifyBearerToken failed: %v", err)
	}
	if got != syncedUserID {
		t.Fatalf("UserID = %v, want %v", got, syncedUserID)
	}
	if syncedInfo.Subject != "supabase-user-id" || syncedInfo.Email != "student@example.com" || syncedInfo.Name != "学生 太郎" {
		t.Fatalf("syncedInfo = %+v", syncedInfo)
	}
}

func TestVerifierVerifyBearerTokenRejectsInvalidClaims(t *testing.T) {
	key := newRSAKey(t)
	jwksServer := newJWKSServer(t, key, "test-key")
	defer jwksServer.Close()

	verifier := newTestVerifier(t, jwksServer.URL, inmemory.NewExternalIdentityRepository())
	rawToken := signSupabaseJWT(t, key, "test-key", "https://evil.example.com/auth/v1", defaultAudience, "supabase-user-id")

	_, err := verifier.VerifyBearerToken(context.Background(), rawToken)
	if !errors.Is(err, value.ErrAuthTokenInvalid) {
		t.Fatalf("error = %v, want ErrAuthTokenInvalid", err)
	}
}

func TestVerifierVerifyBearerTokenRejectsUnknownIdentity(t *testing.T) {
	key := newRSAKey(t)
	jwksServer := newJWKSServer(t, key, "test-key")
	defer jwksServer.Close()

	verifier := newTestVerifier(t, jwksServer.URL, inmemory.NewExternalIdentityRepository())
	rawToken := signSupabaseJWT(t, key, "test-key", testIssuer, defaultAudience, "unknown-user-id")

	_, err := verifier.VerifyBearerToken(context.Background(), rawToken)
	if !errors.Is(err, value.ErrAuthTokenInvalid) {
		t.Fatalf("error = %v, want ErrAuthTokenInvalid", err)
	}
}

func TestVerifierVerifyBearerTokenReturnsInfrastructureErrorForJWKSFailure(t *testing.T) {
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer jwksServer.Close()

	key := newRSAKey(t)
	verifier := newTestVerifier(t, jwksServer.URL, inmemory.NewExternalIdentityRepository())
	rawToken := signSupabaseJWT(t, key, "test-key", testIssuer, defaultAudience, "supabase-user-id")

	_, err := verifier.VerifyBearerToken(context.Background(), rawToken)
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Is(err, value.ErrAuthTokenInvalid) {
		t.Fatalf("error = %v, want infrastructure error", err)
	}
}

func TestVerifierUnknownKidRefreshIsRateLimited(t *testing.T) {
	key := newRSAKey(t)

	var fetches int32
	keySet := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{{
			Key:       &key.PublicKey,
			KeyID:     "server-key",
			Algorithm: jwt.SigningMethodRS256.Alg(),
			Use:       "sig",
		}},
	}
	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&fetches, 1)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(keySet); err != nil {
			t.Fatalf("Encode JWKS failed: %v", err)
		}
	}))
	defer jwksServer.Close()

	verifier := newTestVerifier(t, jwksServer.URL, inmemory.NewExternalIdentityRepository())

	// Every token carries an unknown kid, so each request misses the cache and
	// would fetch JWKS if there were no cooldown.
	for i := 0; i < 5; i++ {
		rawToken := signSupabaseJWT(t, key, "unknown-kid", testIssuer, defaultAudience, "supabase-user-id")
		if _, err := verifier.VerifyBearerToken(context.Background(), rawToken); !errors.Is(err, value.ErrAuthTokenInvalid) {
			t.Fatalf("request %d error = %v, want ErrAuthTokenInvalid", i, err)
		}
	}

	if got := atomic.LoadInt32(&fetches); got != 1 {
		t.Fatalf("JWKS fetches = %d, want 1 (cooldown should suppress repeated unknown-kid refreshes)", got)
	}
}

func newTestVerifier(t *testing.T, jwksURL string, identityRepo *inmemory.ExternalIdentityRepository) *Verifier {
	t.Helper()
	return newTestVerifierWithSync(t, jwksURL, identityRepo, nil)
}

func newTestVerifierWithSync(
	t *testing.T,
	jwksURL string,
	identityRepo *inmemory.ExternalIdentityRepository,
	userSync UserSyncFunc,
) *Verifier {
	t.Helper()
	verifier, err := NewVerifier(Config{
		Issuer:       testIssuer,
		Audience:     defaultAudience,
		JWKSURL:      jwksURL,
		JWKSCacheTTL: time.Minute,
		UserSync:     userSync,
	}, identityRepo)
	if err != nil {
		t.Fatalf("NewVerifier failed: %v", err)
	}
	return verifier
}

func newRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}
	return key
}

func newJWKSServer(t *testing.T, key *rsa.PrivateKey, keyID string) *httptest.Server {
	t.Helper()
	keySet := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:       &key.PublicKey,
				KeyID:     keyID,
				Algorithm: jwt.SigningMethodRS256.Alg(),
				Use:       "sig",
			},
		},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(keySet); err != nil {
			t.Fatalf("Encode JWKS failed: %v", err)
		}
	}))
}

func signSupabaseJWT(t *testing.T, key *rsa.PrivateKey, keyID, issuer, audience, subject string) string {
	t.Helper()
	return signSupabaseJWTWithUser(t, key, keyID, issuer, audience, subject, "", "")
}

func signSupabaseJWTWithUser(
	t *testing.T,
	key *rsa.PrivateKey,
	keyID string,
	issuer string,
	audience string,
	subject string,
	email string,
	name string,
) string {
	t.Helper()
	claims := jwt.MapClaims{
		"iss": issuer,
		"sub": subject,
		"aud": audience,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Add(-time.Minute).Unix(),
	}
	if email != "" {
		claims["email"] = email
	}
	if name != "" {
		claims["user_metadata"] = map[string]any{"name": name}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID
	rawToken, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("SignedString failed: %v", err)
	}
	return rawToken
}
