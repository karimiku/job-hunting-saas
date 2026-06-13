package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestResolveMCPUserWithAPIKey(t *testing.T) {
	email, err := value.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	name, err := value.NewUserName("User")
	if err != nil {
		t.Fatalf("NewUserName failed: %v", err)
	}
	user := entity.NewUser(email, name)
	secret, err := value.GenerateAIAccessTokenSecret()
	if err != nil {
		t.Fatalf("GenerateAIAccessTokenSecret failed: %v", err)
	}
	token := entity.NewAIAccessToken(user.ID(), "Codex", secret.Hash(), secret.Preview())
	userRepo := &fakeMCPUserRepo{user: user}
	tokenRepo := &fakeMCPTokenRepo{token: token}

	t.Setenv("MCP_API_KEY", secret.String())
	t.Setenv("MCP_USER_ID", "")
	t.Setenv("MCP_USER_EMAIL", "ignored@example.com")

	got, err := resolveMCPUser(context.Background(), userRepo, tokenRepo)
	if err != nil {
		t.Fatalf("resolveMCPUser() failed: %v", err)
	}
	if got.ID() != user.ID() {
		t.Errorf("resolved user ID = %s, want %s", got.ID(), user.ID())
	}
	if tokenRepo.touchedID != token.ID() {
		t.Errorf("touchedID = %s, want %s", tokenRepo.touchedID, token.ID())
	}
}

func TestResolveMCPUserRejectsUnknownAPIKey(t *testing.T) {
	secret, err := value.GenerateAIAccessTokenSecret()
	if err != nil {
		t.Fatalf("GenerateAIAccessTokenSecret failed: %v", err)
	}
	t.Setenv("MCP_API_KEY", secret.String())
	t.Setenv("MCP_USER_ID", "")
	t.Setenv("MCP_USER_EMAIL", "")

	_, err = resolveMCPUser(context.Background(), &fakeMCPUserRepo{}, &fakeMCPTokenRepo{})
	if err == nil {
		t.Fatal("resolveMCPUser() succeeded, want error")
	}
}

func TestResolveMCPUserFallsBackToEmail(t *testing.T) {
	email, err := value.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	name, err := value.NewUserName("User")
	if err != nil {
		t.Fatalf("NewUserName failed: %v", err)
	}
	user := entity.NewUser(email, name)

	t.Setenv("MCP_API_KEY", "")
	t.Setenv("MCP_USER_ID", "")
	t.Setenv("MCP_USER_EMAIL", email.String())

	got, err := resolveMCPUser(context.Background(), &fakeMCPUserRepo{user: user}, &fakeMCPTokenRepo{})
	if err != nil {
		t.Fatalf("resolveMCPUser() failed: %v", err)
	}
	if got.ID() != user.ID() {
		t.Errorf("resolved user ID = %s, want %s", got.ID(), user.ID())
	}
}

type fakeMCPUserRepo struct {
	user *entity.User
}

func (r *fakeMCPUserRepo) Save(context.Context, *entity.User) error {
	return nil
}

func (r *fakeMCPUserRepo) FindByID(_ context.Context, id entity.UserID) (*entity.User, error) {
	if r.user != nil && r.user.ID() == id {
		return r.user, nil
	}
	return nil, repository.ErrNotFound
}

func (r *fakeMCPUserRepo) FindByEmail(_ context.Context, email value.Email) (*entity.User, error) {
	if r.user != nil && r.user.Email().String() == email.String() {
		return r.user, nil
	}
	return nil, repository.ErrNotFound
}

func (r *fakeMCPUserRepo) Delete(context.Context, entity.UserID) error {
	return nil
}

type fakeMCPTokenRepo struct {
	token     *entity.AIAccessToken
	touchedID entity.AIAccessTokenID
}

func (r *fakeMCPTokenRepo) Save(context.Context, *entity.AIAccessToken) error {
	return nil
}

func (r *fakeMCPTokenRepo) FindActiveByHash(_ context.Context, tokenHash string) (*entity.AIAccessToken, error) {
	if r.token != nil && r.token.TokenHash() == tokenHash {
		return r.token, nil
	}
	return nil, repository.ErrNotFound
}

func (r *fakeMCPTokenRepo) ListByUserID(context.Context, entity.UserID) ([]*entity.AIAccessToken, error) {
	return nil, nil
}

func (r *fakeMCPTokenRepo) TouchLastUsed(_ context.Context, id entity.AIAccessTokenID, _ time.Time) error {
	if id.IsZero() {
		return errors.New("id is required")
	}
	r.touchedID = id
	return nil
}

func (r *fakeMCPTokenRepo) Revoke(context.Context, entity.UserID, entity.AIAccessTokenID, time.Time) error {
	return nil
}
