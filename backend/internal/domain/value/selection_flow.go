package value

import (
	"errors"
	"strings"
)

var (
	// ErrSelectionFlowSourceEmpty は source が空文字のときに返されるエラー。
	ErrSelectionFlowSourceEmpty = errors.New("selection flow source must not be empty")
	// ErrSelectionFlowSourceInvalid は source が未定義値のときに返されるエラー。
	ErrSelectionFlowSourceInvalid = errors.New("selection flow source is invalid")
	// ErrSelectionStagePosition は stage position が1未満のときに返されるエラー。
	ErrSelectionStagePosition = errors.New("selection stage position must be positive")
	// ErrSelectionConfidenceInvalid は confidence が0-100の範囲外のときに返されるエラー。
	ErrSelectionConfidenceInvalid = errors.New("selection flow confidence must be between 0 and 100")
)

const (
	selectionFlowSourceTemplate = "template"
	selectionFlowSourceManual   = "manual"
	selectionFlowSourceAIInbox  = "ai_inbox"
	selectionFlowSourceAIPaste  = "ai_paste"
)

var validSelectionFlowSources = map[string]bool{
	selectionFlowSourceTemplate: true,
	selectionFlowSourceManual:   true,
	selectionFlowSourceAIInbox:  true,
	selectionFlowSourceAIPaste:  true,
}

// SelectionFlowSource は選考フローがどの入口から作られたかを表す値オブジェクト。
type SelectionFlowSource struct {
	value string
}

// NewSelectionFlowSource は raw から SelectionFlowSource を生成する。
func NewSelectionFlowSource(raw string) (SelectionFlowSource, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return SelectionFlowSource{}, ErrSelectionFlowSourceEmpty
	}
	if !validSelectionFlowSources[raw] {
		return SelectionFlowSource{}, ErrSelectionFlowSourceInvalid
	}
	return SelectionFlowSource{value: raw}, nil
}

// String は source を文字列で返す。
func (s SelectionFlowSource) String() string { return s.value }

// Equals は 2 つの SelectionFlowSource が等しいかを判定する。
func (s SelectionFlowSource) Equals(other SelectionFlowSource) bool {
	return s.value == other.value
}

// SelectionFlowSourceTemplate は標準テンプレート由来の source を返す。
func SelectionFlowSourceTemplate() SelectionFlowSource {
	return SelectionFlowSource{value: selectionFlowSourceTemplate}
}

// SelectionFlowSourceManual はユーザー手入力由来の source を返す。
func SelectionFlowSourceManual() SelectionFlowSource {
	return SelectionFlowSource{value: selectionFlowSourceManual}
}

// SelectionFlowSourceAIInbox は Inbox をAI解析した source を返す。
func SelectionFlowSourceAIInbox() SelectionFlowSource {
	return SelectionFlowSource{value: selectionFlowSourceAIInbox}
}

// SelectionFlowSourceAIPaste は貼り付け本文をAI解析した source を返す。
func SelectionFlowSourceAIPaste() SelectionFlowSource {
	return SelectionFlowSource{value: selectionFlowSourceAIPaste}
}

// NewSelectionStagePosition は position が1以上であることを検証する。
func NewSelectionStagePosition(position int) (int, error) {
	if position <= 0 {
		return 0, ErrSelectionStagePosition
	}
	return position, nil
}

// NewSelectionConfidence は 0-100 の信頼度を検証する。nil は未設定として許容する。
func NewSelectionConfidence(confidence *int) (*int, error) {
	if confidence == nil {
		return nil, nil
	}
	if *confidence < 0 || *confidence > 100 {
		return nil, ErrSelectionConfidenceInvalid
	}
	value := *confidence
	return &value, nil
}
