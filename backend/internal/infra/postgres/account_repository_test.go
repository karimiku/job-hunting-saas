//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func TestAccountRepository_DeleteAccount_CascadesUserData(t *testing.T) {
	pool := setupTestDB(t)
	ctx := context.Background()
	userID := insertCommittedTestUser(t, pool)
	ids := insertCommittedAccountData(t, pool, userID)
	repo := postgres.NewAccountRepository(pool)

	if err := repo.DeleteAccount(ctx, userID); err != nil {
		t.Fatalf("DeleteAccount failed: %v", err)
	}

	assertZeroCount(t, pool, `SELECT count(*) FROM users WHERE id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM companies WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM entries WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM company_aliases WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM external_identities WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM password_credentials WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM inbox_clips WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM es_memos WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM ai_access_tokens WHERE user_id = $1`, uuid.UUID(userID))
	assertZeroCount(t, pool, `SELECT count(*) FROM tasks WHERE entry_id = $1`, ids.entryID)
	assertZeroCount(t, pool, `SELECT count(*) FROM stage_histories WHERE entry_id = $1`, ids.entryID)
}

func TestAccountRepository_DeleteAccount_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewAccountRepository(pool)

	err := repo.DeleteAccount(context.Background(), entity.NewUserID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

type committedAccountDataIDs struct {
	entryID uuid.UUID
}

func insertCommittedAccountData(t *testing.T, pool *pgxpool.Pool, userID entity.UserID) committedAccountDataIDs {
	t.Helper()
	ctx := context.Background()
	companyID := uuid.New()
	entryID := uuid.New()
	taskID := uuid.New()
	stageHistoryID := uuid.New()
	aliasID := uuid.New()
	externalIdentityID := uuid.New()
	passwordCredentialID := uuid.New()
	clipID := uuid.New()
	memoID := uuid.New()
	tokenID := uuid.New()

	statements := []struct {
		sql  string
		args []any
	}{
		{
			sql: `INSERT INTO companies (id, user_id, name, created_at, updated_at)
				  VALUES ($1, $2, '退会企業', now(), now())`,
			args: []any{companyID, uuid.UUID(userID)},
		},
		{
			sql: `INSERT INTO entries (id, user_id, company_id, route, source, stage_label, created_at, updated_at)
				  VALUES ($1, $2, $3, '本選考', 'リクナビ', '応募', now(), now())`,
			args: []any{entryID, uuid.UUID(userID), companyID},
		},
		{
			sql: `INSERT INTO tasks (id, entry_id, title, task_type, created_at, updated_at)
				  VALUES ($1, $2, '締切', 'deadline', now(), now())`,
			args: []any{taskID, entryID},
		},
		{
			sql: `INSERT INTO stage_histories (id, entry_id, stage_kind, stage_label, created_at)
				  VALUES ($1, $2, 'application', '応募', now())`,
			args: []any{stageHistoryID, entryID},
		},
		{
			sql: `INSERT INTO company_aliases (id, user_id, company_id, alias, created_at)
				  VALUES ($1, $2, $3, '別名', now())`,
			args: []any{aliasID, uuid.UUID(userID), companyID},
		},
		{
			sql: `INSERT INTO external_identities (id, user_id, provider, subject, created_at)
				  VALUES ($1, $2, 'google', $3, now())`,
			args: []any{externalIdentityID, uuid.UUID(userID), uuid.New().String()},
		},
		{
			sql: `INSERT INTO password_credentials (id, user_id, password_hash, created_at, updated_at)
				  VALUES ($1, $2, 'hash', now(), now())`,
			args: []any{passwordCredentialID, uuid.UUID(userID)},
		},
		{
			sql: `INSERT INTO inbox_clips (id, user_id, url, title, source, captured_at)
				  VALUES ($1, $2, $3, '求人', 'リクナビ', now())`,
			args: []any{clipID, uuid.UUID(userID), "https://example.com/" + uuid.NewString()},
		},
		{
			sql: `INSERT INTO es_memos (id, user_id, entry_id, title, content, created_at, updated_at)
				  VALUES ($1, $2, $3, 'ES', '内容', now(), now())`,
			args: []any{memoID, uuid.UUID(userID), entryID},
		},
		{
			sql: `INSERT INTO ai_access_tokens (id, user_id, name, token_hash, token_prefix, created_at)
				  VALUES ($1, $2, 'AI', $3, 'entre_ai_', now())`,
			args: []any{tokenID, uuid.UUID(userID), uuid.NewString()},
		},
	}
	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt.sql, stmt.args...); err != nil {
			t.Fatalf("failed to seed account data: %v", err)
		}
	}

	return committedAccountDataIDs{entryID: entryID}
}

func assertZeroCount(t *testing.T, pool *pgxpool.Pool, sql string, arg any) {
	t.Helper()
	var count int
	if err := pool.QueryRow(context.Background(), sql, arg).Scan(&count); err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("count = %d for query %q, want 0", count, sql)
	}
}
