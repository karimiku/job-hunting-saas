package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	aiaccesstokenuc "github.com/karimiku/job-hunting-saas/internal/usecase/ai_access_token"
)

func setupAiAccessTokenHandler() (*AiAccessTokenHandler, *inmemory.AIAccessTokenRepository) {
	repo := inmemory.NewAIAccessTokenRepository()
	return NewAiAccessTokenHandler(
		aiaccesstokenuc.NewCreate(repo),
		aiaccesstokenuc.NewList(repo),
		aiaccesstokenuc.NewRevoke(repo),
	), repo
}

func TestCreateAiAccessToken_ReturnsRawTokenOnce(t *testing.T) {
	h, _ := setupAiAccessTokenHandler()
	userID := entity.NewUserID()
	body, _ := json.Marshal(openapi.CreateAiAccessTokenRequest{Name: "Claude Desktop"})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()
	h.CreateAiAccessToken(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var resp openapi.CreateAiAccessTokenResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("Token is empty")
	}
	if resp.AccessToken.TokenPrefix == "" {
		t.Fatal("TokenPrefix is empty")
	}
}

func TestListAndRevokeAiAccessToken(t *testing.T) {
	h, repo := setupAiAccessTokenHandler()
	userID := entity.NewUserID()
	created, err := aiaccesstokenuc.NewCreate(repo).Execute(
		httptest.NewRequest(http.MethodPost, "/", nil).Context(),
		aiaccesstokenuc.CreateInput{UserID: userID, Name: "Codex"},
	)
	if err != nil {
		t.Fatalf("seed token: %v", err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/", nil)
	listReq = listReq.WithContext(middleware.SetUserID(listReq.Context(), userID))
	listW := httptest.NewRecorder()
	h.ListAiAccessTokens(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("list status = %d", listW.Code)
	}
	var listResp struct {
		Tokens []openapi.AiAccessTokenResponse `json:"tokens"`
	}
	if err := json.NewDecoder(listW.Body).Decode(&listResp); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listResp.Tokens) != 1 {
		t.Fatalf("len(tokens) = %d, want 1", len(listResp.Tokens))
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/", nil)
	deleteReq = deleteReq.WithContext(middleware.SetUserID(deleteReq.Context(), userID))
	deleteW := httptest.NewRecorder()
	h.RevokeAiAccessToken(deleteW, deleteReq, openapi.AiAccessTokenId(created.Token.ID()))
	if deleteW.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", deleteW.Code)
	}
}

func TestAiAccessTokenManagement_RejectsBearerAuth(t *testing.T) {
	h, repo := setupAiAccessTokenHandler()
	userID := entity.NewUserID()
	created, err := aiaccesstokenuc.NewCreate(repo).Execute(
		context.Background(),
		aiaccesstokenuc.CreateInput{UserID: userID, Name: "Existing token"},
	)
	if err != nil {
		t.Fatalf("seed token: %v", err)
	}

	ctx := middleware.SetAuthMethod(
		middleware.SetUserID(context.Background(), userID),
		middleware.AuthMethodBearer,
	)
	body, _ := json.Marshal(openapi.CreateAiAccessTokenRequest{Name: "Should be denied"})

	createReq := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body)).WithContext(ctx)
	createW := httptest.NewRecorder()
	h.CreateAiAccessToken(createW, createReq)
	if createW.Code != http.StatusForbidden {
		t.Fatalf("create status = %d, want 403", createW.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	listW := httptest.NewRecorder()
	h.ListAiAccessTokens(listW, listReq)
	if listW.Code != http.StatusForbidden {
		t.Fatalf("list status = %d, want 403", listW.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/", nil).WithContext(ctx)
	deleteW := httptest.NewRecorder()
	h.RevokeAiAccessToken(deleteW, deleteReq, openapi.AiAccessTokenId(created.Token.ID()))
	if deleteW.Code != http.StatusForbidden {
		t.Fatalf("delete status = %d, want 403", deleteW.Code)
	}
}
