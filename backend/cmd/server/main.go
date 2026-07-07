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
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/handler"
	fbinfra "github.com/karimiku/job-hunting-saas/internal/infra/firebase"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
	"github.com/karimiku/job-hunting-saas/internal/infra/supabaseauth"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	aiaccesstokenuc "github.com/karimiku/job-hunting-saas/internal/usecase/ai_access_token"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
	companyaliasuc "github.com/karimiku/job-hunting-saas/internal/usecase/company_alias"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
	esmemo "github.com/karimiku/job-hunting-saas/internal/usecase/es_memo"
	inboxclipuc "github.com/karimiku/job-hunting-saas/internal/usecase/inbox_clip"
	selectionflowuc "github.com/karimiku/job-hunting-saas/internal/usecase/selection_flow"
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
	devAuthEnabled := os.Getenv("DEV_AUTH_ENABLED") == "true"
	if devAuthEnabled && isProductionRuntime() {
		return errors.New("DEV_AUTH_ENABLED must not be true in production")
	}
	devSessionSecret := ""
	if devAuthEnabled {
		devSessionSecret = strings.TrimSpace(os.Getenv("DEV_AUTH_SECRET"))
		if devSessionSecret == "" {
			devSessionSecret = uuid.NewString()
			log.Println("dev auth enabled with ephemeral session secret")
		}
	}

	ctx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()

	var (
		accountRepo          repository.AccountRepository
		companyRepo          repository.CompanyRepository
		companyAliasRepo     repository.CompanyAliasRepository
		entryRepo            repository.EntryRepository
		entryWithCompanyRepo repository.EntryWithCompanyRepository
		taskRepo             repository.TaskRepository
		stageHistoryRepo     repository.StageHistoryRepository
		selectionFlowRepo    repository.SelectionFlowRepository
		userRepo             repository.UserRepository
		extIDRepo            repository.ExternalIdentityRepository
		inboxClipRepo        repository.InboxClipRepository
		esMemoRepo           repository.ESMemoRepository
		aiTokenRepo          repository.AIAccessTokenRepository
	)

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		pool, err := postgres.NewPool(ctx, dbURL)
		if err != nil {
			return fmt.Errorf("connect to database: %w", err)
		}
		defer pool.Close()

		accountRepo = postgres.NewAccountRepository(pool)
		companyRepo = postgres.NewCompanyRepository(pool)
		companyAliasRepo = postgres.NewCompanyAliasRepository(pool)
		entryRepo = postgres.NewEntryRepository(pool)
		entryWithCompanyRepo = postgres.NewEntryWithCompanyRepository(pool)
		taskRepo = postgres.NewTaskRepository(pool)
		stageHistoryRepo = postgres.NewStageHistoryRepository(pool)
		selectionFlowRepo = postgres.NewSelectionFlowRepository(pool)
		userRepo = postgres.NewUserRepository(pool)
		extIDRepo = postgres.NewExternalIdentityRepository(pool)
		inboxClipRepo = postgres.NewInboxClipRepository(pool)
		esMemoRepo = postgres.NewESMemoRepository(pool)
		aiTokenRepo = postgres.NewAIAccessTokenRepository(pool)
		log.Println("using PostgreSQL repositories")
	} else {
		// DATABASE_URL 未設定 = 開発・ローカルテストモード。auth middleware も配線できないため
		// 全エンドポイントが認証なしで通る。これを誤って本番起動しないよう明示フラグを要求。
		if os.Getenv("ALLOW_INSECURE_NO_AUTH") != "true" {
			return errors.New("DATABASE_URL is not set, which disables authentication and " +
				"causes all clips/entries to share a zero-value UserID; set " +
				"ALLOW_INSECURE_NO_AUTH=true to proceed in dev mode, or configure DATABASE_URL " +
				"and SUPABASE_AUTH_ISSUER for a real environment")
		}
		inMemoryCompanyRepo := inmemory.NewCompanyRepository()
		inMemoryEntryRepo := inmemory.NewEntryRepository()
		inMemoryUserRepo := inmemory.NewUserRepository()

		accountRepo = inmemory.NewAccountRepository(inMemoryUserRepo)
		companyRepo = inMemoryCompanyRepo
		companyAliasRepo = inmemory.NewCompanyAliasRepository()
		entryRepo = inMemoryEntryRepo
		entryWithCompanyRepo = inmemory.NewEntryWithCompanyRepository(inMemoryCompanyRepo, inMemoryEntryRepo)
		taskRepo = inmemory.NewTaskRepository(inMemoryEntryRepo)
		stageHistoryRepo = inmemory.NewStageHistoryRepository()
		selectionFlowRepo = inmemory.NewSelectionFlowRepository()
		inboxClipRepo = inmemory.NewInboxClipRepository()
		esMemoRepo = inmemory.NewESMemoRepository()
		aiTokenRepo = inmemory.NewAIAccessTokenRepository()
		log.Println("using in-memory repositories (ALLOW_INSECURE_NO_AUTH=true) — auth endpoints disabled, all data shared across users")
	}

	// Auth 配線は DB 永続化できる場合のみ有効化する。
	var (
		authHandler             *handler.AuthHandler
		devAuthHandler          *handler.DevAuthHandler
		authMiddleware          func(http.Handler) http.Handler
		authConfig              handler.AuthConfig
		firebaseLoginConfigured bool
	)
	if userRepo != nil && extIDRepo != nil {
		authenticateUC := useruc.NewAuthenticate(userRepo, extIDRepo)
		cookieSameSite, err := parseCookieSameSite(os.Getenv("COOKIE_SAME_SITE"))
		if err != nil {
			return err
		}
		authConfig = handler.AuthConfig{
			CookieDomain:   os.Getenv("COOKIE_DOMAIN"),
			CookieSecure:   os.Getenv("COOKIE_SECURE") == "true",
			CookieSameSite: cookieSameSite,
		}

		aiBearerVerifier := aiaccesstokenuc.NewVerify(aiTokenRepo)
		bearerVerifier := middleware.NewChainedBearerTokenVerifier(aiBearerVerifier)

		supabaseAuthConfigured := false
		if supabaseIssuer := strings.TrimSpace(os.Getenv("SUPABASE_AUTH_ISSUER")); supabaseIssuer != "" {
			supabaseVerifier, err := supabaseauth.NewVerifier(supabaseauth.Config{
				Issuer:   supabaseIssuer,
				Audience: os.Getenv("SUPABASE_JWT_AUDIENCE"),
				JWKSURL:  os.Getenv("SUPABASE_JWKS_URL"),
				UserSync: func(ctx context.Context, info supabaseauth.UserInfo) (entity.UserID, error) {
					out, err := authenticateUC.Execute(ctx, useruc.AuthenticateInput{
						Provider: "supabase",
						Subject:  info.Subject,
						Email:    info.Email,
						Name:     info.Name,
					})
					if err != nil {
						return entity.UserID{}, err
					}
					return out.User.ID(), nil
				},
			}, extIDRepo)
			if err != nil {
				return fmt.Errorf("init supabase auth verifier: %w", err)
			}
			bearerVerifier = middleware.NewChainedBearerTokenVerifier(aiBearerVerifier, supabaseVerifier)
			supabaseAuthConfigured = true
			log.Println("supabase auth bearer verifier wired")
		}

		projectID := strings.TrimSpace(os.Getenv("FIREBASE_PROJECT_ID"))
		var sessionVerifier middleware.FirebaseSessionVerifier
		var sessionCreator handler.FirebaseSessionCreator
		if projectID != "" {
			// GOOGLE_APPLICATION_CREDENTIALS を使うなら credentialsPath を空にして ADC に任せる
			credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_FILE")

			fb, err := fbinfra.NewClient(ctx, credentialsPath, projectID)
			if err != nil {
				return fmt.Errorf("init firebase: %w", err)
			}

			// Firebase SDK 型を handler / middleware に漏らさないため、adapter で DTO に変換する。
			sessionCreator = fbinfra.NewSessionCreator(fb.Auth)
			sessionVerifierCacheTTL, err := firebaseSessionVerifierCacheTTL(os.Getenv("FIREBASE_SESSION_VERIFY_CACHE_TTL"))
			if err != nil {
				return err
			}
			sessionVerifier = fbinfra.NewSessionVerifier(fb.Auth)
			sessionVerifier = middleware.NewCachedSessionVerifier(sessionVerifier, sessionVerifierCacheTTL)
			firebaseLoginConfigured = true
			log.Println("firebase auth wired")
		}
		// authHandler は GET /auth/me・DELETE /auth/session（ログアウト）を提供し、
		// どちらも firebaseAuth を必要としない。POST /auth/session（ログイン）だけが
		// firebaseAuth を使うため、Firebase 未設定時は sessionCreator が nil のまま渡す。
		// この場合でも main.go 側のルーティングで LoginRoute を登録しない限り
		// CreateSession（nil interface 呼び出し）は到達しない。
		authHandler = handler.NewAuthHandler(sessionCreator, authenticateUC, userRepo, authConfig)

		if projectID == "" && !supabaseAuthConfigured && !devAuthEnabled {
			return errors.New("FIREBASE_PROJECT_ID or SUPABASE_AUTH_ISSUER must be set when DATABASE_URL is configured")
		}

		authMiddleware = middleware.NewAuthWithBearerAndDevSession(
			sessionVerifier,
			extIDRepo,
			bearerVerifier,
			devSessionSecret,
		)
		if devAuthEnabled {
			devAuthHandler = handler.NewDevAuthHandler(authenticateUC, authConfig, devSessionSecret)
			log.Println("dev auth wired")
		}
	}

	companyHandler := handler.NewCompanyHandler(
		companyuc.NewCreate(companyRepo),
		companyuc.NewGet(companyRepo),
		companyuc.NewList(companyRepo),
		companyuc.NewUpdate(companyRepo),
		companyuc.NewDelete(companyRepo),
	)

	companyAliasHandler := handler.NewCompanyAliasHandler(
		companyaliasuc.NewCreate(companyAliasRepo, companyRepo),
		companyaliasuc.NewGet(companyAliasRepo),
		companyaliasuc.NewList(companyAliasRepo, companyRepo),
		companyaliasuc.NewDelete(companyAliasRepo),
	)

	entryHandler := handler.NewEntryHandler(
		entryuc.NewCreate(entryRepo, companyRepo),
		entryuc.NewCreateWithCompany(entryWithCompanyRepo),
		entryuc.NewGet(entryRepo),
		entryuc.NewList(entryRepo),
		companyuc.NewList(companyRepo),
		entryuc.NewUpdate(entryRepo),
		entryuc.NewDelete(entryRepo),
	)

	taskHandler := handler.NewTaskHandler(
		taskuc.NewCreate(taskRepo, entryRepo),
		taskuc.NewGet(taskRepo),
		taskuc.NewList(taskRepo),
		taskuc.NewListAll(taskRepo),
		taskuc.NewUpdate(taskRepo),
		taskuc.NewDelete(taskRepo),
	)

	pageDataHandler := handler.NewPageDataHandler(
		userRepo,
		entryuc.NewList(entryRepo),
		companyuc.NewList(companyRepo),
		inboxclipuc.NewList(inboxClipRepo),
		taskuc.NewListAll(taskRepo),
	)

	stageHistoryHandler := handler.NewStageHistoryHandler(
		stagehistoryuc.NewCreate(stageHistoryRepo, entryRepo),
		stagehistoryuc.NewList(stageHistoryRepo, entryRepo),
	)

	selectionFlowHandler := handler.NewSelectionFlowHandler(
		selectionflowuc.NewGet(selectionFlowRepo, entryRepo),
		selectionflowuc.NewUpsert(selectionFlowRepo, entryRepo),
		selectionflowuc.NewUpdateCurrent(selectionFlowRepo, entryRepo),
	)

	inboxClipHandler := handler.NewInboxClipHandler(
		inboxclipuc.NewCreate(inboxClipRepo),
		inboxclipuc.NewList(inboxClipRepo),
		inboxclipuc.NewDelete(inboxClipRepo),
	)

	aiTokenHandler := handler.NewAiAccessTokenHandler(
		aiaccesstokenuc.NewCreate(aiTokenRepo),
		aiaccesstokenuc.NewList(aiTokenRepo),
		aiaccesstokenuc.NewRevoke(aiTokenRepo),
	)

	esMemoHandler := handler.NewESMemoHandler(
		esmemo.NewAppend(esMemoRepo, entryRepo),
		esmemo.NewList(esMemoRepo),
	)

	meHandler := handler.NewMeHandler(
		useruc.NewDeleteAccount(accountRepo),
		authConfig,
	)

	h := &handler.Handler{
		MeHandler:            meHandler,
		CompanyHandler:       companyHandler,
		CompanyAliasHandler:  companyAliasHandler,
		EntryHandler:         entryHandler,
		TaskHandler:          taskHandler,
		PageDataHandler:      pageDataHandler,
		StageHistoryHandler:  stageHistoryHandler,
		SelectionFlowHandler: selectionFlowHandler,
		InboxClipHandler:     inboxClipHandler,
		AiAccessTokenHandler: aiTokenHandler,
		ESMemoHandler:        esMemoHandler,
	}

	router := chi.NewRouter()
	router.Use(stripPathPrefixMiddleware("/backend"))
	// CORS_ALLOWED_ORIGINS (新): カンマ区切りで複数 origin を許可 (chrome 拡張等を追加するときに使う)。
	// CORS_ALLOWED_ORIGIN (旧, 後方互換): 単一 origin。両方セットされていれば新の方を優先。
	corsOriginsRaw := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOriginsRaw == "" {
		corsOriginsRaw = os.Getenv("CORS_ALLOWED_ORIGIN")
	}
	corsOrigins := allowedOrigins(corsOriginsRaw)
	warnOnCredentialedWildcardCORS(corsOrigins, os.Getenv("COOKIE_SAME_SITE"))
	globalRateLimit, err := requestsPerMinuteFromEnv("RATE_LIMIT_GLOBAL_REQUESTS_PER_MINUTE", 30)
	if err != nil {
		return err
	}
	authRateLimit, err := requestsPerMinuteFromEnv("RATE_LIMIT_AUTH_REQUESTS_PER_MINUTE", 5)
	if err != nil {
		return err
	}
	userRateLimit, err := requestsPerMinuteFromEnv("RATE_LIMIT_AUTHENTICATED_REQUESTS_PER_MINUTE", 60)
	if err != nil {
		return err
	}

	router.Use(corsMiddleware(corsOrigins))
	router.Use(middleware.NewIPRateLimiter(globalRateLimit, time.Minute))
	router.Use(middleware.NewServerTiming())

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, "ok"); err != nil {
			log.Printf("health: write response failed: %v", err)
		}
	})

	// 認証不要ルート（ログイン / ログアウト）
	// LogoutRoute は firebaseAuth 不要なので Firebase 未設定でも常に登録する。
	// LoginRoute（Firebase ID Token 検証）は firebaseAuth が設定されている場合のみ登録する。
	if authHandler != nil {
		router.Group(func(r chi.Router) {
			r.Use(middleware.NewIPRateLimiter(authRateLimit, time.Minute))
			r.Use(middleware.NewOriginGuard(corsOrigins))
			authHandler.LogoutRoute(r)
			if firebaseLoginConfigured {
				authHandler.LoginRoute(r)
			}
		})
	}
	if devAuthHandler != nil {
		router.Group(func(r chi.Router) {
			r.Use(middleware.NewIPRateLimiter(authRateLimit, time.Minute))
			r.Use(middleware.NewOriginGuard(corsOrigins))
			devAuthHandler.PublicRoutes(r)
		})
	}

	// 認証必須ルート
	router.Group(func(r chi.Router) {
		if authMiddleware != nil {
			r.Use(authMiddleware)
			r.Use(middleware.NewAuthenticatedUserRateLimiter(userRateLimit, time.Minute))
			r.Use(middleware.NewSessionCSRFProtection(corsOrigins))
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

	// #nosec G706 -- port is deployment configuration, not user-controlled request data.
	log.Printf("server listening on :%s", port)
	serverErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("listen: %w", err)
			return
		}
		serverErr <- nil
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		stopSignals()
		log.Println("server shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown: %w", err)
		}
		return <-serverErr
	}
}

