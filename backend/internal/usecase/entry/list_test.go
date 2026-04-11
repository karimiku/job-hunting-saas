package entry

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestList_Multiple(t *testing.T) {
	userID := entity.NewUserID()
	e1 := newTestEntry(t, userID)
	e2 := newTestEntry(t, userID)
	repo := &mockEntryRepo{
		listByUserFn: func(_ context.Context, _ entity.UserID, _ repository.EntryFilter) ([]*entity.Entry, error) {
			return []*entity.Entry{e1, e2}, nil
		},
	}

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{
		UserID: userID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Fatalf("len(Entries) = %d, want 2", len(out.Entries))
	}
}

func TestList_Empty(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{
		UserID: entity.NewUserID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 0 {
		t.Fatalf("len(Entries) = %d, want 0", len(out.Entries))
	}
}

func TestList_WithFilter(t *testing.T) {
	userID := entity.NewUserID()
	e1 := newTestEntry(t, userID)
	var capturedFilter repository.EntryFilter
	repo := &mockEntryRepo{
		listByUserFn: func(_ context.Context, _ entity.UserID, filter repository.EntryFilter) ([]*entity.Entry, error) {
			capturedFilter = filter
			return []*entity.Entry{e1}, nil
		},
	}

	status := "in_progress"
	stageKind := "interview"
	source := "リクナビ"

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:    userID,
		Status:    &status,
		StageKind: &stageKind,
		Source:    &source,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Entries) != 1 {
		t.Fatalf("len(Entries) = %d, want 1", len(out.Entries))
	}
	if capturedFilter.Status == nil {
		t.Fatal("filter.Status should not be nil")
	}
	if capturedFilter.Status.String() != "in_progress" {
		t.Errorf("filter.Status = %q, want %q", capturedFilter.Status.String(), "in_progress")
	}
	if capturedFilter.StageKind == nil {
		t.Fatal("filter.StageKind should not be nil")
	}
	if capturedFilter.StageKind.String() != "interview" {
		t.Errorf("filter.StageKind = %q, want %q", capturedFilter.StageKind.String(), "interview")
	}
	if capturedFilter.Source == nil {
		t.Fatal("filter.Source should not be nil")
	}
	if capturedFilter.Source.String() != "リクナビ" {
		t.Errorf("filter.Source = %q, want %q", capturedFilter.Source.String(), "リクナビ")
	}
}

func TestList_InvalidStatus(t *testing.T) {
	repo := &mockEntryRepo{}
	status := "invalid_status"

	uc := NewList(repo)
	_, err := uc.Execute(context.Background(), ListInput{
		UserID: entity.NewUserID(),
		Status: &status,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrEntryStatusInvalid) {
		t.Errorf("error = %v, want ErrEntryStatusInvalid", err)
	}
}

func TestList_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockEntryRepo{
		listByUserFn: func(_ context.Context, _ entity.UserID, _ repository.EntryFilter) ([]*entity.Entry, error) {
			return nil, dbErr
		},
	}

	uc := NewList(repo)
	_, err := uc.Execute(context.Background(), ListInput{
		UserID: entity.NewUserID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
