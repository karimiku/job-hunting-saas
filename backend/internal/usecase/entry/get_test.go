package entry

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestEntry(t *testing.T, userID entity.UserID) *entity.Entry {
	t.Helper()
	route, err := value.NewRoute("本選考")
	if err != nil {
		t.Fatalf("NewRoute failed: %v", err)
	}
	source, err := value.NewSource("リクナビ")
	if err != nil {
		t.Fatalf("NewSource failed: %v", err)
	}
	return entity.NewEntry(userID, entity.NewCompanyID(), route, source)
}

func TestGet_Found(t *testing.T) {
	userID := entity.NewUserID()
	expected := newTestEntry(t, userID)
	repo := &mockEntryRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) (*entity.Entry, error) {
			return expected, nil
		},
	}

	uc := NewGet(repo)
	out, err := uc.Execute(context.Background(), GetInput{
		UserID:  userID,
		EntryID: expected.ID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entry != expected {
		t.Error("Entry should be the expected entry")
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := &mockEntryRepo{}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
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
	repo := &mockEntryRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) (*entity.Entry, error) {
			return nil, dbErr
		},
	}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
