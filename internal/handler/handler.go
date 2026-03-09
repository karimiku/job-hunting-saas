package handler

import "github.com/karimiku/job-hunting-saas/internal/gen/openapi"

// Handler は oapi-codegen の ServerInterface を実装する。
// ドメインごとのHandlerを埋め込みで合成し、単一のServerInterfaceとして振る舞う。
type Handler struct {
	*CompanyHandler
	*EntryHandler
}

var _ openapi.ServerInterface = (*Handler)(nil)
