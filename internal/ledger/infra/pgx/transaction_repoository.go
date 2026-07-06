package pgx

import (
	"FinFlow/internal/ledger/domain"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type TransactionRepository struct {
	db dbtx
}

func NewTransactionRepository(db dbtx) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Save(ctx context.Context, txn *domain.Transaction) error {

	txID := string(txn.ID())

	_, err := r.db.Exec(ctx,
		`INSERT INTO transactions (id, key, description, created_at)                                                   
                 VALUES ($1, $2, $3, $4)`,
		txID, nullString(string(txn.Key())), txn.Description(), txn.CreatedAt())
	if err != nil {
		return err
	}

	for _, p := range txn.Postings() {
		_, err = r.db.Exec(ctx,
			`INSERT INTO postings (transaction_id, account_id, amount, currency)                                   
                         VALUES ($1, $2, $3, $4)`,
			txID, string(p.AccountID()), p.Money().Amount(), p.Money().Currency().Code())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, id domain.TransactionID) (*domain.Transaction, error) {

	var (
		txID        domain.TransactionID
		key         *string
		description string
		createdAt   time.Time
	)

	err := r.db.QueryRow(ctx, `SELECT  id, key, description, created_at FROM transactions WHERE id=$1`, string(id)).Scan(
		&txID, &key, &description, &createdAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTransactionNotFound
	}
	if err != nil {
		return nil, err
	}

	postings, err := r.loadPostings(ctx, string(txID))
	if err != nil {
		return nil, err
	}

	txn := domain.ReconstituteTransaction(txID, postings, createdAt, nullKey(key), description)

	return txn, nil
}

func (r *TransactionRepository) FindByIdempotencyKey(ctx context.Context, key domain.IdempotencyKey) (*domain.Transaction, error) {

	var (
		txID        domain.TransactionID
		txKey       *string
		description string
		createdAt   time.Time
	)

	err := r.db.QueryRow(ctx, `SELECT  id, key, description, created_at FROM transactions WHERE key=$1`, string(key)).Scan(
		&txID, &txKey, &description, &createdAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrTransactionNotFound
	}
	if err != nil {
		return nil, err
	}

	postings, err := r.loadPostings(ctx, string(txID))
	if err != nil {
		return nil, err
	}

	txn := domain.ReconstituteTransaction(txID, postings, createdAt, nullKey(txKey), description)

	return txn, nil
}

func (r *TransactionRepository) loadPostings(ctx context.Context, txID string) ([]domain.Posting, error) {

	rows, err := r.db.Query(ctx,
		`SELECT account_id, amount, currency FROM postings WHERE transaction_id = $1 ORDER BY id`,
		txID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postings []domain.Posting
	for rows.Next() {
		var (
			accountID string
			amount    int64
			currency  string
		)
		if err := rows.Scan(&accountID, &amount, &currency); err != nil {
			return nil, err
		}
		ccy, err := domain.NewCurrency(currency)
		if err != nil {
			return nil, err
		}
		postings = append(postings,
			domain.ReconstitutePosting(domain.AccountID(accountID), domain.NewMoney(amount, ccy)))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return postings, nil
}

func nullKey(key *string) domain.IdempotencyKey {

	if key == nil {
		return ""
	}

	return domain.IdempotencyKey(*key)
}
