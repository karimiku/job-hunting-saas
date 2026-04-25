// Package inmemory は Repository インターフェースの sync.Map ベース実装を提供する。
// 開発環境とユニットテストで DB を立ち上げずに動作確認するために使う。
package inmemory

import (
	"context"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// CompanyRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// 本番ではPostgreSQL実装に差し替える。
// 本番実装と同じ userID スコープを守り、テストが認可前提をすり抜けないようにする。
type CompanyRepository struct {
	mu            sync.RWMutex
	companiesByID map[entity.CompanyID]*entity.Company
}

func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{
		companiesByID: make(map[entity.CompanyID]*entity.Company),
	}
}

func (r *CompanyRepository) Save(_ context.Context, company *entity.Company) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.companiesByID[company.ID()] = company
	return nil
}

func (r *CompanyRepository) FindByID(_ context.Context, userID entity.UserID, companyID entity.CompanyID) (*entity.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	storedCompany, exists := r.companiesByID[companyID]
	if !exists || storedCompany.UserID() != userID {
		return nil, repository.ErrNotFound
	}
	return storedCompany, nil
}

func (r *CompanyRepository) ListByUserID(_ context.Context, userID entity.UserID) ([]*entity.Company, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var ownedCompanies []*entity.Company
	for _, company := range r.companiesByID {
		if company.UserID() == userID {
			ownedCompanies = append(ownedCompanies, company)
		}
	}
	return ownedCompanies, nil
}

func (r *CompanyRepository) Delete(_ context.Context, userID entity.UserID, companyID entity.CompanyID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storedCompany, exists := r.companiesByID[companyID]
	if !exists || storedCompany.UserID() != userID {
		return repository.ErrNotFound
	}
	delete(r.companiesByID, companyID)
	return nil
}
