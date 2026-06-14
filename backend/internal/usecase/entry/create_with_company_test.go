package entry

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

func TestCreateWithCompany_Success(t *testing.T) {
	userID := entity.NewUserID()
	var savedCompany *entity.Company
	var savedEntry *entity.Entry
	uc := NewCreateWithCompany(&mockEntryWithCompanyRepo{
		saveFn: func(_ context.Context, company *entity.Company, entry *entity.Entry) error {
			savedCompany = company
			savedEntry = entry
			return nil
		},
	})

	out, err := uc.Execute(context.Background(), CreateWithCompanyInput{
		UserID:      userID,
		CompanyName: "テスト企業",
		Route:       "本選考",
		Source:      "リクナビ",
		SourceURL:   "https://job.rikunabi.com/2027/company/r123/",
		Memo:        "memo",
	})
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if out.Company != savedCompany || out.Entry != savedEntry {
		t.Fatal("output should contain saved company and entry")
	}
	if savedCompany.UserID() != userID {
		t.Errorf("company.UserID = %v, want %v", savedCompany.UserID(), userID)
	}
	if savedEntry.UserID() != userID {
		t.Errorf("entry.UserID = %v, want %v", savedEntry.UserID(), userID)
	}
	if savedEntry.CompanyID() != savedCompany.ID() {
		t.Errorf("entry.CompanyID = %v, want %v", savedEntry.CompanyID(), savedCompany.ID())
	}
	if savedEntry.SourceURL() == nil {
		t.Fatal("entry.SourceURL should be set")
	}
	if savedEntry.Memo() != "memo" {
		t.Errorf("entry.Memo = %q, want memo", savedEntry.Memo())
	}
}

func TestCreateWithCompany_InvalidCompanyName(t *testing.T) {
	called := false
	uc := NewCreateWithCompany(&mockEntryWithCompanyRepo{
		saveFn: func(_ context.Context, _ *entity.Company, _ *entity.Entry) error {
			called = true
			return nil
		},
	})

	_, err := uc.Execute(context.Background(), CreateWithCompanyInput{
		UserID:      entity.NewUserID(),
		CompanyName: "",
		Route:       "本選考",
		Source:      "リクナビ",
	})
	if err == nil {
		t.Fatal("Execute should return validation error")
	}
	if called {
		t.Fatal("repository should not be called on validation error")
	}
}

func TestCreateWithCompany_InvalidSourceURL(t *testing.T) {
	called := false
	uc := NewCreateWithCompany(&mockEntryWithCompanyRepo{
		saveFn: func(_ context.Context, _ *entity.Company, _ *entity.Entry) error {
			called = true
			return nil
		},
	})

	_, err := uc.Execute(context.Background(), CreateWithCompanyInput{
		UserID:      entity.NewUserID(),
		CompanyName: "テスト企業",
		Route:       "本選考",
		Source:      "リクナビ",
		SourceURL:   "not a url",
	})
	if err == nil {
		t.Fatal("Execute should return validation error")
	}
	if called {
		t.Fatal("repository should not be called on validation error")
	}
}

func TestCreateWithCompany_RepositoryError(t *testing.T) {
	expected := errors.New("db failed")
	uc := NewCreateWithCompany(&mockEntryWithCompanyRepo{
		saveFn: func(_ context.Context, _ *entity.Company, _ *entity.Entry) error {
			return expected
		},
	})

	_, err := uc.Execute(context.Background(), CreateWithCompanyInput{
		UserID:      entity.NewUserID(),
		CompanyName: "テスト企業",
		Route:       "本選考",
		Source:      "リクナビ",
	})
	if !errors.Is(err, expected) {
		t.Fatalf("err = %v, want %v", err, expected)
	}
}
