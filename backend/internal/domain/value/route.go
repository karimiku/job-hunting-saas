package value

import (
	"errors"
	"strings"
)

var (
	ErrRouteEmpty   = errors.New("route must not be empty")
	ErrRouteInvalid = errors.New("route format is invalid")
)

// Route は応募経路（本選考・インターン等）を表す値オブジェクト。
type Route struct {
	value string
}

func NewRoute(raw string) (Route, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Route{}, ErrRouteEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Route{}, ErrRouteInvalid
	}
	return Route{value: raw}, nil
}

func (r Route) String() string {
	return r.value
}

func (r Route) Equals(other Route) bool {
	return r.value == other.value
}
