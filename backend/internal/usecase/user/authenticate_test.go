package user

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// --- mocks ---

type mockUserRepo struct {
	findByIDFn    func(ctx context.Context, id entity.UserID) (*entity.User, error)
	findByEmailFn func(ctx context.Context, email value.Email) (*entity.User, error)
	saveFn        func(ctx context.Context, user *entity.User) error
}

func (m *mockUserRepo) Save(ctx context.Context, user *entity.User) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepo) Delete(_ context.Context, _ entity.UserID) error {
	return nil
}

type mockExternalIdentityRepo struct {
	findFn func(ctx context.Context, provider value.AuthProvider, subject string) (*entity.ExternalIdentity, error)
	saveFn func(ctx context.Context, identity *entity.ExternalIdentity) error
}

func (m *mockExternalIdentityRepo) Save(ctx context.Context, identity *entity.ExternalIdentity) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, identity)
	}
	return nil
}

func (m *mockExternalIdentityRepo) FindByProviderAndSubject(ctx context.Context, provider value.AuthProvider, subject string) (*entity.ExternalIdentity, error) {
	if m.findFn != nil {
		return m.findFn(ctx, provider, subject)
	}
	return nil, repository.ErrNotFound
}

// --- helpers ---

func newTestEmail(t *testing.T, raw string) value.Email {
	t.Helper()
	email, err := value.NewEmail(raw)
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	return email
}

func newTestUserName(t *testing.T, raw string) value.UserName {
	t.Helper()
	name, err := value.NewUserName(raw)
	if err != nil {
		t.Fatalf("NewUserName failed: %v", err)
	}
	return name
}

func newExistingUser(t *testing.T) *entity.User {
	t.Helper()
	return entity.NewUser(
		newTestEmail(t, "existing@example.com"),
		newTestUserName(t, "既存ユーザー"),
	)
}

func validInput() AuthenticateInput {
	return AuthenticateInput{
		Provider: "google",
		Subject:  "firebase-uid-123",
		Email:    "test@example.com",
		Name:     "テスト",
	}
}

// --- tests ---

// 1. (provider, subject) でヒット → 対応 User を返す
func TestAuthenticate_FoundByExternalIdentity(t *testing.T) {
	existing := newExistingUser(t)
	identity := entity.NewExternalIdentity(existing.ID(), value.AuthProviderGoogle(), "firebase-uid-123")

	userSaveCalled := false
	idSaveCalled := false
	repo := &mockUserRepo{
		findByIDFn: func(_ context.Context, id entity.UserID) (*entity.User, error) {
			if id != existing.ID() {
				t.Errorf("FindByID called with unexpected ID: %v", id)
			}
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.User) error {
			userSaveCalled = true
			return nil
		},
	}
	idRepo := &mockExternalIdentityRepo{
		findFn: func(_ context.Context, _ value.AuthProvider, _ string) (*entity.ExternalIdentity, error) {
			return identity, nil
		},
		saveFn: func(_ context.Context, _ *entity.ExternalIdentity) error {
			idSaveCalled = true
			return nil
		},
	}

	uc := NewAuthenticate(repo, idRepo)
	out, err := uc.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Created {
		t.Error("Created should be false")
	}
	if out.User != existing {
		t.Error("User should be the existing user")
	}
	if userSaveCalled || idSaveCalled {
		t.Error("no Save should be called when identity found")
	}
}

