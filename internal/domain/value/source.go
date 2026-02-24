package value

import (
	"errors"
	"strings"
)

var (
	ErrSourceEmpty   = errors.New("source must not be empty")
	ErrSourceInvalid = errors.New("source format is invalid")
)

type Source struct {
	value string
}

func NewSource(raw string) (Source, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Source{}, ErrSourceEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Source{}, ErrSourceInvalid
	}
	return Source{value: raw}, nil
}

func (s Source) String() string {
	return s.value
}

func (s Source) Equals(other Source) bool {
	return s.value == other.value
}
