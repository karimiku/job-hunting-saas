package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	esmemo "github.com/karimiku/job-hunting-saas/internal/usecase/es_memo"
	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

func TestServiceAppendESMemo_RequiresConfirmation(t *testing.T) {
	memoRepo := &recordingMemoRepo{}
	service := newTestService(entity.NewUserID(), &fakeContextQuery{}, memoRepo, &recordingTaskRepo{})

	out, err := service.AppendESMemo(context.Background(), AppendESMemoInput{
		Title:   "面接で話す改善経験",
		Content: "顧客課題を分解して改善した",
	})
	if err != nil {
		t.Fatalf("AppendESMemo() failed: %v", err)
	}
	if memoRepo.saved != nil {
		t.Fatal("memo should not be saved without confirm=true")
	}
	got := out.(map[string]any)
	if got["confirmationRequired"] != true {
		t.Errorf("confirmationRequired = %v, want true", got["confirmationRequired"])
	}
}

func TestServiceAppendESMemo_WithConfirmationSaves(t *testing.T) {
	memoRepo := &recordingMemoRepo{}
	service := newTestService(entity.NewUserID(), &fakeContextQuery{}, memoRepo, &recordingTaskRepo{})

	out, err := service.AppendESMemo(context.Background(), AppendESMemoInput{
		Title:   "面接で話す改善経験",
		Content: "顧客課題を分解して改善した",
		Confirm: true,
	})
	if err != nil {
		t.Fatalf("AppendESMemo() failed: %v", err)
	}
	if memoRepo.saved == nil {
		t.Fatal("memo should be saved with confirm=true")
	}
	got := out.(map[string]any)
	if got["created"] != true {
		t.Errorf("created = %v, want true", got["created"])
	}
}

func TestServiceCreateTask_RequiresConfirmation(t *testing.T) {
	taskRepo := &recordingTaskRepo{}
	entryID := entity.NewEntryID()
	service := newTestService(entity.NewUserID(), &fakeContextQuery{entryID: entryID, company: "Example Inc."}, &recordingMemoRepo{}, taskRepo)

	out, err := service.CreateTask(context.Background(), CreateTaskInput{
		EntryID: entryID.String(),
		Title:   "ES提出",
		Type:    "deadline",
		DueDate: "2026-06-10",
	})
	if err != nil {
		t.Fatalf("CreateTask() failed: %v", err)
	}
	if taskRepo.saved != nil {
		t.Fatal("task should not be saved without confirm=true")
	}
	got := out.(map[string]any)
	if got["confirmationRequired"] != true {
		t.Errorf("confirmationRequired = %v, want true", got["confirmationRequired"])
	}
	task := got["task"].(map[string]any)
	if task["company"] != "Example Inc." {
		t.Errorf("company = %v, want Example Inc.", task["company"])
	}
}

func TestServiceCreateTask_WithConfirmationSaves(t *testing.T) {
	taskRepo := &recordingTaskRepo{}
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	service := newTestService(userID, &fakeContextQuery{entryID: entryID, company: "Example Inc."}, &recordingMemoRepo{}, taskRepo)

	out, err := service.CreateTask(context.Background(), CreateTaskInput{
		EntryID: entryID.String(),
		Title:   "ES提出",
		Type:    "deadline",
		Confirm: true,
		Notify:  true,
	})
	if err != nil {
		t.Fatalf("CreateTask() failed: %v", err)
	}
	if taskRepo.saved == nil {
		t.Fatal("task should be saved with confirm=true")
	}
	if !taskRepo.saved.Notify() {
		t.Error("Notify() = false, want true")
	}
	got := out.(map[string]any)
	if got["created"] != true {
		t.Errorf("created = %v, want true", got["created"])
	}
}

func newTestService(userID entity.UserID, query ContextQuery, memoRepo *recordingMemoRepo, taskRepo *recordingTaskRepo) *Service {
	entryRepo := &fakeEntryRepo{userID: userID}
	return NewService(
		userID,
		query,
		esmemo.NewAppend(memoRepo, entryRepo),
		taskuc.NewCreate(taskRepo, entryRepo),
		jobemail.NewExtract(),
	)
}

type fakeContextQuery struct {
	entryID entity.EntryID
	company string
}

func (q *fakeContextQuery) ListEntries(context.Context, entity.UserID) ([]EntryDTO, error) {
	return []EntryDTO{}, nil
}

func (q *fakeContextQuery) GetEntryContext(_ context.Context, _ entity.UserID, entryID entity.EntryID) (*EntryContextDTO, error) {
	company := q.company
	if company == "" {
		company = "Example Inc."
	}
	return &EntryContextDTO{
		Entry: EntryDTO{
			ID:      entryID.String(),
			Company: company,
		},
		Tasks: []TaskDTO{},
	}, nil
}

func (q *fakeContextQuery) ListOpenTasks(context.Context, entity.UserID) ([]TaskDTO, error) {
	return []TaskDTO{}, nil
}

func (q *fakeContextQuery) ListInboxClips(context.Context, entity.UserID) ([]InboxClipDTO, error) {
	return []InboxClipDTO{}, nil
}

type recordingMemoRepo struct {
	saved *entity.ESMemo
}

func (r *recordingMemoRepo) Save(_ context.Context, memo *entity.ESMemo) error {
	r.saved = memo
	return nil
}

func (r *recordingMemoRepo) ListByUserID(context.Context, entity.UserID, int32) ([]*entity.ESMemo, error) {
	return []*entity.ESMemo{}, nil
}

type recordingTaskRepo struct {
	saved *entity.Task
}

func (r *recordingTaskRepo) Save(_ context.Context, task *entity.Task) error {
	r.saved = task
	return nil
}

func (r *recordingTaskRepo) FindByID(context.Context, entity.UserID, entity.TaskID) (*entity.Task, error) {
	return nil, repository.ErrNotFound
}

func (r *recordingTaskRepo) ListByEntryID(context.Context, entity.UserID, entity.EntryID) ([]*entity.Task, error) {
	return []*entity.Task{}, nil
}

func (r *recordingTaskRepo) ListByUserIDWithDueBefore(context.Context, entity.UserID, time.Time) ([]*entity.Task, error) {
	return []*entity.Task{}, nil
}

func (r *recordingTaskRepo) Delete(context.Context, entity.UserID, entity.TaskID) error {
	return nil
}

type fakeEntryRepo struct {
	userID entity.UserID
}

func (r *fakeEntryRepo) Save(context.Context, *entity.Entry) error {
	return nil
}

func (r *fakeEntryRepo) FindByID(_ context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	if userID != r.userID {
		return nil, repository.ErrNotFound
	}
	route, err := value.NewRoute("direct")
	if err != nil {
		return nil, err
	}
	source, err := value.NewSource("manual")
	if err != nil {
		return nil, err
	}
	return entity.ReconstructEntry(
		id,
		userID,
		entity.NewCompanyID(),
		route,
		source,
		value.EntryStatusInProgress(),
		nil,
		value.MustNewStage(value.StageKindApplication(), "応募"),
		"",
		time.Now(),
		time.Now(),
	), nil
}

func (r *fakeEntryRepo) ListByUserID(context.Context, entity.UserID, repository.EntryFilter) ([]*entity.Entry, error) {
	return []*entity.Entry{}, nil
}

func (r *fakeEntryRepo) Delete(context.Context, entity.UserID, entity.EntryID) error {
	return nil
}
