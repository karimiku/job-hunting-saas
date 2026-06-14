package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	useruc "github.com/karimiku/job-hunting-saas/internal/usecase/user"
)

// MeHandler は現在のユーザー自身に対する操作を扱う。
type MeHandler struct {
	deleteAccountUseCase *useruc.DeleteAccount
	cfg                  AuthConfig
}

// NewMeHandler は MeHandler に必要な依存を DI して新しい MeHandler を返す。
func NewMeHandler(deleteAccountUseCase *useruc.DeleteAccount, cfg AuthConfig) *MeHandler {
	return &MeHandler{deleteAccountUseCase: deleteAccountUseCase, cfg: cfg}
}

// DeleteMe は DELETE /me の handler。退会後に session cookie を削除する。
func (h *MeHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID.IsZero() {
		writeJSON(w, http.StatusUnauthorized, openapi.ErrorResponse{Message: "unauthenticated"})
		return
	}
	if middleware.GetAuthMethod(r.Context()) != middleware.AuthMethodSession {
		writeJSON(w, http.StatusUnauthorized, openapi.ErrorResponse{Message: "session authentication required"})
		return
	}
	if h.deleteAccountUseCase == nil {
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "internal error"})
		return
	}

	if err := h.deleteAccountUseCase.Execute(r.Context(), useruc.DeleteAccountInput{UserID: userID}); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			clearSessionCookie(w, h.cfg)
			writeJSON(w, http.StatusUnauthorized, openapi.ErrorResponse{Message: "unauthenticated"})
			return
		}
		log.Printf("me: DeleteAccount failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "internal error"})
		return
	}

	clearSessionCookie(w, h.cfg)
	w.WriteHeader(http.StatusNoContent)
}