func requestsPerMinuteFromEnv(name string, defaultValue int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s %q (want non-negative integer requests per minute): %w", name, raw, err)
	}
	if value < 0 {
		return 0, fmt.Errorf("invalid %s %q (must be non-negative)", name, raw)
	}
	return value, nil
}

// isProductionRuntime is a best-effort guard used only to hard-fail startup when
// DEV_AUTH_ENABLED is left on in an environment that self-identifies as
// production. It intentionally does NOT flip unset envs to "production": doing so
// would block local dev (where these envs are commonly unset) and would trade a
// startup failure for a false sense of safety. The real protection against dev
// sessions in production is (a) never setting DEV_AUTH_ENABLED=true there, and
// (b) isLocalDevRequest trusting only the server-observed r.Host, which is not a
// loopback value on a public deployment. Keep both in place.
func isProductionRuntime() bool {
	for _, name := range []string{"APP_ENV", "GO_ENV", "ENV", "GIN_MODE"} {
		if strings.EqualFold(strings.TrimSpace(os.Getenv(name)), "production") {
			return true
		}
	}
	return false
}

func parseCookieSameSite(raw string) (http.SameSite, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "lax":
		return http.SameSiteLaxMode, nil
	case "strict":
		return http.SameSiteStrictMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return 0, fmt.Errorf("invalid COOKIE_SAME_SITE %q (want lax, strict, or none)", raw)
	}
}

