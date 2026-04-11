package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
)

// writeJSON はHTTPレスポンスをJSON形式で書き出す。
// 全handlerで共通して使う低レベルなレスポンス関数。
func writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
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
