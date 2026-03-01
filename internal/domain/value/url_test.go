package value

import (
	"errors"
	"testing"
)

func TestNewURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"basic", "https://example.com", nil},
		{"mynavi", "https://www.mynavi.jp/company/123", nil},
		{"rikunabi", "https://job.rikunabi.com/2027/company/r123/", nil},
		{"onecareer", "https://www.onecareer.jp/companies/12345", nil},
		{"i-web ats", "https://mypage.i-web.jpn.com/company2027/", nil},
		{"query params", "https://example.com/path?q=test&page=1", nil},
		{"fragment", "https://example.com/event#seminar", nil},
		{"encoded japanese", "https://example.com/search?q=%E5%B0%B1%E6%B4%BB", nil},
		{"deep path", "https://careers.company.co.jp/graduate/2027/mypage/login", nil},

		// 異常系
		{"empty", "", ErrURLEmpty},
		{"whitespace only", "   ", ErrURLEmpty},
		{"http", "http://example.com", ErrURLInvalid},
		{"ftp", "ftp://example.com", ErrURLInvalid},
		{"no scheme", "example.com", ErrURLInvalid},
		{"leading space", " https://example.com", ErrURLInvalid},
		{"trailing space", "https://example.com ", ErrURLInvalid},
		{"no host", "https://", ErrURLInvalid},
		{"no host with path", "https:///path", ErrURLInvalid},
		{"space in host", "https://exa mple.com", ErrURLInvalid},
		{"uppercase scheme", "HTTPS://example.com", ErrURLInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewURL(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestURL_String(t *testing.T) {
	u, err := NewURL("https://example.com/path?q=test")
	if err != nil {
		t.Fatalf("NewURL failed: %v", err)
	}
	if u.String() != "https://example.com/path?q=test" {
		t.Errorf("String() = %q, want %q", u.String(), "https://example.com/path?q=test")
	}
}

func TestURL_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same url", "https://example.com", "https://example.com", true},
		{"different path", "https://example.com/a", "https://example.com/b", false},
		{"different query", "https://example.com?a=1", "https://example.com?a=2", false},
		{"with and without slash", "https://example.com", "https://example.com/", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlA, err := NewURL(tt.a)
			if err != nil {
				t.Fatalf("NewURL(%q) failed: %v", tt.a, err)
			}
			urlB, err := NewURL(tt.b)
			if err != nil {
				t.Fatalf("NewURL(%q) failed: %v", tt.b, err)
			}
			if got := urlA.Equals(urlB); got != tt.want {
				t.Errorf("URL(%q).Equals(URL(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
