package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	mcphandler "github.com/karimiku/job-hunting-saas/internal/handler/mcp"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
	esmemo "github.com/karimiku/job-hunting-saas/internal/usecase/es_memo"
	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

func main() {
	log.SetOutput(os.Stderr)

	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	user, err := resolveMCPUser(ctx, userRepo)
	if err != nil {
		log.Fatalf("resolve MCP user: %v", err)
	}

	entryRepo := postgres.NewEntryRepository(pool)
	taskRepo := postgres.NewTaskRepository(pool)
	memoRepo := postgres.NewESMemoRepository(pool)
	query := postgres.NewMCPQuery(pool)
	app := mcpuc.NewService(
		user.ID(),
		query,
		esmemo.NewAppend(memoRepo, entryRepo),
		taskuc.NewCreate(taskRepo, entryRepo),
		jobemail.NewExtract(),
	)

	server := mcphandler.NewServer(app)
	if err := mcphandler.ServeStdio(ctx, os.Stdin, os.Stdout, server); err != nil && !errors.Is(err, io.EOF) {
		log.Fatalf("serve MCP: %v", err)
	}
}

func resolveMCPUser(ctx context.Context, repo repository.UserRepository) (*entity.User, error) {
	if raw := strings.TrimSpace(os.Getenv("MCP_USER_ID")); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid MCP_USER_ID: %w", err)
		}
		return repo.FindByID(ctx, entity.UserID(id))
	}
	if raw := strings.TrimSpace(os.Getenv("MCP_USER_EMAIL")); raw != "" {
		email, err := value.NewEmail(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid MCP_USER_EMAIL: %w", err)
		}
		return repo.FindByEmail(ctx, email)
	}
	return nil, errors.New("set MCP_USER_EMAIL or MCP_USER_ID")
}
