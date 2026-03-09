package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
)

func TestWriteError_NotFound(t *testing.T) {
	w := httptest.NewRecorder()

	writeError(w, repository.ErrNotFound)

	if w.Code != 404 {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestWriteError_OtherError(t *testing.T) {
	w := httptest.NewRecorder()

	writeError(w, repository.ErrNotFound)

	// 他のエラーは400になることを確認
	w2 := httptest.NewRecorder()
	writeError(w2, java_err("validation failed"))

	if w2.Code != 400 {
		t.Errorf("status = %d, want 400", w2.Code)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	writeJSON(w, 201, openapi.ErrorResponse{Message: "created"})

	if w.Code != 201 {
		t.Errorf("status = %d, want 201", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}

type java_err string

func (e java_err) Error() string { return string(e) }
