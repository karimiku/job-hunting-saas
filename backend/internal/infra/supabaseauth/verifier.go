// Package supabaseauth verifies Supabase Auth access tokens for the backend.
package supabaseauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

const (
	defaultAudience     = "authenticated"
	defaultJWKSCacheTTL = 10 * time.Minute
	maxJWKSCacheTTL     = 10 * time.Minute
	defaultHTTPTimeout  = 5 * time.Second
)

var (
	errInvalidJWT = errors.New("invalid supabase jwt")

	supportedJWTAlgorithms = []string{
		jwt.SigningMethodRS256.Alg(),
		jwt.SigningMethodES256.Alg(),
		jwt.SigningMethodEdDSA.Alg(),
	}
)

// Config contains the Supabase Auth JWT verification settings.
type Config struct {
	Issuer       string
	Audience     string
	JWKSURL      string
	JWKSCacheTTL time.Duration
	HTTPClient   *http.Client
}

// Verifier resolves a valid Supabase access token to this app's UserID.
type Verifier struct {
	issuer       string
	audience     string
	jwksURL      string
	cacheTTL     time.Duration
	httpClient   *http.Client
	extIDRepo    repository.ExternalIdentityRepository
	mu           sync.RWMutex
	cachedJWKS   jose.JSONWebKeySet
	jwksExpireAt time.Time
}

// NewVerifier creates a Supabase JWT verifier.
func NewVerifier(cfg Config, extIDRepo repository.ExternalIdentityRepository) (*Verifier, error) {
	if extIDRepo == nil {
		return nil, errors.New("supabaseauth: external identity repository is required")
	}

	issuer := strings.TrimRight(strings.TrimSpace(cfg.Issuer), "/")
	if issuer == "" {
		return nil, errors.New("supabaseauth: issuer is required")
	}
	jwksURL := strings.TrimSpace(cfg.JWKSURL)
	if jwksURL == "" {
		jwksURL = issuer + "/.well-known/jwks.json"
	}
	if _, err := url.ParseRequestURI(jwksURL); err != nil {
		return nil, fmt.Errorf("supabaseauth: invalid JWKS URL: %w", err)
	}

	audience := strings.TrimSpace(cfg.Audience)
	if audience == "" {
		audience = defaultAudience
	}
	cacheTTL := cfg.JWKSCacheTTL
	if cacheTTL == 0 {
		cacheTTL = defaultJWKSCacheTTL
	}
	if cacheTTL < 0 || cacheTTL > maxJWKSCacheTTL {
		return nil, fmt.Errorf("supabaseauth: JWKS cache TTL must be between 0 and %s", maxJWKSCacheTTL)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	}

	return &Verifier{
		issuer:     issuer,
		audience:   audience,
		jwksURL:    jwksURL,
		cacheTTL:   cacheTTL,
		httpClient: httpClient,
		extIDRepo:  extIDRepo,
	}, nil
}

// VerifyBearerToken implements middleware.BearerTokenVerifier.
func (v *Verifier) VerifyBearerToken(ctx context.Context, rawToken string) (entity.UserID, error) {
	claims, err := v.verifyJWT(ctx, strings.TrimSpace(rawToken))
	if err != nil {
		if errors.Is(err, errInvalidJWT) {
			return entity.UserID{}, value.ErrAuthTokenInvalid
		}
		return entity.UserID{}, err
	}

	identity, err := v.extIDRepo.FindByProviderAndSubject(ctx, value.AuthProviderSupabase(), claims.Subject)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return entity.UserID{}, value.ErrAuthTokenInvalid
		}
		return entity.UserID{}, err
	}

	return identity.UserID(), nil
}

func (v *Verifier) verifyJWT(ctx context.Context, rawToken string) (*jwt.RegisteredClaims, error) {
	if rawToken == "" {
		return nil, errInvalidJWT
	}

	claims := &jwt.RegisteredClaims{}
	parser := jwt.Parser{ValidMethods: supportedJWTAlgorithms}
	token, err := parser.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
		return v.publicKey(ctx, token)
	})
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorUnverifiable != 0 &&
				validationErr.Inner != nil &&
				!errors.Is(validationErr.Inner, errInvalidJWT) {
				return nil, fmt.Errorf("supabaseauth: verify JWT: %w", validationErr.Inner)
			}
			return nil, errInvalidJWT
		}
		return nil, errInvalidJWT
	}
	if token == nil || !token.Valid {
		return nil, errInvalidJWT
	}
	if claims.Subject == "" ||
		claims.ExpiresAt == nil ||
		!claims.VerifyIssuer(v.issuer, true) ||
		!claims.VerifyAudience(v.audience, true) {
		return nil, errInvalidJWT
	}

	return claims, nil
}

func (v *Verifier) publicKey(ctx context.Context, token *jwt.Token) (interface{}, error) {
	if token == nil || token.Method == nil {
		return nil, errInvalidJWT
	}
	alg := token.Method.Alg()
	if !isSupportedJWTAlgorithm(alg) {
		return nil, errInvalidJWT
	}

	kid, ok := token.Header["kid"].(string)
	if !ok || strings.TrimSpace(kid) == "" {
		return nil, errInvalidJWT
	}

	if key, ok := v.cachedKey(kid, alg); ok {
		return key, nil
	}
	if err := v.refreshJWKS(ctx); err != nil {
		return nil, err
	}
	if key, ok := v.cachedKey(kid, alg); ok {
		return key, nil
	}
	return nil, errInvalidJWT
}

func (v *Verifier) cachedKey(kid, alg string) (interface{}, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if time.Now().After(v.jwksExpireAt) {
		return nil, false
	}
	return selectJWK(v.cachedJWKS, kid, alg)
}

func (v *Verifier) refreshJWKS(ctx context.Context) error {
	keySet, err := v.fetchJWKS(ctx)
	if err != nil {
		return err
	}

	v.mu.Lock()
	v.cachedJWKS = keySet
	v.jwksExpireAt = time.Now().Add(v.cacheTTL)
	v.mu.Unlock()
	return nil
}

func (v *Verifier) fetchJWKS(ctx context.Context) (jose.JSONWebKeySet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("supabaseauth: create JWKS request: %w", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("supabaseauth: fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return jose.JSONWebKeySet{}, fmt.Errorf("supabaseauth: fetch JWKS: unexpected status %s", resp.Status)
	}

	var keySet jose.JSONWebKeySet
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&keySet); err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("supabaseauth: decode JWKS: %w", err)
	}
	if len(keySet.Keys) == 0 {
		return jose.JSONWebKeySet{}, errors.New("supabaseauth: JWKS contains no keys")
	}
	return keySet, nil
}

func selectJWK(keySet jose.JSONWebKeySet, kid, alg string) (interface{}, bool) {
	for _, key := range keySet.Key(kid) {
		if key.Algorithm != "" && key.Algorithm != alg {
			continue
		}
		if !key.Valid() {
			continue
		}
		return key.Key, true
	}
	return nil, false
}

func isSupportedJWTAlgorithm(alg string) bool {
	for _, supported := range supportedJWTAlgorithms {
		if alg == supported {
			return true
		}
	}
	return false
}
