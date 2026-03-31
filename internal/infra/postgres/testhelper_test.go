//go:build integration

package postgres_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

const defaultTestDatabaseURL = "postgres://postgres:postgres@localhost:15432/job_hunting_test?sslmode=disable"

var (
	testPool    *pgxpool.Pool
	testInitErr error
	testOnce    sync.Once
)

// initTestDB はテスト用の pgxpool.Pool を初期化する。
// sync.Once 内から呼ばれるため、t.Fatalf ではなくエラーを返す。
func initTestDB() (*pgxpool.Pool, error) {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultTestDatabaseURL
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	schema, err := os.ReadFile("../../../sql/schema.sql")
	if err != nil {
		pool.Close()
		return nil, err
	}

	dropSQL := `
		DROP TABLE IF EXISTS password_credentials CASCADE;
		DROP TABLE IF EXISTS external_identities CASCADE;
		DROP TABLE IF EXISTS company_aliases CASCADE;
		DROP TABLE IF EXISTS stage_histories CASCADE;
		DROP TABLE IF EXISTS tasks CASCADE;
		DROP TABLE IF EXISTS entries CASCADE;
		DROP TABLE IF EXISTS companies CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
		DROP TYPE IF EXISTS entry_status CASCADE;
		DROP TYPE IF EXISTS stage_kind CASCADE;
		DROP TYPE IF EXISTS task_type CASCADE;
		DROP TYPE IF EXISTS task_status CASCADE;
		DROP TYPE IF EXISTS auth_provider CASCADE;
	`
	if _, err := pool.Exec(ctx, dropSQL); err != nil {
		pool.Close()
		return nil, err
	}
	if _, err := pool.Exec(ctx, string(schema)); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

// setupTestDB はテスト用の pgxpool.Pool を返す。
// 初回呼び出し時にスキーマを適用する。初期化エラーは全テストに伝播する。
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	testOnce.Do(func() {
		testPool, testInitErr = initTestDB()
	})

	if testInitErr != nil {
		t.Fatalf("test db init failed: %v", testInitErr)
	}

	return testPool
}

// beginTx はテスト用のトランザクションを開始する。
// t.Cleanup でロールバックが登録されるため、テスト後に自動でクリーンアップされる。
func beginTx(t *testing.T, pool *pgxpool.Pool) pgx.Tx {
	t.Helper()
	ctx := context.Background()

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})

	return tx
}

// insertTestUser はテスト用の User を直接 SQL INSERT する。
// FK 前提データとして使う。リポジトリアダプタではなく直接 INSERT。
func insertTestUser(t *testing.T, tx pgx.Tx) entity.UserID {
	t.Helper()
	ctx := context.Background()

	userID := entity.NewUserID()

	_, err := tx.Exec(ctx,
		`INSERT INTO users (id, email, name, created_at, updated_at)
		 VALUES ($1, $2, $3, now(), now())`,
		uuid.UUID(userID),
		uuid.New().String()+"@test.example.com",
		"テストユーザー",
	)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	return userID
}

// insertTestCompany はテスト用の Company を直接 SQL INSERT する。
// Entry/Task テストの FK 前提データとして使う。
func insertTestCompany(t *testing.T, tx pgx.Tx, userID entity.UserID) entity.CompanyID {
	t.Helper()
	ctx := context.Background()

	companyID := entity.NewCompanyID()
	_, err := tx.Exec(ctx,
		`INSERT INTO companies (id, user_id, name, memo, created_at, updated_at)
		 VALUES ($1, $2, $3, '', now(), now())`,
		uuid.UUID(companyID),
		uuid.UUID(userID),
		"テスト企業",
	)
	if err != nil {
		t.Fatalf("failed to insert test company: %v", err)
	}

	return companyID
}

// insertTestEntry はテスト用の Entry を直接 SQL INSERT する。
// Task テストの FK 前提データとして使う。
func insertTestEntry(t *testing.T, tx pgx.Tx, userID entity.UserID, companyID entity.CompanyID) entity.EntryID {
	t.Helper()
	ctx := context.Background()

	entryID := entity.NewEntryID()
	_, err := tx.Exec(ctx,
		`INSERT INTO entries (id, user_id, company_id, route, source, status, stage_kind, stage_label, memo, created_at, updated_at)
		 VALUES ($1, $2, $3, 'テストルート', 'テストソース', 'in_progress', 'application', '応募', '', now(), now())`,
		uuid.UUID(entryID),
		uuid.UUID(userID),
		uuid.UUID(companyID),
	)
	if err != nil {
		t.Fatalf("failed to insert test entry: %v", err)
	}

	return entryID
}