func firebaseSessionVerifierCacheTTL(raw string) (time.Duration, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 30 * time.Second, nil
	}
	ttl, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid FIREBASE_SESSION_VERIFY_CACHE_TTL %q (want duration like 30s or 1m): %w", raw, err)
	}
	if ttl < 0 {
		return 0, fmt.Errorf("invalid FIREBASE_SESSION_VERIFY_CACHE_TTL %q (must not be negative)", raw)
	}
	return ttl, nil
}

// corsMiddleware はフロントエンドとの Cookie 付き通信を許可する最小 CORS 実装。
// 受け取るのはカンマ区切りの allowlist。空なら http://localhost:3000 のみ。
//
// Chrome 拡張から呼びたい場合は `chrome-extension://<extension-id>` を allowlist に追加する。
// 拡張の popup/background から host_permissions 付きで呼ぶ場合は SameSite=Strict の
// session cookie 共有を検証する。動かない場合は拡張専用 token 方式に切り替える。
func allowedOrigins(allowedOriginsRaw string) []string {
	if allowedOriginsRaw == "" {
		allowedOriginsRaw = "http://localhost:3000"
	}
	origins := make([]string, 0)
	for _, o := range strings.Split(allowedOriginsRaw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			o = strings.TrimRight(o, "/")
			origins = append(origins, o)
		}
	}
	return origins
}

