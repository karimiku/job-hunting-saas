package service

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestInboxClipValues(t *testing.T) (value.URL, value.InboxClipTitle, value.Source, value.InboxClipGuess, value.InboxClipContentText) {
	t.Helper()

	url, err := value.NewURL("https://example.com/jobs/1")
	if err != nil {
		t.Fatalf("NewURL: %v", err)
	}
	title, err := value.NewInboxClipTitle("株式会社サンプル 募集要項")
	if err != nil {
		t.Fatalf("NewInboxClipTitle: %v", err)
	}
	source, err := value.NewSource("マイナビ")
	if err != nil {
		t.Fatalf("NewSource: %v", err)
	}
	guess, err := value.NewInboxClipGuess("株式会社サンプル")
	if err != nil {
		t.Fatalf("NewInboxClipGuess: %v", err)
	}
	contentText, err := value.NewInboxClipContentText("選考フロー: ES提出、一次面接")
	if err != nil {
		t.Fatalf("NewInboxClipContentText: %v", err)
	}
	return url, title, source, guess, contentText
}

func TestInboxClipRegistrationService_Register_ReturnsExisting(t *testing.T) {
	userID := entity.NewUserID()
	url, title, source, guess, contentText := newTestInboxClipValues(t)
	existing := entity.NewInboxClip(userID, url, title, source, guess, contentText)
	repo := &recordingInboxClipRepo{existing: existing}
	service := NewInboxClipRegistrationService(repo)

	got, err := service.Register(context.Background(), userID, url, title, source, guess, contentText)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got != existing {
		t.Fatal("Register should return existing clip")
	}
	if repo.createCalled {
		t.Fatal("Create should not be called when existing clip is found")
	}
}

func TestInboxClipRegistrationService_Register_CreatesNew(t *testing.T) {
	userID := entity.NewUserID()
	url, title, source, guess, contentText := newTestInboxClipValues(t)
	repo := &recordingInboxClipRepo{}
	service := NewInboxClipRegistrationService(repo)

	got, err := service.Register(context.Background(), userID, url, title, source, guess, contentText)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got == nil {
		t.Fatal("Register should return created clip")
	}
	if repo.created != got {
		t.Fatal("created clip should be returned")
	}
	if got.UserID() != userID {
		t.Errorf("UserID = %v, want %v", got.UserID(), userID)
	}
}

func TestInboxClipRegistrationService_Register_AlreadyExistsFetchesExisting(t *testing.T) {
	userID := entity.NewUserID()
	url, title, source, guess, contentText := newTestInboxClipValues(t)
	existing := entity.NewInboxClip(userID, url, title, source, guess, contentText)
	repo := &recordingInboxClipRepo{
		createErr:           repository.ErrAlreadyExists,
		existingAfterCreate: existing,
	}
	service := NewInboxClipRegistrationService(repo)

	got, err := service.Register(context.Background(), userID, url, title, source, guess, contentText)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got != existing {
		t.Fatal("Register should fetch existing clip after ErrAlreadyExists")
	}
}

func TestInboxClipRegistrationService_Register_PropagatesFindError(t *testing.T) {
	userID := entity.NewUserID()
	url, title, source, guess, contentText := newTestInboxClipValues(t)
	expected := errors.New("db failed")
	service := NewInboxClipRegistrationService(&recordingInboxClipRepo{findErr: expected})

	_, err := service.Register(context.Background(), userID, url, title, source, guess, contentText)
	if !errors.Is(err, expected) {
		t.Fatalf("error = %v, want %v", err, expected)
	}
}

type recordingInboxClipRepo struct {
	existing            *entity.InboxClip
	existingAfterCreate *entity.InboxClip
	findErr             error
	createErr           error
	createCalled        bool
	created             *entity.InboxClip
}

func (r *recordingInboxClipRepo) Create(_ context.Context, clip *entity.InboxClip) error {
	r.createCalled = true
	if r.createErr != nil {
		return r.createErr
	}
	r.created = clip
	return nil
}

func (r *recordingInboxClipRepo) FindByID(context.Context, entity.UserID, entity.InboxClipID) (*entity.InboxClip, error) {
	return nil, repository.ErrNotFound
}

func (r *recordingInboxClipRepo) FindByUserIDAndURL(context.Context, entity.UserID, value.URL) (*entity.InboxClip, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if r.existing != nil {
		return r.existing, nil
	}
	if r.createCalled && r.existingAfterCreate != nil {
		return r.existingAfterCreate, nil
	}
	return nil, repository.ErrNotFound
}

func (r *recordingInboxClipRepo) ListByUserID(context.Context, entity.UserID) ([]*entity.InboxClip, error) {
	return nil, nil
}

func (r *recordingInboxClipRepo) Delete(context.Context, entity.UserID, entity.InboxClipID) error {
	return nil
}
