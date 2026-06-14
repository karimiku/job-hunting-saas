package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// EntryWithCompanyRepository は Company と Entry を同一 PostgreSQL transaction で保存する。
type EntryWithCompanyRepository struct {
	pool *pgxpool.Pool
}

// NewEntryWithCompanyRepository は EntryWithCompanyRepository を新規生成する。
func NewEntryWithCompanyRepository(pool *pgxpool.Pool) *EntryWithCompanyRepository {
	return &EntryWithCompanyRepository{pool: pool}
}

// SaveEntryWithCompany は Company 作成と Entry 作成を同一 transaction に閉じ込める。
func (r *EntryWithCompanyRepository) SaveEntryWithCompany(ctx context.Context, company *entity.Company, entry *entity.Entry) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("postgres: begin SaveEntryWithCompany tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := NewCompanyRepository(tx).Save(ctx, company); err != nil {
		return fmt.Errorf("postgres: SaveEntryWithCompany company: %w", err)
	}
	if err := NewEntryRepository(tx).Save(ctx, entry); err != nil {
		return fmt.Errorf("postgres: SaveEntryWithCompany entry: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("postgres: commit SaveEntryWithCompany tx: %w", err)
	}
	return nil
}
