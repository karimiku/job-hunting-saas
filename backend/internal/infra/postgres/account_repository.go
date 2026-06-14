package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// AccountRepository はアカウント単位の PostgreSQL 削除を提供する。
type AccountRepository struct {
	pool *pgxpool.Pool
}

// NewAccountRepository は AccountRepository を新規生成する。
func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{pool: pool}
}

// DeleteAccount は users 削除を transaction 内で実行し、関連データは FK cascade に委ねる。
func (r *AccountRepository) DeleteAccount(ctx context.Context, userID entity.UserID) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres: begin DeleteAccount tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := NewUserRepository(tx).Delete(ctx, userID); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("postgres: commit DeleteAccount tx: %w", err)
	}
	return nil
}
