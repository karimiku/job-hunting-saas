// Command mcp-server exposes job-hunting context operations over stdio MCP.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

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

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	if usesRemoteAPIEnv() {
		app, err := newRemoteApplicationFromEnv()
		if err != nil {
			return err
		}
		return serveMCP(ctx, app)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	tokenRepo := postgres.NewAIAccessTokenRepository(pool)
	user, err := resolveMCPUser(ctx, userRepo, tokenRepo)
	if err != nil {
		return fmt.Errorf("resolve MCP user: %w", err)
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

	return serveMCP(ctx, app)
}

func serveMCP(ctx context.Context, app mcphandler.Application) error {
	server := mcphandler.NewServer(app)
	if err := mcphandler.ServeStdio(ctx, os.Stdin, os.Stdout, server); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("serve MCP: %w", err)
	}
	return nil
}

func resolveMCPUser(
	ctx context.Context,
	userRepo repository.UserRepository,
	tokenRepo repository.AIAccessTokenRepository,
) (*entity.User, error) {
	if raw := strings.TrimSpace(os.Getenv("MCP_API_KEY")); raw != "" {
		secret, err := value.NewAIAccessTokenSecret(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid MCP_API_KEY: %w", err)
		}
		token, err := tokenRepo.FindActiveByHash(ctx, secret.Hash())
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, errors.New("invalid MCP_API_KEY")
			}
			return nil, fmt.Errorf("find MCP_API_KEY: %w", err)
		}
		if err := tokenRepo.TouchLastUsed(ctx, token.ID(), time.Now()); err != nil {
			return nil, fmt.Errorf("touch MCP_API_KEY: %w", err)
		}
		return userRepo.FindByID(ctx, token.UserID())
	}

	if raw := strings.TrimSpace(os.Getenv("MCP_USER_ID")); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid MCP_USER_ID: %w", err)
		}
		return userRepo.FindByID(ctx, entity.UserID(id))
	}
	if raw := strings.TrimSpace(os.Getenv("MCP_USER_EMAIL")); raw != "" {
		email, err := value.NewEmail(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid MCP_USER_EMAIL: %w", err)
		}
		return userRepo.FindByEmail(ctx, email)
	}
	return nil, errors.New("set MCP_API_KEY, MCP_USER_EMAIL, or MCP_USER_ID")
}
