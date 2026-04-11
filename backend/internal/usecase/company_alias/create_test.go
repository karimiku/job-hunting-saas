package companyalias

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func companyFound() *mockCompanyRepo {
	return &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return &entity.Company{}, nil
		},
	}
}

func TestCreate_Success(t *testing.T) {
	createCalled := false
	aliasRepo := &mockAliasRepo{
		createFn: func(_ context.Context, _ *entity.CompanyAlias) error {
			createCalled = true
			return nil
		},
	}

	uc := NewCreate(aliasRepo, companyFound())
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Alias:     "グーグル",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.CompanyAlias == nil {
		t.Fatal("CompanyAlias should not be nil")
	}
	if out.CompanyAlias.Alias().String() != "グーグル" {
		t.Errorf("Alias = %q, want %q", out.CompanyAlias.Alias().String(), "グーグル")
	}
	if !createCalled {
		t.Error("Create should be called")
	}
}

func TestCreate_CompanyNotFound(t *testing.T) {
	uc := NewCreate(&mockAliasRepo{}, &mockCompanyRepo{})

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Alias:     "グーグル",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestCreate_EmptyAlias(t *testing.T) {
	uc := NewCreate(&mockAliasRepo{}, companyFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Alias:     "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrAliasEmpty) {
		t.Errorf("error = %v, want ErrAliasEmpty", err)
	}
}

func TestCreate_CreateError(t *testing.T) {
	createErr := errors.New("db write failed")
	aliasRepo := &mockAliasRepo{
		createFn: func(_ context.Context, _ *entity.CompanyAlias) error {
			return createErr
		},
	}

	uc := NewCreate(aliasRepo, companyFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Alias:     "グーグル",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, createErr) {
		t.Errorf("error = %v, want createErr", err)
	}
}