// warnOnCredentialedWildcardCORS logs a startup warning when a wildcard origin is
// configured for credentialed CORS. Any matching subdomain (e.g. evil.vercel.app)
// could then send cookie/credentialed requests; combined with SameSite=none the
// legacy session cookie becomes cross-site theftable. Behavior is unchanged — this
// only surfaces a misconfiguration risk that must be resolved at the ops layer.
func warnOnCredentialedWildcardCORS(origins []string, cookieSameSite string) {
	var wildcards []string
	for _, o := range origins {
		if strings.Contains(o, "*.") {
			wildcards = append(wildcards, o)
		}
	}
	if len(wildcards) == 0 {
		return
	}
	log.Printf("security warning: credentialed CORS allows wildcard origin(s) %s; "+
		"any matching subdomain can send cookie/credentialed requests. Prefer exact "+
		"origins in CORS_ALLOWED_ORIGINS for credentialed paths.", strings.Join(wildcards, ", "))
	if strings.EqualFold(strings.TrimSpace(cookieSameSite), "none") {
		log.Printf("security warning: COOKIE_SAME_SITE=none combined with wildcard CORS " +
			"origins exposes legacy session cookies to cross-site theft from any matching " +
			"subdomain; use exact origins or SameSite=lax/strict for the session cookie.")
	}
}

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, o := range allowedOrigins {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed[o] = struct{}{}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowed[origin]; ok && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Expose-Headers", "Server-Timing")
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

func stripPathPrefixMiddleware(prefix string) func(http.Handler) http.Handler {
	prefix = "/" + strings.Trim(strings.TrimSpace(prefix), "/")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == prefix {
				r = cloneRequestWithPath(r, "/")
			} else if strings.HasPrefix(path, prefix+"/") {
				r = cloneRequestWithPath(r, strings.TrimPrefix(path, prefix))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func cloneRequestWithPath(r *http.Request, path string) *http.Request {
	clone := r.Clone(r.Context())
	clone.URL.Path = path
	clone.URL.RawPath = ""
	return clone
}
