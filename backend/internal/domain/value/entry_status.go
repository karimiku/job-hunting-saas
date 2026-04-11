package value

import "errors"

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

func NewEntryStatus(raw string) (EntryStatus, error) {
	if raw == "" {
		return EntryStatus{}, ErrEntryStatusEmpty
	}
	if !validEntryStatuses[raw] {
		return EntryStatus{}, ErrEntryStatusInvalid
	}
	return EntryStatus{value: raw}, nil
}

func (s EntryStatus) String() string {
	return s.value
}

func (s EntryStatus) Equals(other EntryStatus) bool {
	return s.value == other.value
}

func (s EntryStatus) IsOpen() bool {
	return s.value == entryStatusInProgress || s.value == entryStatusOffered
}

// --- 定数コンストラクタ ---
// ハードコードされた既知の値に対して、エラーなしでインスタンスを返す。
// エンティティのファクトリ関数やメソッド内で `_, _ :=` のエラー握りつぶしを避けるために使う。

func EntryStatusInProgress() EntryStatus { return EntryStatus{value: entryStatusInProgress} }
func EntryStatusOffered() EntryStatus    { return EntryStatus{value: entryStatusOffered} }
func EntryStatusAccepted() EntryStatus   { return EntryStatus{value: entryStatusAccepted} }
func EntryStatusRejected() EntryStatus   { return EntryStatus{value: entryStatusRejected} }
func EntryStatusWithdrawn() EntryStatus  { return EntryStatus{value: entryStatusWithdrawn} }
