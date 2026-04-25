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

// UserRepository は UserRepository インターフェースの PostgreSQL 実装。
type UserRepository struct {
	q *sqlc.Queries
}

// NewUserRepository は UserRepository を新規生成する。db には pgxpool.Pool もしくは tx を渡す。
func NewUserRepository(db sqlc.DBTX) *UserRepository {
	return &UserRepository{q: sqlc.New(db)}
}

// Save は User を upsert する。email の一意制約に違反する場合は repository.ErrAlreadyExists を返す。
func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	if err := r.q.UpsertUser(ctx, sqlc.UpsertUserParams{
		ID:        uuid.UUID(user.ID()),
		Email:     user.Email().String(),
		Name:      user.Name().String(),
		CreatedAt: pgtype.Timestamptz{Time: user.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: user.UpdatedAt(), Valid: true},
	}); err != nil {
		if isUniqueViolation(err) {
			return repository.ErrAlreadyExists
		}
		return fmt.Errorf("postgres: UpsertUser: %w", err)
	}
	return nil
}

// FindByID は ID から User を取得する。存在しない場合は repository.ErrNotFound を返す。
func (r *UserRepository) FindByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	row, err := r.q.FindUserByID(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindUserByID: %w", err)
	}
	return reconstructUser(row)
}

// FindByEmail は email から User を取得する。存在しない場合は repository.ErrNotFound を返す。
func (r *UserRepository) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	row, err := r.q.FindUserByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindUserByEmail: %w", err)
	}
	return reconstructUser(row)
}

// Delete は ID から User を削除する。存在しない場合は repository.ErrNotFound を返す。
func (r *UserRepository) Delete(ctx context.Context, id entity.UserID) error {
	n, err := r.q.DeleteUser(ctx, uuid.UUID(id))
	if err != nil {
		return fmt.Errorf("postgres: DeleteUser: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructUser(row sqlc.User) (*entity.User, error) {
	email, err := value.NewEmail(row.Email)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: user email: %w", err)
	}
	name, err := value.NewUserName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: user name: %w", err)
	}
	return entity.ReconstructUser(
		entity.UserID(row.ID),
		email,
		name,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}
