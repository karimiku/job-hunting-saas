// Package devsession signs local-development session cookie values.
package devsession

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

const prefix = "dev_"

// Sign returns a signed dev session value for userID.
func Sign(userID entity.UserID, secret string) string {
	id := userID.String()
	return prefix + id + "." + signature(id, secret)
}

// Verify checks a dev session value and returns the embedded userID.
func Verify(raw string, secret string) (entity.UserID, bool) {
	if secret == "" || !strings.HasPrefix(raw, prefix) {
		return entity.UserID{}, false
	}
	payload := strings.TrimPrefix(raw, prefix)
	id, sig, ok := strings.Cut(payload, ".")
	if !ok || id == "" || sig == "" {
		return entity.UserID{}, false
	}
	if !hmac.Equal([]byte(sig), []byte(signature(id, secret))) {
		return entity.UserID{}, false
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return entity.UserID{}, false
	}
	return entity.UserID(parsed), true
}

func signature(id string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(id))
	return hex.EncodeToString(mac.Sum(nil))
}
