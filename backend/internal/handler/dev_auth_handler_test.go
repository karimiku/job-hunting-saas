package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsLocalDevRequestTrustsOnlyHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		origin  string
		referer string
		want    bool
	}{
		{
			name: "loopback host allowed",
			host: "localhost:3000",
			want: true,
		},
		{
			name: "loopback IPv4 host allowed",
			host: "127.0.0.1:8080",
			want: true,
		},
		{
			name: "loopback IPv6 host allowed",
			host: "[::1]:8080",
			want: true,
		},
		{
			name: "public host is rejected",
			host: "api.example.com",
			want: false,
		},
		{
			name:    "spoofed loopback Origin does not bypass a public host",
			host:    "api.example.com",
			origin:  "http://localhost:3000",
			referer: "http://localhost:3000/login",
			want:    false,
		},
		{
			name:    "spoofed loopback Referer does not bypass a public host",
			host:    "victim.vercel.app",
			referer: "http://127.0.0.1/",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/dev/session", nil)
			req.Host = tt.host
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}
			if got := isLocalDevRequest(req); got != tt.want {
				t.Fatalf("isLocalDevRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
