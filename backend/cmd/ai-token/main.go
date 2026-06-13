// Command ai-token registers or generates an AI access token for MCP clients.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ai-token: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		userIDRaw string
		emailRaw  string
		name      string
		tokenRaw  string
	)
	flag.StringVar(&userIDRaw, "user-id", "", "target user UUID")
	flag.StringVar(&emailRaw, "email", "", "target user email")
	flag.StringVar(&name, "name", "AI client", "token label")
	flag.StringVar(&tokenRaw, "token", "", "existing entre_ai_ token to register; defaults to AI_ACCESS_TOKEN env or generates a new one")
	flag.Parse()

	if strings.TrimSpace(tokenRaw) == "" {
		tokenRaw = os.Getenv("AI_ACCESS_TOKEN")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databaseURL) == "" {
		return errors.New("DATABASE_URL is required")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepository(pool)
	user, err := resolveTargetUser(ctx, userRepo, userIDRaw, emailRaw)
	if err != nil {
		return err
	}

	secret, generated, err := resolveTokenSecret(tokenRaw)
	if err != nil {
		return err
	}

	tokenName := strings.TrimSpace(name)
	if tokenName == "" {
		tokenName = "AI client"
	}
	token := entity.NewAIAccessToken(user.ID(), tokenName, secret.Hash(), secret.Preview())
	tokenRepo := postgres.NewAIAccessTokenRepository(pool)
	if err := tokenRepo.Save(ctx, token); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return errors.New("AI access token is already registered")
		}
		return err
	}

	if generated {
		fmt.Printf("AI access token created for %s\n", user.Email().String())
		fmt.Printf("Token: %s\n", secret.String())
		fmt.Println("Store this value now; only its hash is saved.")
		return nil
	}
	fmt.Printf("AI access token registered for %s (%s)\n", user.Email().String(), secret.Preview())
	return nil
}

func resolveTargetUser(ctx context.Context, repo repository.UserRepository, userIDRaw, emailRaw string) (*entity.User, error) {
	if raw := strings.TrimSpace(userIDRaw); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid -user-id: %w", err)
		}
		return repo.FindByID(ctx, entity.UserID(id))
	}
	if raw := strings.TrimSpace(emailRaw); raw != "" {
		email, err := value.NewEmail(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid -email: %w", err)
		}
		return repo.FindByEmail(ctx, email)
	}
	return nil, errors.New("set -email or -user-id")
}

func resolveTokenSecret(raw string) (value.AIAccessTokenSecret, bool, error) {
	if strings.TrimSpace(raw) == "" {
		secret, err := value.GenerateAIAccessTokenSecret()
		return secret, true, err
	}
	secret, err := value.NewAIAccessTokenSecret(raw)
	return secret, false, err
}
