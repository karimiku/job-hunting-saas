package repository

import "errors"

// ErrNotFound はインフラ層の差異(DB, InMemory等)を隠蔽し、
// 上位層が統一的にHTTP 404等へ変換するための共通エラー。
var ErrNotFound = errors.New("entity not found")
