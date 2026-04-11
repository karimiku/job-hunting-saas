package entry

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
	existing := newTestEntry(t, userID)
	saveCalled := false
	repo := &mockEntryRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) (*entity.Entry, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Entry) error {
			saveCalled = true
			return nil
		},
	}

	uc := NewUpdate(repo)
	out, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     userID,
		EntryID:    existing.ID(),
		Source:     "マイナビ",
		Status:     "offered",
		StageKind:  "interview",
		StageLabel: "一次面接",
		Memo:       "新しいメモ",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entry.Source().String() != "マイナビ" {
		t.Errorf("Source = %q, want %q", out.Entry.Source().String(), "マイナビ")
	}
	if out.Entry.Status().String() != "offered" {
		t.Errorf("Status = %q, want %q", out.Entry.Status().String(), "offered")
	}
	if out.Entry.Stage().Kind().String() != "interview" {
		t.Errorf("Stage.Kind = %q, want %q", out.Entry.Stage().Kind().String(), "interview")
	}
	if out.Entry.Stage().Label() != "一次面接" {
		t.Errorf("Stage.Label = %q, want %q", out.Entry.Stage().Label(), "一次面接")
	}
	if out.Entry.Memo() != "新しいメモ" {
		t.Errorf("Memo = %q, want %q", out.Entry.Memo(), "新しいメモ")
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     entity.NewUserID(),
		EntryID:    entity.NewEntryID(),
		Source:     "リクナビ",
		Status:     "in_progress",
		StageKind:  "application",
		StageLabel: "応募",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestUpdate_InvalidSource(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     entity.NewUserID(),
		EntryID:    entity.NewEntryID(),
		Source:     "",
		Status:     "in_progress",
		StageKind:  "application",
		StageLabel: "応募",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrSourceEmpty) {
		t.Errorf("error = %v, want ErrSourceEmpty", err)
	}
}

func TestUpdate_InvalidStatus(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     entity.NewUserID(),
		EntryID:    entity.NewEntryID(),
		Source:     "リクナビ",
		Status:     "invalid",
		StageKind:  "application",
		StageLabel: "応募",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrEntryStatusInvalid) {
		t.Errorf("error = %v, want ErrEntryStatusInvalid", err)
	}
}

func TestUpdate_InvalidStageKind(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     entity.NewUserID(),
		EntryID:    entity.NewEntryID(),
		Source:     "リクナビ",
		Status:     "in_progress",
		StageKind:  "invalid",
		StageLabel: "応募",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrStageKindInvalid) {
		t.Errorf("error = %v, want ErrStageKindInvalid", err)
	}
}

func TestUpdate_EmptyStageLabel(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     entity.NewUserID(),
		EntryID:    entity.NewEntryID(),
		Source:     "リクナビ",
		Status:     "in_progress",
		StageKind:  "application",
		StageLabel: "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrStageLabelEmpty) {
		t.Errorf("error = %v, want ErrStageLabelEmpty", err)
	}
}

func TestUpdate_SaveError(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestEntry(t, userID)
	saveErr := errors.New("db write failed")
	repo := &mockEntryRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) (*entity.Entry, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Entry) error {
			return saveErr
		},
	}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID:     userID,
		EntryID:    existing.ID(),
		Source:     "リクナビ",
		Status:     "in_progress",
		StageKind:  "application",
		StageLabel: "応募",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
