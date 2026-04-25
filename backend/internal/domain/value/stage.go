package value

import (
	"errors"
	"strings"
)

// ErrStageKindEmpty は stage kind が空文字のときに返されるエラー。
// ErrStageKindInvalid は stage kind が未定義の値のときに返されるエラー。
// ErrStageLabelEmpty は stage label が空文字のときに返されるエラー。
// ErrStageLabelInvalid は stage label の形式が不正なときに返されるエラー。
var (
	ErrStageKindEmpty    = errors.New("stage kind must not be empty")
	ErrStageKindInvalid  = errors.New("stage kind is invalid")
	ErrStageLabelEmpty   = errors.New("stage label must not be empty")
	ErrStageLabelInvalid = errors.New("stage label format is invalid")
)

const (
	stageKindApplication = "application"
	stageKindDocument    = "document"
	stageKindTest        = "test"
	stageKindInterview   = "interview"
	stageKindGroup       = "group"
	stageKindOffer       = "offer"
	stageKindOther       = "other"
)

var validStageKinds = map[string]bool{
	stageKindApplication: true,
	stageKindDocument:    true,
	stageKindTest:        true,
	stageKindInterview:   true,
	stageKindGroup:       true,
	stageKindOffer:       true,
	stageKindOther:       true,
}

// StageKind は選考フェーズの種別を表す値オブジェクト。
// application / document / test / interview / group / offer / other のいずれか。
type StageKind struct {
	value string
}

// NewStageKind は文字列から StageKind を生成する。
// 空文字列や未定義の値が渡された場合はエラーを返す。
func NewStageKind(raw string) (StageKind, error) {
	if raw == "" {
		return StageKind{}, ErrStageKindEmpty
	}
	if !validStageKinds[raw] {
		return StageKind{}, ErrStageKindInvalid
	}
	return StageKind{value: raw}, nil
}

// String は stage kind を文字列で返す。
func (k StageKind) String() string {
	return k.value
}

// Equals は 2 つの StageKind が等しいかを判定する。
func (k StageKind) Equals(other StageKind) bool {
	return k.value == other.value
}

// --- 定数コンストラクタ ---
// ハードコードされた既知の値に対して、エラーなしでインスタンスを返す。

// StageKindApplication は application 種別の StageKind を返す。
func StageKindApplication() StageKind { return StageKind{value: stageKindApplication} }

// StageKindDocument は document 種別の StageKind を返す。
func StageKindDocument() StageKind { return StageKind{value: stageKindDocument} }

// StageKindTest は test 種別の StageKind を返す。
func StageKindTest() StageKind { return StageKind{value: stageKindTest} }

// StageKindInterview は interview 種別の StageKind を返す。
func StageKindInterview() StageKind { return StageKind{value: stageKindInterview} }

// StageKindGroup は group 種別の StageKind を返す。
func StageKindGroup() StageKind { return StageKind{value: stageKindGroup} }

// StageKindOffer は offer 種別の StageKind を返す。
func StageKindOffer() StageKind { return StageKind{value: stageKindOffer} }

// StageKindOther は other 種別の StageKind を返す。
func StageKindOther() StageKind { return StageKind{value: stageKindOther} }

// Stage は選考フェーズを表す値オブジェクト。
// kind（応募・ES・テスト・面接等のカテゴリ）と label（表示名）の組で構成される。
type Stage struct {
	kind  StageKind
	label string
}

// NewStage は kind と label から Stage を生成する。空文字や不正な label は対応するエラーを返す。
func NewStage(kind StageKind, label string) (Stage, error) {
	if label == "" {
		return Stage{}, ErrStageLabelEmpty
	}
	if label != strings.TrimSpace(label) {
		return Stage{}, ErrStageLabelInvalid
	}
	return Stage{kind: kind, label: label}, nil
}

// Kind は Stage の種別 (StageKind) を返す。
func (s Stage) Kind() StageKind {
	return s.kind
}

// Label は Stage の表示名を返す。
func (s Stage) Label() string {
	return s.label
}

// Equals は 2 つの Stage が等しいかを判定する。
func (s Stage) Equals(other Stage) bool {
	return s.kind.Equals(other.kind) && s.label == other.label
}

// MustNewStage は NewStage のパニック版。
// ハードコードされた既知の値に対して使う。
// 不正な値が渡された場合はプログラマのバグなのでパニックする。
func MustNewStage(kind StageKind, label string) Stage {
	s, err := NewStage(kind, label)
	if err != nil {
		panic("invalid stage: " + err.Error())
	}
	return s
}
