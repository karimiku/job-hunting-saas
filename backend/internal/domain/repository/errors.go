package repository

import "errors"

// ErrNotFound はリポジトリで対象エンティティが見つからない場合に返すセンチネルエラー。
var ErrNotFound = errors.New("entity not found")

// ErrAlreadyExists はユニーク制約違反など、エンティティが既に存在する場合に返すセンチネルエラー。
var ErrAlreadyExists = errors.New("entity already exists")
