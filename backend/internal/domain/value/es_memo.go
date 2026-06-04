package value

import (
	"errors"
	"strings"
)

const (
	maxESMemoCategoryLen = 64
	maxESMemoTitleLen    = 120
	maxESMemoSourceLen   = 64
)

// ESMemoCategory は ES / 自己PR / 面接ネタ用メモの分類を表す。
// LLMやMCPクライアントがカテゴリを増やせるよう、列挙型ではなく制約付き文字列にしている。
type ESMemoCategory struct {
	value string
}

func NewESMemoCategory(raw string) (ESMemoCategory, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		v = "general"
	}
	if len([]rune(v)) > maxESMemoCategoryLen {
		return ESMemoCategory{}, errors.New("es memo category must be 64 characters or less")
	}
	return ESMemoCategory{value: v}, nil
}

func (c ESMemoCategory) String() string { return c.value }

// ESMemoTitle は ESメモの見出しを表す。
type ESMemoTitle struct {
	value string
}

func NewESMemoTitle(raw string) (ESMemoTitle, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ESMemoTitle{}, errors.New("es memo title is required")
	}
	if len([]rune(v)) > maxESMemoTitleLen {
		return ESMemoTitle{}, errors.New("es memo title must be 120 characters or less")
	}
	return ESMemoTitle{value: v}, nil
}

func (t ESMemoTitle) String() string { return t.value }

// ESMemoContent は ESメモ本文を表す。
type ESMemoContent struct {
	value string
}

func NewESMemoContent(raw string) (ESMemoContent, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return ESMemoContent{}, errors.New("es memo content is required")
	}
	return ESMemoContent{value: v}, nil
}

func (c ESMemoContent) String() string { return c.value }

// ESMemoSource はメモの入力元を表す。
type ESMemoSource struct {
	value string
}

func NewESMemoSource(raw string) (ESMemoSource, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		v = "mcp"
	}
	if len([]rune(v)) > maxESMemoSourceLen {
		return ESMemoSource{}, errors.New("es memo source must be 64 characters or less")
	}
	return ESMemoSource{value: v}, nil
}

func (s ESMemoSource) String() string { return s.value }
