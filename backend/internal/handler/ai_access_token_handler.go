package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
)

const defaultAIIntegrationLabel = "AI連携トークン"

// AIAccessTokenHandler はAI連携用アクセストークンのHTTP handler。
type AIAccessTokenHandler struct {
	repo repository.AIAccessTokenRepository
}

// NewAIAccessTokenHandler はAIAccessTokenHandlerを生成する。
func NewAIAccessTokenHandler(repo repository.AIAccessTokenRepository) *AIAccessTokenHandler {
	return &AIAccessTokenHandler{repo: repo}
}

// ListAIAccessTokens はログインユーザーのAI連携token一覧を返す。
func (h *AIAccessTokenHandler) ListAIAccessTokens(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireSession(w, r)
	if !ok {
		return
	}

	tokens, err := h.repo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to list tokens"})
		return
	}
	items := make([]openapi.AIAccessTokenResponse, 0, len(tokens))
	for _, token := range tokens {
		items = append(items, toAIAccessTokenResponse(token))
	}
	writeJSON(w, http.StatusOK, map[string]any{"tokens": items})
}

// CreateAIAccessToken は新しいAI連携tokenを発行する。
func (h *AIAccessTokenHandler) CreateAIAccessToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.requireSession(w, r)
	if !ok {
		return
	}

	var req openapi.CreateAIAccessTokenRequest
	if r.Body != nil {
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil && !errors.Is(err, io.EOF) {
			writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
			return
		}
	}
	name := strings.TrimSpace(valueOrEmpty(req.Name))
	if name == "" {
		name = defaultAIIntegrationLabel
	}
	if len([]rune(name)) > 80 {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "name must be 80 characters or less"})
		return
	}

	secret, err := value.GenerateAIAccessTokenSecret()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to generate token"})
		return
	}
	token := entity.NewAIAccessToken(userID, name, secret.Hash(), secret.Preview())
	if err := h.repo.Save(r.Context(), token); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			writeJSON(w, http.StatusConflict, openapi.ErrorResponse{Message: "token already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to save token"})
		return
	}

	writeJSON(w, http.StatusCreated, openapi.CreateAIAccessTokenResponse{
		Token: secret.String(),
		Item:  toAIAccessTokenResponse(token),
	})
}

// RevokeAIAccessToken は指定tokenを失効する。
func (h *AIAccessTokenHandler) RevokeAIAccessToken(w http.ResponseWriter, r *http.Request, tokenId openapi.AIAccessTokenId) {
	userID, ok := h.requireSession(w, r)
	if !ok {
		return
	}

	err := h.repo.Revoke(r.Context(), userID, entity.AIAccessTokenID(tokenId), time.Now())
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, openapi.ErrorResponse{Message: "token not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "failed to revoke token"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AIAccessTokenHandler) requireSession(w http.ResponseWriter, r *http.Request) (entity.UserID, bool) {
	if h == nil || h.repo == nil {
		writeJSON(w, http.StatusServiceUnavailable, openapi.ErrorResponse{Message: "token management is unavailable"})
		return entity.UserID{}, false
	}
	userID := middleware.GetUserID(r.Context())
	if userID.IsZero() {
		writeJSON(w, http.StatusUnauthorized, openapi.ErrorResponse{Message: "unauthenticated"})
		return entity.UserID{}, false
	}
	if middleware.GetAuthMethod(r.Context()) != middleware.AuthMethodSession {
		writeJSON(w, http.StatusForbidden, openapi.ErrorResponse{Message: "session authentication is required"})
		return entity.UserID{}, false
	}
	return userID, true
}

func toAIAccessTokenResponse(token *entity.AIAccessToken) openapi.AIAccessTokenResponse {
	resp := openapi.AIAccessTokenResponse{
		Id:           uuid.UUID(token.ID()),
		Name:         token.Name(),
		TokenPreview: token.TokenPreview(),
		CreatedAt:    token.CreatedAt(),
	}
	if token.LastUsedAt() != nil {
		resp.LastUsedAt = token.LastUsedAt()
	}
	if token.RevokedAt() != nil {
		resp.RevokedAt = token.RevokedAt()
	}
	return resp
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
