//go:build integration

package postgres_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func TestDomainUniqueConstraintsMigration_IdempotentAndNormalizesLegacyInboxClips(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	if _, err := tx.Exec(ctx, `ALTER TABLE inbox_clips DROP CONSTRAINT inbox_clips_user_id_url_key`); err != nil {
		t.Fatalf("drop inbox_clips constraint: %v", err)
	}
	if _, err := tx.Exec(ctx, `ALTER TABLE company_aliases DROP CONSTRAINT company_aliases_user_id_company_id_alias_key`); err != nil {
		t.Fatalf("drop company_aliases constraint: %v", err)
	}

	userID := insertTestUser(t, tx)
	trimmedID := entity.NewInboxClipID()
	fallbackID := entity.NewInboxClipID()
	fallbackURL := "https://example.com/jobs/legacy-empty-title"

	if _, err := tx.Exec(ctx, `
		INSERT INTO inbox_clips (id, user_id, url, title, source, guess, captured_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
	`, uuid.UUID(trimmedID), uuid.UUID(userID), "https://example.com/jobs/legacy-trim", "  Legacy title  ", "マイナビ", "  Legacy Co  "); err != nil {
		t.Fatalf("insert legacy trimmed clip: %v", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO inbox_clips (id, user_id, url, title, source, guess, captured_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
	`, uuid.UUID(fallbackID), uuid.UUID(userID), fallbackURL, " \t\n ", "マイナビ", " "+strings.Repeat("あ", 300)+" "); err != nil {
		t.Fatalf("insert legacy fallback clip: %v", err)
	}

	migration, err := os.ReadFile("../../../sql/migrations/20260614_add_domain_unique_constraints.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	for i := 0; i < 2; i++ {
		if _, err := tx.Exec(ctx, string(migration)); err != nil {
			t.Fatalf("migration execution %d failed: %v", i+1, err)
		}
	}

	repo := postgres.NewInboxClipRepository(tx)
	trimmed, err := repo.FindByID(ctx, userID, trimmedID)
	if err != nil {
		t.Fatalf("FindByID trimmed clip: %v", err)
	}
	if got := trimmed.Title().String(); got != "Legacy title" {
		t.Errorf("trimmed title = %q, want %q", got, "Legacy title")
	}
	if got := trimmed.Guess().String(); got != "Legacy Co" {
		t.Errorf("trimmed guess = %q, want %q", got, "Legacy Co")
	}

	fallback, err := repo.FindByID(ctx, userID, fallbackID)
	if err != nil {
		t.Fatalf("FindByID fallback clip: %v", err)
	}
	if got := fallback.Title().String(); got != fallbackURL {
		t.Errorf("fallback title = %q, want %q", got, fallbackURL)
	}
	if got := utf8.RuneCountInString(fallback.Guess().String()); got != value.InboxClipGuessMaxLength {
		t.Errorf("fallback guess length = %d, want %d", got, value.InboxClipGuessMaxLength)
	}
}
