package value

import (
	"testing"
)

func TestNewURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// 正常系
		{"basic", "https://example.com", false},
		{"mynavi", "https://www.mynavi.jp/company/123", false},
		{"rikunabi", "https://job.rikunabi.com/2027/company/r123/", false},
		{"onecareer", "https://www.onecareer.jp/companies/12345", false},
		{"i-web ats", "https://mypage.i-web.jpn.com/company2027/", false},
		{"query params", "https://example.com/path?q=test&page=1", false},
		{"fragment", "https://example.com/event#seminar", false},
		{"encoded japanese", "https://example.com/search?q=%E5%B0%B1%E6%B4%BB", false},
		{"deep path", "https://careers.company.co.jp/graduate/2027/mypage/login", false},

		// 異常系
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"http", "http://example.com", true},
		{"ftp", "ftp://example.com", true},
		{"no scheme", "example.com", true},
		{"leading space", " https://example.com", true},
		{"trailing space", "https://example.com ", true},
		{"no host", "https://", true},
		{"no host with path", "https:///path", true},
		{"space in host", "https://exa mple.com", true},
		{"uppercase scheme", "HTTPS://example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewURL(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("NewURL(%q) should return error, but got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("NewURL(%q) should succeed, but got error: %v", tt.input, err)
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
