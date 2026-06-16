package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
)

// SelectionFlowRepository は SelectionFlowRepository インターフェースの PostgreSQL 実装。
type SelectionFlowRepository struct {
	q *sqlc.Queries
}

// NewSelectionFlowRepository は SelectionFlowRepository を新規生成する。
func NewSelectionFlowRepository(db sqlc.DBTX) *SelectionFlowRepository {
	return &SelectionFlowRepository{q: sqlc.New(db)}
}

// Upsert はEntryごとの選考フローを置換保存する。
func (r *SelectionFlowRepository) Upsert(ctx context.Context, flow *entity.SelectionFlow) (*entity.SelectionFlow, error) {
	currentStagePosition, err := intToInt32(flow.CurrentStagePosition())
	if err != nil {
		return nil, fmt.Errorf("postgres: selection flow current stage position: %w", err)
	}
	confidence, err := intPtrToPgInt4(flow.Confidence())
	if err != nil {
		return nil, fmt.Errorf("postgres: selection flow confidence: %w", err)
	}
	row, err := r.q.UpsertSelectionFlow(ctx, sqlc.UpsertSelectionFlowParams{
		ID:                   uuid.UUID(flow.ID()),
		EntryID:              uuid.UUID(flow.EntryID()),
		Source:               flow.Source().String(),
		CurrentStagePosition: currentStagePosition,
		Confidence:           confidence,
		InboxClipID:          inboxClipIDToPgUUID(flow.InboxClipID()),
		CreatedAt:            pgtype.Timestamptz{Time: flow.CreatedAt(), Valid: true},
		UpdatedAt:            pgtype.Timestamptz{Time: flow.UpdatedAt(), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: UpsertSelectionFlow: %w", err)
	}
	if err := r.q.DeleteSelectionStagesByFlowID(ctx, row.ID); err != nil {
		return nil, fmt.Errorf("postgres: DeleteSelectionStagesByFlowID: %w", err)
	}
	for _, stage := range flow.Stages() {
		position, err := intToInt32(stage.Position())
		if err != nil {
			return nil, fmt.Errorf("postgres: selection stage position: %w", err)
		}
		if err := r.q.CreateSelectionStage(ctx, sqlc.CreateSelectionStageParams{
			ID:           uuid.UUID(stage.ID()),
			FlowID:       row.ID,
			Position:     position,
			StageKind:    sqlc.StageKind(stage.Stage().Kind().String()),
			StageLabel:   stage.Stage().Label(),
			EvidenceText: stage.EvidenceText(),
			CreatedAt:    pgtype.Timestamptz{Time: stage.CreatedAt(), Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("postgres: CreateSelectionStage: %w", err)
		}
	}
	return r.reconstruct(ctx, row)
}

// FindByEntryID は userID 所有のEntryに紐づく選考フローを返す。
func (r *SelectionFlowRepository) FindByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) (*entity.SelectionFlow, error) {
	row, err := r.q.FindSelectionFlowByEntryID(ctx, sqlc.FindSelectionFlowByEntryIDParams{
		UserID:  uuid.UUID(userID),
		EntryID: uuid.UUID(entryID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindSelectionFlowByEntryID: %w", err)
	}
	return r.reconstruct(ctx, row)
}

func (r *SelectionFlowRepository) reconstruct(ctx context.Context, row sqlc.SelectionFlow) (*entity.SelectionFlow, error) {
	stageRows, err := r.q.ListSelectionStagesByFlowID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("postgres: ListSelectionStagesByFlowID: %w", err)
	}
	stages := make([]*entity.SelectionStage, 0, len(stageRows))
	for _, stageRow := range stageRows {
		stage, err := reconstructSelectionStage(stageRow)
		if err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}
	source, err := value.NewSelectionFlowSource(row.Source)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: selection_flow source: %w", err)
	}
	return entity.ReconstructSelectionFlow(
		entity.SelectionFlowID(row.ID),
		entity.EntryID(row.EntryID),
		source,
		int(row.CurrentStagePosition),
		pgInt4ToIntPtr(row.Confidence),
		pgUUIDToInboxClipIDPtr(row.InboxClipID),
		stages,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}

func reconstructSelectionStage(row sqlc.SelectionStage) (*entity.SelectionStage, error) {
	stageKind, err := value.NewStageKind(string(row.StageKind))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: selection_stage kind: %w", err)
	}
	stage, err := value.NewStage(stageKind, row.StageLabel)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: selection_stage stage: %w", err)
	}
	return entity.ReconstructSelectionStage(
		entity.SelectionStageID(row.ID),
		entity.SelectionFlowID(row.FlowID),
		int(row.Position),
		stage,
		row.EvidenceText,
		row.CreatedAt.Time,
	), nil
}

func intPtrToPgInt4(value *int) (pgtype.Int4, error) {
	if value == nil {
		return pgtype.Int4{}, nil
	}
	out, err := intToInt32(*value)
	if err != nil {
		return pgtype.Int4{}, err
	}
	return pgtype.Int4{Int32: out, Valid: true}, nil
}

func pgInt4ToIntPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	out := int(value.Int32)
	return &out
}

func inboxClipIDToPgUUID(id *entity.InboxClipID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: uuid.UUID(*id), Valid: true}
}

func pgUUIDToInboxClipIDPtr(value pgtype.UUID) *entity.InboxClipID {
	if !value.Valid {
		return nil
	}
	id := entity.InboxClipID(value.Bytes)
	return &id
}

func intToInt32(value int) (int32, error) {
	if value < math.MinInt32 || value > math.MaxInt32 {
		return 0, fmt.Errorf("value %d overflows int32", value)
	}
	return int32(value), nil
}
