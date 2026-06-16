package mcp

import (
	"context"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
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

func TestServiceDeleteEntry_RequiresConfirmation(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	entryRepo := &fakeEntryRepo{userID: userID}
	service := newTestServiceWithEntryRepo(userID, &fakeContextQuery{entryID: entryID, company: "Example Inc."}, entryRepo, &recordingMemoRepo{}, &recordingTaskRepo{})

	out, err := service.DeleteEntry(context.Background(), DeleteEntryInput{
		EntryID: entryID.String(),
	})
	if err != nil {
		t.Fatalf("DeleteEntry() failed: %v", err)
	}
	if entryRepo.deleted {
		t.Fatal("entry should not be deleted without confirm=true")
	}
	got := out.(map[string]any)
	if got["confirmationRequired"] != true {
		t.Errorf("confirmationRequired = %v, want true", got["confirmationRequired"])
	}
	if got["relatedTaskCount"] != 1 {
		t.Errorf("relatedTaskCount = %v, want 1", got["relatedTaskCount"])
	}
}

func TestServiceDeleteEntry_WithConfirmationDeletes(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	entryRepo := &fakeEntryRepo{userID: userID}
	service := newTestServiceWithEntryRepo(userID, &fakeContextQuery{entryID: entryID, company: "Example Inc."}, entryRepo, &recordingMemoRepo{}, &recordingTaskRepo{})

	out, err := service.DeleteEntry(context.Background(), DeleteEntryInput{
		EntryID: entryID.String(),
		Confirm: true,
	})
	if err != nil {
		t.Fatalf("DeleteEntry() failed: %v", err)
	}
	if !entryRepo.deleted {
		t.Fatal("entry should be deleted with confirm=true")
	}
	if entryRepo.deletedEntryID != entryID {
		t.Errorf("deletedEntryID = %v, want %v", entryRepo.deletedEntryID, entryID)
	}
	got := out.(map[string]any)
	if got["deleted"] != true {
		t.Errorf("deleted = %v, want true", got["deleted"])
	}
}

func TestServiceListESMemos_ReturnsCompanyForEntryMemo(t *testing.T) {
	userID := entity.NewUserID()
	entryID := entity.NewEntryID()
	memoRepo := &recordingMemoRepo{
		items: []*entity.ESMemo{
			newTestESMemo(t, userID, &entryID, "interview", "面接で話す改善経験", "顧客課題を分解して改善した"),
		},
	}
	service := newTestService(userID, &fakeContextQuery{entryID: entryID, company: "Example Inc."}, memoRepo, &recordingTaskRepo{})

	out, err := service.ListESMemos(context.Background(), 10)
	if err != nil {
		t.Fatalf("ListESMemos() failed: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("len(out) = %d, want 1", len(out))
	}
	if out[0].Company == nil || *out[0].Company != "Example Inc." {
		t.Fatalf("Company = %v, want Example Inc.", out[0].Company)
	}
	if out[0].Title != "面接で話す改善経験" {
		t.Errorf("Title = %q, want 面接で話す改善経験", out[0].Title)
	}
}

func newTestService(userID entity.UserID, query ContextQuery, memoRepo *recordingMemoRepo, taskRepo *recordingTaskRepo) *Service {
	entryRepo := &fakeEntryRepo{userID: userID}
	return newTestServiceWithEntryRepo(userID, query, entryRepo, memoRepo, taskRepo)
}

func newTestServiceWithEntryRepo(userID entity.UserID, query ContextQuery, entryRepo *fakeEntryRepo, memoRepo *recordingMemoRepo, taskRepo *recordingTaskRepo) *Service {
	return NewService(
		userID,
		query,
		esmemo.NewAppend(memoRepo, entryRepo),
		esmemo.NewList(memoRepo),
		taskuc.NewCreate(taskRepo, entryRepo),
		jobemail.NewExtract(),
		nil,
		entryuc.NewDelete(entryRepo),
		nil,
		nil,
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
		Tasks: []TaskDTO{{ID: "task-1", EntryID: entryID.String(), Title: "ES提出"}},
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
	items []*entity.ESMemo
}

func (r *recordingMemoRepo) Save(_ context.Context, memo *entity.ESMemo) error {
	r.saved = memo
	return nil
}

func (r *recordingMemoRepo) ListByUserID(context.Context, entity.UserID, int32) ([]*entity.ESMemo, error) {
	return r.items, nil
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

func (r *recordingTaskRepo) ListByUserID(context.Context, entity.UserID) ([]*entity.Task, error) {
	return []*entity.Task{}, nil
}

func (r *recordingTaskRepo) ListByUserIDWithDueBefore(context.Context, entity.UserID, time.Time) ([]*entity.Task, error) {
	return []*entity.Task{}, nil
}

func (r *recordingTaskRepo) Delete(context.Context, entity.UserID, entity.TaskID) error {
	return nil
}

type fakeEntryRepo struct {
	userID         entity.UserID
	deleted        bool
	deletedEntryID entity.EntryID
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

func (r *fakeEntryRepo) Delete(_ context.Context, userID entity.UserID, id entity.EntryID) error {
	if userID != r.userID {
		return repository.ErrNotFound
	}
	r.deleted = true
	r.deletedEntryID = id
	return nil
}

func newTestESMemo(t *testing.T, userID entity.UserID, entryID *entity.EntryID, category string, title string, content string) *entity.ESMemo {
	t.Helper()
	cat, err := value.NewESMemoCategory(category)
	if err != nil {
		t.Fatal(err)
	}
	memoTitle, err := value.NewESMemoTitle(title)
	if err != nil {
		t.Fatal(err)
	}
	memoContent, err := value.NewESMemoContent(content)
	if err != nil {
		t.Fatal(err)
	}
	source, err := value.NewESMemoSource("mcp")
	if err != nil {
		t.Fatal(err)
	}
	return entity.NewESMemo(userID, entryID, cat, memoTitle, memoContent, source)
}
