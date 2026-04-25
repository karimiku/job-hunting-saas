package value

import (
	"errors"
	"strings"
)

// ErrRouteEmpty は route が空文字のときに返されるエラー。
// ErrRouteInvalid は route の形式が不正なときに返されるエラー。
var (
	ErrRouteEmpty   = errors.New("route must not be empty")
	ErrRouteInvalid = errors.New("route format is invalid")
)

// Route は応募経路（本選考・インターン等）を表す値オブジェクト。
type Route struct {
	value string
}

// NewRoute は raw から Route を生成する。空文字や不正値は対応するエラーを返す。
func NewRoute(raw string) (Route, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Route{}, ErrRouteEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Route{}, ErrRouteInvalid
	}
	return Route{value: raw}, nil
}

// String は route を文字列で返す。
func (r Route) String() string {
	return r.value
}

// Equals は 2 つの Route が等しいかを判定する。
func (r Route) Equals(other Route) bool {
	return r.value == other.value
}
