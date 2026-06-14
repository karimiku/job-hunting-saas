package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOriginGuard_AllowsAllowedOrigin(t *testing.T) {
	handler := NewOriginGuard([]string{"https://entre.example"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Origin", "https://entre.example")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestOriginGuard_AllowsAllowedReferer(t *testing.T) {
	handler := NewOriginGuard([]string{"https://entre.example"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Referer", "https://entre.example/dashboard")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestOriginGuard_RejectsMissingOriginForUnsafeMethod(t *testing.T) {
	handler := NewOriginGuard([]string{"https://entre.example"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestSessionCSRFProtection_RejectsSessionUnsafeWithoutOrigin(t *testing.T) {
	handler := NewSessionCSRFProtection([]string{"https://entre.example"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	req := httptest.NewRequest(http.MethodDelete, "/", nil).WithContext(
		SetAuthMethod(context.Background(), AuthMethodSession),
	)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestSessionCSRFProtection_SkipsBearerAndSafeMethods(t *testing.T) {
	for _, tt := range []struct {
		name       string
		method     string
		methodAuth AuthMethod
	}{
		{name: "bearer unsafe", method: http.MethodPost, methodAuth: AuthMethodBearer},
		{name: "session safe", method: http.MethodGet, methodAuth: AuthMethodSession},
	} {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewSessionCSRFProtection([]string{"https://entre.example"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			req := httptest.NewRequest(tt.method, "/", nil).WithContext(
				SetAuthMethod(context.Background(), tt.methodAuth),
			)

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			if w.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
			}
		})
	}
}
