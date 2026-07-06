package pgx

import (
	"FinFlow/internal/ledger/app"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UnitOfWork struct {
	pool *pgxpool.Pool
}

func NewUnitOfWork(pool *pgxpool.Pool) *UnitOfWork {
	return &UnitOfWork{pool: pool}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(r app.Repositories) error) error {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	repos := app.Repositories{
		Accounts:     NewAccountRepository(tx),
		Transactions: NewTransactionRepository(tx),
	}

	if err = fn(repos); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
