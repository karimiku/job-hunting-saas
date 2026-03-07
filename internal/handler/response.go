package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
)

// writeJSON はHTTPレスポンスをJSON形式で書き出す。
// 全handlerで共通して使う低レベルなレスポンス関数。
func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// writeError はドメイン/リポジトリ層のエラーをHTTPステータスに変換する。
// handlerごとにエラー分岐が重複するのを防ぐための集約ポイント。
// 既知のドメインエラーのみクライアントに返し、それ以外は500+汎用メッセージで内部詳細を隠す。
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeJSON(w, http.StatusNotFound, openapi.ErrorResponse{Message: "not found"})
	case errors.Is(err, value.ErrCompanyNameEmpty),
		errors.Is(err, value.ErrCompanyNameInvalid):
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: err.Error()})
	default:
		log.Printf("unexpected error: %v", err)
		writeJSON(w, http.StatusInternalServerError, openapi.ErrorResponse{Message: "internal server error"})
	}
}
