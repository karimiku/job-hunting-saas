package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

var (
	ErrEntryRouteEmpty = errors.New("entry route must not be empty")
)

// Entry は企業への応募（インターン・本選考等）を表すエンティティ。
// サービスの中心的な管理単位であり、Task や StageHistory が紐づく。
type Entry struct {
	id        EntryID
	userID    UserID
	companyID CompanyID
	route     string
	source    value.Source
	status    value.EntryStatus
	stage     value.Stage
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

func NewEntry(userID UserID, companyID CompanyID, route string, source value.Source) (*Entry, error) {
	trimmed := strings.TrimSpace(route)
	if trimmed == "" {
		return nil, ErrEntryRouteEmpty
	}

	// 定数コンストラクタを使い、エラー握りつぶし（_, _）を回避する
	status := value.EntryStatusInProgress()
	stage := value.MustNewStage(value.StageKindApplication(), "応募")

	now := time.Now()
	return &Entry{
		id:        NewEntryID(),
		userID:    userID,
		companyID: companyID,
		route:     trimmed,
		source:    source,
		status:    status,
		stage:     stage,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructEntry はDBから読み取ったデータでEntryを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructEntry(
	id EntryID, userID UserID, companyID CompanyID,
	route string, source value.Source, status value.EntryStatus,
	stage value.Stage, memo string, createdAt, updatedAt time.Time,
) *Entry {
	return &Entry{
		id:        id,
		userID:    userID,
		companyID: companyID,
		route:     route,
		source:    source,
		status:    status,
		stage:     stage,
		memo:      memo,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (e *Entry) ID() EntryID              { return e.id }
func (e *Entry) UserID() UserID            { return e.userID }
func (e *Entry) CompanyID() CompanyID      { return e.companyID }
func (e *Entry) Route() string             { return e.route }
func (e *Entry) Source() value.Source       { return e.source }
func (e *Entry) Status() value.EntryStatus { return e.status }
func (e *Entry) Stage() value.Stage        { return e.stage }
func (e *Entry) Memo() string              { return e.memo }
func (e *Entry) CreatedAt() time.Time      { return e.createdAt }
func (e *Entry) UpdatedAt() time.Time      { return e.updatedAt }

func (e *Entry) UpdateSource(source value.Source) {
	e.source = source
	e.updatedAt = time.Now()
}

func (e *Entry) UpdateStage(stage value.Stage) {
	e.stage = stage
	e.updatedAt = time.Now()
}

func (e *Entry) UpdateStatus(status value.EntryStatus) {
	e.status = status
	e.updatedAt = time.Now()
}

func (e *Entry) UpdateMemo(memo string) {
	e.memo = memo
	e.updatedAt = time.Now()
}
