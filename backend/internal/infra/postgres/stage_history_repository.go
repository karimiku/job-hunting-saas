package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
)

// StageHistoryRepository は StageHistoryRepository インターフェースの PostgreSQL 実装。
type StageHistoryRepository struct {
	q *sqlc.Queries
}

func NewStageHistoryRepository(db sqlc.DBTX) *StageHistoryRepository {
	return &StageHistoryRepository{q: sqlc.New(db)}
}

func (r *StageHistoryRepository) Create(ctx context.Context, history *entity.StageHistory) error {
	if err := r.q.CreateStageHistory(ctx, sqlc.CreateStageHistoryParams{
		ID:         uuid.UUID(history.ID()),
		EntryID:    uuid.UUID(history.EntryID()),
		StageKind:  sqlc.StageKind(history.Stage().Kind().String()),
		StageLabel: history.Stage().Label(),
		Note:       history.Note(),
		CreatedAt:  pgtype.Timestamptz{Time: history.CreatedAt(), Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: CreateStageHistory: %w", err)
	}
	return nil
}

func (r *StageHistoryRepository) ListByEntryID(ctx context.Context, entryID entity.EntryID) ([]*entity.StageHistory, error) {
	rows, err := r.q.ListStageHistoriesByEntryID(ctx, uuid.UUID(entryID))
	if err != nil {
		return nil, fmt.Errorf("postgres: ListStageHistoriesByEntryID: %w", err)
	}

	histories := make([]*entity.StageHistory, 0, len(rows))
	for _, row := range rows {
		h, err := reconstructStageHistory(row)
		if err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}

	return histories, nil
}

func reconstructStageHistory(row sqlc.StageHistory) (*entity.StageHistory, error) {
	stageKind, err := value.NewStageKind(string(row.StageKind))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: stage_history stage_kind: %w", err)
	}

	stage, err := value.NewStage(stageKind, row.StageLabel)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: stage_history stage: %w", err)
	}

	return entity.ReconstructStageHistory(
		entity.StageHistoryID(row.ID),
		entity.EntryID(row.EntryID),
		stage,
		row.Note,
		row.CreatedAt.Time,
	), nil
}
