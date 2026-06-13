package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
)

func TestCreateAIAccessTokenRequiresSessionAuth(t *testing.T) {
	userID := entity.NewUserID()
	h := NewAIAccessTokenHandler(&fakeAIAccessTokenHandlerRepo{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai-access-tokens", strings.NewReader(`{}`))
	req = req.WithContext(middleware.SetAuth(req.Context(), userID, middleware.AuthMethodAIAccessToken))
	w := httptest.NewRecorder()

	h.CreateAIAccessToken(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestCreateAIAccessTokenSavesHashAndReturnsPlainTokenOnce(t *testing.T) {
	userID := entity.NewUserID()
	repo := &fakeAIAccessTokenHandlerRepo{}
	h := NewAIAccessTokenHandler(repo)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ai-access-tokens", strings.NewReader(`{"name":"自分用"}`))
	req = req.WithContext(middleware.SetAuth(req.Context(), userID, middleware.AuthMethodSession))
	w := httptest.NewRecorder()

	h.CreateAIAccessToken(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", w.Code, w.Body.String())
	}
	if repo.saved == nil {
		t.Fatal("token was not saved")
	}
	if repo.saved.UserID() != userID {
		t.Errorf("saved userID = %s, want %s", repo.saved.UserID(), userID)
	}
	if repo.saved.Name() != "自分用" {
		t.Errorf("saved name = %q, want 自分用", repo.saved.Name())
	}

	var body struct {
		Token string `json:"token"`
		Item  struct {
			TokenPreview string `json:"tokenPreview"`
		} `json:"item"`
	}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Token == "" || !strings.HasPrefix(body.Token, "entre_ai_") {
		t.Fatalf("token = %q, want entre_ai_ prefix", body.Token)
	}
	if repo.saved.TokenHash() == body.Token {
		t.Error("saved token hash must not be plain token")
	}
	if body.Item.TokenPreview == body.Token {
		t.Error("token preview must not expose the full token")
	}
}

type fakeAIAccessTokenHandlerRepo struct {
	saved *entity.AIAccessToken
	items []*entity.AIAccessToken
}

func (r *fakeAIAccessTokenHandlerRepo) Save(_ context.Context, token *entity.AIAccessToken) error {
	r.saved = token
	r.items = append([]*entity.AIAccessToken{token}, r.items...)
	return nil
}

func (r *fakeAIAccessTokenHandlerRepo) FindActiveByHash(context.Context, string) (*entity.AIAccessToken, error) {
	return nil, repository.ErrNotFound
}

func (r *fakeAIAccessTokenHandlerRepo) ListByUserID(context.Context, entity.UserID) ([]*entity.AIAccessToken, error) {
	return r.items, nil
}

func (r *fakeAIAccessTokenHandlerRepo) TouchLastUsed(context.Context, entity.AIAccessTokenID, time.Time) error {
	return nil
}

func (r *fakeAIAccessTokenHandlerRepo) Revoke(_ context.Context, _ entity.UserID, id entity.AIAccessTokenID, at time.Time) error {
	for _, token := range r.items {
		if token.ID() == id {
			token.Revoke(at)
			return nil
		}
	}
	return repository.ErrNotFound
}
