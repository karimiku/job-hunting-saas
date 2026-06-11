package aiaccesstoken

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

func TestCreateListVerify_RevealsRawTokenOnlyOnCreate(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewAIAccessTokenRepository()
	userID := entity.NewUserID()

	created, err := NewCreate(repo).Execute(ctx, CreateInput{
		UserID: userID,
		Name:   " Claude Desktop ",
	})
	if err != nil {
		t.Fatalf("Create.Execute() failed: %v", err)
	}
	if !strings.HasPrefix(created.RawToken, value.AIAccessTokenRawPrefix) {
		t.Fatalf("RawToken prefix = %q", created.RawToken)
	}
	if created.Token.Name().String() != "Claude Desktop" {
		t.Errorf("Name = %q", created.Token.Name().String())
	}

	listed, err := NewList(repo).Execute(ctx, ListInput{UserID: userID})
	if err != nil {
		t.Fatalf("List.Execute() failed: %v", err)
	}
	if len(listed.Tokens) != 1 {
		t.Fatalf("len(tokens) = %d, want 1", len(listed.Tokens))
	}
	if listed.Tokens[0].TokenHash().String() == created.RawToken {
		t.Fatal("repository must not expose raw token as hash")
	}

	verified, err := NewVerify(repo).Execute(ctx, VerifyInput{RawToken: created.RawToken})
	if err != nil {
		t.Fatalf("Verify.Execute() failed: %v", err)
	}
	if verified.UserID != userID {
		t.Errorf("UserID = %s, want %s", verified.UserID, userID)
	}
	afterVerify, err := NewList(repo).Execute(ctx, ListInput{UserID: userID})
	if err != nil {
		t.Fatalf("List after verify failed: %v", err)
	}
	if afterVerify.Tokens[0].LastUsedAt() == nil {
		t.Fatal("LastUsedAt() = nil, want timestamp after verify")
	}
}

func TestVerify_RevokedTokenIsInvalid(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewAIAccessTokenRepository()
	userID := entity.NewUserID()

	created, err := NewCreate(repo).Execute(ctx, CreateInput{UserID: userID, Name: "Codex"})
	if err != nil {
		t.Fatalf("Create.Execute() failed: %v", err)
	}
	if err := NewRevoke(repo).Execute(ctx, RevokeInput{UserID: userID, TokenID: created.Token.ID()}); err != nil {
		t.Fatalf("Revoke.Execute() failed: %v", err)
	}
	_, err = NewVerify(repo).Execute(ctx, VerifyInput{RawToken: created.RawToken})
	if !errors.Is(err, value.ErrAIAccessTokenInvalid) {
		t.Fatalf("Verify revoked token error = %v, want ErrAIAccessTokenInvalid", err)
	}
}

func TestCreate_RequiresUserID(t *testing.T) {
	_, err := NewCreate(inmemory.NewAIAccessTokenRepository()).Execute(context.Background(), CreateInput{
		Name: "Claude",
	})
	if !errors.Is(err, ErrUserIDRequired) {
		t.Fatalf("error = %v, want ErrUserIDRequired", err)
	}
}
