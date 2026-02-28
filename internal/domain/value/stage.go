package value

import (
	"errors"
	"strings"
)

var (
	ErrStageKindEmpty   = errors.New("stage kind must not be empty")
	ErrStageKindInvalid = errors.New("stage kind is invalid")
	ErrStageLabelEmpty  = errors.New("stage label must not be empty")
	ErrStageLabelInvalid = errors.New("stage label format is invalid")
)

const (
	StageKindApplication = "application"
	StageKindDocument    = "document"
	StageKindTest        = "test"
	StageKindInterview   = "interview"
	StageKindGroup       = "group"
	StageKindOffer       = "offer"
	StageKindOther       = "other"
)

var validStageKinds = map[string]bool{
	StageKindApplication: true,
	StageKindDocument:    true,
	StageKindTest:        true,
	StageKindInterview:   true,
	StageKindGroup:       true,
	StageKindOffer:       true,
	StageKindOther:       true,
}

// Stage は選考フェーズを表す値オブジェクト。
// kind（応募・ES・テスト・面接等のカテゴリ）と label（表示名）の組で構成される。
type Stage struct {
	kind  string
	label string
}

func NewStage(kind string, label string) (Stage, error) {
	if kind == "" {
		return Stage{}, ErrStageKindEmpty
	}
	if !validStageKinds[kind] {
		return Stage{}, ErrStageKindInvalid
	}
	if label == "" {
		return Stage{}, ErrStageLabelEmpty
	}
	if label != strings.TrimSpace(label) {
		return Stage{}, ErrStageLabelInvalid
	}
	return Stage{kind: kind, label: label}, nil
}

func (s Stage) Kind() string {
	return s.kind
}

func (s Stage) Label() string {
	return s.label
}

func (s Stage) Equals(other Stage) bool {
	return s.kind == other.kind && s.label == other.label
}

// MustNewStage は NewStage のパニック版。
// ハードコードされた既知の値に対して使う。
// 不正な値が渡された場合はプログラマのバグなのでパニックする。
func MustNewStage(kind string, label string) Stage {
	s, err := NewStage(kind, label)
	if err != nil {
		panic("invalid stage: " + err.Error())
	}
	return s
}
