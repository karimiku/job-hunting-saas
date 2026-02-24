package value

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmailEmpty   = errors.New("email must not be empty")
	ErrEmailInvalid = errors.New("email format is invalid")

	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([.+_-][a-zA-Z0-9]+)*@[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
)

type Email struct {
	value string
}

func NewEmail(raw string) (Email, error) {
	if raw == "" {
		return Email{}, ErrEmailEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Email{}, ErrEmailInvalid
	}
	if !emailRegex.MatchString(raw) {
		return Email{}, ErrEmailInvalid
	}
	return Email{value: strings.ToLower(raw)}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
