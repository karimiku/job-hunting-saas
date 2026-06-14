package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	useruc "github.com/karimiku/job-hunting-saas/internal/usecase/user"
)

type fakeAccountRepo struct {
	deleteFn func(ctx context.Context, userID entity.UserID) error
}

func (r *fakeAccountRepo) DeleteAccount(ctx context.Context, userID entity.UserID) error {
	if r.deleteFn != nil {
		return r.deleteFn(ctx, userID)
	}
	return nil
}

func TestDeleteMe_Success(t *testing.T) {
	userRepo := inmemory.NewUserRepository()
	user := seedUser(t, userRepo, "delete-me@example.com", "退会ユーザー")
	h := NewMeHandler(
		useruc.NewDeleteAccount(inmemory.NewAccountRepository(userRepo)),
		AuthConfig{CookieDomain: "example.com", CookieSecure: true, CookieSameSite: http.SameSiteNoneMode},
	)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	ctx := middleware.SetAuthMethod(
		middleware.SetUserID(req.Context(), user.ID()),
		middleware.AuthMethodSession,
	)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.DeleteMe(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204, body = %s", w.Code, w.Body.String())
	}
	if _, err := userRepo.FindByID(context.Background(), user.ID()); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("user should be deleted, err = %v", err)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies = %d, want 1", len(cookies))
	}
	c := cookies[0]
	if c.Name != "session" || c.Value != "" || c.MaxAge >= 0 {
		t.Fatalf("cookie = %+v, want cleared session cookie", c)
	}
	if c.Domain != "example.com" || !c.Secure || c.SameSite != http.SameSiteNoneMode {
		t.Fatalf("cookie config = %+v, want domain/secure/samesite preserved", c)
	}
}

func TestDeleteMe_Unauthenticated(t *testing.T) {
	called := false
	h := NewMeHandler(
		useruc.NewDeleteAccount(&fakeAccountRepo{
			deleteFn: func(_ context.Context, _ entity.UserID) error {
				called = true
				return nil
			},
		}),
		AuthConfig{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	w := httptest.NewRecorder()

	h.DeleteMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if called {
		t.Fatal("repository should not be called")
	}
}

func TestDeleteMe_RejectsBearerAuth(t *testing.T) {
	called := false
	h := NewMeHandler(
		useruc.NewDeleteAccount(&fakeAccountRepo{
			deleteFn: func(_ context.Context, _ entity.UserID) error {
				called = true
				return nil
			},
		}),
		AuthConfig{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	ctx := middleware.SetAuthMethod(
		middleware.SetUserID(req.Context(), entity.NewUserID()),
		middleware.AuthMethodBearer,
	)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.DeleteMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if called {
		t.Fatal("repository should not be called")
	}
}

func TestDeleteMe_UserNotFoundClearsCookie(t *testing.T) {
	h := NewMeHandler(
		useruc.NewDeleteAccount(&fakeAccountRepo{
			deleteFn: func(_ context.Context, _ entity.UserID) error {
				return repository.ErrNotFound
			},
		}),
		AuthConfig{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	ctx := middleware.SetAuthMethod(
		middleware.SetUserID(req.Context(), entity.NewUserID()),
		middleware.AuthMethodSession,
	)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.DeleteMe(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if got := w.Result().Cookies(); len(got) != 1 || got[0].MaxAge >= 0 {
		t.Fatalf("cookie should be cleared, got %+v", got)
	}
}

func TestDeleteMe_InternalError(t *testing.T) {
	h := NewMeHandler(
		useruc.NewDeleteAccount(&fakeAccountRepo{
			deleteFn: func(_ context.Context, _ entity.UserID) error {
				return errors.New("db failed")
			},
		}),
		AuthConfig{},
	)

	req := httptest.NewRequest(http.MethodDelete, "/me", nil)
	ctx := middleware.SetAuthMethod(
		middleware.SetUserID(req.Context(), entity.NewUserID()),
		middleware.AuthMethodSession,
	)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.DeleteMe(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}
