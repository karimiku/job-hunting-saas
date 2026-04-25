package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
)

// writeJSON はHTTPレスポンスをJSON形式で書き出す。
// 全handlerで共通して使う低レベルなレスポンス関数。
// Encode のエラーは ResponseWriter が既に書き込み開始している以上回復不能のため、
// 明示的に握りつぶしてログに出すに留める。
func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writeJSON: encode failed: %v", err)
	}
}

// writeError はドメイン/リポジトリ層のエラーをHTTPステータスに変換する。
// handlerごとにエラー分岐が重複するのを防ぐための集約ポイント。
// TODO: 予期しないエラーは500を返すように改善する
func writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, openapi.ErrorResponse{Message: "not found"})
		return
	}
	if errors.Is(err, repository.ErrAlreadyExists) {
		writeJSON(w, http.StatusConflict, openapi.ErrorResponse{Message: "already exists"})
		return
	}
	writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: err.Error()})
}