// 2. ExternalIdentity なし / email でヒット → ExternalIdentity を紐付け保存
func TestAuthenticate_LinkExistingUserByEmail(t *testing.T) {
	existing := newExistingUser(t)

	var savedIdentity *entity.ExternalIdentity
	userSaveCalled := false
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ value.Email) (*entity.User, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.User) error {
			userSaveCalled = true
			return nil
		},
	}
	idRepo := &mockExternalIdentityRepo{
		saveFn: func(_ context.Context, id *entity.ExternalIdentity) error {
			savedIdentity = id
			return nil
		},
	}

	uc := NewAuthenticate(repo, idRepo)
	out, err := uc.Execute(context.Background(), AuthenticateInput{
		Provider: "google",
		Subject:  "firebase-uid-123",
		Email:    "existing@example.com",
		Name:     "既存ユーザー",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Created {
		t.Error("Created should be false for linked existing user")
	}
	if out.User != existing {
		t.Error("User should be the existing user")
	}
	if userSaveCalled {
		t.Error("User Save should not be called")
	}
	if savedIdentity == nil {
		t.Fatal("ExternalIdentity should be saved")
	}
	if savedIdentity.UserID() != existing.ID() {
		t.Errorf("Identity linked to wrong user: got %v, want %v", savedIdentity.UserID(), existing.ID())
	}
	if savedIdentity.Subject() != "firebase-uid-123" {
		t.Errorf("Identity subject = %q, want %q", savedIdentity.Subject(), "firebase-uid-123")
	}
}

// 3. User なし → User + ExternalIdentity の両方を新規保存
func TestAuthenticate_NewUser(t *testing.T) {
	var savedUser *entity.User
	var savedIdentity *entity.ExternalIdentity
	repo := &mockUserRepo{
		saveFn: func(_ context.Context, u *entity.User) error {
			savedUser = u
			return nil
		},
	}
	idRepo := &mockExternalIdentityRepo{
		saveFn: func(_ context.Context, id *entity.ExternalIdentity) error {
			savedIdentity = id
			return nil
		},
	}

	uc := NewAuthenticate(repo, idRepo)
	out, err := uc.Execute(context.Background(), validInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Created {
		t.Error("Created should be true")
	}
	if savedUser == nil {
		t.Fatal("User should be saved")
	}
	if savedUser.Email().String() != "test@example.com" {
		t.Errorf("Email = %q", savedUser.Email().String())
	}
	if savedIdentity == nil {
		t.Fatal("ExternalIdentity should be saved")
	}
	if savedIdentity.UserID() != savedUser.ID() {
		t.Error("Identity should be linked to new user")
	}
}

func TestAuthenticate_InvalidProvider(t *testing.T) {
	uc := NewAuthenticate(&mockUserRepo{}, &mockExternalIdentityRepo{})
	_, err := uc.Execute(context.Background(), AuthenticateInput{
		Provider: "facebook",
		Subject:  "x",
		Email:    "test@example.com",
		Name:     "a",
	})
	if !errors.Is(err, value.ErrAuthProviderInvalid) {
		t.Errorf("error = %v, want ErrAuthProviderInvalid", err)
	}
}

func TestAuthenticate_EmptySubject(t *testing.T) {
	uc := NewAuthenticate(&mockUserRepo{}, &mockExternalIdentityRepo{})
	_, err := uc.Execute(context.Background(), AuthenticateInput{
		Provider: "google",
		Subject:  "",
		Email:    "test@example.com",
		Name:     "a",
	})
	if err == nil {
		t.Fatal("expected error for empty subject")
	}
}

func TestAuthenticate_EmailValidationError(t *testing.T) {
	uc := NewAuthenticate(&mockUserRepo{}, &mockExternalIdentityRepo{})
	_, err := uc.Execute(context.Background(), AuthenticateInput{
		Provider: "google",
		Subject:  "x",
		Email:    "",
		Name:     "a",
	})
	if !errors.Is(err, value.ErrEmailEmpty) {
		t.Errorf("error = %v, want ErrEmailEmpty", err)
	}
}

func TestAuthenticate_FindByIdentityError(t *testing.T) {
	dbErr := errors.New("db down")
	idRepo := &mockExternalIdentityRepo{
		findFn: func(_ context.Context, _ value.AuthProvider, _ string) (*entity.ExternalIdentity, error) {
			return nil, dbErr
		},
	}
	uc := NewAuthenticate(&mockUserRepo{}, idRepo)
	_, err := uc.Execute(context.Background(), validInput())
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}

func TestAuthenticate_SaveUserError(t *testing.T) {
	saveErr := errors.New("write failed")
	repo := &mockUserRepo{
		saveFn: func(_ context.Context, _ *entity.User) error {
			return saveErr
		},
	}
	uc := NewAuthenticate(repo, &mockExternalIdentityRepo{})
	_, err := uc.Execute(context.Background(), validInput())
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
