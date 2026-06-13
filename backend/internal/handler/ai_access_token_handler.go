package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	aiaccesstokenuc "github.com/karimiku/job-hunting-saas/internal/usecase/ai_access_token"
)

// AiAccessTokenHandler は AI / MCP 連携用アクセストークン管理の HTTP リクエストを受ける。
type AiAccessTokenHandler struct {
	createUseCase *aiaccesstokenuc.Create
	listUseCase   *aiaccesstokenuc.List
	revokeUseCase *aiaccesstokenuc.Revoke
}

// NewAiAccessTokenHandler は AiAccessTokenHandler を DI 構築する。
func NewAiAccessTokenHandler(
	createUC *aiaccesstokenuc.Create,
	listUC *aiaccesstokenuc.List,
	revokeUC *aiaccesstokenuc.Revoke,
) *AiAccessTokenHandler {
	return &AiAccessTokenHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
		revokeUseCase: revokeUC,
	}
}

const maxAiAccessTokenBodyBytes = 8 * 1024

// CreateAiAccessToken は POST /api/v1/ai/tokens のハンドラ。
func (h *AiAccessTokenHandler) CreateAiAccessToken(w http.ResponseWriter, r *http.Request) {
	if rejectBearerTokenManagement(w, r) {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxAiAccessTokenBodyBytes)

	var req openapi.CreateAiAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeJSON(w, http.StatusRequestEntityTooLarge, openapi.ErrorResponse{Message: "request body too large"})
			return
		}
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	out, err := h.createUseCase.Execute(r.Context(), aiaccesstokenuc.CreateInput{
		UserID: middleware.GetUserID(r.Context()),
		Name:   stringValue(req.Name),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, openapi.CreateAiAccessTokenResponse{
		Token:       out.RawToken,
		AccessToken: toAiAccessTokenResponse(out.Token),
	})
}

func stringValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

// ListAiAccessTokens は GET /api/v1/ai/tokens のハンドラ。
func (h *AiAccessTokenHandler) ListAiAccessTokens(w http.ResponseWriter, r *http.Request) {
	if rejectBearerTokenManagement(w, r) {
		return
	}

	out, err := h.listUseCase.Execute(r.Context(), aiaccesstokenuc.ListInput{
		UserID: middleware.GetUserID(r.Context()),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	items := make([]openapi.AiAccessTokenResponse, len(out.Tokens))
	for i, t := range out.Tokens {
		items[i] = toAiAccessTokenResponse(t)
	}
	writeJSON(w, http.StatusOK, map[string]any{"tokens": items})
}

// RevokeAiAccessToken は DELETE /api/v1/ai/tokens/{aiAccessTokenId} のハンドラ。
func (h *AiAccessTokenHandler) RevokeAiAccessToken(w http.ResponseWriter, r *http.Request, aiAccessTokenId openapi.AiAccessTokenId) {
	if rejectBearerTokenManagement(w, r) {
		return
	}

	err := h.revokeUseCase.Execute(r.Context(), aiaccesstokenuc.RevokeInput{
		UserID:  middleware.GetUserID(r.Context()),
		TokenID: entity.AIAccessTokenID(aiAccessTokenId),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func rejectBearerTokenManagement(w http.ResponseWriter, r *http.Request) bool {
	if middleware.GetAuthMethod(r.Context()) != middleware.AuthMethodBearer {
		return false
	}
	writeJSON(w, http.StatusForbidden, openapi.ErrorResponse{
		Message: "AI access token management requires session authentication",
	})
	return true
}

func toAiAccessTokenResponse(t *entity.AIAccessToken) openapi.AiAccessTokenResponse {
	return openapi.AiAccessTokenResponse{
		Id:          uuid.UUID(t.ID()),
		Name:        t.Name().String(),
		TokenPrefix: t.Prefix().String(),
		CreatedAt:   t.CreatedAt(),
		LastUsedAt:  t.LastUsedAt(),
		RevokedAt:   t.RevokedAt(),
	}
}
