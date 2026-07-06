package domain

import "context"

type AccountRepository interface {
	Save(ctx context.Context, acc *Account) error
	FindByID(ctx context.Context, id AccountID) (*Account, error)
	FindSystemAccount(ctx context.Context, currency Currency) (*Account, error)
}

type TransactionRepository interface {
	Save(ctx context.Context, tx *Transaction) error
	FindByID(ctx context.Context, id TransactionID) (*Transaction, error)
	FindByIdempotencyKey(ctx context.Context, key IdempotencyKey) (*Transaction, error)
}

type IDGenerator interface {
	Generate() TransactionID
}
