package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	fbauth "firebase.google.com/go/v4/auth"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

// --- SetUserID / GetUserID ---

func TestGetUserID_WhenSet_ReturnsValue(t *testing.T) {
	expected := entity.NewUserID()
	ctx := SetUserID(context.Background(), expected)

	got := GetUserID(ctx)
	if got != expected {
		t.Errorf("GetUserID = %v, want %v", got, expected)
	}
}

func TestGetUserID_WhenNotSet_ReturnsZero(t *testing.T) {
	got := GetUserID(context.Background())
	if !got.IsZero() {
		t.Errorf("GetUserID = %v, want zero", got)
	}
}

func TestGetUserID_WhenWrongType_ReturnsZero(t *testing.T) {
	// 別パッケージから同名キーで誤った型をセットされても拾わないこと（ユーザー由来の改ざん耐性）
	ctx := context.WithValue(context.Background(), contextKey("userID"), "not-a-userid")

	got := GetUserID(ctx)
	if !got.IsZero() {
		t.Errorf("GetUserID = %v, want zero (wrong type assertion should fail)", got)
	}
}

func TestSetUserID_OverwritesPrevious(t *testing.T) {
	first := entity.NewUserID()
	second := entity.NewUserID()

	ctx := SetUserID(context.Background(), first)
	ctx = SetUserID(ctx, second)

	if got := GetUserID(ctx); got != second {
		t.Errorf("GetUserID = %v, want %v (latest wins)", got, second)
	}
}

// --- NewAuth ---

// mockSessionVerifier は FirebaseSessionVerifier のテスト実装。
type mockSessionVerifier struct {
	verifyFn func(ctx context.Context, cookie string) (*fbauth.Token, error)
}

func (m *mockSessionVerifier) VerifySessionCookie(ctx context.Context, cookie string) (*fbauth.Token, error) {
	return m.verifyFn(ctx, cookie)
}

// nextAssertingUserID は ServeHTTP まで到達したかと、
// context にセットされた userID が期待値と一致するかを確認する next handler を返す。
func nextAssertingUserID(t *testing.T, called *bool, expected entity.UserID) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		got := GetUserID(r.Context())
		if got != expected {
			t.Errorf("GetUserID in next handler = %v, want %v", got, expected)
		}
		w.WriteHeader(http.StatusOK)
	})
}

func TestNewAuth_Success(t *testing.T) {
	userID := entity.NewUserID()
	extIDRepo := inmemory.NewExternalIdentityRepository()
	identity := entity.NewExternalIdentity(userID, value.AuthProviderGoogle(), "firebase-uid")
	if err := extIDRepo.Save(context.Background(), identity); err != nil {
		t.Fatalf("save identity: %v", err)
	}

	fb := &mockSessionVerifier{
		verifyFn: func(_ context.Context, cookie string) (*fbauth.Token, error) {
			if cookie != "valid-cookie" {
				t.Errorf("cookie = %q, want valid-cookie", cookie)
			}
			return &fbauth.Token{UID: "firebase-uid"}, nil
		},
	}

	called := false
	handler := NewAuth(fb, extIDRepo)(nextAssertingUserID(t, &called, userID))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "valid-cookie"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if !called {
		t.Error("next handler should be called")
	}
}

func TestNewAuth_NoCookie(t *testing.T) {
	called := false
	handler := NewAuth(&mockSessionVerifier{}, inmemory.NewExternalIdentityRepository())(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true }),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
	if called {
		t.Error("next handler should NOT be called when unauthenticated")
	}
}

func TestNewAuth_EmptyCookie(t *testing.T) {
	called := false
	handler := NewAuth(&mockSessionVerifier{}, inmemory.NewExternalIdentityRepository())(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true }),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: ""})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
	if called {
		t.Error("next handler should NOT be called")
	}
}

func TestNewAuth_InvalidSessionCookie(t *testing.T) {
	fb := &mockSessionVerifier{
		verifyFn: func(_ context.Context, _ string) (*fbauth.Token, error) {
			return nil, errors.New("expired cookie")
		},
	}
	called := false
	handler := NewAuth(fb, inmemory.NewExternalIdentityRepository())(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true }),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "expired"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
	if called {
		t.Error("next handler should NOT be called")
	}
}

func TestNewAuth_IdentityNotFound(t *testing.T) {
	// Session は有効だが DB に該当 ExternalIdentity がない異常系
	fb := &mockSessionVerifier{
		verifyFn: func(_ context.Context, _ string) (*fbauth.Token, error) {
			return &fbauth.Token{UID: "unknown-uid"}, nil
		},
	}
	called := false
	handler := NewAuth(fb, inmemory.NewExternalIdentityRepository())(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true }),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "valid"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", w.Code)
	}
	if called {
		t.Error("next handler should NOT be called")
	}
}

// failingExtIDRepo は FindByProviderAndSubject が ErrNotFound 以外のエラーを返すリポジトリ。
type failingExtIDRepo struct {
	err error
}

func (r *failingExtIDRepo) Save(_ context.Context, _ *entity.ExternalIdentity) error {
	return nil
}

func (r *failingExtIDRepo) FindByProviderAndSubject(_ context.Context, _ value.AuthProvider, _ string) (*entity.ExternalIdentity, error) {
	return nil, r.err
}

func TestNewAuth_RepoUnexpectedError(t *testing.T) {
	fb := &mockSessionVerifier{
		verifyFn: func(_ context.Context, _ string) (*fbauth.Token, error) {
			return &fbauth.Token{UID: "uid"}, nil
		},
	}
	repo := &failingExtIDRepo{err: errors.New("db unreachable")}

	// 念のため repository.ErrNotFound と区別されているか確認のためのアサーション
	if errors.Is(repo.err, repository.ErrNotFound) {
		t.Fatal("test setup error: error must NOT be ErrNotFound")
	}

	called := false
	handler := NewAuth(fb, repo)(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) { called = true }))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: "valid"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
	if called {
		t.Error("next handler should NOT be called on internal error")
	}
}
