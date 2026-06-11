package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

type mockBearerVerifier struct {
	userID entity.UserID
	err    error
	token  string
}

func (m *mockBearerVerifier) VerifyBearerToken(_ context.Context, rawToken string) (entity.UserID, error) {
	m.token = rawToken
	if m.err != nil {
		return entity.UserID{}, m.err
	}
	return m.userID, nil
}

func TestNewAuthWithBearer_Success(t *testing.T) {
	userID := entity.NewUserID()
	bearer := &mockBearerVerifier{userID: userID}
	called := false
	handler := NewAuthWithBearer(&mockSessionVerifier{}, inmemory.NewExternalIdentityRepository(), bearer)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			if got := GetUserID(r.Context()); got != userID {
				t.Errorf("GetUserID = %v, want %v", got, userID)
			}
			if got := GetAuthMethod(r.Context()); got != AuthMethodBearer {
				t.Errorf("GetAuthMethod = %q, want %q", got, AuthMethodBearer)
			}
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer entre_ai_test")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if !called {
		t.Fatal("next handler should be called")
	}
	if bearer.token != "entre_ai_test" {
		t.Errorf("token = %q, want entre_ai_test", bearer.token)
	}
}

func TestNewAuthWithBearer_InvalidToken(t *testing.T) {
	bearer := &mockBearerVerifier{err: value.ErrAIAccessTokenInvalid}
	called := false
	handler := NewAuthWithBearer(&mockSessionVerifier{}, inmemory.NewExternalIdentityRepository(), bearer)(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true }),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if called {
		t.Fatal("next handler should not be called")
	}
}

func TestNewAuthWithBearer_UnexpectedVerifierError(t *testing.T) {
	bearer := &mockBearerVerifier{err: errors.New("db down")}
	handler := NewAuthWithBearer(&mockSessionVerifier{}, inmemory.NewExternalIdentityRepository(), bearer)(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer entre_ai_test")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}
