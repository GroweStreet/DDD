package memory

import (
	"FinFlow/internal/ledger/domain"
	"context"
	"sync"
)

type TransactionRepository struct {
	byID  map[domain.TransactionID]*domain.Transaction
	byKey map[domain.IdempotencyKey]*domain.Transaction
	mu    sync.RWMutex
}

var _ domain.TransactionRepository = (*TransactionRepository)(nil)

func NewTransactionRepository() *TransactionRepository {

	var repo TransactionRepository

	repo.byID = make(map[domain.TransactionID]*domain.Transaction)
	repo.byKey = make(map[domain.IdempotencyKey]*domain.Transaction)

	return &repo
}

func (t *TransactionRepository) Save(ctx context.Context, tx *domain.Transaction) error {

	t.mu.Lock()
	defer t.mu.Unlock()

	t.byID[tx.ID()] = tx
	if tx.Key() != "" {
		t.byKey[tx.Key()] = tx
	}

	return nil
}

func (t *TransactionRepository) FindByID(ctx context.Context, id domain.TransactionID) (*domain.Transaction, error) {

	t.mu.RLock()
	defer t.mu.RUnlock()

	tx, ok := t.byID[id]
	if !ok {
		return nil, domain.ErrTransactionNotFound
	}

	return tx, nil
}

func (t *TransactionRepository) FindByIdempotencyKey(ctx context.Context, key domain.IdempotencyKey) (*domain.Transaction, error) {

	t.mu.RLock()
	defer t.mu.RUnlock()

	tx, ok := t.byKey[key]
	if !ok {
		return nil, domain.ErrTransactionNotFound
	}

	return tx, nil
}
