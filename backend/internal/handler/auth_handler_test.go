package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	useruc "github.com/karimiku/job-hunting-saas/internal/usecase/user"
)

// failingUserRepo は FindByID が non-NotFound エラーを返すリポジトリ。
type failingUserRepo struct{}

func (r *failingUserRepo) Save(_ context.Context, _ *entity.User) error { return nil }
func (r *failingUserRepo) FindByID(_ context.Context, _ entity.UserID) (*entity.User, error) {
	return nil, errors.New("db unreachable")
}
func (r *failingUserRepo) FindByEmail(_ context.Context, _ value.Email) (*entity.User, error) {
	return nil, repository.ErrNotFound
}
func (r *failingUserRepo) Delete(_ context.Context, _ entity.UserID) error { return nil }

// mockFirebaseAuth は FirebaseSessionCreator の差し替え可能なテスト実装。
type mockFirebaseAuth struct {
	verifyIDTokenFn func(ctx context.Context, idToken string) (*IDTokenClaims, error)
	sessionCookieFn func(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
}

func (m *mockFirebaseAuth) VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error) {
	return m.verifyIDTokenFn(ctx, idToken)
}

func (m *mockFirebaseAuth) SessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	return m.sessionCookieFn(ctx, idToken, expiresIn)
}

// setupAuthHandler は Me / DeleteSession テスト用に AuthHandler を組み立てる。
// CreateSession を呼ばないテストでは firebase 引数は nil で良い。
func setupAuthHandler() (*AuthHandler, *inmemory.UserRepository) {
	userRepo := inmemory.NewUserRepository()
	h := NewAuthHandler(nil, nil, userRepo, AuthConfig{})
	return h, userRepo
}

// setupAuthHandlerWithFirebase は CreateSession テスト用にモック Firebase + 本物の Authenticate UC を組む。
func setupAuthHandlerWithFirebase(fb *mockFirebaseAuth) (*AuthHandler, *inmemory.UserRepository, *inmemory.ExternalIdentityRepository) {
	userRepo := inmemory.NewUserRepository()
	extIDRepo := inmemory.NewExternalIdentityRepository()
	authUC := useruc.NewAuthenticate(userRepo, extIDRepo)
	h := NewAuthHandler(fb, authUC, userRepo, AuthConfig{})
	return h, userRepo, extIDRepo
}

func seedUser(t *testing.T, userRepo *inmemory.UserRepository, email, name string) *entity.User {
	t.Helper()

	emailVO, err := value.NewEmail(email)
	if err != nil {
		t.Fatalf("NewEmail: %v", err)
	}
	nameVO, err := value.NewUserName(name)
	if err != nil {
		t.Fatalf("NewUserName: %v", err)
	}
	user := entity.NewUser(emailVO, nameVO)
	if err := userRepo.Save(context.Background(), user); err != nil {
		t.Fatalf("save user: %v", err)
	}
	return user
}

// freshClaims は鮮度チェックを通る IDTokenClaims を返す（AuthTime = 直前）。
func freshClaims(uid, email, name string) *IDTokenClaims {
	return &IDTokenClaims{
		UID:      uid,
		Email:    email,
		Name:     name,
		AuthTime: time.Now(),
	}
}

// --- CreateSession ---

func TestCreateSession_Success_NewUser(t *testing.T) {
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return freshClaims("firebase-uid-1", "new@example.com", "新規ユーザー"), nil
		},
		sessionCookieFn: func(_ context.Context, _ string, _ time.Duration) (string, error) {
			return "session-cookie-value", nil
		},
	}
	h, userRepo, _ := setupAuthHandlerWithFirebase(fb)

	body, _ := json.Marshal(createSessionRequest{IDToken: "valid-id-token"})
	req := httptest.NewRequest(http.MethodPost, "/auth/session", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateSession(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	// User が新規作成されていること
	email, _ := value.NewEmail("new@example.com")
	if _, err := userRepo.FindByEmail(context.Background(), email); err != nil {
		t.Errorf("user should be created: %v", err)
	}

	// Session Cookie が発行されること
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies = %d, want 1", len(cookies))
	}
	if cookies[0].Name != "session" || cookies[0].Value != "session-cookie-value" {
		t.Errorf("cookie = %+v, want session=session-cookie-value", cookies[0])
	}
	if !cookies[0].HttpOnly {
		t.Error("cookie should be HttpOnly")
	}
}

func TestCreateSession_Success_ExistingUserByEmail(t *testing.T) {
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return freshClaims("firebase-uid-2", "existing@example.com", "既存ユーザー"), nil
		},
		sessionCookieFn: func(_ context.Context, _ string, _ time.Duration) (string, error) {
			return "cookie", nil
		},
	}
	h, userRepo, _ := setupAuthHandlerWithFirebase(fb)
	existing := seedUser(t, userRepo, "existing@example.com", "既存ユーザー")

	body, _ := json.Marshal(createSessionRequest{IDToken: "valid"})
	req := httptest.NewRequest(http.MethodPost, "/auth/session", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateSession(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp authUserResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.ID != existing.ID().String() {
		t.Errorf("ID = %q, want %q (existing user reused)", resp.ID, existing.ID().String())
	}
}

func TestCreateSession_FallbackNameToEmail(t *testing.T) {
	// Google アカウントで displayName が空のケース
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return &IDTokenClaims{
				UID:      "uid-3",
				Email:    "noname@example.com",
				AuthTime: time.Now(),
			}, nil
		},
		sessionCookieFn: func(_ context.Context, _ string, _ time.Duration) (string, error) {
			return "cookie", nil
		},
	}
	h, _, _ := setupAuthHandlerWithFirebase(fb)

	body, _ := json.Marshal(createSessionRequest{IDToken: "valid"})
	req := httptest.NewRequest(http.MethodPost, "/auth/session", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}
	var resp authUserResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Name != "noname@example.com" {
		t.Errorf("Name = %q, want fallback to email", resp.Name)
	}
}

