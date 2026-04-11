package company

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestCreate_NameOnly(t *testing.T) {
	saveCalled := false
	repo := &mockCompanyRepo{
		saveFn: func(_ context.Context, _ *entity.Company) error {
			saveCalled = true
			return nil
		},
	}

	uc := NewCreate(repo)
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		Name:   "株式会社テスト",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Company == nil {
		t.Fatal("Company should not be nil")
	}
	if out.Company.Name().String() != "株式会社テスト" {
		t.Errorf("Name = %q, want %q", out.Company.Name().String(), "株式会社テスト")
	}
	if out.Company.Memo() != "" {
		t.Errorf("Memo = %q, want empty", out.Company.Memo())
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestCreate_WithMemo(t *testing.T) {
	repo := &mockCompanyRepo{}

	uc := NewCreate(repo)
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		Name:   "株式会社テスト",
		Memo:   "メモ内容",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Company.Memo() != "メモ内容" {
		t.Errorf("Memo = %q, want %q", out.Company.Memo(), "メモ内容")
	}
}

func TestCreate_EmptyName(t *testing.T) {
	repo := &mockCompanyRepo{}

	uc := NewCreate(repo)
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		Name:   "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrCompanyNameEmpty) {
		t.Errorf("error = %v, want ErrCompanyNameEmpty", err)
	}
}

func TestCreate_SaveError(t *testing.T) {
	saveErr := errors.New("db write failed")
	repo := &mockCompanyRepo{
		saveFn: func(_ context.Context, _ *entity.Company) error {
			return saveErr
		},
	}

	uc := NewCreate(repo)
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		Name:   "株式会社テスト",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
