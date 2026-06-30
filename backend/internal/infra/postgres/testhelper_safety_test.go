//go:build integration

package postgres_test

import "testing"

func TestAssertSafeDestructiveTestDatabaseAllowsLocalHosts(t *testing.T) {
	t.Setenv(destructiveIntegrationTestOverrideEnv, "")

	testCases := []string{
		"postgres://postgres:postgres@localhost:15432/job_hunting_test?sslmode=disable",
		"postgres://postgres:postgres@127.0.0.1:15432/job_hunting_test?sslmode=disable",
		"postgres://postgres:postgres@[::1]:15432/job_hunting_test?sslmode=disable",
		"postgres://postgres:postgres@db:5432/job_hunting_test?sslmode=disable",
		"postgres://postgres:postgres@postgres:5432/job_hunting_test?sslmode=disable",
	}

	for _, databaseURL := range testCases {
		t.Run(databaseURL, func(t *testing.T) {
			if err := assertSafeDestructiveTestDatabase(databaseURL); err != nil {
				t.Fatalf("expected local database URL to be allowed: %v", err)
			}
		})
	}
}

func TestAssertSafeDestructiveTestDatabaseRejectsRemoteHosts(t *testing.T) {
	t.Setenv(destructiveIntegrationTestOverrideEnv, "")

	testCases := []string{
		"postgres://postgres:secret@db.example.com:5432/postgres?sslmode=require",
		"postgres://postgres:secret@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres?sslmode=require",
		"postgres://postgres:secret@db.project-ref.supabase.co:5432/postgres?sslmode=require",
	}

	for _, databaseURL := range testCases {
		t.Run(databaseURL, func(t *testing.T) {
			if err := assertSafeDestructiveTestDatabase(databaseURL); err == nil {
				t.Fatal("expected remote database URL to be rejected")
			}
		})
	}
}

func TestAssertSafeDestructiveTestDatabaseAllowsExplicitOverride(t *testing.T) {
	t.Setenv(destructiveIntegrationTestOverrideEnv, "true")

	databaseURL := "postgres://postgres:secret@db.project-ref.supabase.co:5432/postgres?sslmode=require"
	if err := assertSafeDestructiveTestDatabase(databaseURL); err != nil {
		t.Fatalf("expected override to allow remote database URL: %v", err)
	}
}
