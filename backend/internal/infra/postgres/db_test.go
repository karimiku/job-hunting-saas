package postgres

import (
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

const testDatabaseURL = "postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable"

func TestNewPoolConfigDefaultsForVercelSupabaseRuntime(t *testing.T) {
	clearPoolEnv(t)

	config, err := NewPoolConfig(testDatabaseURL)
	if err != nil {
		t.Fatalf("NewPoolConfig returned error: %v", err)
	}

	if config.MaxConns != defaultMaxConns {
		t.Errorf("MaxConns = %d, want %d", config.MaxConns, defaultMaxConns)
	}
	if config.MinConns != defaultMinConns {
		t.Errorf("MinConns = %d, want %d", config.MinConns, defaultMinConns)
	}
	if config.MaxConnIdleTime != defaultMaxConnIdleTime {
		t.Errorf("MaxConnIdleTime = %s, want %s", config.MaxConnIdleTime, defaultMaxConnIdleTime)
	}
	if config.MaxConnLifetime != defaultMaxConnLifetime {
		t.Errorf("MaxConnLifetime = %s, want %s", config.MaxConnLifetime, defaultMaxConnLifetime)
	}
	if config.HealthCheckPeriod != defaultHealthCheckPeriod {
		t.Errorf("HealthCheckPeriod = %s, want %s", config.HealthCheckPeriod, defaultHealthCheckPeriod)
	}
	if config.ConnConfig.DefaultQueryExecMode != pgx.QueryExecModeExec {
		t.Errorf("DefaultQueryExecMode = %v, want %v", config.ConnConfig.DefaultQueryExecMode, pgx.QueryExecModeExec)
	}
	if got := config.ConnConfig.RuntimeParams["application_name"]; got != defaultApplicationName {
		t.Errorf("application_name = %q, want %q", got, defaultApplicationName)
	}
}

func TestNewPoolConfigReadsEnvOverrides(t *testing.T) {
	clearPoolEnv(t)
	t.Setenv("PGPOOL_MAX_CONNS", "8")
	t.Setenv("PGPOOL_MIN_CONNS", "2")
	t.Setenv("PGPOOL_MAX_CONN_IDLE_TIME", "45s")
	t.Setenv("PGPOOL_MAX_CONN_LIFETIME", "10m")
	t.Setenv("PGPOOL_HEALTH_CHECK_PERIOD", "15s")
	t.Setenv("PGX_DEFAULT_QUERY_EXEC_MODE", "simple_protocol")
	t.Setenv("PGAPPNAME", "entre-api-preview")

	config, err := NewPoolConfig(testDatabaseURL)
	if err != nil {
		t.Fatalf("NewPoolConfig returned error: %v", err)
	}

	if config.MaxConns != 8 {
		t.Errorf("MaxConns = %d, want 8", config.MaxConns)
	}
	if config.MinConns != 2 {
		t.Errorf("MinConns = %d, want 2", config.MinConns)
	}
	if config.MaxConnIdleTime != 45*time.Second {
		t.Errorf("MaxConnIdleTime = %s, want 45s", config.MaxConnIdleTime)
	}
	if config.MaxConnLifetime != 10*time.Minute {
		t.Errorf("MaxConnLifetime = %s, want 10m", config.MaxConnLifetime)
	}
	if config.HealthCheckPeriod != 15*time.Second {
		t.Errorf("HealthCheckPeriod = %s, want 15s", config.HealthCheckPeriod)
	}
	if config.ConnConfig.DefaultQueryExecMode != pgx.QueryExecModeSimpleProtocol {
		t.Errorf("DefaultQueryExecMode = %v, want %v", config.ConnConfig.DefaultQueryExecMode, pgx.QueryExecModeSimpleProtocol)
	}
	if got := config.ConnConfig.RuntimeParams["application_name"]; got != "entre-api-preview" {
		t.Errorf("application_name = %q, want entre-api-preview", got)
	}
}

func TestNewPoolConfigValidatesEnv(t *testing.T) {
	tests := []struct {
		name        string
		envName     string
		envValue    string
		wantMessage string
	}{
		{
			name:        "negative max conns",
			envName:     "PGPOOL_MAX_CONNS",
			envValue:    "-1",
			wantMessage: "PGPOOL_MAX_CONNS",
		},
		{
			name:        "invalid idle duration",
			envName:     "PGPOOL_MAX_CONN_IDLE_TIME",
			envValue:    "soon",
			wantMessage: "PGPOOL_MAX_CONN_IDLE_TIME",
		},
		{
			name:        "invalid query mode",
			envName:     "PGX_DEFAULT_QUERY_EXEC_MODE",
			envValue:    "prepared",
			wantMessage: "PGX_DEFAULT_QUERY_EXEC_MODE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearPoolEnv(t)
			t.Setenv(tt.envName, tt.envValue)

			_, err := NewPoolConfig(testDatabaseURL)
			if err == nil {
				t.Fatal("NewPoolConfig returned nil error")
			}
			if !strings.Contains(err.Error(), tt.wantMessage) {
				t.Fatalf("error = %q, want it to mention %q", err.Error(), tt.wantMessage)
			}
		})
	}
}

func TestNewPoolConfigRejectsMinConnsGreaterThanMaxConns(t *testing.T) {
	clearPoolEnv(t)
	t.Setenv("PGPOOL_MAX_CONNS", "2")
	t.Setenv("PGPOOL_MIN_CONNS", "3")

	_, err := NewPoolConfig(testDatabaseURL)
	if err == nil {
		t.Fatal("NewPoolConfig returned nil error")
	}
	if !strings.Contains(err.Error(), "PGPOOL_MIN_CONNS") {
		t.Fatalf("error = %q, want it to mention PGPOOL_MIN_CONNS", err.Error())
	}
}

func clearPoolEnv(t *testing.T) {
	t.Helper()

	for _, name := range []string{
		"PGPOOL_MAX_CONNS",
		"PGPOOL_MIN_CONNS",
		"PGPOOL_MAX_CONN_IDLE_TIME",
		"PGPOOL_MAX_CONN_LIFETIME",
		"PGPOOL_HEALTH_CHECK_PERIOD",
		"PGX_DEFAULT_QUERY_EXEC_MODE",
		"PGAPPNAME",
	} {
		t.Setenv(name, "")
	}
}
