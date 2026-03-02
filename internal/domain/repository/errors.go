package repository

import "errors"

// ErrNotFound はリポジトリで対象エンティティが見つからない場合に返すセンチネルエラー。
var ErrNotFound = errors.New("entity not found")
