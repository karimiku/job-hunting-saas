package companyalias

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestAlias(t *testing.T) *entity.CompanyAlias {
	t.Helper()
	alias, err := value.NewAlias("グーグル")
	if err != nil {
		t.Fatalf("NewAlias failed: %v", err)
	}
	return entity.NewCompanyAlias(entity.NewUserID(), entity.NewCompanyID(), alias)
}

func TestGet_Found(t *testing.T) {
	expected := newTestAlias(t)
	repo := &mockAliasRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyAliasID) (*entity.CompanyAlias, error) {
			return expected, nil
		},
	}

	uc := NewGet(repo)
	out, err := uc.Execute(context.Background(), GetInput{
		UserID:         entity.NewUserID(),
		CompanyAliasID: expected.ID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.CompanyAlias != expected {
		t.Error("CompanyAlias should be the expected alias")
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := &mockAliasRepo{}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID:         entity.NewUserID(),
		CompanyAliasID: entity.NewCompanyAliasID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestGet_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockAliasRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyAliasID) (*entity.CompanyAlias, error) {
			return nil, dbErr
		},
	}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID:         entity.NewUserID(),
		CompanyAliasID: entity.NewCompanyAliasID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
