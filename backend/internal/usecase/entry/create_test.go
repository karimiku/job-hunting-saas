package entry

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
	saveCalled := false
	entryRepo := &mockEntryRepo{
		saveFn: func(_ context.Context, _ *entity.Entry) error {
			saveCalled = true
			return nil
		},
	}

	uc := NewCreate(entryRepo, companyFound())
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Route:     "本選考",
		Source:    "リクナビ",
		Memo:      "メモ内容",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entry == nil {
		t.Fatal("Entry should not be nil")
	}
	if out.Entry.Route().String() != "本選考" {
		t.Errorf("Route = %q, want %q", out.Entry.Route().String(), "本選考")
	}
	if out.Entry.Source().String() != "リクナビ" {
		t.Errorf("Source = %q, want %q", out.Entry.Source().String(), "リクナビ")
	}
	if out.Entry.Memo() != "メモ内容" {
		t.Errorf("Memo = %q, want %q", out.Entry.Memo(), "メモ内容")
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestCreate_WithoutMemo(t *testing.T) {
	uc := NewCreate(&mockEntryRepo{}, companyFound())
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Route:     "本選考",
		Source:    "リクナビ",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entry.Memo() != "" {
		t.Errorf("Memo = %q, want empty", out.Entry.Memo())
	}
}

func TestCreate_CompanyNotFound(t *testing.T) {
	uc := NewCreate(&mockEntryRepo{}, &mockCompanyRepo{})

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Route:     "本選考",
		Source:    "リクナビ",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestCreate_EmptyRoute(t *testing.T) {
	uc := NewCreate(&mockEntryRepo{}, companyFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Route:     "",
		Source:    "リクナビ",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrRouteEmpty) {
		t.Errorf("error = %v, want ErrRouteEmpty", err)
	}
}

func TestCreate_EmptySource(t *testing.T) {
	uc := NewCreate(&mockEntryRepo{}, companyFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Route:     "本選考",
		Source:    "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrSourceEmpty) {
		t.Errorf("error = %v, want ErrSourceEmpty", err)
	}
}

func TestCreate_SaveError(t *testing.T) {
	saveErr := errors.New("db write failed")
	entryRepo := &mockEntryRepo{
		saveFn: func(_ context.Context, _ *entity.Entry) error {
			return saveErr
		},
	}

	uc := NewCreate(entryRepo, companyFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Route:     "本選考",
		Source:    "リクナビ",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
