package company

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestUpdate_Success(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestCompany(t, userID)
	saveCalled := false
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Company) error {
			saveCalled = true
			return nil
		},
	}

	newName := "新しい社名"
	newMemo := "新しいメモ"
	uc := NewUpdate(repo)
	out, err := uc.Execute(context.Background(), UpdateInput{
		UserID:    userID,
		CompanyID: existing.ID(),
		Name:      &newName,
		Memo:      &newMemo,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Company.Name().String() != "新しい社名" {
		t.Errorf("Name = %q, want %q", out.Company.Name().String(), "新しい社名")
	}
	if out.Company.Memo() != "新しいメモ" {
		t.Errorf("Memo = %q, want %q", out.Company.Memo(), "新しいメモ")
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestUpdate_PartialUpdate_NameOnly(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestCompany(t, userID)
	originalMemo := existing.Memo()
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Company) error {
			return nil
		},
	}

	newName := "新しい社名"
	uc := NewUpdate(repo)
	out, err := uc.Execute(context.Background(), UpdateInput{
		UserID:    userID,
		CompanyID: existing.ID(),
		Name:      &newName,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Company.Name().String() != "新しい社名" {
		t.Errorf("Name = %q, want %q", out.Company.Name().String(), "新しい社名")
	}
	if out.Company.Memo() != originalMemo {
		t.Errorf("Memo = %q, want %q (unchanged)", out.Company.Memo(), originalMemo)
	}
}

func TestUpdate_EmptyName(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestCompany(t, userID)
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return existing, nil
		},
	}

	emptyName := ""
	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:    userID,
		CompanyID: existing.ID(),
		Name:      &emptyName,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrCompanyNameEmpty) {
		t.Errorf("error = %v, want ErrCompanyNameEmpty", err)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockCompanyRepo{}

	newName := "株式会社テスト"
	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Name:      &newName,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestUpdate_SaveError(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestCompany(t, userID)
	saveErr := errors.New("db write failed")
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Company) error {
			return saveErr
		},
	}

	newName := "株式会社テスト"
	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:    userID,
		CompanyID: existing.ID(),
		Name:      &newName,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
