package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// Entry は企業への応募（インターン・本選考等）を表すエンティティ。
// サービスの中心的な管理単位であり、Task や StageHistory が紐づく。
type Entry struct {
	id        EntryID
	userID    UserID
	companyID CompanyID
	route     value.Route
	source    value.Source
	status    value.EntryStatus
	stage     value.Stage
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

// NewEntry は Entry を新規作成する。各値オブジェクトのバリデーションは呼び出し側で済んでいる前提。
func NewEntry(userID UserID, companyID CompanyID, route value.Route, source value.Source) *Entry {
	// 定数コンストラクタを使い、エラー握りつぶし（_, _）を回避する
	status := value.EntryStatusInProgress()
	stage := value.MustNewStage(value.StageKindApplication(), "応募")

	now := time.Now()
	return &Entry{
		id:        NewEntryID(),
		userID:    userID,
		companyID: companyID,
		route:     route,
		source:    source,
		status:    status,
		stage:     stage,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructEntry はDBから読み取ったデータでEntryを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructEntry(
	id EntryID, userID UserID, companyID CompanyID,
	route value.Route, source value.Source, status value.EntryStatus,
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

// ID は Entry の ID を返す。
func (e *Entry) ID() EntryID { return e.id }

// UserID は Entry を所有するユーザの ID を返す。
func (e *Entry) UserID() UserID { return e.userID }

// CompanyID は Entry に紐づく Company の ID を返す。
func (e *Entry) CompanyID() CompanyID { return e.companyID }

// Route は応募経路を返す。
func (e *Entry) Route() value.Route { return e.route }

// Source は流入元を返す。
func (e *Entry) Source() value.Source { return e.source }

// Status は Entry の選考ステータスを返す。
func (e *Entry) Status() value.EntryStatus { return e.status }

// Stage は Entry の選考フェーズを返す。
func (e *Entry) Stage() value.Stage { return e.stage }

// Memo は Entry のメモを返す。
func (e *Entry) Memo() string { return e.memo }

// CreatedAt は Entry の作成日時を返す。
func (e *Entry) CreatedAt() time.Time { return e.createdAt }

// UpdatedAt は Entry の更新日時を返す。
func (e *Entry) UpdatedAt() time.Time { return e.updatedAt }

// UpdateSource は流入元を更新し、UpdatedAt を現在時刻にする。
func (e *Entry) UpdateSource(source value.Source) {
	e.source = source
	e.updatedAt = time.Now()
}

// UpdateStage は選考フェーズを更新し、UpdatedAt を現在時刻にする。
func (e *Entry) UpdateStage(stage value.Stage) {
	e.stage = stage
	e.updatedAt = time.Now()
}

// UpdateStatus は選考ステータスを更新し、UpdatedAt を現在時刻にする。
func (e *Entry) UpdateStatus(status value.EntryStatus) {
	e.status = status
	e.updatedAt = time.Now()
}

// UpdateMemo は memo を更新し、UpdatedAt を現在時刻にする。
func (e *Entry) UpdateMemo(memo string) {
	e.memo = memo
	e.updatedAt = time.Now()
}
