package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prevOut := log.Writer()
	prevFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(prevOut)
		log.SetFlags(prevFlags)
	}()
	fn()
	return buf.String()
}

func TestWarnOnCredentialedWildcardCORS(t *testing.T) {
	t.Run("exact origins emit no warning", func(t *testing.T) {
		out := captureLog(t, func() {
			warnOnCredentialedWildcardCORS([]string{"https://app.example.com", "http://localhost:3000"}, "lax")
		})
		if out != "" {
			t.Fatalf("expected no warning, got %q", out)
		}
	})

	t.Run("wildcard origin warns", func(t *testing.T) {
		out := captureLog(t, func() {
			warnOnCredentialedWildcardCORS([]string{"https://*.vercel.app"}, "lax")
		})
		if !strings.Contains(out, "wildcard origin") {
			t.Fatalf("expected wildcard warning, got %q", out)
		}
		if strings.Contains(out, "SameSite") {
			t.Fatalf("did not expect SameSite=none warning for lax, got %q", out)
		}
	})

	t.Run("wildcard origin with SameSite=none escalates", func(t *testing.T) {
		out := captureLog(t, func() {
			warnOnCredentialedWildcardCORS([]string{"https://*.vercel.app"}, "none")
		})
		if !strings.Contains(out, "wildcard origin") || !strings.Contains(out, "COOKIE_SAME_SITE=none") {
			t.Fatalf("expected both wildcard and SameSite=none warnings, got %q", out)
		}
	})
}
