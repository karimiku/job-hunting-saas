package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxConns          int32 = 4
	defaultMinConns          int32 = 0
	defaultMaxConnIdleTime         = 60 * time.Second
	defaultMaxConnLifetime         = 30 * time.Minute
	defaultHealthCheckPeriod       = 30 * time.Second
	defaultApplicationName         = "job-hunting-saas-api"
)

// NewPool は databaseURL から pgxpool.Pool を作成する。
// 呼び出し側で pool.Close() を defer すること。
//
// 本番環境では sslmode=require 以上を推奨。
// databaseURL にはパスワードが含まれるためログに出力しないこと。
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := NewPoolConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres: failed to ping: %w", err)
	}

	return pool, nil
}

// NewPoolConfig は DATABASE_URL と環境変数から pgxpool.Config を作成する。
//
// Vercel など auto-scale する実行環境から Supabase に接続する場合、各 container の
// app-side pool が積み上がるため、MaxConns は小さめに抑える。
// また Supabase transaction pooler は prepared statement をサポートしないため、
// pgx の statement cache を使わない exec mode を既定にする。
func NewPoolConfig(databaseURL string) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse database URL: %w", err)
	}

	maxConns, err := int32FromEnv("PGPOOL_MAX_CONNS", defaultMaxConns)
	if err != nil {
		return nil, err
	}
	minConns, err := int32FromEnv("PGPOOL_MIN_CONNS", defaultMinConns)
	if err != nil {
		return nil, err
	}
	if minConns > maxConns {
		return nil, fmt.Errorf("postgres: PGPOOL_MIN_CONNS (%d) must be <= PGPOOL_MAX_CONNS (%d)", minConns, maxConns)
	}

	maxConnIdleTime, err := durationFromEnv("PGPOOL_MAX_CONN_IDLE_TIME", defaultMaxConnIdleTime)
	if err != nil {
		return nil, err
	}
	maxConnLifetime, err := durationFromEnv("PGPOOL_MAX_CONN_LIFETIME", defaultMaxConnLifetime)
	if err != nil {
		return nil, err
	}
	healthCheckPeriod, err := durationFromEnv("PGPOOL_HEALTH_CHECK_PERIOD", defaultHealthCheckPeriod)
	if err != nil {
		return nil, err
	}
	queryExecMode, err := queryExecModeFromEnv("PGX_DEFAULT_QUERY_EXEC_MODE", pgx.QueryExecModeExec)
	if err != nil {
		return nil, err
	}

	config.MaxConns = maxConns
	config.MinConns = minConns
	config.MaxConnIdleTime = maxConnIdleTime
	config.MaxConnLifetime = maxConnLifetime
	config.HealthCheckPeriod = healthCheckPeriod
	config.ConnConfig.DefaultQueryExecMode = queryExecMode

	if config.ConnConfig.RuntimeParams == nil {
		config.ConnConfig.RuntimeParams = make(map[string]string)
	}
	if config.ConnConfig.RuntimeParams["application_name"] == "" {
		appName := strings.TrimSpace(os.Getenv("PGAPPNAME"))
		if appName == "" {
			appName = defaultApplicationName
		}
		config.ConnConfig.RuntimeParams["application_name"] = appName
	}

	return config, nil
}

func int32FromEnv(name string, defaultValue int32) (int32, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("postgres: invalid %s %q (want non-negative integer): %w", name, raw, err)
	}
	if value < 0 {
		return 0, fmt.Errorf("postgres: invalid %s %q (must be non-negative)", name, raw)
	}
	return int32(value), nil
}

func durationFromEnv(name string, defaultValue time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("postgres: invalid %s %q (want duration like 30s or 5m): %w", name, raw, err)
	}
	if value < 0 {
		return 0, fmt.Errorf("postgres: invalid %s %q (must be non-negative)", name, raw)
	}
	return value, nil
}

func queryExecModeFromEnv(name string, defaultValue pgx.QueryExecMode) (pgx.QueryExecMode, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}
	switch raw {
	case "cache_statement":
		return pgx.QueryExecModeCacheStatement, nil
	case "cache_describe":
		return pgx.QueryExecModeCacheDescribe, nil
	case "describe_exec":
		return pgx.QueryExecModeDescribeExec, nil
	case "exec":
		return pgx.QueryExecModeExec, nil
	case "simple_protocol":
		return pgx.QueryExecModeSimpleProtocol, nil
	default:
		return 0, fmt.Errorf("postgres: invalid %s %q (want cache_statement, cache_describe, describe_exec, exec, or simple_protocol)", name, raw)
	}
}