func TestCreateSession_InvalidJSON(t *testing.T) {
	h, _, _ := setupAuthHandlerWithFirebase(&mockFirebaseAuth{})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateSession_EmptyIDToken(t *testing.T) {
	h, _, _ := setupAuthHandlerWithFirebase(&mockFirebaseAuth{})

	body, _ := json.Marshal(createSessionRequest{IDToken: ""})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateSession_InvalidIDToken(t *testing.T) {
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return nil, errors.New("invalid token")
		},
	}
	h, _, _ := setupAuthHandlerWithFirebase(fb)

	body, _ := json.Marshal(createSessionRequest{IDToken: "bad"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestCreateSession_StaleIDToken(t *testing.T) {
	// 5分以上前に発行されたトークンは拒否される（Session Fixation 防止）
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return &IDTokenClaims{
				UID:      "uid",
				Email:    "stale@example.com",
				Name:     "name",
				AuthTime: time.Now().Add(-10 * time.Minute),
			}, nil
		},
	}
	h, _, _ := setupAuthHandlerWithFirebase(fb)

	body, _ := json.Marshal(createSessionRequest{IDToken: "stale"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestCreateSession_AuthenticateFails(t *testing.T) {
	// email が空のクレームで Authenticate UC が value.NewEmail で失敗するケース
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return &IDTokenClaims{
				UID:      "uid",
				AuthTime: time.Now(),
			}, nil
		},
	}
	h, _, _ := setupAuthHandlerWithFirebase(fb)

	body, _ := json.Marshal(createSessionRequest{IDToken: "valid"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

func TestCreateSession_SessionCookieFails(t *testing.T) {
	fb := &mockFirebaseAuth{
		verifyIDTokenFn: func(_ context.Context, _ string) (*IDTokenClaims, error) {
			return freshClaims("uid", "ok@example.com", "name"), nil
		},
		sessionCookieFn: func(_ context.Context, _ string, _ time.Duration) (string, error) {
			return "", errors.New("firebase down")
		},
	}
	h, _, _ := setupAuthHandlerWithFirebase(fb)

	body, _ := json.Marshal(createSessionRequest{IDToken: "valid"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.CreateSession(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

// --- Me ---

func TestMe_Success(t *testing.T) {
	h, userRepo := setupAuthHandler()
	user := seedUser(t, userRepo, "test@example.com", "テストユーザー")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), user.ID()))
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	var resp authUserResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != user.ID().String() {
		t.Errorf("ID = %q, want %q", resp.ID, user.ID().String())
	}
	if resp.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", resp.Email, "test@example.com")
	}
	if resp.Name != "テストユーザー" {
		t.Errorf("Name = %q, want %q", resp.Name, "テストユーザー")
	}
}

func TestMe_Unauthenticated(t *testing.T) {
	h, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

func TestMe_UserNotFound(t *testing.T) {
	h, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
}

// --- DeleteSession ---

func TestDeleteSession_ClearsCookie(t *testing.T) {
	h, _ := setupAuthHandler()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()

	h.DeleteSession(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookies = %d, want 1", len(cookies))
	}
	c := cookies[0]
	if c.Name != "session" {
		t.Errorf("cookie.Name = %q, want %q", c.Name, "session")
	}
	if c.Value != "" {
		t.Errorf("cookie.Value = %q, want empty", c.Value)
	}
	if c.MaxAge >= 0 {
		t.Errorf("cookie.MaxAge = %d, want negative (clearing)", c.MaxAge)
	}
	if !c.HttpOnly {
		t.Error("cookie should be HttpOnly")
	}
}

func TestMe_InternalError(t *testing.T) {
	h := NewAuthHandler(nil, nil, &failingUserRepo{}, AuthConfig{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}

// --- PublicRoutes / ProtectedRoutes ---

func TestPublicRoutes_Registered(t *testing.T) {
	h, _ := setupAuthHandler()
	r := chi.NewRouter()
	h.PublicRoutes(r)

	// POST /auth/session と DELETE /auth/session が登録されていることを確認
	if !routeExists(r, http.MethodPost, "/auth/session") {
		t.Error("POST /auth/session should be registered")
	}
	if !routeExists(r, http.MethodDelete, "/auth/session") {
		t.Error("DELETE /auth/session should be registered")
	}
}

func TestProtectedRoutes_Registered(t *testing.T) {
	h, _ := setupAuthHandler()
	r := chi.NewRouter()
	h.ProtectedRoutes(r)

	if !routeExists(r, http.MethodGet, "/auth/me") {
		t.Error("GET /auth/me should be registered")
	}
}

func routeExists(r chi.Router, method, path string) bool {
	tctx := chi.NewRouteContext()
	return r.Match(tctx, method, path)
}

func TestDeleteSession_RespectsConfig(t *testing.T) {
	userRepo := inmemory.NewUserRepository()
	h := NewAuthHandler(nil, nil, userRepo, AuthConfig{
		CookieDomain: "example.com",
		CookieSecure: true,
	})

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	h.DeleteSession(w, req)

	c := w.Result().Cookies()[0]
	if c.Domain != "example.com" {
		t.Errorf("Domain = %q, want %q", c.Domain, "example.com")
	}
	if !c.Secure {
		t.Error("Secure should be true")
	}
}
