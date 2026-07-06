package memory

import (
	"FinFlow/internal/ledger/domain"
	"context"
	"sync"
)

type AccountRepository struct {
	accounts map[domain.AccountID]*domain.Account
	mu       sync.RWMutex
}

var _ domain.AccountRepository = (*AccountRepository)(nil)

func NewAccountRepository() *AccountRepository {
	var repo AccountRepository
	repo.accounts = make(map[domain.AccountID]*domain.Account)

	return &repo
}

func (a *AccountRepository) Save(ctx context.Context, acc *domain.Account) error {

	a.mu.Lock()
	defer a.mu.Unlock()

	a.accounts[acc.ID()] = acc
	return nil
}

func (a *AccountRepository) FindByID(ctx context.Context, id domain.AccountID) (*domain.Account, error) {

	a.mu.RLock()
	defer a.mu.RUnlock()

	acc, ok := a.accounts[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	clone := domain.ReconstituteAccount(acc.ID(), acc.UserID(), acc.Currency(), acc.Balance(), acc.Status(), acc.Version(), acc.Type(), acc.Created())

	return clone, nil
}

func (a *AccountRepository) FindSystemAccount(ctx context.Context, currency domain.Currency) (*domain.Account, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, v := range a.accounts {
		if v.Currency() == currency && v.Type() == domain.AccountSystemType {
			clone := domain.ReconstituteAccount(v.ID(), v.UserID(), v.Currency(), v.Balance(), v.Status(), v.Version(), v.Type(), v.Created())
			return clone, nil
		}
	}

	return nil, domain.ErrSystemAccountNotFound
}
