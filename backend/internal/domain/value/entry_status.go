package value

import "errors"

// ErrEntryStatusEmpty は entry status が空文字のときに返されるエラー。
// ErrEntryStatusInvalid は entry status が未定義の値のときに返されるエラー。
var (
	ErrEntryStatusEmpty   = errors.New("entry status must not be empty")
	ErrEntryStatusInvalid = errors.New("entry status is invalid")
)

const (
	entryStatusInProgress = "in_progress"
	entryStatusOffered    = "offered"
	entryStatusAccepted   = "accepted"
	entryStatusRejected   = "rejected"
	entryStatusWithdrawn  = "withdrawn"
)

var validEntryStatuses = map[string]bool{
	entryStatusInProgress: true,
	entryStatusOffered:    true,
	entryStatusAccepted:   true,
	entryStatusRejected:   true,
	entryStatusWithdrawn:  true,
}

// EntryStatus は応募の進捗状態を表す値オブジェクト。
// in_progress / offered / accepted / rejected / withdrawn のいずれか。
type EntryStatus struct {
	value string
}

// NewEntryStatus は raw から EntryStatus を生成する。空文字や未定義値は対応するエラーを返す。
func NewEntryStatus(raw string) (EntryStatus, error) {
	if raw == "" {
		return EntryStatus{}, ErrEntryStatusEmpty
	}
	if !validEntryStatuses[raw] {
		return EntryStatus{}, ErrEntryStatusInvalid
	}
	return EntryStatus{value: raw}, nil
}

// String は entry status を文字列で返す。
func (s EntryStatus) String() string {
	return s.value
}

// Equals は 2 つの EntryStatus が等しいかを判定する。
func (s EntryStatus) Equals(other EntryStatus) bool {
	return s.value == other.value
}

// IsOpen は 進行中 (in_progress / offered) かを返す。
func (s EntryStatus) IsOpen() bool {
	return s.value == entryStatusInProgress || s.value == entryStatusOffered
}

// --- 定数コンストラクタ ---
// ハードコードされた既知の値に対して、エラーなしでインスタンスを返す。
// エンティティのファクトリ関数やメソッド内で `_, _ :=` のエラー握りつぶしを避けるために使う。

// EntryStatusInProgress は in_progress 状態の EntryStatus を返す。
func EntryStatusInProgress() EntryStatus { return EntryStatus{value: entryStatusInProgress} }

// EntryStatusOffered は offered 状態の EntryStatus を返す。
func EntryStatusOffered() EntryStatus { return EntryStatus{value: entryStatusOffered} }

// EntryStatusAccepted は accepted 状態の EntryStatus を返す。
func EntryStatusAccepted() EntryStatus { return EntryStatus{value: entryStatusAccepted} }

// EntryStatusRejected は rejected 状態の EntryStatus を返す。
func EntryStatusRejected() EntryStatus { return EntryStatus{value: entryStatusRejected} }

// EntryStatusWithdrawn は withdrawn 状態の EntryStatus を返す。
func EntryStatusWithdrawn() EntryStatus { return EntryStatus{value: entryStatusWithdrawn} }
