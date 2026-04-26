// Package main は HTTP サーバのエントリポイント。
// 依存解決 (DI) と HTTP ルーティングの最終配線をここで行う。
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/handler"
	fbinfra "github.com/karimiku/job-hunting-saas/internal/infra/firebase"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
	inboxclipuc "github.com/karimiku/job-hunting-saas/internal/usecase/inbox_clip"
	stagehistoryuc "github.com/karimiku/job-hunting-saas/internal/usecase/stage_history"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
	useruc "github.com/karimiku/job-hunting-saas/internal/usecase/user"
)

func main() {
	// run() が return すれば defer (pool.Close 等) が実行されてから os.Exit に到達する。
	// log.Fatal を直接呼ぶと defer がスキップされて DB 接続が閉じられないため避ける。
	if err := run(); err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	var (
		companyRepo      repository.CompanyRepository
		entryRepo        repository.EntryRepository
		taskRepo         repository.TaskRepository
		stageHistoryRepo repository.StageHistoryRepository
		userRepo         repository.UserRepository
		extIDRepo        repository.ExternalIdentityRepository
		inboxClipRepo    repository.InboxClipRepository
	)

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		pool, err := postgres.NewPool(ctx, dbURL)
		if err != nil {
			return fmt.Errorf("connect to database: %w", err)
		}
		defer pool.Close()

		companyRepo = postgres.NewCompanyRepository(pool)
		entryRepo = postgres.NewEntryRepository(pool)
		taskRepo = postgres.NewTaskRepository(pool)
		stageHistoryRepo = postgres.NewStageHistoryRepository(pool)
		userRepo = postgres.NewUserRepository(pool)
		extIDRepo = postgres.NewExternalIdentityRepository(pool)
		// Postgres 版が未実装のため、Inbox は in-memory のまま運用する。
		// （Chrome 拡張用のクリップ保存。永続化は後続 PR で対応予定。）
		inboxClipRepo = inmemory.NewInboxClipRepository()
		log.Println("using PostgreSQL repositories (Inbox clips: in-memory only)")
	} else {
		inMemoryCompanyRepo := inmemory.NewCompanyRepository()
		inMemoryEntryRepo := inmemory.NewEntryRepository()

		companyRepo = inMemoryCompanyRepo
		entryRepo = inMemoryEntryRepo
		taskRepo = inmemory.NewTaskRepository(inMemoryEntryRepo)
		stageHistoryRepo = inmemory.NewStageHistoryRepository()
		inboxClipRepo = inmemory.NewInboxClipRepository()
		log.Println("using in-memory repositories (DATABASE_URL not set) — auth endpoints disabled")
	}

	// Firebase 初期化 / Auth 配線は DB 永続化できる場合のみ有効化する
	var (
		authHandler    *handler.AuthHandler
		authMiddleware func(http.Handler) http.Handler
	)
	if userRepo != nil && extIDRepo != nil {
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		if projectID == "" {
			return errors.New("FIREBASE_PROJECT_ID must be set when DATABASE_URL is configured")
		}
		// GOOGLE_APPLICATION_CREDENTIALS を使うなら credentialsPath を空にして ADC に任せる
		credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_FILE")

		fb, err := fbinfra.NewClient(ctx, credentialsPath, projectID)
		if err != nil {
			return fmt.Errorf("init firebase: %w", err)
		}

		authenticateUC := useruc.NewAuthenticate(userRepo, extIDRepo)
		// Firebase SDK 型を handler / middleware に漏らさないため、adapter で DTO に変換する。
		sessionCreator := fbinfra.NewSessionCreator(fb.Auth)
		sessionVerifier := fbinfra.NewSessionVerifier(fb.Auth)
		authHandler = handler.NewAuthHandler(sessionCreator, authenticateUC, userRepo, handler.AuthConfig{
			CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
		})
		authMiddleware = middleware.NewAuth(sessionVerifier, extIDRepo)
		log.Println("firebase auth wired")
	}

	companyHandler := handler.NewCompanyHandler(
		companyuc.NewCreate(companyRepo),
		companyuc.NewGet(companyRepo),
		companyuc.NewList(companyRepo),
		companyuc.NewUpdate(companyRepo),
		companyuc.NewDelete(companyRepo),
	)

	entryHandler := handler.NewEntryHandler(
		entryuc.NewCreate(entryRepo, companyRepo),
		entryuc.NewGet(entryRepo),
		entryuc.NewList(entryRepo),
		entryuc.NewUpdate(entryRepo),
		entryuc.NewDelete(entryRepo),
	)

	taskHandler := handler.NewTaskHandler(
		taskuc.NewCreate(taskRepo, entryRepo),
		taskuc.NewGet(taskRepo),
		taskuc.NewList(taskRepo),
		taskuc.NewUpdate(taskRepo),
		taskuc.NewDelete(taskRepo),
	)

	stageHistoryHandler := handler.NewStageHistoryHandler(
		stagehistoryuc.NewCreate(stageHistoryRepo, entryRepo),
		stagehistoryuc.NewList(stageHistoryRepo, entryRepo),
	)

	inboxClipHandler := handler.NewInboxClipHandler(
		inboxclipuc.NewCreate(inboxClipRepo),
		inboxclipuc.NewList(inboxClipRepo),
		inboxclipuc.NewDelete(inboxClipRepo),
	)

	h := &handler.Handler{
		CompanyHandler:      companyHandler,
		EntryHandler:        entryHandler,
		TaskHandler:         taskHandler,
		StageHistoryHandler: stageHistoryHandler,
		InboxClipHandler:    inboxClipHandler,
	}

	router := chi.NewRouter()
	router.Use(corsMiddleware(os.Getenv("CORS_ALLOWED_ORIGIN")))

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, "ok"); err != nil {
			log.Printf("health: write response failed: %v", err)
		}
	})

	// 認証不要ルート（ログイン / ログアウト）
	if authHandler != nil {
		authHandler.PublicRoutes(router)
	}

	// 認証必須ルート
	router.Group(func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware)
		}
		if authHandler != nil {
			authHandler.ProtectedRoutes(r)
		}
		openapi.HandlerFromMux(h, r)
	})

	// http.ListenAndServe には timeout が無く Slowloris 等の DoS に弱いため、
	// http.Server を明示してヘッダ・ボディ・アイドル各種 timeout を設定する。
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

// corsMiddleware はフロントエンドとの Cookie 付き通信を許可する最小 CORS 実装。
// allowedOrigin が空の場合は http://localhost:3000 をデフォルトとする。
func corsMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
					http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
				}, ", "))
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
