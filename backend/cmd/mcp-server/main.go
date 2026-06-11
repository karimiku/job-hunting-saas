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

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	mcphandler "github.com/karimiku/job-hunting-saas/internal/handler/mcp"
	"github.com/karimiku/job-hunting-saas/internal/infra/entreapi"
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
	if token := strings.TrimSpace(os.Getenv("ENTRE_API_TOKEN")); token != "" {
		app, err := entreapi.NewMCPApplication(os.Getenv("ENTRE_API_BASE_URL"), token, nil)
		if err != nil {
			return err
		}
		log.Println("mcp-server using Entré API bridge")
		return serve(ctx, app)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("set ENTRE_API_TOKEN for API bridge mode, or DATABASE_URL for direct database mode")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	user, err := resolveMCPUser(ctx, userRepo)
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

	log.Println("mcp-server using direct database mode")
	return serve(ctx, app)
}

func serve(ctx context.Context, app mcphandler.Application) error {
	server := mcphandler.NewServer(app)
	if err := mcphandler.ServeStdio(ctx, os.Stdin, os.Stdout, server); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("serve MCP: %w", err)
	}
	return nil
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
