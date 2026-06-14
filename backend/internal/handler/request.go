package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
)

const maxDefaultJSONBodyBytes = 256 * 1024

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) bool {
	return decodeJSONBodyInternal(w, r, dst, maxBytes, false)
}

func decodeOptionalJSONBody(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) bool {
	return decodeJSONBodyInternal(w, r, dst, maxBytes, true)
}

func decodeJSONBodyInternal(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64, allowEmpty bool) bool {
	if maxBytes <= 0 {
		maxBytes = maxDefaultJSONBodyBytes
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		if allowEmpty && errors.Is(err, io.EOF) {
			return true
		}
		writeDecodeError(w, err)
		return false
	}

	var extra any
	if err := decoder.Decode(&extra); err != nil && !errors.Is(err, io.EOF) {
		writeDecodeError(w, err)
		return false
	} else if err == nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return false
	}

	return true
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		writeJSON(w, http.StatusRequestEntityTooLarge, openapi.ErrorResponse{Message: "request body too large"})
		return
	}
	writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
}
