package companyalias

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

func TestList_Multiple(t *testing.T) {
	a1 := newTestAlias(t)
	a2 := newTestAlias(t)
	repo := &mockAliasRepo{
		listByCompanyFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) ([]*entity.CompanyAlias, error) {
			return []*entity.CompanyAlias{a1, a2}, nil
		},
	}

	uc := NewList(repo, companyFound())
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.CompanyAliases) != 2 {
		t.Fatalf("len(CompanyAliases) = %d, want 2", len(out.CompanyAliases))
	}
}

func TestList_Empty(t *testing.T) {
	repo := &mockAliasRepo{}

	uc := NewList(repo, companyFound())
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.CompanyAliases) != 0 {
		t.Fatalf("len(CompanyAliases) = %d, want 0", len(out.CompanyAliases))
	}
}

func TestList_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockAliasRepo{
		listByCompanyFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) ([]*entity.CompanyAlias, error) {
			return nil, dbErr
		},
	}

	uc := NewList(repo, companyFound())
	_, err := uc.Execute(context.Background(), ListInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}

func TestList_CompanyNotFound(t *testing.T) {
	// 存在しない/他人の company では alias を引かず ErrNotFound を返す。
	listed := false
	repo := &mockAliasRepo{
		listByCompanyFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) ([]*entity.CompanyAlias, error) {
			listed = true
			return nil, nil
		},
	}

	uc := NewList(repo, &mockCompanyRepo{})
	_, err := uc.Execute(context.Background(), ListInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
	if listed {
		t.Error("ListByCompanyID should not be called when company is not found")
	}
}
