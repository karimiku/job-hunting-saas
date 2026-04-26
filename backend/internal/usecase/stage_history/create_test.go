package stagehistory

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestCreate_Success(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	createCalled := false
	historyRepo := &mockHistoryRepo{
		createFn: func(_ context.Context, _ *entity.StageHistory) error {
			createCalled = true
			return nil
		},
	}

	uc := NewCreate(historyRepo, expectFindByID(t, userID, entryID))
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "interview",
		Label:     "一次面接",
		Note:      "オンライン",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.StageHistory == nil {
		t.Fatal("StageHistory should not be nil")
	}
	if out.StageHistory.Stage().Kind().String() != "interview" {
		t.Errorf("Kind = %q, want %q", out.StageHistory.Stage().Kind().String(), "interview")
	}
	if out.StageHistory.Stage().Label() != "一次面接" {
		t.Errorf("Label = %q, want %q", out.StageHistory.Stage().Label(), "一次面接")
	}
	if out.StageHistory.Note() != "オンライン" {
		t.Errorf("Note = %q, want %q", out.StageHistory.Note(), "オンライン")
	}
	if !createCalled {
		t.Error("Create should be called")
	}
}

func TestCreate_WithoutNote(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	uc := NewCreate(&mockHistoryRepo{}, expectFindByID(t, userID, entryID))
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "document",
		Label:     "ES提出",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.StageHistory.Note() != "" {
		t.Errorf("Note = %q, want empty", out.StageHistory.Note())
	}
}

func TestCreate_EntryNotFound(t *testing.T) {
	// Entry が見つからないとき、historyRepo が呼ばれないことも担保する
	uc := NewCreate(failOnCallHistoryRepo(t), &mockEntryRepo{})

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    entity.NewUserID(),
		EntryID:   entity.NewEntryID(),
		StageKind: "interview",
		Label:     "一次面接",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestCreate_EmptyStageKind(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	uc := NewCreate(failOnCallHistoryRepo(t), expectFindByID(t, userID, entryID))
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "",
		Label:     "一次面接",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrStageKindEmpty) {
		t.Errorf("error = %v, want ErrStageKindEmpty", err)
	}
}

func TestCreate_InvalidStageKind(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	uc := NewCreate(failOnCallHistoryRepo(t), expectFindByID(t, userID, entryID))
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "unknown_kind",
		Label:     "一次面接",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrStageKindInvalid) {
		t.Errorf("error = %v, want ErrStageKindInvalid", err)
	}
}

func TestCreate_EmptyLabel(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	uc := NewCreate(failOnCallHistoryRepo(t), expectFindByID(t, userID, entryID))
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "interview",
		Label:     "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrStageLabelEmpty) {
		t.Errorf("error = %v, want ErrStageLabelEmpty", err)
	}
}

func TestCreate_LabelWithSurroundingSpace(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	uc := NewCreate(failOnCallHistoryRepo(t), expectFindByID(t, userID, entryID))
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "interview",
		Label:     " 一次面接 ",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrStageLabelInvalid) {
		t.Errorf("error = %v, want ErrStageLabelInvalid", err)
	}
}

func TestCreate_RepoError(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	createErr := errors.New("db write failed")
	historyRepo := &mockHistoryRepo{
		createFn: func(_ context.Context, _ *entity.StageHistory) error {
			return createErr
		},
	}

	uc := NewCreate(historyRepo, expectFindByID(t, userID, entryID))
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:    userID,
		EntryID:   entryID,
		StageKind: "interview",
		Label:     "一次面接",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, createErr) {
		t.Errorf("error = %v, want createErr", err)
	}
}
