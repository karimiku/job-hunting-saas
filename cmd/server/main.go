package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/handler"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// InMemory実装はプロセス再起動でデータが消える。本番ではPostgreSQL実装に差し替える。
	companyRepo := inmemory.NewCompanyRepository()
	entryRepo := inmemory.NewEntryRepository()

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

	h := &handler.Handler{
		CompanyHandler: companyHandler,
		EntryHandler:   entryHandler,
	}

	router := chi.NewRouter()
	router.Use(middleware.Auth)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
	// oapi-codegen が生成した ServerInterface のルーティングを登録する
	openapi.HandlerFromMux(h, router)

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
