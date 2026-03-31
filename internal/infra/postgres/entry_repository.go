package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
)

type EntryRepository struct {
	q *sqlc.Queries
}

func NewEntryRepository(db sqlc.DBTX) *EntryRepository {
	return &EntryRepository{q: sqlc.New(db)}
}

func (r *EntryRepository) Save(ctx context.Context, entry *entity.Entry) error {
	if err := r.q.UpsertEntry(ctx, sqlc.UpsertEntryParams{
		ID:         uuid.UUID(entry.ID()),
		UserID:     uuid.UUID(entry.UserID()),
		CompanyID:  uuid.UUID(entry.CompanyID()),
		Route:      entry.Route().String(),
		Source:     entry.Source().String(),
		Status:     sqlc.EntryStatus(entry.Status().String()),
		StageKind:  sqlc.StageKind(entry.Stage().Kind().String()),
		StageLabel: entry.Stage().Label(),
		Memo:       entry.Memo(),
		CreatedAt:  pgtype.Timestamptz{Time: entry.CreatedAt(), Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: entry.UpdatedAt(), Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: UpsertEntry: %w", err)
	}
	return nil
}

func (r *EntryRepository) FindByID(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	row, err := r.q.FindEntryByID(ctx, sqlc.FindEntryByIDParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindEntryByID: %w", err)
	}

	return reconstructEntry(row)
}

func (r *EntryRepository) ListByUserID(ctx context.Context, userID entity.UserID, filter repository.EntryFilter) ([]*entity.Entry, error) {
	params := sqlc.ListEntriesByUserIDParams{
		UserID: uuid.UUID(userID),
	}

	if filter.Status != nil {
		params.Status = sqlc.NullEntryStatus{
			EntryStatus: sqlc.EntryStatus(filter.Status.String()),
			Valid:       true,
		}
	}
	if filter.StageKind != nil {
		params.StageKind = sqlc.NullStageKind{
			StageKind: sqlc.StageKind(filter.StageKind.String()),
			Valid:     true,
		}
	}
	if filter.Source != nil {
		params.Source = pgtype.Text{
			String: filter.Source.String(),
			Valid:  true,
		}
	}

	rows, err := r.q.ListEntriesByUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("postgres: ListEntriesByUserID: %w", err)
	}

	entries := make([]*entity.Entry, 0, len(rows))
	for _, row := range rows {
		e, err := reconstructEntry(row)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}

func (r *EntryRepository) Delete(ctx context.Context, userID entity.UserID, id entity.EntryID) error {
	n, err := r.q.DeleteEntry(ctx, sqlc.DeleteEntryParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		return fmt.Errorf("postgres: DeleteEntry: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructEntry(row sqlc.Entry) (*entity.Entry, error) {
	route, err := value.NewRoute(row.Route)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: entry route: %w", err)
	}

	source, err := value.NewSource(row.Source)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: entry source: %w", err)
	}

	status, err := value.NewEntryStatus(string(row.Status))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: entry status: %w", err)
	}

	stageKind, err := value.NewStageKind(string(row.StageKind))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: entry stage_kind: %w", err)
	}

	stage, err := value.NewStage(stageKind, row.StageLabel)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: entry stage: %w", err)
	}

	return entity.ReconstructEntry(
		entity.EntryID(row.ID),
		entity.UserID(row.UserID),
		entity.CompanyID(row.CompanyID),
		route,
		source,
		status,
		stage,
		row.Memo,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}
