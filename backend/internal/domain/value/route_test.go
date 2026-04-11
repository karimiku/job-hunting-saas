package value

import (
	"errors"
	"testing"
)

func TestNewRoute(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"japanese route", "本選考", nil},
		{"internship", "インターン", nil},

		// 異常系
		{"empty", "", ErrRouteEmpty},
		{"whitespace only", " ", ErrRouteEmpty},
		{"leading space", " 本選考", ErrRouteInvalid},
		{"trailing space", "本選考 ", ErrRouteInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRoute(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewRoute(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestRoute_String(t *testing.T) {
	r, err := NewRoute("本選考")
	if err != nil {
		t.Fatalf("NewRoute failed: %v", err)
	}
	if r.String() != "本選考" {
		t.Errorf("String() = %q, want %q", r.String(), "本選考")
	}
}

func TestRoute_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "本選考", "本選考", true},
		{"different value", "本選考", "インターン", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routeA, err := NewRoute(tt.a)
			if err != nil {
				t.Fatalf("NewRoute(%q) failed: %v", tt.a, err)
			}
			routeB, err := NewRoute(tt.b)
			if err != nil {
				t.Fatalf("NewRoute(%q) failed: %v", tt.b, err)
			}
			if got := routeA.Equals(routeB); got != tt.want {
				t.Errorf("Route(%q).Equals(Route(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
