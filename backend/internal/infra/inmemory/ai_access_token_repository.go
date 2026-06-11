package inmemory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// AIAccessTokenRepository はメモリ上に保存する AI 連携トークン用リポジトリ。
type AIAccessTokenRepository struct {
	mu         sync.RWMutex
	tokensByID map[entity.AIAccessTokenID]*entity.AIAccessToken
}

// NewAIAccessTokenRepository は AIAccessTokenRepository を新規生成する。
func NewAIAccessTokenRepository() *AIAccessTokenRepository {
	return &AIAccessTokenRepository{
		tokensByID: make(map[entity.AIAccessTokenID]*entity.AIAccessToken),
	}
}

// Create はトークンを保存する。
func (r *AIAccessTokenRepository) Create(_ context.Context, token *entity.AIAccessToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.tokensByID {
		if existing.TokenHash().String() == token.TokenHash().String() {
			return repository.ErrAlreadyExists
		}
	}
	r.tokensByID[token.ID()] = token
	return nil
}

// ListByUserID は userID 所有のトークンを作成日時の新しい順で返す。
func (r *AIAccessTokenRepository) ListByUserID(_ context.Context, userID entity.UserID) ([]*entity.AIAccessToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	owned := make([]*entity.AIAccessToken, 0)
	for _, t := range r.tokensByID {
		if t.UserID() == userID {
			owned = append(owned, t)
		}
	}
	sort.Slice(owned, func(i, j int) bool {
		return owned[i].CreatedAt().After(owned[j].CreatedAt())
	})
	return owned, nil
}

// FindActiveByHash は未失効トークンを hash から取得する。
func (r *AIAccessTokenRepository) FindActiveByHash(_ context.Context, hash value.AIAccessTokenHash) (*entity.AIAccessToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tokensByID {
		if t.TokenHash().String() == hash.String() && !t.IsRevoked() {
			return t, nil
		}
	}
	return nil, repository.ErrNotFound
}

// Revoke は userID 所有のトークンを失効する。
func (r *AIAccessTokenRepository) Revoke(_ context.Context, userID entity.UserID, id entity.AIAccessTokenID, revokedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	token, ok := r.tokensByID[id]
	if !ok || token.UserID() != userID || token.IsRevoked() {
		return repository.ErrNotFound
	}
	r.tokensByID[id] = reconstructAIAccessTokenWithUsage(token, token.LastUsedAt(), &revokedAt)
	return nil
}

// TouchLastUsed はトークンの最終利用日時を更新する。
func (r *AIAccessTokenRepository) TouchLastUsed(_ context.Context, id entity.AIAccessTokenID, usedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	token, ok := r.tokensByID[id]
	if !ok || token.IsRevoked() {
		return repository.ErrNotFound
	}
	r.tokensByID[id] = reconstructAIAccessTokenWithUsage(token, &usedAt, token.RevokedAt())
	return nil
}

func reconstructAIAccessTokenWithUsage(
	token *entity.AIAccessToken,
	lastUsedAt *time.Time,
	revokedAt *time.Time,
) *entity.AIAccessToken {
	return entity.ReconstructAIAccessToken(
		token.ID(),
		token.UserID(),
		token.Name(),
		token.TokenHash(),
		token.Prefix(),
		token.CreatedAt(),
		lastUsedAt,
		revokedAt,
	)
}
