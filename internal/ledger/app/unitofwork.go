package app

import (
	"FinFlow/internal/ledger/domain"
	"context"
)

type Repositories struct {
	Accounts     domain.AccountRepository
	Transactions domain.TransactionRepository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(r Repositories) error) error
}
